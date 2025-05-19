package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/juju/juju/api/client/application"
	"github.com/openhdc/otterscale/internal/domain/model"
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

func Test_filterPersistentVolumeClaim(t *testing.T) {
	pvcs := []corev1.PersistentVolumeClaim{
		{ObjectMeta: metav1.ObjectMeta{Name: "pvc1"}, Spec: corev1.PersistentVolumeClaimSpec{StorageClassName: strPtr("sc1")}},
	}
	vs := []corev1.Volume{
		{Name: "v1", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "pvc1"}}},
	}
	scm := map[string]storagev1.StorageClass{"sc1": {ObjectMeta: metav1.ObjectMeta{Name: "sc1"}}}
	namespace := "metadata"
	out := filterPersistentVolumeClaim(pvcs, vs, namespace, scm)
	if len(out) != 1 || out[0].StorageClass == nil || out[0].StorageClass.Name != "sc1" {
		t.Errorf("unexpected output: %+v", out)
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
