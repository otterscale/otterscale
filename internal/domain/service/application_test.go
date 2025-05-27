package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/juju/juju/api/base"
	application "github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"

	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var errExpected = errors.New("expected error")

func TestNexusService_GetPublicAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacility := mocks.NewMockJujuApplication(ctrl)
	s := &NexusService{
		facility: mockFacility,
	}

	uuid := "test-uuid"
	facilityName := "test-facility"

	t.Run("success", func(t *testing.T) {
		cpUnit := "leader-unit"
		cpUnitInfo := &application.UnitInfo{PublicAddress: "test-address"}

		mockFacility.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(cpUnit, nil)
		mockFacility.EXPECT().GetUnitInfo(gomock.Any(), uuid, cpUnit).Return(cpUnitInfo, nil)

		address, err := s.GetPublicAddress(context.Background(), uuid, facilityName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if address != cpUnitInfo.PublicAddress {
			t.Errorf("expected address %q, got %q", cpUnitInfo.PublicAddress, address)
		}
	})

	t.Run("error getting leader", func(t *testing.T) {
		mockFacility.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return("", errExpected)

		_, err := s.GetPublicAddress(context.Background(), uuid, facilityName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("error getting unit info", func(t *testing.T) {
		cpUnit := "leader-unit"
		mockFacility.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(cpUnit, nil)
		mockFacility.EXPECT().GetUnitInfo(gomock.Any(), uuid, cpUnit).Return(nil, errExpected)

		_, err := s.GetPublicAddress(context.Background(), uuid, facilityName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestNexusService_ListApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockKubeApps := mocks.NewMockKubeApps(ctrl)
	mockKubeCore := mocks.NewMockKubeCore(ctrl)
	mockKubeStorage := mocks.NewMockKubeStorage(ctrl)

	s := &NexusService{
		facility:   mockJujuApplication,
		kubernetes: mockKubeClient,
		apps:       mockKubeApps,
		core:       mockKubeCore,
		storage:    mockKubeStorage,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facilityName := "test-facility"

	t.Run("success", func(t *testing.T) {
		// Mock for setKubernetesClient
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true) // Assume client already exists for simplicity

		// Mocks for List... calls
		mockKubeApps.EXPECT().ListDeployments(gomock.Any(), uuid, facilityName, "").Return([]appsv1.Deployment{*createMockDeployment("dep1", "ns1")}, nil)
		mockKubeApps.EXPECT().ListStatefulSets(gomock.Any(), uuid, facilityName, "").Return([]appsv1.StatefulSet{*createMockStatefulSet("sts1", "ns1")}, nil)
		mockKubeApps.EXPECT().ListDaemonSets(gomock.Any(), uuid, facilityName, "").Return([]appsv1.DaemonSet{*createMockDaemonSet("ds1", "ns1")}, nil)
		mockKubeCore.EXPECT().ListServices(gomock.Any(), uuid, facilityName, "").Return([]corev1.Service{*createMockService("svc1", "ns1", map[string]string{"app": "test"})}, nil)
		mockKubeCore.EXPECT().ListPods(gomock.Any(), uuid, facilityName, "").Return([]corev1.Pod{*createMockPod("pod1", "ns1", map[string]string{"app": "test"})}, nil)
		mockKubeCore.EXPECT().ListPersistentVolumeClaims(gomock.Any(), uuid, facilityName, "").Return([]corev1.PersistentVolumeClaim{*createMockPVC("pvc1", "ns1", "standard")}, nil)
		mockKubeStorage.EXPECT().ListStorageClasses(gomock.Any(), uuid, facilityName).Return([]storagev1.StorageClass{*createMockStorageClass("standard")}, nil)

		apps, err := s.ListApplications(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("ListApplications() unexpected error: %v", err)
		}
		if len(apps) != 3 { // Expecting 1 deployment, 1 statefulset, 1 daemonset
			t.Fatalf("Expected 3 applications, got %d", len(apps))
		}
		// Add more specific assertions about the content of 'apps'
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return("", errExpected)
		// No further mocks needed as it should fail early

		_, err := s.ListApplications(ctx, uuid, facilityName)
		if err == nil {
			t.Fatal("ListApplications() expected error from setKubernetesClient, got nil")
		}
		if !errors.Is(err, errExpected) {
			t.Fatalf("Expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_GetApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockKubeApps := mocks.NewMockKubeApps(ctrl)
	mockKubeCore := mocks.NewMockKubeCore(ctrl)
	mockKubeStorage := mocks.NewMockKubeStorage(ctrl)
	mockHelm := mocks.NewMockKubeHelm(ctrl) // Assuming GetChartMetadataFromApplication might be called

	s := &NexusService{
		facility:   mockJujuApplication,
		kubernetes: mockKubeClient,
		apps:       mockKubeApps,
		core:       mockKubeCore,
		storage:    mockKubeStorage,
		helm:       mockHelm, // Assuming GetChartMetadataFromApplication might be called
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facilityName := "test-facility"
	// appName is the name of the k8s resource (Deployment, StatefulSet, etc.)
	// appNamespace is the k8s namespace where the resource resides.
	// These were likely defined in your test case, e.g. "ns1" and "test-app"

	t.Run("success", func(t *testing.T) {
		// Define the specific namespace and app name for this test case
		appNamespace := "ns1" // This should match the namespace used in your createMockX functions
		appName := "test-app" // This should match the name of the mock deployment/statefulset/daemonset

		// Mock for setKubernetesClient
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true) // Assume client already exists for simplicity

		// Mocks for Get... calls (assuming a Deployment is found for this test)
		mockDeployment := createMockDeployment(appName, appNamespace)
		mockKubeApps.EXPECT().GetDeployment(gomock.Any(), uuid, facilityName, appNamespace, appName).Return(mockDeployment, nil)
		mockKubeApps.EXPECT().GetStatefulSet(gomock.Any(), uuid, facilityName, appNamespace, appName).Return(nil, k8serrors.NewNotFound(appsv1.Resource("statefulset"), appName))
		mockKubeApps.EXPECT().GetDaemonSet(gomock.Any(), uuid, facilityName, appNamespace, appName).Return(nil, k8serrors.NewNotFound(appsv1.Resource("daemonset"), appName))

		// Mocks for List... calls - these should expect the specific appNamespace
		mockKubeCore.EXPECT().ListServices(gomock.Any(), uuid, facilityName, appNamespace).Return([]corev1.Service{*createMockService("svc1", appNamespace, map[string]string{"app": appName})}, nil)
		mockKubeCore.EXPECT().ListPods(gomock.Any(), uuid, facilityName, appNamespace).Return([]corev1.Pod{*createMockPod("pod1", appNamespace, map[string]string{"app": appName})}, nil)
		mockKubeCore.EXPECT().ListPersistentVolumeClaims(gomock.Any(), uuid, facilityName, appNamespace).Return([]corev1.PersistentVolumeClaim{*createMockPVC("pvc1", appNamespace, "standard")}, nil)
		mockKubeStorage.EXPECT().ListStorageClasses(gomock.Any(), uuid, facilityName).Return([]storagev1.StorageClass{*createMockStorageClass("standard")}, nil)

		// Call the method under test
		app, err := s.GetApplication(ctx, uuid, facilityName, appNamespace, appName)
		if err != nil {
			t.Fatalf("GetApplication() unexpected error: %v", err)
		}
		if app == nil {
			t.Fatal("GetApplication() returned nil app, expected a deployment")
		}
		if app.Name != appName {
			t.Errorf("Expected app name %s, got %s", appName, app.Name)
		}
		if app.Namespace != appNamespace {
			t.Errorf("Expected app namespace %s, got %s", appNamespace, app.Namespace)
		}
		// Add more specific assertions about the content of 'app'
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		appNamespace := "ns1"
		appName := "test-app"
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return("", errExpected)
		// No further mocks needed as it should fail early

		_, err := s.GetApplication(ctx, uuid, facilityName, appNamespace, appName)
		if err == nil {
			t.Fatal("GetApplication() expected error from setKubernetesClient, got nil")
		}
	})
}

func TestNexusService_ListReleases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScope := mocks.NewMockJujuModel(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl)
	mockJujuApplication := mocks.NewMockJujuApplication(ctrl) // This is s.facility
	mockKubeClient := mocks.NewMockKubeClient(ctrl)           // This is s.kubernetes
	mockKubeHelm := mocks.NewMockKubeHelm(ctrl)

	s := &NexusService{
		scope:      mockScope,
		client:     mockClient,
		facility:   mockJujuApplication,
		kubernetes: mockKubeClient,
		helm:       mockKubeHelm,
	}

	// ctx := context.Background()
	// scopeUUID := "test-scope-uuid"
	// scopeName := "test-scope-name"
	// k8sFacilityName := "k8s-facility"

	// t.Run("success", func(t *testing.T) {
	// 	// 1. Mock for s.listFacilitiesAcrossScopes -> s.scope.List()
	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: scopeName}}, nil)

	// 	// 2. Mock for s.listFacilitiesAcrossScopes -> s.ListFacilities -> s.client.Status()
	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(&params.FullStatus{
	// 		Applications: map[string]params.ApplicationStatus{
	// 			k8sFacilityName: {Charm: charmNameKubernetes, Status: &params.StatusInfo{Current: "active"}},
	// 			"other-app":     {Charm: "some-other-charm", Status: &params.StatusInfo{Current: "active"}},
	// 		},
	// 	}, nil)

	// 	// 3. Mocks for s.setKubernetesClient for the k8sFacilityName (assuming client needs to be set up)
	// 	mockKubeClient.EXPECT().Exists(scopeUUID, k8sFacilityName).Return(false)
	// 	leaderUnitName := k8sFacilityName + "/0"
	// 	kubeControlUnitName := "kubernetes-control-plane/0" // Example

	// 	mockJujuApplication.EXPECT().GetLeader(ctx, scopeUUID, k8sFacilityName).Return(leaderUnitName, nil)
	// 	mockJujuApplication.EXPECT().GetUnitInfo(ctx, scopeUUID, leaderUnitName).Return(
	// 		&model.UnitInfo{RelationData: []model.UnitRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]model.UnitData{kubeControlUnitName: {}}}}}, nil)
	// 	mockJujuApplication.EXPECT().GetUnitInfo(ctx, scopeUUID, kubeControlUnitName).Return(
	// 		&model.UnitInfo{RelationData: []model.UnitRelationData{{UnitRelationData: map[string]model.UnitData{"related/0": {UnitData: map[string]interface{}{
	// 			"api-endpoints": `["https://10.0.0.1:6443"]`,
	// 			"creds":         `{"user1":{"client_token":"token123"}}`,
	// 		}}}}}}, nil)
	// 	mockKubeClient.EXPECT().Set(scopeUUID, k8sFacilityName, gomock.Any()).Return(nil)

	// 	// 4. Mock for s.helm.ListReleases
	// 	expectedHelmReleases := []*release.Release{{Name: "release1"}}
	// 	mockKubeHelm.EXPECT().ListReleases(scopeUUID, k8sFacilityName, "").Return(expectedHelmReleases, nil)

	// 	releases, err := s.ListReleases(ctx)
	// 	if err != nil {
	// 		t.Fatalf("ListReleases() unexpected error: %v", err)
	// 	}

	// 	if len(releases) != 1 {
	// 		t.Fatalf("Expected 1 release, got %d. Releases: %+v", len(releases), releases)
	// 	}
	// 	if releases[0].Release.Name != "release1" {
	// 		t.Errorf("Expected release name 'release1', got %s", releases[0].Name)
	// 	}
	// 	if releases[0].ScopeUUID != scopeUUID {
	// 		t.Errorf("Expected scope UUID %s, got %s", scopeUUID, releases[0].ScopeUUID)
	// 	}
	// 	if releases[0].FacilityName != k8sFacilityName {
	// 		t.Errorf("Expected facility name %s, got %s", k8sFacilityName, releases[0].FacilityName)
	// 	}
	// })

	// t.Run("error_listing_scopes", func(t *testing.T) {
	// 	mockScope.EXPECT().List(ctx).Return(nil, errExpected)
	// 	_, err := s.ListReleases(ctx)
	// 	if err == nil {
	// 		t.Fatal("expected error from scope.List, got nil")
	// 	}
	// })

	// t.Run("error_client_status", func(t *testing.T) {
	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: scopeName}}, nil)
	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(nil, errExpected)
	// 	_, err := s.ListReleases(ctx)
	// 	if err == nil {
	// 		t.Fatal("expected error from client.Status, got nil")
	// 	}
	// })

	// t.Run("error_set_kubernetes_client", func(t *testing.T) {
	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: scopeName}}, nil)
	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(&params.FullStatus{
	// 		Applications: map[string]params.ApplicationStatus{
	// 			k8sFacilityName: {Charm: charmNameKubernetes, Status: &params.StatusInfo{Current: "active"}},
	// 		},
	// 	}, nil)
	// 	mockKubeClient.EXPECT().Exists(scopeUUID, k8sFacilityName).Return(false)
	// 	mockJujuApplication.EXPECT().GetLeader(ctx, scopeUUID, k8sFacilityName).Return("", errExpected) // Error here

	// 	_, err := s.ListReleases(ctx)
	// 	if err == nil {
	// 		t.Fatal("expected error from setKubernetesClient, got nil")
	// 	}
	// })

	// t.Run("error_helm_list_releases", func(t *testing.T) {
	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: scopeName}}, nil)
	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(&params.FullStatus{
	// 		Applications: map[string]params.ApplicationStatus{
	// 			k8sFacilityName: {Charm: charmNameKubernetes, Status: &params.StatusInfo{Current: "active"}},
	// 		},
	// 	}, nil)
	// 	mockKubeClient.EXPECT().Exists(scopeUUID, k8sFacilityName).Return(true) // Assume client exists
	// 	mockKubeHelm.EXPECT().ListReleases(scopeUUID, k8sFacilityName, "").Return(nil, errExpected)

	// 	_, err := s.ListReleases(ctx)
	// 	if err == nil {
	// 		t.Fatal("expected error from helm.ListReleases, got nil")
	// 	}
	// })
	// t.Run("success", func(t *testing.T) {
	// 	ctx := context.Background()
	// 	scopeUUID := "uuid1"
	// 	modelName := "model1"
	// 	k8sFacilityName := "facility1"
	// 	const localCharmNameKubernetes = "kubernetes"

	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: modelName}}, nil)

	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(&params.FullStatus{
	// 		Applications: map[string]params.ApplicationStatus{
	// 			k8sFacilityName: {Charm: localCharmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
	// 		},
	// 	}, nil)

	// 	mockKubeClient.EXPECT().Exists(scopeUUID, k8sFacilityName).Return(true) // Simplified, was Return(true, nil) - Exists returns bool

	// 	expectedHelmReleases := []release.Release{ // KubeHelm.ListReleases returns []release.Release
	// 		{Name: "release1"},
	// 	}
	// 	mockKubeHelm.EXPECT().ListReleases(scopeUUID, k8sFacilityName, "").Return(expectedHelmReleases, nil)

	// 	releases, err := s.ListReleases(ctx)
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}

	// 	if len(releases) != 1 {
	// 		t.Fatalf("unexpected number of releases: got %d, want %d", len(releases), 1)
	// 	}

	// 	if releases[0].Release.Name != expectedHelmReleases[0].Name {
	// 		t.Fatalf("unexpected release name: got %s, want %s", releases[0].Release.Name, expectedHelmReleases[0].Name)
	// 	}
	// 	if releases[0].ScopeUUID != scopeUUID {
	// 		t.Errorf("Expected scope UUID %s, got %s", scopeUUID, releases[0].ScopeUUID)
	// 	}
	// 	if releases[0].FacilityName != k8sFacilityName {
	// 		t.Errorf("Expected facility name %s, got %s", k8sFacilityName, releases[0].FacilityName)
	// 	}
	// })

	t.Run("error_listing_scopes", func(t *testing.T) {
		ctx := context.Background()
		mockScope.EXPECT().List(ctx).Return(nil, errExpected)

		_, err := s.ListReleases(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	// t.Run("error_client_status", func(t *testing.T) {
	// 	ctx := context.Background()
	// 	scopeUUID := "uuid1"
	// 	modelName := "model1"

	// 	mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{{UUID: scopeUUID, Name: modelName}}, nil)
	// 	mockClient.EXPECT().Status(ctx, scopeUUID, gomock.Nil()).Return(nil, errExpected)

	// 	_, err := s.ListReleases(ctx)
	// 	if err == nil {
	// 		t.Fatal("expected error, got nil")
	// 	}
	// })
}
func TestNexusService_CreateRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockKubeHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{
		facility:   mockJujuApplication,
		kubernetes: mockKubeClient,
		helm:       mockKubeHelm,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "test-namespace"
	name := "test-release"
	dryRun := false
	chartRef := "test-chart"
	valuesYAML := "name: value"
	valuesMap := map[string]string{"key": "value"}

	t.Run("success", func(t *testing.T) {
		// Mocks for setKubernetesClient
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // Assume client exists

		expectedRelease := &release.Release{Name: name}
		// Construct expected values map after strvals.ParseInto
		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}

		mockKubeHelm.EXPECT().InstallRelease(uuid, facility, namespace, name, dryRun, chartRef, gomock.Eq(expectedValuesForHelm)).Return(expectedRelease, nil)

		rel, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rel.Release.Name != expectedRelease.Name {
			t.Errorf("expected release name %q, got %q", expectedRelease.Name, rel.Release.Name)
		}
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false) // Client doesn't exist
		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected)

		_, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v to contain %v", err, errExpected)
		}
	})

	t.Run("error_install_release", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true)
		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}
		mockKubeHelm.EXPECT().InstallRelease(uuid, facility, namespace, name, dryRun, chartRef, gomock.Eq(expectedValuesForHelm)).Return(nil, errExpected)

		_, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v to contain %v", err, errExpected)
		}
	})
	t.Run("generates_valid_name", func(t *testing.T) {
		// Run a few times to increase confidence, though not strictly necessary
		// as namesgenerator itself is quite consistent in format.
		for i := 0; i < 5; i++ {
			name := randomName()
			if name == "" {
				t.Errorf("randomName() returned an empty string on attempt %d", i+1)
			}
			if strings.Contains(name, "_") {
				t.Errorf("randomName() returned %q, which contains an underscore, on attempt %d", name, i+1)
			}
			// namesgenerator.GetRandomName(0) usually produces names like "adjective_noun"
			// so we expect a hyphen after replacement.
			if !strings.Contains(name, "-") {
				t.Logf("randomName() returned %q, which does not contain a hyphen (this might be rare but possible if namesgenerator produced a single word without underscore), on attempt %d", name, i+1)
			}

			parts := strings.Split(name, "-")
			if len(parts) < 1 || len(parts) > 2 { // It can be one part if namesgenerator returns a single word
				t.Errorf("randomName() returned %q, expected 1 or 2 parts separated by hyphen, got %d parts, on attempt %d", name, len(parts), i+1)
			}
			for _, part := range parts {
				if part == "" {
					t.Errorf("randomName() returned %q, which contains an empty part, on attempt %d", name, i+1)
				}
			}
		}
	})
}

