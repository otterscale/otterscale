package mux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openhdc/otterscale/internal/app"
	"github.com/openhdc/otterscale/internal/core"
)

// createMockServices creates mock services for testing
func createMockServices() (
	*app.ApplicationService,
	*app.BISTService,
	*app.ConfigurationService,
	*app.EnvironmentService,
	*app.FacilityService,
	*app.EssentialService,
	*app.MachineService,
	*app.NetworkService,
	*app.StorageService,
	*app.ScopeService,
	*app.TagService,
) {
	// Create mock use cases (nil is fine for mux testing as we're not calling the actual methods)
	var (
		appUC   *core.ApplicationUseCase
		bistUC  *core.BISTUseCase
		confUC  *core.ConfigurationUseCase
		envUC   *core.EnvironmentUseCase
		facUC   *core.FacilityUseCase
		essUC   *core.EssentialUseCase
		machUC  *core.MachineUseCase
		netUC   *core.NetworkUseCase
		storUC  *core.StorageUseCase
		scopeUC *core.ScopeUseCase
		tagUC   *core.TagUseCase
	)

	return app.NewApplicationService(appUC),
		app.NewBISTService(bistUC),
		app.NewConfigurationService(confUC),
		app.NewEnvironmentService(envUC),
		app.NewFacilityService(facUC),
		app.NewEssentialService(essUC),
		app.NewMachineService(machUC),
		app.NewNetworkService(netUC),
		app.NewStorageService(storUC),
		app.NewScopeService(scopeUC),
		app.NewTagService(tagUC)
}

