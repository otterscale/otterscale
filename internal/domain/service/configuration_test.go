package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/base"
	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNexusService_GetConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockMAASServer(ctrl)
	mockPackageRepo := mocks.NewMockMAASPackageRepository(ctrl)
	mockBootResource := mocks.NewMockMAASBootResource(ctrl)
	mockBootSource := mocks.NewMockMAASBootSource(ctrl)
	mockBootSourceSelection := mocks.NewMockMAASBootSourceSelection(ctrl)

	s := &NexusService{
		server:              mockServer,
		packageRepository:   mockPackageRepo,
		bootResource:        mockBootResource,
		bootSource:          mockBootSource,
		bootSourceSelection: mockBootSourceSelection,
	}

	ctx := context.Background()

	mockServer.EXPECT().Get(ctx, maasConfigNTPServers).Return([]byte(`"ntp1.example.com ntp2.example.com"`), nil)
	mockPackageRepo.EXPECT().List(ctx).Return([]entity.PackageRepository{{URL: "http://repo1.example.com/"}}, nil)
	mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil)
	mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{}, nil)
	mockBootSourceSelection.EXPECT().List(ctx, gomock.Any()).Return([]entity.BootSourceSelection{}, nil).AnyTimes()
	mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"jammy"`), nil)

	cfg, err := s.GetConfiguration(ctx)
	if err != nil {
		t.Fatalf("GetConfiguration failed: %v", err)
	}

	expectedNTPServers := []string{"ntp1.example.com", "ntp2.example.com"}
	if !reflect.DeepEqual(cfg.NTPServers, expectedNTPServers) { // 使用 cfg 變數
		t.Errorf("NTPServers mismatch: expected %v, got %v", expectedNTPServers, cfg.NTPServers)
	}

	expectedPackageRepo := []model.PackageRepository{{URL: "http://repo1.example.com/"}}
	if !reflect.DeepEqual(cfg.PackageRepositories, expectedPackageRepo) { // 使用 cfg 變數
		t.Errorf("PackageRepositories mismatch: expected %v, got %v", expectedPackageRepo, cfg.PackageRepositories)
	}

	//  新增 BootImages 的斷言，即使它是空的，也要驗證
	if len(cfg.BootImages) != 0 { // 使用 cfg 變數
		t.Errorf("BootImages mismatch: expected empty slice, got %v", cfg.BootImages)
	}
}
func TestUpdateNTPServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockMAASServer(ctrl)
	s := &NexusService{server: mockServer}
	ctx := context.Background()
	addresses := []string{"ntp3.example.com", "ntp4.example.com"}

	t.Run("success", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigNTPServers, "ntp3.example.com ntp4.example.com").Return(nil)
		mockServer.EXPECT().Get(ctx, maasConfigNTPServers).Return([]byte(`"ntp3.example.com ntp4.example.com"`), nil)

		updatedAddresses, err := s.UpdateNTPServer(ctx, addresses)
		if err != nil {
			t.Fatalf("UpdateNTPServer failed: %v", err)
		}
		if !reflect.DeepEqual(updatedAddresses, addresses) {
			t.Errorf("expected addresses %v, got %v", addresses, updatedAddresses)
		}
	})

	t.Run("update error", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigNTPServers, "ntp3.example.com ntp4.example.com").Return(errors.New("update failed"))

		_, err := s.UpdateNTPServer(ctx, addresses)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "update failed" {
			t.Errorf("expected 'update failed', got %v", err)
		}
	})

	t.Run("list after update error", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigNTPServers, "ntp3.example.com ntp4.example.com").Return(nil)
		mockServer.EXPECT().Get(ctx, maasConfigNTPServers).Return(nil, errors.New("list failed"))

		_, err := s.UpdateNTPServer(ctx, addresses)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "list failed" {
			t.Errorf("expected 'list failed', got %v", err)
		}
	})
}

func TestCreateBootImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBootSourceSelection := mocks.NewMockMAASBootSourceSelection(ctrl)
	s := &NexusService{bootSourceSelection: mockBootSourceSelection}
	ctx := context.Background()
	distroSeries := "jammy"
	architectures := []string{"amd64", "arm64"}

	t.Run("success", func(t *testing.T) {
		mockBootSourceSelection.EXPECT().CreateFromMAASIO(ctx, distroSeries, architectures).Return(&entity.BootSourceSelection{}, nil)
		_, err := s.CreateBootImage(ctx, distroSeries, architectures)
		if err != nil {
			t.Fatalf("CreateBootImage failed: %v", err)
		}
	})

	t.Run("create error", func(t *testing.T) {
		mockBootSourceSelection.EXPECT().CreateFromMAASIO(ctx, distroSeries, architectures).Return(nil, errors.New("create failed"))
		_, err := s.CreateBootImage(ctx, distroSeries, architectures)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "create failed" {
			t.Errorf("expected 'create failed', got %v", err)
		}
	})
}
func TestSetDefaultBootImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockMAASServer(ctrl)
	s := &NexusService{server: mockServer}
	ctx := context.Background()
	distroSeries := "jammy"

	t.Run("success", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem).Return(nil)
		mockServer.EXPECT().Update(ctx, maasConfigDefaultDistroSeries, distroSeries).Return(nil)
		mockServer.EXPECT().Update(ctx, maasConfigCommissioningDistroSeries, distroSeries).Return(nil)

		if err := s.SetDefaultBootImage(ctx, distroSeries); err != nil {
			t.Fatalf("SetDefaultBootImage failed: %v", err)
		}
	})

	t.Run("error updating default os system", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem).Return(errors.New("os system update error"))
		// No other mocks should be called

		err := s.SetDefaultBootImage(ctx, distroSeries)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "os system update error" {
			t.Errorf("expected 'os system update error', got %v", err)
		}
	})

	t.Run("error updating default distro series", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem).Return(nil)
		mockServer.EXPECT().Update(ctx, maasConfigDefaultDistroSeries, distroSeries).Return(errors.New("distro series update error"))
		// No other mocks should be called

		err := s.SetDefaultBootImage(ctx, distroSeries)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "distro series update error" {
			t.Errorf("expected 'distro series update error', got %v", err)
		}
	})

	t.Run("error updating commissioning distro series", func(t *testing.T) {
		mockServer.EXPECT().Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem).Return(nil)
		mockServer.EXPECT().Update(ctx, maasConfigDefaultDistroSeries, distroSeries).Return(nil)
		mockServer.EXPECT().Update(ctx, maasConfigCommissioningDistroSeries, distroSeries).Return(errors.New("commissioning update error"))

		err := s.SetDefaultBootImage(ctx, distroSeries)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "commissioning update error" {
			t.Errorf("expected 'commissioning update error', got %v", err)
		}
	})
}

func TestUpdatePackageRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMAASPackageRepository(ctrl)
	s := &NexusService{packageRepository: mockRepo}
	ctx := context.Background()
	id := 1
	url := "http://newrepo.example.com/"
	skipJuju := true // skipJuju is not used by the service method, but kept for consistency if it was intended for other layers

	t.Run("success", func(t *testing.T) {
		expectedRepo := &entity.PackageRepository{
			ID:   id,
			URL:  url,
			Name: "newrepo.example.com", // Assuming name is derived or set
		}
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(expectedRepo, nil)

		repo, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err != nil {
			t.Fatalf("UpdatePackageRepository failed: %v", err)
		}
		if repo == nil {
			t.Fatal("expected repository, got nil")
		}
		if repo.ID != expectedRepo.ID {
			t.Errorf("expected ID %d, got %d", expectedRepo.ID, repo.ID)
		}
		if repo.URL != expectedRepo.URL {
			t.Errorf("expected URL %s, got %s", expectedRepo.URL, repo.URL)
		}
		if repo.Name != expectedRepo.Name {
			t.Errorf("expected Name %s, got %s", expectedRepo.Name, repo.Name)
		}
	})

	t.Run("update error", func(t *testing.T) {
		expectedErr := errors.New("update failed")
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(nil, expectedErr)

		_, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
	t.Run("update error with nil repository", func(t *testing.T) {
		expectedErr := errors.New("update failed")
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(nil, expectedErr)

		repo, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if repo != nil {
			t.Fatal("expected nil repository, got non-nil")
		}
	})
	t.Run("update error with empty URL", func(t *testing.T) {
		id := 1
		url := "" // Empty URL
		skipJuju := true

		expectedErr := errors.New("URL cannot be empty")
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(nil, expectedErr)

		_, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
	t.Run("success with juju update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockMAASPackageRepository(ctrl)
		mockScope := mocks.NewMockJujuModel(ctrl)             // Changed from NewMockMAASScope
		mockScopeConfig := mocks.NewMockJujuModelConfig(ctrl) // Changed from NewMockMAASScopeConfig

		s := &NexusService{
			packageRepository: mockRepo,
			scope:             mockScope,
			scopeConfig:       mockScopeConfig,
		}
		ctx := context.Background()
		id := 1
		url := "http://newrepo.example.com/"
		skipJuju := false // Test this path

		expectedRepoEntity := &entity.PackageRepository{
			ID:   id,
			URL:  url,
			Name: "newrepo.example.com",
		}
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(expectedRepoEntity, nil)

		scopes := []base.UserModelSummary{ // Changed from model.UserMAASScope
			{UUID: "uuid1", Name: "scope1"},
			{UUID: "uuid2", Name: "scope2"},
		}
		mockScope.EXPECT().List(ctx).Return(scopes, nil)

		expectedJujuConfig := map[string]any{jujuConfigAPTMirror: url}
		mockScopeConfig.EXPECT().Set(ctx, scopes[0].UUID, expectedJujuConfig).Return(nil)
		mockScopeConfig.EXPECT().Set(ctx, scopes[1].UUID, expectedJujuConfig).Return(nil)

		repo, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err != nil {
			t.Fatalf("UpdatePackageRepository failed: %v", err)
		}
		if repo == nil {
			t.Fatal("expected repository, got nil")
		}
		if repo.ID != expectedRepoEntity.ID {
			t.Errorf("expected ID %d, got %d", expectedRepoEntity.ID, repo.ID)
		}
		if repo.URL != expectedRepoEntity.URL {
			t.Errorf("expected URL %s, got %s", expectedRepoEntity.URL, repo.URL)
		}
	})

	t.Run("error listing scopes during juju update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockMAASPackageRepository(ctrl)
		mockScope := mocks.NewMockJujuModel(ctrl) // Changed from NewMockMAASScope
		// mockScopeConfig is not strictly needed here as the error happens before its use

		s := &NexusService{
			packageRepository: mockRepo,
			scope:             mockScope,
			// scopeConfig: mockScopeConfig,
		}
		ctx := context.Background()
		id := 1
		url := "http://newrepo.example.com/"
		skipJuju := false

		expectedRepoEntity := &entity.PackageRepository{
			ID:   id,
			URL:  url,
			Name: "newrepo.example.com",
		}
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(expectedRepoEntity, nil)

		expectedErr := errors.New("scope list error")
		mockScope.EXPECT().List(ctx).Return(nil, expectedErr)

		_, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("error setting scope config during juju update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockMAASPackageRepository(ctrl)
		mockScope := mocks.NewMockJujuModel(ctrl)             // Changed from NewMockMAASScope
		mockScopeConfig := mocks.NewMockJujuModelConfig(ctrl) // Changed from NewMockMAASScopeConfig

		s := &NexusService{
			packageRepository: mockRepo,
			scope:             mockScope,
			scopeConfig:       mockScopeConfig,
		}
		ctx := context.Background()
		id := 1
		url := "http://newrepo.example.com/"
		skipJuju := false

		expectedRepoEntity := &entity.PackageRepository{
			ID:   id,
			URL:  url,
			Name: "newrepo.example.com",
		}
		mockRepo.EXPECT().Update(ctx, id, &entity.PackageRepositoryParams{URL: url}).Return(expectedRepoEntity, nil)

		scopes := []base.UserModelSummary{ // Changed from model.UserMAASScope
			{UUID: "uuid1", Name: "scope1"},
		}
		mockScope.EXPECT().List(ctx).Return(scopes, nil)

		expectedJujuConfig := map[string]any{jujuConfigAPTMirror: url}
		expectedErr := errors.New("scope config set error")
		mockScopeConfig.EXPECT().Set(ctx, scopes[0].UUID, expectedJujuConfig).Return(expectedErr)

		_, err := s.UpdatePackageRepository(ctx, id, url, skipJuju)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestImportBootImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBootResource := mocks.NewMockMAASBootResource(ctrl)
	s := &NexusService{bootResource: mockBootResource}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBootResource.EXPECT().Import(ctx).Return(nil)
		if err := s.ImportBootImages(ctx); err != nil {
			t.Fatalf("ImportBootImages failed: %v", err)
		}
	})

	t.Run("import error", func(t *testing.T) {
		mockBootResource.EXPECT().Import(ctx).Return(errors.New("import failed"))
		if err := s.ImportBootImages(ctx); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestIsImportingBootImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBootResource := mocks.NewMockMAASBootResource(ctrl)
	s := &NexusService{bootResource: mockBootResource}
	ctx := context.Background()

	t.Run("is importing", func(t *testing.T) {
		mockBootResource.EXPECT().IsImporting(ctx).Return(true, nil)
		importing, err := s.IsImportingBootImages(ctx)
		if err != nil {
			t.Fatalf("IsImportingBootImages failed: %v", err)
		}
		if !importing {
			t.Errorf("IsImportingBootImages expected true, got false")
		}
	})

	t.Run("not importing", func(t *testing.T) {
		mockBootResource.EXPECT().IsImporting(ctx).Return(false, nil)
		importing, err := s.IsImportingBootImages(ctx)
		if err != nil {
			t.Fatalf("IsImportingBootImages failed: %v", err)
		}
		if importing {
			t.Errorf("IsImportingBootImages expected false, got true")
		}
	})

	t.Run("is importing error", func(t *testing.T) {
		mockBootResource.EXPECT().IsImporting(ctx).Return(false, errors.New("is importing error"))
		_, err := s.IsImportingBootImages(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "is importing error" {
			t.Errorf("expected 'is importing error', got %v", err)
		}
	})
}

func TestListBootImageSelections(t *testing.T) {
	s := &NexusService{} // This function does not have external dependencies on mocks for this method
	ctx := context.Background()

	selections, err := s.ListBootImageSelections(ctx)
	if err != nil {
		t.Fatalf("ListBootImageSelections failed: %v", err)
	}

	// Check if the number of selections matches the number of entries in ubuntuDistroSeriesMap
	if len(selections) != len(ubuntuDistroSeriesMap) {
		t.Errorf("Expected %d selections, got %d. Selections: %+v", len(ubuntuDistroSeriesMap), len(selections), selections)
	}

	// Verify that each entry in ubuntuDistroSeriesMap is present in the selections
	for series, expected := range ubuntuDistroSeriesMap {
		found := false
		for _, actual := range selections {
			if actual.DistroSeries == series {
				// Compare the actual selection with the expected model.BootImageSelection
				// The 'expected' variable from the map is already of type model.BootImageSelection.
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Mismatch for DistroSeries %s:\nExpected: %+v\nActual:   %+v",
						series, expected, actual)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected selection for DistroSeries %s not found in results: %+v", series, selections)
		}
	}

	// Optional: Verify that there are no unexpected selections in the results
	// This is implicitly covered if len(selections) == len(ubuntuDistroSeriesMap)
	// and all expected items are found. However, an explicit check can be added for robustness.
	for _, actual := range selections {
		if _, ok := ubuntuDistroSeriesMap[actual.DistroSeries]; !ok {
			t.Errorf("Found unexpected selection in results: %+v", actual)
		}
	}
}

func TestNexusService_listBootImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockMAASServer(ctrl)
	mockBootResource := mocks.NewMockMAASBootResource(ctrl)
	mockBootSource := mocks.NewMockMAASBootSource(ctrl)
	mockBootSourceSelection := mocks.NewMockMAASBootSourceSelection(ctrl)

	s := &NexusService{
		server:              mockServer,
		bootResource:        mockBootResource,
		bootSource:          mockBootSource,
		bootSourceSelection: mockBootSourceSelection,
	}
	ctx := context.Background()

	t.Run("success with data", func(t *testing.T) {
		defaultSeries := "focal"
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"`+defaultSeries+`"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{
			{Name: "ubuntu/focal", Architecture: "amd64", Type: "Synced"},
			{Name: "ubuntu/focal", Architecture: "arm64", Type: "Synced"},
			{Name: "ubuntu/jammy", Architecture: "amd64", Type: "Synced"},
		}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{
			{ID: 1, URL: "http://images.maas.io/ephemeral-v3/daily/"},
		}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return([]entity.BootSourceSelection{
			{Release: "focal", OS: "ubuntu", Arches: []string{"amd64", "arm64"}},
			{Release: "jammy", OS: "ubuntu", Arches: []string{"amd64"}},
		}, nil)

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}

		if len(images) != 2 {
			t.Fatalf("expected 2 images, got %d: %+v", len(images), images)
		}

		focalImageFound := false
		jammyImageFound := false
		for _, img := range images {
			if img.DistroSeries == "focal" {
				focalImageFound = true
				if img.Name != ubuntuDistroSeriesMap["focal"].Name {
					t.Errorf("expected focal name %s, got %s", ubuntuDistroSeriesMap["focal"].Name, img.Name)
				}
				if !img.Default {
					t.Error("expected focal to be default")
				}
				if len(img.ArchitectureStatusMap) != 2 {
					t.Errorf("expected 2 architectures for focal, got %d", len(img.ArchitectureStatusMap))
				}
				if _, ok := img.ArchitectureStatusMap["amd64"]; !ok {
					t.Error("expected amd64 for focal")
				}
				if _, ok := img.ArchitectureStatusMap["arm64"]; !ok {
					t.Error("expected arm64 for focal")
				}
			}
			if img.DistroSeries == "jammy" {
				jammyImageFound = true
				if img.Name != ubuntuDistroSeriesMap["jammy"].Name {
					t.Errorf("expected jammy name %s, got %s", ubuntuDistroSeriesMap["jammy"].Name, img.Name)
				}
				if img.Default {
					t.Error("expected jammy not to be default")
				}
				if len(img.ArchitectureStatusMap) != 1 {
					t.Errorf("expected 1 architecture for jammy, got %d", len(img.ArchitectureStatusMap))
				}
				if _, ok := img.ArchitectureStatusMap["amd64"]; !ok {
					t.Error("expected amd64 for jammy")
				}
			}
		}
		if !focalImageFound {
			t.Error("focal image not found")
		}
		if !jammyImageFound {
			t.Error("jammy image not found")
		}
	})

	t.Run("error getting default distro series", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return(nil, errors.New("server get error"))
		// Other mocks should not be called

		_, err := s.listBootImages(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "server get error" {
			t.Errorf("expected 'server get error', got %v", err)
		}
	})

	t.Run("error listing boot resources", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return(nil, errors.New("boot resource error"))
		// Other mocks should not be called

		_, err := s.listBootImages(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "boot resource error" {
			t.Errorf("expected 'boot resource error', got %v", err)
		}
	})

	t.Run("error listing boot sources", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil)
		mockBootSource.EXPECT().List(ctx).Return(nil, errors.New("boot source error"))
		// Other mocks should not be called

		_, err := s.listBootImages(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "boot source error" {
			t.Errorf("expected 'boot source error', got %v", err)
		}
	})

	t.Run("error listing boot source selections", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{{ID: 1}}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return(nil, errors.New("bss error"))

		_, err := s.listBootImages(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "bss error" {
			t.Errorf("expected 'bss error', got %v", err)
		}
	})

	t.Run("empty boot resources", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil) // Empty
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{{ID: 1}}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return([]entity.BootSourceSelection{
			{Release: "focal", OS: "ubuntu", Arches: []string{"amd64"}},
		}, nil)

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}
		if len(images) != 1 {
			t.Fatalf("expected 1 image, got %d", len(images))
		}
		if len(images[0].ArchitectureStatusMap) != 0 { // No matching boot resources
			t.Errorf("expected 0 architectures, got %d", len(images[0].ArchitectureStatusMap))
		}
	})

	t.Run("empty boot sources", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{}, nil) // Empty
		// bootSourceSelection.List should not be called

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}
		if len(images) != 0 {
			t.Fatalf("expected 0 images, got %d", len(images))
		}
	})

	t.Run("boot resource name without slash", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{
			{Name: "focal-custom", Architecture: "amd64"}, // No slash
		}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{{ID: 1}}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return([]entity.BootSourceSelection{
			{Release: "focal", OS: "ubuntu", Arches: []string{"amd64"}},
		}, nil)

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}
		if len(images) != 1 {
			t.Fatalf("expected 1 image, got %d", len(images))
		}
		// The architecture map will be empty because "focal-custom" won't match "focal" key in brm
		if len(images[0].ArchitectureStatusMap) != 0 {
			t.Errorf("expected 0 architectures due to name mismatch, got %d: %+v", len(images[0].ArchitectureStatusMap), images[0].ArchitectureStatusMap)
		}
	})

	t.Run("boot source selection list returns empty", func(t *testing.T) {
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"focal"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{{ID: 1}}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return([]entity.BootSourceSelection{}, nil) // Empty

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}
		if len(images) != 0 {
			t.Fatalf("expected 0 images, got %d", len(images))
		}
	})

	t.Run("distro series not in ubuntuDistroSeriesMap", func(t *testing.T) {
		customSeries := "mycustomos"
		mockServer.EXPECT().Get(ctx, maasConfigDefaultDistroSeries).Return([]byte(`"`+customSeries+`"`), nil)
		mockBootResource.EXPECT().List(ctx).Return([]entity.BootResource{
			{Name: "vendor/" + customSeries, Architecture: "amd64", Type: "Synced"},
		}, nil)
		mockBootSource.EXPECT().List(ctx).Return([]entity.BootSource{{ID: 1}}, nil)
		mockBootSourceSelection.EXPECT().List(ctx, 1).Return([]entity.BootSourceSelection{
			{Release: customSeries, OS: "vendor", Arches: []string{"amd64"}},
		}, nil)

		images, err := s.listBootImages(ctx)
		if err != nil {
			t.Fatalf("listBootImages failed: %v", err)
		}
		if len(images) != 1 {
			t.Fatalf("expected 1 image, got %d", len(images))
		}
		if images[0].Name != customSeries { // Should use brss[j].Release as Name
			t.Errorf("expected name %s, got %s", customSeries, images[0].Name)
		}
		if !images[0].Default {
			t.Error("expected customSeries to be default")
		}
		if _, ok := images[0].ArchitectureStatusMap["amd64"]; !ok {
			t.Error("expected amd64 for customSeries")
		}
	})
}