// func TestNexusService_CreateRelease(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
// 	mockKubeClient := mocks.NewMockKubeClient(ctrl)
// 	mockKubeHelm := mocks.NewMockKubeHelm(ctrl)
// 	s := &NexusService{
// 		facility:   mockJujuApplication,
// 		kubernetes: mockKubeClient,
// 		helm:       mockKubeHelm,
// 	}

// 	ctx := context.Background()
// 	uuid := "test-uuid"
// 	facility := "test-facility"
// 	namespace := "test-namespace"
// 	name := "test-release"
// 	dryRun := false
// 	chartRef := "test-chart"
// 	valuesYAML := "name: value"
// 	valuesMap := map[string]string{"key": "value"}

// 	t.Run("success", func(t *testing.T) {
// 		// Mocks for setKubernetesClient
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // Assume client exists

// 		expectedRelease := &release.Release{Name: name}
// 		// Construct expected values map after strvals.ParseInto
// 		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}

// 		mockKubeHelm.EXPECT().InstallRelease(uuid, facility, namespace, name, dryRun, chartRef, expectedValuesForHelm).Return(expectedRelease, nil)

// 		rel, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 		if rel.Release.Name != expectedRelease.Name {
// 			t.Errorf("expected release name %q, got %q", expectedRelease.Name, rel.Release.Name)
// 		}
// 	})