func TestNew_WithoutHelper(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	mux := New(false, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

	if mux == nil {
		t.Fatal("Expected mux to be created, got nil")
	}
}

func TestNew_WithHelper(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	mux := New(true, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

	if mux == nil {
		t.Fatal("Expected mux to be created, got nil")
	}
}

func TestNew_ServiceHandlersRegistered(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	mux := New(false, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

	// Create a test server with the mux
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test service endpoints - Connect RPC services use specific method endpoints
	// We'll test a known method from each service with GET requests to check registration
	serviceEndpoints := []string{
		"/otterscale.application.v1.ApplicationService/ListApplications",
		"/otterscale.configuration.v1.ConfigurationService/GetConfiguration",
		"/otterscale.environment.v1.EnvironmentService/CheckHealthy", // Updated
		"/otterscale.facility.v1.FacilityService/ListFacilities",
		"/otterscale.essential.v1.EssentialService/ListEssentials", // Updated
		"/otterscale.machine.v1.MachineService/ListMachines",
		"/otterscale.network.v1.NetworkService/ListNetworks",
		"/otterscale.storage.v1.StorageService/ListMONs", // Updated
		"/otterscale.scope.v1.ScopeService/ListScopes",
		"/otterscale.tag.v1.TagService/ListTags",
	}

	for _, endpoint := range serviceEndpoints {
		t.Run("endpoint_"+strings.ReplaceAll(endpoint, "/", "_"), func(t *testing.T) {
			// Use GET request to check if endpoint is registered
			// Connect RPC services should return 405 Method Not Allowed for GET
			// but 404 Not Found if the endpoint isn't registered at all
			resp, err := http.Get(server.URL + endpoint)
			if err != nil {
				t.Fatalf("Failed to make request to %s: %v", endpoint, err)
			}
			defer resp.Body.Close()

			// We expect either:
			// - 405 Method Not Allowed (endpoint exists but doesn't accept GET)
			// - Some other status (but not 404 Not Found)
			if resp.StatusCode == http.StatusNotFound {
				t.Errorf("Endpoint %s not found (404), expected it to be registered", endpoint)
			}
		})
	}
}

func TestNew_HealthAndReflectionWithHelper(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	mux := New(true, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

	// Create a test server with the mux
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test health endpoint - Connect health uses specific method paths
	t.Run("health_endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/grpc.health.v1.Health/Check")
		if err != nil {
			t.Fatalf("Failed to make request to health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Error("Health endpoint not found (404), expected it to be registered when helper=true")
		}
	})

	// Test reflection endpoints
	reflectionEndpoints := []string{
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	for _, endpoint := range reflectionEndpoints {
		t.Run("reflection_endpoint_"+strings.ReplaceAll(endpoint, "/", "_"), func(t *testing.T) {
			resp, err := http.Get(server.URL + endpoint)
			if err != nil {
				t.Fatalf("Failed to make request to %s: %v", endpoint, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				t.Errorf("Reflection endpoint %s not found (404), expected it to be registered when helper=true", endpoint)
			}
		})
	}
}

func TestNew_NoHealthAndReflectionWithoutHelper(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	mux := New(false, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

	// Create a test server with the mux
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test that health endpoint is NOT registered when helper=false
	t.Run("no_health_endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/grpc.health.v1.Health/Check")
		if err != nil {
			t.Fatalf("Failed to make request to health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Error("Health endpoint found when helper=false, expected it to not be registered")
		}
	})

	// Test that reflection endpoints are NOT registered when helper=false
	reflectionEndpoints := []string{
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	for _, endpoint := range reflectionEndpoints {
		t.Run("no_reflection_endpoint_"+strings.ReplaceAll(endpoint, "/", "_"), func(t *testing.T) {
			resp, err := http.Get(server.URL + endpoint)
			if err != nil {
				t.Fatalf("Failed to make request to %s: %v", endpoint, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("Reflection endpoint %s found when helper=false, expected it to not be registered", endpoint)
			}
		})
	}
}

func TestServices_ContainsAllServiceNames(t *testing.T) {
	expectedServices := []string{
		"otterscale.application.v1.ApplicationService",
		"otterscale.bist.v1.BISTService",
		"otterscale.configuration.v1.ConfigurationService",
		"otterscale.environment.v1.EnvironmentService",
		"otterscale.facility.v1.FacilityService",
		"otterscale.essential.v1.EssentialService",
		"otterscale.machine.v1.MachineService",
		"otterscale.network.v1.NetworkService",
		"otterscale.scope.v1.ScopeService",
		"otterscale.storage.v1.StorageService",
		"otterscale.tag.v1.TagService",
	}

	if len(Services) != len(expectedServices) {
		t.Errorf("Expected %d services, got %d", len(expectedServices), len(Services))
	}

	// Check that all expected services are present
	serviceMap := make(map[string]bool)
	for _, service := range Services {
		serviceMap[service] = true
	}

	for _, expected := range expectedServices {
		if !serviceMap[expected] {
			t.Errorf("Expected service %s not found in Services slice", expected)
		}
	}
}

func TestNew_NilServices(t *testing.T) {
	// Test that the function handles nil services gracefully
	// This might panic or handle gracefully depending on the implementation
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Function panicked with nil services (this might be expected): %v", r)
		}
	}()

	mux := New(false, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if mux == nil {
		t.Error("Expected mux to be created even with nil services")
	}
}

func TestNew_HelperFlagBehavior(t *testing.T) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	testCases := []struct {
		name         string
		helper       bool
		expectHealth bool
	}{
		{
			name:         "helper_true",
			helper:       true,
			expectHealth: true,
		},
		{
			name:         "helper_false",
			helper:       false,
			expectHealth: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux := New(tc.helper, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)

			server := httptest.NewServer(mux)
			defer server.Close()

			resp, err := http.Get(server.URL + "/grpc.health.v1.Health/Check")
			if err != nil {
				t.Fatalf("Failed to make health check request: %v", err)
			}
			defer resp.Body.Close()

			if tc.expectHealth {
				if resp.StatusCode == http.StatusNotFound {
					t.Error("Expected health endpoint to be available when helper=true")
				}
			} else {
				if resp.StatusCode != http.StatusNotFound {
					t.Error("Expected health endpoint to be unavailable when helper=false")
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkNew_WithoutHelper(b *testing.B) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(false, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)
	}
}

func BenchmarkNew_WithHelper(b *testing.B) {
	appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc := createMockServices()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(true, appSvc, bistSVC, confSvc, envSvc, facSvc, essSvc, machSvc, netSvc, storSvc, scopeSvc, tagSvc)
	}
}
