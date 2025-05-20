package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/juju/juju/api/client/application"
	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_extractKubeControl(t *testing.T) {
	unitInfo := &model.UnitInfo{
		RelationData: []application.EndpointRelationData{
			{
				Endpoint: "kube-control",
				UnitRelationData: map[string]application.RelationData{
					"unit-1": {},
				},
			},
		},
	}
	name, err := extractKubeControl(unitInfo)
	if err != nil || name != "unit-1" {
		t.Errorf("expected unit-1, got %v, err %v", name, err)
	}
	_, err = extractKubeControl(&model.UnitInfo{})
	if err == nil {
		t.Error("expected error for missing kube-control")
	}
}

func Test_extractEndpoint(t *testing.T) {
	endpoints := []string{"https://1.2.3.4"}
	endpointsJSON, _ := json.Marshal(endpoints)
	unitInfo := &model.UnitInfo{
		RelationData: []application.EndpointRelationData{
			{
				UnitRelationData: map[string]application.RelationData{
					"u": {
						UnitData: map[string]interface{}{
							"api-endpoints": string(endpointsJSON),
						},
					},
				},
			},
		},
	}
	ep, err := extractEndpoint(unitInfo)
	if err != nil || ep != "https://1.2.3.4" {
		t.Errorf("expected endpoint, got %v, err %v", ep, err)
	}
	_, err = extractEndpoint(&model.UnitInfo{})
	if err == nil {
		t.Error("expected error for missing endpoint")
	}
}

func Test_extractClientToken(t *testing.T) {
	creds := map[string]model.ControlPlaneCredential{
		"foo": {ClientToken: "tok"},
	}
	credsJSON, _ := json.Marshal(creds)
	unitInfo := &model.UnitInfo{
		RelationData: []application.EndpointRelationData{
			{
				UnitRelationData: map[string]application.RelationData{
					"u": {
						UnitData: map[string]interface{}{
							"creds": string(credsJSON),
						},
					},
				},
			},
		},
	}
	tok, err := extractClientToken(unitInfo)
	if err != nil || tok != "tok" {
		t.Errorf("expected tok, got %v, err %v", tok, err)
	}
	_, err = extractClientToken(&model.UnitInfo{})
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func Test_toStorageClassMap(t *testing.T) {
	scs := []storagev1.StorageClass{
		{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
		{ObjectMeta: metav1.ObjectMeta{Name: "bar"}},
	}
	m := toStorageClassMap(scs)
	if len(m) != 2 || m["foo"].Name != "foo" || m["bar"].Name != "bar" {
		t.Errorf("unexpected map: %v", m)
	}
}

func Test_filterServices(t *testing.T) {
	svcs := []corev1.Service{
		{Spec: corev1.ServiceSpec{Selector: map[string]string{"app": "a"}}},
	}
	sel := metav1.SetAsLabelSelector(map[string]string{"app": "a"})
	selector, _ := metav1.LabelSelectorAsSelector(sel)
	namespace := "metadata"
	out := filterServices(svcs, namespace, selector)
	if len(out) != 1 {
		t.Errorf("expected 1, got %d", len(out))
	}
}

func Test_filterPods(t *testing.T) {
	pods := []corev1.Pod{
		{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "a"}}},
	}
	sel := metav1.SetAsLabelSelector(map[string]string{"app": "a"})
	selector, _ := metav1.LabelSelectorAsSelector(sel)
	namespace := "metadata"
	out := filterPods(pods, namespace, selector)
	if len(out) != 1 {
		t.Errorf("expected 1, got %d", len(out))
	}
}

func Test_randomName(t *testing.T) {
	name := randomName()
	if name == "" || name != fmt.Sprintf("%v", name) {
		t.Errorf("unexpected random name: %v", name)
	}
}

func Test_isKeyNotFoundError(t *testing.T) {
	err := &fakeStatusError{code: http.StatusNotFound}
	if !isKeyNotFoundError(err) {
		t.Error("expected true for not found error")
	}
	if isKeyNotFoundError(errors.New("other")) {
		t.Error("expected false for non-status error")
	}
}

type fakeStatusError struct{ code int }

func (f *fakeStatusError) Error() string { return "status error" }
func (f *fakeStatusError) Status() *metav1.Status {
	return &metav1.Status{Code: int32(f.code)}
}

func strPtr(s string) *string { return &s }