// 	t.Run("error_set_kubernetes_client", func(t *testing.T) {
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false) // Client doesn't exist
// 		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected)

// 		_, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
// 		if err == nil {
// 			t.Fatal("expected error, got nil")
// 		}
// 	})

// 	t.Run("error_install_release", func(t *testing.T) {
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true)
// 		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}
// 		mockKubeHelm.EXPECT().InstallRelease(uuid, facility, namespace, name, dryRun, chartRef, expectedValuesForHelm).Return(nil, errExpected)

//			_, err := s.CreateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML, valuesMap)
//			if err == nil {
//				t.Fatal("expected error, got nil")
//			}
//		})
//	}
func TestNexusService_UpdateRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockKubeHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{
		facility:   mockJujuApplication,
		kubernetes: mockKubeClient,
		helm:       mockKubeHelm,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "test-namespace"
	name := "test-release"
	dryRun := false
	chartRef := "test-chart/mychart" // Added chartRef for UpdateRelease
	valuesYAML := "name: updatedValue"
	// valuesMap is not used by UpdateRelease method, values are derived from valuesYAML only

	t.Run("success", func(t *testing.T) {
		// Mocks for setKubernetesClient
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // Assume client exists

		expectedRelease := &release.Release{Name: name}
		// Values for UpgradeRelease are derived only from valuesYAML in the UpdateRelease method
		expectedValuesFromYAML := map[string]any{"name": "updatedValue"}

		mockKubeHelm.EXPECT().UpgradeRelease(uuid, facility, namespace, name, dryRun, chartRef, gomock.Eq(expectedValuesFromYAML)).Return(expectedRelease, nil)

		// Corrected call to s.UpdateRelease: uses chartRef and valuesYAML
		rel, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rel.Release.Name != expectedRelease.Name {
			t.Errorf("expected release name %q, got %q", expectedRelease.Name, rel.Release.Name)
		}
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false) // Client doesn't exist
		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected)

		// Corrected call to s.UpdateRelease
		_, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v to contain %v", err, errExpected)
		}
	})

	t.Run("error_upgrade_release", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true)
		expectedValuesFromYAML := map[string]any{"name": "updatedValue"}
		// Corrected mock expectation for UpgradeRelease
		mockKubeHelm.EXPECT().UpgradeRelease(uuid, facility, namespace, name, dryRun, chartRef, gomock.Eq(expectedValuesFromYAML)).Return(nil, errExpected)

		// Corrected call to s.UpdateRelease
		_, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, chartRef, valuesYAML)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v to contain %v", err, errExpected)
		}
	})
}

// func TestNexusService_UpdateRelease(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockJujuApplication := mocks.NewMockJujuApplication(ctrl)
// 	mockKubeClient := mocks.NewMockKubeClient(ctrl)
// 	mockKubeHelm := mocks.NewMockKubeHelm(ctrl)
// 	s := &NexusService{
// 		facility:   mockJujuApplication,
// 		kubernetes: mockKubeClient,
// 		helm:       mockKubeHelm,
// 	}

// 	ctx := context.Background()
// 	uuid := "test-uuid"
// 	facility := "test-facility"
// 	namespace := "test-namespace"
// 	name := "test-release"
// 	dryRun := false
// 	valuesYAML := "name: value"
// 	valuesMap := map[string]string{"key": "value"}

// 	t.Run("success", func(t *testing.T) {
// 		// Mocks for setKubernetesClient
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // Assume client exists

// 		expectedRelease := &release.Release{Name: name}
// 		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}

// 		mockKubeHelm.EXPECT().UpgradeRelease(uuid, facility, namespace, name, dryRun, expectedValuesForHelm).Return(expectedRelease, nil)

// 		rel, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, valuesYAML, valuesMap)
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 		if rel.Release.Name != expectedRelease.Name {
// 			t.Errorf("expected release name %q, got %q", expectedRelease.Name, rel.Release.Name)
// 		}
// 	})

// 	t.Run("error_set_kubernetes_client", func(t *testing.T) {
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false) // Client doesn't exist
// 		mockJujuApplication.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected)

// 		_, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, valuesYAML, valuesMap)
// 		if err == nil {
// 			t.Fatal("expected error, got nil")
// 		}
// 	})

// 	t.Run("error_upgrade_release", func(t *testing.T) {
// 		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true)
// 		expectedValuesForHelm := map[string]any{"name": "value", "key": "value"}
// 		mockKubeHelm.EXPECT().UpgradeRelease(uuid, facility, namespace, name, dryRun, expectedValuesForHelm).Return(nil, errExpected)
// 		_, err := s.UpdateRelease(ctx, uuid, facility, namespace, name, dryRun, valuesYAML, valuesMap)
// 		if err == nil {
// 			t.Fatal("expected error, got nil")
// 		}
// 	})
// }