func Test_toApplication(t *testing.T) {
	sel := metav1.SetAsLabelSelector(map[string]string{"app": "a"})
	_, err1 := metav1.LabelSelectorAsSelector(sel)
	if err1 != nil {
		fmt.Printf("failed to create selector: %v", err1)
	}
	app, err := toApplication(sel, "Deployment", "foo", "ns", &metav1.ObjectMeta{Name: "foo"}, map[string]string{"app": "a"}, str32Ptr(1), []corev1.Container{}, []corev1.Volume{}, []corev1.Service{}, []corev1.Pod{}, []corev1.PersistentVolumeClaim{}, map[string]storagev1.StorageClass{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.Name != "foo" || app.Type != "Deployment" {
		t.Errorf("unexpected app: %+v", app)
	}
	// error path
	_, err = toApplication(&metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{{Key: "a", Operator: "invalid"}}}, "Deployment", "foo", "ns", &metav1.ObjectMeta{}, nil, nil, nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for invalid selector")
	}
}

func str32Ptr(i int32) *int32 { return &i }

func Test_newKubernetesConfig(t *testing.T) {
	unitInfo := &model.UnitInfo{
		RelationData: []application.EndpointRelationData{
			{
				UnitRelationData: map[string]application.RelationData{
					"u": {
						UnitData: map[string]interface{}{
							"api-endpoints": `["https://1.2.3.4"]`,
							"creds":         `{"foo":{"client-token":"tok"}}`,
						},
					},
				},
			},
		},
	}
	cfg, err := newKubernetesConfig(unitInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "https://1.2.3.4" || cfg.BearerToken != "tok" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestNexusService_ListApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApps := mocks.NewMockKubeApps(ctrl)
	mockCore := mocks.NewMockKubeCore(ctrl)
	mockStorage := mocks.NewMockKubeStorage(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)

	ns := &NexusService{
		apps:       mockApps,
		core:       mockCore,
		storage:    mockStorage,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"

	// Set up mocks for setKubernetesClient
	mockKube.EXPECT().Exists(uuid, facility).Return(true)

	// Set up mocks for ListDeployments, ListStatefulSets, ListDaemonSets
	mockApps.EXPECT().ListDeployments(ctx, uuid, facility, "").Return([]appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: "ns1"},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "a"}},
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{Name: "c1"}},
						Volumes:    []corev1.Volume{},
					},
				},
			},
		},
	}, nil)
	mockApps.EXPECT().ListStatefulSets(ctx, uuid, facility, "").Return([]appsv1.StatefulSet{}, nil)
	mockApps.EXPECT().ListDaemonSets(ctx, uuid, facility, "").Return([]appsv1.DaemonSet{}, nil)

	// Set up mocks for ListServices, ListPods, ListPersistentVolumeClaims, ListStorageClasses
	mockCore.EXPECT().ListServices(ctx, uuid, facility, "").Return([]corev1.Service{}, nil)
	mockCore.EXPECT().ListPods(ctx, uuid, facility, "").Return([]corev1.Pod{}, nil)
	mockCore.EXPECT().ListPersistentVolumeClaims(ctx, uuid, facility, "").Return([]corev1.PersistentVolumeClaim{}, nil)
	mockStorage.EXPECT().ListStorageClasses(ctx, uuid, facility).Return([]storagev1.StorageClass{}, nil)

	apps, err := ns.ListApplications(ctx, uuid, facility)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("expected 1 application, got %d", len(apps))
	}
	if apps[0].Name != "dep1" || apps[0].Type != "Deployment" {
		t.Errorf("unexpected application: %+v", apps[0])
	}
}