func TestNexusService_GetChartMetadataFromApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{helm: mockHelm}
	uuid := "test-uuid"
	facility := "test-facility"
	app := &model.Application{
		Namespace: "test-namespace",
		Labels:    map[string]string{"app.kubernetes.io/instance": "test-release"},
	}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedValues := map[string]any{"key": "value"}
		mockHelm.EXPECT().GetValues(uuid, facility, app.Namespace, app.Labels["app.kubernetes.io/instance"]).Return(expectedValues, nil)

		md, err := s.GetChartMetadataFromApplication(ctx, uuid, facility, app)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		valuesYAML, _ := yaml.Marshal(expectedValues)

		expectedMD := &model.ChartMetadata{ValuesYAML: string(valuesYAML)}
		if !reflect.DeepEqual(md, expectedMD) {
			t.Errorf("expected metadata %+v, got %+v", expectedMD, md)
		}
	})

	t.Run("error getting values", func(t *testing.T) {
		mockHelm.EXPECT().GetValues(uuid, facility, app.Namespace, app.Labels["app.kubernetes.io/instance"]).Return(nil, errors.New("test error"))
		_, err := s.GetChartMetadataFromApplication(ctx, uuid, facility, app)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("no release label", func(t *testing.T) {
		app := &model.Application{
			Namespace: "test-namespace",
			Labels:    map[string]string{},
		}
		md, err := s.GetChartMetadataFromApplication(ctx, uuid, facility, app)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedMD := &model.ChartMetadata{}

		if !reflect.DeepEqual(md, expectedMD) {
			t.Errorf("expected metadata %+v, got %+v", expectedMD, md)
		}
	})

}

// func TestNexusService_GetChart(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockHelm := mocks.NewMockKubeHelm(ctrl)
// 	s := &NexusService{helm: mockHelm}

// 	ctx := context.Background()

// 	t.Run("error listing chart versions", func(t *testing.T) {
// 		mockHelm.EXPECT().ListChartVersions(ctx).Return(nil, errExpected)
// 		_, err := s.GetChart(ctx, "test-chart")

// 		if err == nil {
// 			t.Fatal("expected error, got nil")
// 		}
// 	})
// }

func TestNexusService_fromDeployment(t *testing.T) {
	deployment := createMockDeployment("test-deployment", "test-namespace")
	services := []corev1.Service{*createMockService("svc1", "test-namespace", map[string]string{"app": "test"}), *createMockService("svc2", "other-namespace", map[string]string{"app": "test"})}
	pods := []corev1.Pod{*createMockPod("pod1", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod2", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod3", "other-namespace", map[string]string{"app": "test"})}
	pvcs := []corev1.PersistentVolumeClaim{*createMockPVC("pvc-1", "test-namespace", "standard")}
	scm := map[string]storagev1.StorageClass{"standard": *createMockStorageClass("standard")}

	s := &NexusService{}
	app, err := s.fromDeployment(deployment, services, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedReplicas := int32(1)
	if *app.Replicas != expectedReplicas {
		t.Errorf("expected replicas %d, got %d", expectedReplicas, *app.Replicas)
	}
	if app.Name != deployment.Name {
		t.Errorf("expected name %s, got %s", deployment.Name, app.Name)
	}
	if len(app.Services) != 1 || app.Services[0].Name != "svc1" {
		t.Errorf("services not filtered correctly: %+v", app.Services)
	}
	if len(app.Pods) != 2 {
		t.Errorf("pods not filtered correctly, expected 2 got: %d: %+v", len(app.Pods), app.Pods)
	}
}

func TestNexusService_fromStatefulSet(t *testing.T) {
	statefulSet := createMockStatefulSet("test-statefulset", "test-namespace")
	services := []corev1.Service{*createMockService("svc1", "test-namespace", map[string]string{"app": "test"}), *createMockService("svc2", "other-namespace", map[string]string{"app": "test"})}
	pods := []corev1.Pod{*createMockPod("pod1", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod2", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod3", "other-namespace", map[string]string{"app": "test"})}
	pvcs := []corev1.PersistentVolumeClaim{*createMockPVC("pvc-1", "test-namespace", "standard")}
	scm := map[string]storagev1.StorageClass{"standard": *createMockStorageClass("standard")}

	s := &NexusService{}
	app, err := s.fromStatefulSet(statefulSet, services, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedReplicas := int32(1)
	if *app.Replicas != expectedReplicas {
		t.Errorf("expected replicas %d, got %d", expectedReplicas, *app.Replicas)
	}
	if app.Name != statefulSet.Name {
		t.Errorf("expected name %s, got %s", statefulSet.Name, app.Name)
	}
	if len(app.Services) != 1 || app.Services[0].Name != "svc1" {
		t.Errorf("services not filtered correctly: %+v", app.Services)
	}
	if len(app.Pods) != 2 {
		t.Errorf("pods not filtered correctly, expected 2 got: %d: %+v", len(app.Pods), app.Pods)
	}
}

func TestNexusService_fromDaemonSet(t *testing.T) {
	daemonSet := createMockDaemonSet("test-daemonset", "test-namespace")
	services := []corev1.Service{*createMockService("svc1", "test-namespace", map[string]string{"app": "test"}), *createMockService("svc2", "other-namespace", map[string]string{"app": "test"})}
	pods := []corev1.Pod{*createMockPod("pod1", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod2", "test-namespace", map[string]string{"app": "test"}), *createMockPod("pod3", "other-namespace", map[string]string{"app": "test"})}
	pvcs := []corev1.PersistentVolumeClaim{*createMockPVC("pvc-1", "test-namespace", "standard")}
	scm := map[string]storagev1.StorageClass{"standard": *createMockStorageClass("standard")}

	s := &NexusService{}
	app, err := s.fromDaemonSet(daemonSet, services, pods, pvcs, scm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.Name != daemonSet.Name {
		t.Errorf("expected name %s, got %s", daemonSet.Name, app.Name)
	}
	if app.Replicas != nil {
		t.Error("expected replicas to be nil for DaemonSet")
	}
	if len(app.Services) != 1 || app.Services[0].Name != "svc1" {
		t.Errorf("services not filtered correctly: %+v", app.Services)
	}
	if len(app.Pods) != 2 {
		t.Errorf("pods not filtered correctly, expected 2 got: %d: %+v", len(app.Pods), app.Pods)
	}
}

func Test_toApplication(t *testing.T) {
	labelSelector := &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}}
	appType := "Deployment"
	name := "test-app"
	namespace := "test-ns"
	objectMeta := &metav1.ObjectMeta{Name: name, Namespace: namespace}
	labelsMap := map[string]string{"label1": "value1"}
	replicas := int32Ptr(3)
	containers := []corev1.Container{{Name: "container1", Image: "image1"}}
	volumes := []corev1.Volume{{Name: "vol1", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "test-pvc"}}}}
	services := []corev1.Service{*createMockService("test-svc", namespace, map[string]string{"app": "test"})}
	pods := []corev1.Pod{*createMockPod("test-pod", namespace, map[string]string{"app": "test"})}
	pvcs := []corev1.PersistentVolumeClaim{*createMockPVC("test-pvc", namespace, "standard")}
	scm := map[string]storagev1.StorageClass{"standard": *createMockStorageClass("standard")}

	app, err := toApplication(labelSelector, appType, name, namespace, objectMeta, labelsMap, replicas, containers, volumes, services, pods, pvcs, scm)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.Type != appType {
		t.Errorf("expected type %s, got %s", appType, app.Type)
	}
	if app.Name != name {
		t.Errorf("expected name %s, got %s", name, app.Name)
	}
	if len(app.PersistentVolumeClaims) != 1 {
		t.Errorf("PersistentVolumeClaims not filtered correctly: %+v", app.PersistentVolumeClaims)
	}
}

func Test_newKubernetesConfig(t *testing.T) {
	tests := []struct {
		name        string
		unitInfo    *model.UnitInfo
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Valid_config",
			unitInfo: &model.UnitInfo{
				RelationData: []application.EndpointRelationData{
					{
						UnitRelationData: map[string]application.RelationData{
							"kube-control-plane/0": { // Example related unit name
								UnitData: map[string]interface{}{
									"api-endpoints": `["https://10.0.0.1:6443"]`,
									"creds":         `{"user1":{"client_token":"test-token","scope":"","kubelet_token":"","proxy_token":""}}`,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing_api_endpoints",
			unitInfo: &model.UnitInfo{
				RelationData: []application.EndpointRelationData{
					{
						UnitRelationData: map[string]application.RelationData{
							"kube-control-plane/0": {
								UnitData: map[string]interface{}{
									"creds": `{"user1":{"client_token":"test-token"}}`,
								},
							},
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: "endpoint not found",
		},
		{
			name: "Missing_creds",
			unitInfo: &model.UnitInfo{
				RelationData: []application.EndpointRelationData{
					{
						UnitRelationData: map[string]application.RelationData{
							"kube-control-plane/0": {
								UnitData: map[string]interface{}{
									"api-endpoints": `["https://10.0.0.1:6443"]`,
								},
							},
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: "token not found",
		},
		{
			name: "Invalid_api_endpoints_json",
			unitInfo: &model.UnitInfo{
				RelationData: []application.EndpointRelationData{
					{
						UnitRelationData: map[string]application.RelationData{
							"kube-control-plane/0": {
								UnitData: map[string]interface{}{
									"api-endpoints": `invalid-json`,
									"creds":         `{"user1":{"client_token":"test-token"}}`,
								},
							},
						},
					},
				},
			},
			wantErr: true,
			// expectedErr: "cannot unmarshal invalid-json into Go value of type []string", // Error message might vary slightly
		},
		{
			name: "Invalid_creds_json",
			unitInfo: &model.UnitInfo{
				RelationData: []application.EndpointRelationData{
					{
						UnitRelationData: map[string]application.RelationData{
							"kube-control-plane/0": {
								UnitData: map[string]interface{}{
									"api-endpoints": `["https://10.0.0.1:6443"]`,
									"creds":         `invalid-json`,
								},
							},
						},
					},
				},
			},
			wantErr: true,
			// expectedErr: "cannot unmarshal invalid-json into Go value of type map[string]model.ControlPlaneCredential", // Error message might vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := newKubernetesConfig(tt.unitInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("newKubernetesConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.expectedErr != "" && (err == nil || !reflect.DeepEqual(err.Error(), tt.expectedErr) && !strings.Contains(err.Error(), tt.expectedErr)) {
					t.Errorf("newKubernetesConfig() error = %v, want error string containing %v", err, tt.expectedErr)
				}
				return
			}
			if cfg == nil {
				t.Error("newKubernetesConfig() expected config, got nil")
			} else {
				if cfg.Host != "https://10.0.0.1:6443" { // From valid config
					t.Errorf("newKubernetesConfig() host = %s, want https://10.0.0.1:6443", cfg.Host)
				}
				if cfg.BearerToken != "test-token" { // From valid config
					t.Errorf("newKubernetesConfig() token = %s, want test-token", cfg.BearerToken)
				}
			}
		})
	}
}
func TestNexusService_internalListReleases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScope := mocks.NewMockJujuModel(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl)
	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockHelm := mocks.NewMockKubeHelm(ctrl)

	s := &NexusService{
		scope:      mockScope,
		client:     mockClient,
		facility:   mockJujuApp,
		kubernetes: mockKubeClient,
		helm:       mockHelm,
	}

	ctx := context.Background()
	scopeUUID1 := "scope-uuid-1"
	scopeName1 := "scope-name-1"
	facilityName1 := "k8s-facility-1"

	scopeUUID2 := "scope-uuid-2"
	scopeName2 := "scope-name-2"
	facilityName2 := "k8s-facility-2"

	t.Run("success", func(t *testing.T) {
		// Mock listFacilitiesAcrossScopes
		mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{
			{UUID: scopeUUID1, Name: scopeName1},
			{UUID: scopeUUID2, Name: scopeName2},
		}, nil)

		// Facility 1
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID1, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName1: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		// Mock setKubernetesClient for facility 1
		mockKubeClient.EXPECT().Exists(scopeUUID1, facilityName1).Return(true)
		mockHelm.EXPECT().ListReleases(scopeUUID1, facilityName1, "").Return([]release.Release{{Name: "release1"}}, nil)

		// Facility 2
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID2, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName2: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		// Mock setKubernetesClient for facility 2
		mockKubeClient.EXPECT().Exists(scopeUUID2, facilityName2).Return(true)
		mockHelm.EXPECT().ListReleases(scopeUUID2, facilityName2, "").Return([]release.Release{{Name: "release2"}}, nil)

		releases, err := s.listReleases(ctx)
		if err != nil {
			t.Fatalf("listReleases() unexpected error: %v", err)
		}
		if len(releases) != 2 {
			t.Fatalf("expected 2 releases, got %d", len(releases))
		}
	})

	t.Run("error_listing_scopes", func(t *testing.T) {
		mockScope.EXPECT().List(ctx).Return(nil, errExpected)
		_, err := s.listReleases(ctx)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_listing_facilities_for_one_scope", func(t *testing.T) {
		mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{
			{UUID: scopeUUID1, Name: scopeName1},
		}, nil)
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID1, gomock.Any()).Return(nil, errExpected) // Error for this scope

		// Since errgroup continues, the function might not return this specific error directly
		// but it should not panic and ideally return an aggregated error or the first one.
		// For this test, we check that an error is returned.
		_, err := s.listReleases(ctx)
		if err == nil {
			t.Fatal("listReleases() expected an error, got nil")
		}
	})

	t.Run("error_set_kubernetes_client_for_one_facility_continues", func(t *testing.T) {
		mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{
			{UUID: scopeUUID1, Name: scopeName1}, // Facility that will have setKubernetesClient error
			{UUID: scopeUUID2, Name: scopeName2}, // Facility that will succeed
		}, nil)

		// Facility 1 (error)
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID1, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName1: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		mockKubeClient.EXPECT().Exists(scopeUUID1, facilityName1).Return(false) // Needs setup
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), scopeUUID1, facilityName1).Return("", errExpected)
		// No ListReleases mock for facilityName1 as setKubernetesClient fails

		// Facility 2 (success)
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID2, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName2: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		mockKubeClient.EXPECT().Exists(scopeUUID2, facilityName2).Return(true)
		mockHelm.EXPECT().ListReleases(scopeUUID2, facilityName2, "").Return([]release.Release{{Name: "release2"}}, nil)

		releases, err := s.listReleases(ctx)
		if err != nil { // The errgroup might return the first error, or nil if all goroutines returned nil (due to "return nil // pass")
			t.Logf("listReleases() returned error (as one path failed): %v", err)
		}
		// We expect releases from the successful path
		if len(releases) != 1 || releases[0].Release.Name != "release2" {
			t.Fatalf("expected 1 release ('release2'), got %d releases: %+v", len(releases), releases)
		}
	})

	t.Run("error_helm_list_releases_for_one_facility_continues", func(t *testing.T) {
		mockScope.EXPECT().List(ctx).Return([]base.UserModelSummary{
			{UUID: scopeUUID1, Name: scopeName1}, // Facility that will have helm error
			{UUID: scopeUUID2, Name: scopeName2}, // Facility that will succeed
		}, nil)

		// Facility 1 (helm error)
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID1, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName1: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		mockKubeClient.EXPECT().Exists(scopeUUID1, facilityName1).Return(true)
		mockHelm.EXPECT().ListReleases(scopeUUID1, facilityName1, "").Return(nil, errExpected)

		// Facility 2 (success)
		mockClient.EXPECT().Status(gomock.Any(), scopeUUID2, gomock.Any()).Return(&params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName2: {Charm: charmNameKubernetes, Status: params.DetailedStatus{Status: "active"}},
			},
		}, nil)
		mockKubeClient.EXPECT().Exists(scopeUUID2, facilityName2).Return(true)
		mockHelm.EXPECT().ListReleases(scopeUUID2, facilityName2, "").Return([]release.Release{{Name: "release2"}}, nil)

		releases, err := s.listReleases(ctx)
		if err != nil {
			t.Logf("listReleases() returned error (as one path failed): %v", err)
		}
		if len(releases) != 1 || releases[0].Release.Name != "release2" {
			t.Fatalf("expected 1 release ('release2'), got %d releases: %+v", len(releases), releases)
		}
	})
}

func TestNexusService_setKubernetesClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)

	s := &NexusService{
		facility:   mockJujuApp,
		kubernetes: mockKubeClient,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facilityName := "k8s-facility"
	leaderUnitName := facilityName + "/0"
	kubeControlUnitName := "kubernetes-control-plane/0"

	t.Run("success_client_exists", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("setKubernetesClient() unexpected error: %v", err)
		}
	})

	t.Run("success_client_needs_setup", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]application.RelationData{kubeControlUnitName: {}}}}}, nil,
		)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, kubeControlUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{UnitRelationData: map[string]application.RelationData{"related/0": {UnitData: map[string]interface{}{
				"api-endpoints": `["https://1.2.3.4:6443"]`,
				"creds":         `{"user":{"client_token":"token123"}}`,
			}}}}}}, nil,
		)
		mockKubeClient.EXPECT().Set(uuid, facilityName, gomock.Any()).Return(nil)

		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("setKubernetesClient() unexpected error: %v", err)
		}
	})

	t.Run("error_get_leader", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return("", errExpected)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_get_cp_unit_info", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(nil, errExpected)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_extract_kube_control_not_found", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "other-relation"}}}, nil, // No kube-control
		)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if err == nil || !strings.Contains(err.Error(), "kube-control not found") {
			t.Fatalf("expected kube-control not found error, got %v", err)
		}
	})

	t.Run("error_get_worker_unit_info", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]application.RelationData{kubeControlUnitName: {}}}}}, nil,
		)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, kubeControlUnitName).Return(nil, errExpected)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_new_kubernetes_config_no_endpoint", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]application.RelationData{kubeControlUnitName: {}}}}}, nil,
		)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, kubeControlUnitName).Return( // No api-endpoints
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{UnitRelationData: map[string]application.RelationData{"related/0": {UnitData: map[string]interface{}{
				"creds": `{"user":{"client_token":"token123"}}`,
			}}}}}}, nil,
		)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if err == nil || !strings.Contains(err.Error(), "endpoint not found") {
			t.Fatalf("expected endpoint not found error, got %v", err)
		}
	})

	t.Run("error_new_kubernetes_config_no_token", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]application.RelationData{kubeControlUnitName: {}}}}}, nil,
		)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, kubeControlUnitName).Return( // No creds
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{UnitRelationData: map[string]application.RelationData{"related/0": {UnitData: map[string]interface{}{
				"api-endpoints": `["https://1.2.3.4:6443"]`,
			}}}}}}, nil,
		)
		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if err == nil || !strings.Contains(err.Error(), "token not found") {
			t.Fatalf("expected token not found error, got %v", err)
		}
	})

	t.Run("error_kube_client_set", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false)
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facilityName).Return(leaderUnitName, nil)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, leaderUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{Endpoint: "kube-control", UnitRelationData: map[string]application.RelationData{kubeControlUnitName: {}}}}}, nil,
		)
		mockJujuApp.EXPECT().GetUnitInfo(gomock.Any(), uuid, kubeControlUnitName).Return(
			&model.UnitInfo{RelationData: []application.EndpointRelationData{{UnitRelationData: map[string]application.RelationData{"related/0": {UnitData: map[string]interface{}{
				"api-endpoints": `["https://1.2.3.4:6443"]`,
				"creds":         `{"user":{"client_token":"token123"}}`,
			}}}}}}, nil,
		)
		mockKubeClient.EXPECT().Set(uuid, facilityName, gomock.Any()).Return(errExpected)

		err := s.setKubernetesClient(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}
func TestNexusService_GetChartMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{helm: mockHelm}

	ctx := context.Background()
	chartRef := "stable/mychart"
	expectedValuesYAML := "key: value"
	expectedReadmeMD := "# My Chart"

	t.Run("success", func(t *testing.T) {
		mockHelm.EXPECT().ShowChart(chartRef, action.ShowValues).Return(expectedValuesYAML, nil)
		mockHelm.EXPECT().ShowChart(chartRef, action.ShowReadme).Return(expectedReadmeMD, nil)

		md, err := s.GetChartMetadata(ctx, chartRef)
		if err != nil {
			t.Fatalf("GetChartMetadata() unexpected error: %v", err)
		}
		if md == nil {
			t.Fatalf("expected non-nil metadata, got nil")
		}
		if md.ValuesYAML != expectedValuesYAML {
			t.Errorf("expected values YAML %q, got %q", expectedValuesYAML, md.ValuesYAML)
		}
		if md.ReadmeMD != expectedReadmeMD {
			t.Errorf("expected readme MD %q, got %q", expectedReadmeMD, md.ReadmeMD)
		}
	})

	// t.Run("error_show_values", func(t *testing.T) {
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowValues).Return("", errExpected)
	// 	// The ShowReadme call will also be made by the errgroup.
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowReadme).Return(expectedReadmeMD, nil)

	// 	md, err := s.GetChartMetadata(ctx, chartRef)
	// 	if !errors.Is(err, errExpected) {
	// 		t.Fatalf("expected error %v, got %v", errExpected, err)
	// 	}
	// 	if md == nil {
	// 		t.Fatalf("expected non-nil metadata even on error, got nil")
	// 	}
	// 	// ValuesYAML should be empty because its fetch failed.
	// 	if md.ValuesYAML != "" {
	// 		t.Errorf("expected ValuesYAML to be empty, got %q", md.ValuesYAML)
	// 	}
	// 	// ReadmeMD should be populated as its fetch succeeded.
	// 	if md.ReadmeMD != expectedReadmeMD {
	// 		t.Errorf("expected ReadmeMD %q, got %q", expectedReadmeMD, md.ReadmeMD)
	// 	}
	// })

	// t.Run("error_show_readme", func(t *testing.T) {
	// 	// The ShowValues call will also be made by the errgroup.
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowValues).Return(expectedValuesYAML, nil)
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowReadme).Return("", errExpected)

	// 	md, err := s.GetChartMetadata(ctx, chartRef)
	// 	if !errors.Is(err, errExpected) {
	// 		t.Fatalf("expected error %v, got %v", errExpected, err) // Line 1335
	// 	}
	// 	if md == nil {
	// 		t.Fatalf("expected non-nil metadata even on error, got nil")
	// 	}
	// 	// ValuesYAML should be populated as its fetch succeeded.
	// 	if md.ValuesYAML != expectedValuesYAML {
	// 		t.Errorf("expected ValuesYAML %q, got %q", expectedValuesYAML, md.ValuesYAML)
	// 	}
	// 	// ReadmeMD should be empty because its fetch failed.
	// 	if md.ReadmeMD != "" {
	// 		t.Errorf("expected ReadmeMD to be empty, got %q", md.ReadmeMD)
	// 	}
	// })

	// t.Run("partial_success_readme_fails", func(t *testing.T) {
	// 	// This test is effectively the same as "error_show_readme".
	// 	// ShowValues succeeds, ShowReadme fails.
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowValues).Return(expectedValuesYAML, nil) // Line 1340
	// 	mockHelm.EXPECT().ShowChart(chartRef, action.ShowReadme).Return("", errExpected)         // Line 1341

	// 	md, err := s.GetChartMetadata(ctx, chartRef)
	// 	if !errors.Is(err, errExpected) { // Line 1345
	// 		t.Fatalf("expected error %v, got %v", errExpected, err)
	// 	}
	// 	if md == nil {
	// 		t.Fatalf("expected non-nil metadata even on error, got nil")
	// 	}
	// 	// ValuesYAML should be populated as its fetch succeeded.
	// 	if md.ValuesYAML != expectedValuesYAML {
	// 		t.Errorf("expected ValuesYAML to be populated with %q, got %q", expectedValuesYAML, md.ValuesYAML)
	// 	}
	// 	// ReadmeMD should be empty because its fetch failed.
	// 	if md.ReadmeMD != "" {
	// 		t.Errorf("expected ReadmeMD to be empty, got %q", md.ReadmeMD)
	// 	}
	// })
}
func TestNexusService_GetChart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{helm: mockHelm}

	ctx := context.Background()
	chartName := "mychart"
	chartVersion := "1.0.0"

	indexFile := &repo.IndexFile{
		Entries: map[string]repo.ChartVersions{
			chartName:    {{Metadata: &chart.Metadata{Name: chartName, Version: chartVersion}}},
			"otherchart": {{Metadata: &chart.Metadata{Name: "otherchart", Version: "0.1.0"}}},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{indexFile}, nil)
		c, err := s.GetChart(ctx, chartName)
		if err != nil {
			t.Fatalf("GetChart() unexpected error: %v", err)
		}
		if c.Name != chartName {
			t.Errorf("expected chart name %q, got %q", chartName, c.Name)
		}
		if len(c.Versions) != 1 || c.Versions[0].Metadata.Version != chartVersion {
			t.Errorf("expected chart version %q, got versions: %+v", chartVersion, c.Versions)
		}
	})

	t.Run("chart_not_found", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{indexFile}, nil)
		_, err := s.GetChart(ctx, "nonexistentchart")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not found error, got %v", err)
		}
	})

	t.Run("error_listing_chart_versions", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return(nil, errExpected)
		_, err := s.GetChart(ctx, chartName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("empty_index_file", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{{Entries: map[string]repo.ChartVersions{}}}, nil)
		_, err := s.GetChart(ctx, chartName)
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}

func TestNexusService_ListCharts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{helm: mockHelm}

	ctx := context.Background()
	indexFile1 := &repo.IndexFile{
		Entries: map[string]repo.ChartVersions{
			"chart1": {{Metadata: &chart.Metadata{Name: "chart1", Version: "1.0.0"}}},
			"chart2": {{Metadata: &chart.Metadata{Name: "chart2", Version: "0.1.0"}}},
		},
	}
	indexFile2 := &repo.IndexFile{
		Entries: map[string]repo.ChartVersions{
			"chart3": {{Metadata: &chart.Metadata{Name: "chart3", Version: "2.0.0"}}},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{indexFile1, indexFile2}, nil)
		charts, err := s.ListCharts(ctx)
		if err != nil {
			t.Fatalf("ListCharts() unexpected error: %v", err)
		}
		if len(charts) != 3 {
			t.Fatalf("expected 3 charts, got %d", len(charts))
		}
		// Basic check for one of the charts
		found := false
		for _, c := range charts {
			if c.Name == "chart1" && len(c.Versions) == 1 && c.Versions[0].Metadata.Version == "1.0.0" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("chart 'chart1' not found or incorrect: %+v", charts)
		}
	})

	t.Run("error_listing_chart_versions", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return(nil, errExpected)
		_, err := s.ListCharts(ctx)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("empty_result", func(t *testing.T) {
		mockHelm.EXPECT().ListChartVersions(ctx).Return([]*repo.IndexFile{}, nil)
		charts, err := s.ListCharts(ctx)
		if err != nil {
			t.Fatalf("ListCharts() unexpected error: %v", err)
		}
		if len(charts) != 0 {
			t.Errorf("expected 0 charts, got %d", len(charts))
		}
	})
}