func TestNexusService_GetApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApps := mocks.NewMockKubeApps(ctrl)
	mockCore := mocks.NewMockKubeCore(ctrl)
	mockStorage := mocks.NewMockKubeStorage(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)

	ns := &NexusService{
		apps:       mockApps,
		core:       mockCore,
		storage:    mockStorage,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "ns1"
	name := "dep1"

	// Set up mocks for setKubernetesClient
	mockKube.EXPECT().Exists(uuid, facility).Return(true)

	// Happy path: Deployment found
	mockApps.EXPECT().GetDeployment(ctx, uuid, facility, namespace, name).Return(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "a"}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "c1"}},
					Volumes:    []corev1.Volume{},
				},
			},
		},
	}, nil)
	mockApps.EXPECT().GetStatefulSet(ctx, uuid, facility, namespace, name).Return(nil, nil)
	mockApps.EXPECT().GetDaemonSet(ctx, uuid, facility, namespace, name).Return(nil, nil)
	mockCore.EXPECT().ListServices(ctx, uuid, facility, namespace).Return([]corev1.Service{}, nil)
	mockCore.EXPECT().ListPods(ctx, uuid, facility, namespace).Return([]corev1.Pod{}, nil)
	mockCore.EXPECT().ListPersistentVolumeClaims(ctx, uuid, facility, namespace).Return([]corev1.PersistentVolumeClaim{}, nil)
	mockStorage.EXPECT().ListStorageClasses(ctx, uuid, facility).Return([]storagev1.StorageClass{}, nil)

	app, err := ns.GetApplication(ctx, uuid, facility, namespace, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil || app.Name != name || app.Type != "Deployment" {
		t.Errorf("unexpected application: %+v", app)
	}

	// Error path: All resources not found
	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	mockApps.EXPECT().GetDeployment(ctx, uuid, facility, "ns2", "notfound").Return(nil, nil)
	mockApps.EXPECT().GetStatefulSet(ctx, uuid, facility, "ns2", "notfound").Return(nil, nil)
	mockApps.EXPECT().GetDaemonSet(ctx, uuid, facility, "ns2", "notfound").Return(nil, nil)
	mockCore.EXPECT().ListServices(ctx, uuid, facility, "ns2").Return([]corev1.Service{}, nil)
	mockCore.EXPECT().ListPods(ctx, uuid, facility, "ns2").Return([]corev1.Pod{}, nil)
	mockCore.EXPECT().ListPersistentVolumeClaims(ctx, uuid, facility, "ns2").Return([]corev1.PersistentVolumeClaim{}, nil)
	mockStorage.EXPECT().ListStorageClasses(ctx, uuid, facility).Return([]storagev1.StorageClass{}, nil)

	_, err = ns.GetApplication(ctx, uuid, facility, "ns2", "notfound")
	if err == nil {
		t.Error("expected error for not found application")
	}
}

func TestNexusService_GetChartMetadataFromApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	ns := &NexusService{
		helm: mockHelm,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	app := &model.Application{
		Namespace: "ns1",
		Labels:    map[string]string{"app.kubernetes.io/instance": "rel1"},
	}

	// Happy path: GetValues returns values, ShowChart not called
	mockHelm.EXPECT().GetValues(uuid, facility, "ns1", "rel1").Return(map[string]any{"foo": "bar"}, nil)

	md, err := ns.GetChartMetadataFromApplication(ctx, uuid, facility, app)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md == nil || md.ValuesYAML == "" {
		t.Errorf("expected values YAML in metadata, got: %+v", md)
	}
}