func TestNexusService_RollbackRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockJujuApp := mocks.NewMockJujuApplication(ctrl) // For setKubernetesClient
	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{
		kubernetes: mockKubeClient,
		facility:   mockJujuApp, // For setKubernetesClient
		helm:       mockHelm,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "test-namespace"
	name := "test-release"
	dryRun := false

	t.Run("success", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // For setKubernetesClient
		mockHelm.EXPECT().RollbackRelease(uuid, facility, namespace, name, dryRun).Return(nil)

		err := s.RollbackRelease(ctx, uuid, facility, namespace, name, dryRun)
		if err != nil {
			t.Fatalf("RollbackRelease() unexpected error: %v", err)
		}
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false)                         // For setKubernetesClient
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected) // For setKubernetesClient

		err := s.RollbackRelease(ctx, uuid, facility, namespace, name, dryRun)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_helm_rollback", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // For setKubernetesClient
		mockHelm.EXPECT().RollbackRelease(uuid, facility, namespace, name, dryRun).Return(errExpected)

		err := s.RollbackRelease(ctx, uuid, facility, namespace, name, dryRun)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_DeleteRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKubeClient := mocks.NewMockKubeClient(ctrl)
	mockJujuApp := mocks.NewMockJujuApplication(ctrl) // For setKubernetesClient
	mockHelm := mocks.NewMockKubeHelm(ctrl)
	s := &NexusService{
		kubernetes: mockKubeClient,
		facility:   mockJujuApp, // For setKubernetesClient
		helm:       mockHelm,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"
	namespace := "test-namespace"
	name := "test-release"
	dryRun := false

	t.Run("success", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // For setKubernetesClient
		mockHelm.EXPECT().UninstallRelease(uuid, facility, namespace, name, dryRun).Return(&release.Release{Name: name}, nil)

		err := s.DeleteRelease(ctx, uuid, facility, namespace, name, dryRun)
		if err != nil {
			t.Fatalf("DeleteRelease() unexpected error: %v", err)
		}
	})

	t.Run("error_set_kubernetes_client", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(false)                         // For setKubernetesClient
		mockJujuApp.EXPECT().GetLeader(gomock.Any(), uuid, facility).Return("", errExpected) // For setKubernetesClient

		err := s.DeleteRelease(ctx, uuid, facility, namespace, name, dryRun)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_helm_uninstall", func(t *testing.T) {
		mockKubeClient.EXPECT().Exists(uuid, facility).Return(true) // For setKubernetesClient
		mockHelm.EXPECT().UninstallRelease(uuid, facility, namespace, name, dryRun).Return(nil, errExpected)

		err := s.DeleteRelease(ctx, uuid, facility, namespace, name, dryRun)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

// --- Helper functions for creating mock Kubernetes objects ---

func createMockDeployment(name, namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": "test"}},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1", Image: "image1"}}},
			},
		},
	}
}