func TestNexusService_CreateRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)
	ns := &NexusService{
		helm:       mockHelm,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "ns1"
	name := "rel1"
	chartRef := "mychart"
	valuesYAML := "foo: bar"
	valuesMap := map[string]string{"baz": "qux"}

	// Kubernetes client already exists
	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	mockHelm.EXPECT().InstallRelease(uuid, facility, namespace, name, false, chartRef, gomock.Any()).Return(model.Release{}.Release, nil)

	rel, err := ns.CreateRelease(ctx, uuid, facility, namespace, name, false, chartRef, valuesYAML, valuesMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rel == nil {
		t.Error("expected release, got nil")
	}

	// Error path: invalid YAML
	_, err = ns.CreateRelease(ctx, uuid, facility, namespace, name, false, chartRef, ":", valuesMap)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestNexusService_UpdateRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)
	ns := &NexusService{
		helm:       mockHelm,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "ns1"
	name := "rel1"
	chartRef := "mychart"
	valuesYAML := "foo: bar"

	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	release := model.Release{}
	mockHelm.EXPECT().UpgradeRelease(uuid, facility, namespace, name, false, chartRef, gomock.Any()).Return(release.Release, nil)

	rel, err := ns.UpdateRelease(ctx, uuid, facility, namespace, name, false, chartRef, valuesYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rel == nil {
		t.Error("expected release, got nil")
	}

	// Error path: invalid YAML
	_, err = ns.UpdateRelease(ctx, uuid, facility, namespace, name, false, chartRef, ":")
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestNexusService_DeleteRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)
	ns := &NexusService{
		helm:       mockHelm,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "ns1"
	name := "rel1"

	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	mockHelm.EXPECT().UninstallRelease(uuid, facility, namespace, name, false).Return(model.Release{}.Release, nil)

	err := ns.DeleteRelease(ctx, uuid, facility, namespace, name, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestNexusService_RollbackRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)
	ns := &NexusService{
		helm:       mockHelm,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "ns1"
	name := "rel1"

	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	mockHelm.EXPECT().RollbackRelease(uuid, facility, namespace, name, false).Return(nil)

	err := ns.RollbackRelease(ctx, uuid, facility, namespace, name, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNexusService_ListCharts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	ns := &NexusService{
		helm: mockHelm,
	}

	ctx := context.Background()
	mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{
		{
			APIVersion: "APIVersion_test",
		},
	}, nil)

	charts, err := ns.ListCharts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(charts) == 0 || charts[0].Name != "foo" {
		t.Errorf("expected chart foo, got %+v", charts)
	}
}

func TestNexusService_GetChart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	ns := &NexusService{
		helm: mockHelm,
	}

	ctx := context.Background()
	mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{
		{
			APIVersion: "APIVersion_test",
		},
	}, nil)

	chart, err := ns.GetChart(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chart == nil || chart.Name != "foo" {
		t.Errorf("expected chart foo, got %+v", chart)
	}
}

func TestNexusService_GetChartMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	ns := &NexusService{
		helm: mockHelm,
	}

	ctx := context.Background()
	mockHelm.EXPECT().ShowChart("foo", gomock.Any()).Return("readme content", nil)

	md, err := ns.GetChartMetadata(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md == nil || md.ReadmeMD != "readme content" {
		t.Errorf("expected readme content, got %+v", md)
	}
}

func Test_filterPersistentVolumeClaim(t *testing.T) {
	pvcs := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pvc1", Namespace: "ns"},
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: strPtr("sc1"),
			},
		},
	}
	vs := []corev1.Volume{
		{
			Name: "v1",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "pvc1"},
			},
		},
	}
	scm := map[string]storagev1.StorageClass{
		"sc1": {ObjectMeta: metav1.ObjectMeta{Name: "sc1"}},
	}
	namespace := "ns"
	result := filterPersistentVolumeClaim(pvcs, vs, namespace, scm)
	if len(result) != 1 || result[0].StorageClass == nil || result[0].StorageClass.Name != "sc1" {
		t.Errorf("unexpected result: %+v", result)
	}

	// Test with no matching PVC
	result = filterPersistentVolumeClaim([]corev1.PersistentVolumeClaim{}, vs, namespace, scm)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %+v", result)
	}
}

func TestNexusService_setKubernetesClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKube := mocks.NewMockKubeClient(ctrl)
	ns := &NexusService{kubernetes: mockKube}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"

	// Already exists
	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	err := ns.setKubernetesClient(ctx, uuid, facility)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNexusService_fromDeployment(t *testing.T) {
	ns := &NexusService{}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: "ns"},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "a"}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "c1"}},
					Volumes:    []corev1.Volume{},
				},
			},
		},
	}
	svcs := []corev1.Service{}
	pods := []corev1.Pod{}
	pvcs := []corev1.PersistentVolumeClaim{}
	scm := map[string]storagev1.StorageClass{}

	app, err := ns.fromDeployment(dep, svcs, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil || app.Name != "dep1" || app.Type != "Deployment" {
		t.Errorf("unexpected application: %+v", app)
	}
}

func TestNexusService_fromStatefulSet(t *testing.T) {
	ns := &NexusService{}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "sts1", Namespace: "ns"},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "a"}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "c1"}},
					Volumes:    []corev1.Volume{},
				},
			},
		},
	}
	svcs := []corev1.Service{}
	pods := []corev1.Pod{}
	pvcs := []corev1.PersistentVolumeClaim{}
	scm := map[string]storagev1.StorageClass{}

	app, err := ns.fromStatefulSet(sts, svcs, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil || app.Name != "sts1" || app.Type != "StatefulSet" {
		t.Errorf("unexpected application: %+v", app)
	}
}

func TestNexusService_fromDaemonSet(t *testing.T) {
	ns := &NexusService{}
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: "ds1", Namespace: "ns"},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "a"}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "c1"}},
					Volumes:    []corev1.Volume{},
				},
			},
		},
	}
	svcs := []corev1.Service{}
	pods := []corev1.Pod{}
	pvcs := []corev1.PersistentVolumeClaim{}
	scm := map[string]storagev1.StorageClass{}

	app, err := ns.fromDaemonSet(ds, svcs, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil || app.Name != "ds1" || app.Type != "DaemonSet" {
		t.Errorf("unexpected application: %+v", app)
	}
}