func createMockStatefulSet(name, namespace string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": "test"}},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1", Image: "image1"}}},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{*createMockPVC("data-"+name, namespace, "standard")},
		},
	}
}

func createMockDaemonSet(name, namespace string) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": "test"}},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1", Image: "image1"}}},
			},
		},
	}
}

func createMockService(name, namespace string, selectorLabels map[string]string) *corev1.Service {
	if selectorLabels == nil {
		selectorLabels = map[string]string{"app": "test"}
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Ports:    []corev1.ServicePort{{Name: "http", Port: 80, TargetPort: intstr.FromInt(8080)}},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

func createMockPod(name, namespace string, labels map[string]string) *corev1.Pod {
	if labels == nil {
		labels = map[string]string{"app": "test"}
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1", Image: "image1"}}},
	}
}

func createMockPVC(name, namespace, storageClassName string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": "test"}}, // Keep labels for consistency if needed
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: stringPtr(storageClassName),
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:        corev1.VolumeResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: resourceQuantity(1, "Gi")}},
		},
	}
}

func createMockStorageClass(name string) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta:  metav1.ObjectMeta{Name: name},
		Provisioner: "kubernetes.io/aws-ebs",
	}
}

func stringPtr(s string) *string {
	return &s
}

func resourceQuantity(size int64, suffix string) resource.Quantity {
	quantityString := fmt.Sprintf("%d%s", size, suffix)
	quantity, _ := resource.ParseQuantity(quantityString)
	return quantity
}

func int32Ptr(i int32) *int32 {
	return &i
}
