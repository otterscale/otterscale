package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/base"
	"github.com/openhdc/otterscale/internal/domain/model"
)

// --- Mocks ---

type mockServer struct {
	getFn    func(ctx context.Context, key string) ([]byte, error)
	updateFn func(ctx context.Context, key string, value string) error
}

func (m *mockServer) Get(ctx context.Context, key string) ([]byte, error) {
	return m.getFn(ctx, key)
}
func (m *mockServer) Update(ctx context.Context, key string, value string) error {
	return m.updateFn(ctx, key, value)
}

type mockPackageRepository struct {
	listFn   func(ctx context.Context) ([]model.PackageRepository, error)
	updateFn func(ctx context.Context, id int, params *model.PackageRepositoryParams) (*model.PackageRepository, error)
}

func (m *mockPackageRepository) List(ctx context.Context) ([]model.PackageRepository, error) {
	return m.listFn(ctx)
}
func (m *mockPackageRepository) Update(ctx context.Context, id int, params *model.PackageRepositoryParams) (*model.PackageRepository, error) {
	return m.updateFn(ctx, id, params)
}

type mockBootResource struct {
	listFn        func(ctx context.Context) ([]model.BootImage, error)
	importFn      func(ctx context.Context) error
	isImportingFn func(ctx context.Context) (bool, error)
}

func (m *mockBootResource) List(ctx context.Context) ([]entity.BootResource, error) {
	bootImages, err := m.listFn(ctx)
	if err != nil {
		return nil, err
	}
	bootResources := make([]entity.BootResource, len(bootImages))
	for i, img := range bootImages {
		bootResources[i] = entity.BootResource{
			Name: img.Name,
			// Add more field mappings if needed
		}
	}
	return bootResources, nil
}
func (m *mockBootResource) Import(ctx context.Context) error {
	return m.importFn(ctx)
}
func (m *mockBootResource) IsImporting(ctx context.Context) (bool, error) {
	return m.isImportingFn(ctx)
}

type mockBootSource struct {
	listFn func(ctx context.Context) ([]entity.BootSource, error)
}

func (m *mockBootSource) List(ctx context.Context) ([]entity.BootSource, error) {
	return m.listFn(ctx)
}

type mockBootSourceSelection struct {
	createFn func(ctx context.Context, distro string, arch []string) (*entity.BootSourceSelection, error)
	listFn   func(ctx context.Context, id int) ([]entity.BootSourceSelection, error)
}

func (m *mockBootSourceSelection) CreateFromMAASIO(ctx context.Context, distro string, arch []string) (*entity.BootSourceSelection, error) {
	return m.createFn(ctx, distro, arch)
}
func (m *mockBootSourceSelection) List(ctx context.Context, id int) ([]entity.BootSourceSelection, error) {
	return m.listFn(ctx, id)
}

type mockScope struct {
	listFn   func(ctx context.Context) ([]base.UserModelSummary, error)
	createFn func(ctx context.Context, name string) (*base.ModelInfo, error)
}

func (m *mockScope) List(ctx context.Context) ([]base.UserModelSummary, error) {
	return m.listFn(ctx)
}

func (m *mockScope) Create(ctx context.Context, name string) (*base.ModelInfo, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name)
	}
	return nil, nil
}

type mockScopeConfig struct {
	setFn   func(ctx context.Context, uuid string, cfg map[string]any) error
	listFn  func(ctx context.Context, uuid string) (map[string]any, error)
	unsetFn func(ctx context.Context, uuid string, keys ...string) error
}

func (m *mockScopeConfig) Set(ctx context.Context, uuid string, cfg map[string]any) error {
	return m.setFn(ctx, uuid, cfg)
}

func (m *mockScopeConfig) List(ctx context.Context, uuid string) (map[string]any, error) {
	if m.listFn != nil {
		return m.listFn(ctx, uuid)
	}
	return map[string]any{}, nil
}

func (m *mockScopeConfig) Unset(ctx context.Context, uuid string, keys ...string) error {
	if m.unsetFn != nil {
		return m.unsetFn(ctx, uuid, keys...)
	}
	return nil
}

// --- NexusService factory for tests ---

// newTestNexusServiceConfiguration creates and returns a new NexusService instance with mock dependencies for testing purposes.
func newTestNexusServiceConfiguration() *NexusService {
	return &NexusService{
		server:              &mockServer{},
		packageRepository:   &mockPackageRepository{},
		bootResource:        &mockBootResource{},
		bootSource:          &mockBootSource{},
		bootSourceSelection: &mockBootSourceSelection{},
		scope:               &mockScope{},
		scopeConfig:         &mockScopeConfig{},
	}
}

// --- Tests ---

func TestRemoveQuotes(t *testing.T) {
	s := `"quoted"`
	want := "quoted"
	got := removeQuotes(s)
	if got != want {
		t.Errorf("removeQuotes(%q) = %q, want %q", s, got, want)
	}
}

func TestListNTPServers_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"ntp1 ntp2"`), nil
		},
	}
	got, err := ns.listNTPServers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ntp1", "ntp2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestListNTPServers_Error(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := ns.listNTPServers(context.Background())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateNTPServer_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	called := false
	ns.server = &mockServer{
		updateFn: func(ctx context.Context, key, value string) error {
			called = true
			if key != maasConfigNTPServers {
				t.Errorf("unexpected key: %s", key)
			}
			if value != "ntp1 ntp2" {
				t.Errorf("unexpected value: %s", value)
			}
			return nil
		},
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"ntp1 ntp2"`), nil
		},
	}
	got, err := ns.UpdateNTPServer(context.Background(), []string{"ntp1", "ntp2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected update to be called")
	}
	if !reflect.DeepEqual(got, []string{"ntp1", "ntp2"}) {
		t.Errorf("got %v, want %v", got, []string{"ntp1", "ntp2"})
	}
}

func TestUpdateNTPServer_UpdateError(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		updateFn: func(ctx context.Context, key, value string) error {
			return errors.New("fail")
		},
	}
	_, err := ns.UpdateNTPServer(context.Background(), []string{"ntp1"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdatePackageRepository_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.packageRepository = &mockPackageRepository{
		updateFn: func(ctx context.Context, id int, params *model.PackageRepositoryParams) (*model.PackageRepository, error) {
			return &model.PackageRepository{ID: id, URL: params.URL}, nil
		},
	}
	ns.scope = &mockScope{
		listFn: func(ctx context.Context) ([]base.UserModelSummary, error) {
			return []model.Scope{{UUID: "uuid1"}}, nil
		},
	}
	ns.scopeConfig = &mockScopeConfig{
		setFn: func(ctx context.Context, uuid string, cfg map[string]any) error {
			if uuid != "uuid1" {
				t.Errorf("unexpected uuid: %s", uuid)
			}
			if cfg[jujuConfigAPTMirror] != "http://repo" {
				t.Errorf("unexpected cfg: %v", cfg)
			}
			return nil
		},
	}
	pr, err := ns.UpdatePackageRepository(context.Background(), 1, "http://repo", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.ID != 1 || pr.URL != "http://repo" {
		t.Errorf("unexpected result: %+v", pr)
	}
}

func TestUpdatePackageRepository_SkipJuju(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.packageRepository = &mockPackageRepository{
		updateFn: func(ctx context.Context, id int, params *model.PackageRepositoryParams) (*model.PackageRepository, error) {
			return &model.PackageRepository{ID: id, URL: params.URL}, nil
		},
	}
	pr, err := ns.UpdatePackageRepository(context.Background(), 2, "http://repo2", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.ID != 2 || pr.URL != "http://repo2" {
		t.Errorf("unexpected result: %+v", pr)
	}
}

func TestCreateBootImage_DefaultArch(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.bootSourceSelection = &mockBootSourceSelection{
		createFn: func(ctx context.Context, distro string, arch []string) (*entity.BootSourceSelection, error) {
			if distro != "focal" {
				t.Errorf("unexpected distro: %s", distro)
			}
			if !reflect.DeepEqual(arch, []string{"amd64"}) {
				t.Errorf("unexpected arch: %v", arch)
			}
			return &entity.BootSourceSelection{
				Release: "focal",
				OS:      "Ubuntu",
				Arches:  []string{"amd64", "arm64"},
			}, nil
		},
	}
	img, err := ns.CreateBootImage(context.Background(), "focal", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if img.DistroSeries != "focal" || img.Name != "Ubuntu" {
		t.Errorf("unexpected image: %+v", img)
	}
	if len(img.ArchitectureStatusMap) != 2 {
		t.Errorf("unexpected arch map: %+v", img.ArchitectureStatusMap)
	}
}

func TestSetDefaultBootImage_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	calls := []string{}
	ns.server = &mockServer{
		updateFn: func(ctx context.Context, key, value string) error {
			calls = append(calls, key+":"+value)
			return nil
		},
	}
	err := ns.SetDefaultBootImage(context.Background(), "focal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 3 {
		t.Errorf("expected 3 updates, got %d", len(calls))
	}
}

func TestSetDefaultBootImage_Error(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		updateFn: func(ctx context.Context, key, value string) error {
			return errors.New("fail")
		},
	}
	err := ns.SetDefaultBootImage(context.Background(), "focal")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestImportBootImages(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	called := false
	ns.bootResource = &mockBootResource{
		importFn: func(ctx context.Context) error {
			called = true
			return nil
		},
	}
	err := ns.ImportBootImages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected import to be called")
	}
}

func TestIsImportingBootImages(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.bootResource = &mockBootResource{
		isImportingFn: func(ctx context.Context) (bool, error) {
			return true, nil
		},
	}
	val, err := ns.IsImportingBootImages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !val {
		t.Error("expected true")
	}
}

func TestListBootImageSelections(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	got, err := ns.ListBootImageSelections(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(ubuntuDistroSeriesMap) {
		t.Errorf("expected %d, got %d", len(ubuntuDistroSeriesMap), len(got))
	}
}

func TestListPackageRepositories(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.packageRepository = &mockPackageRepository{
		listFn: func(ctx context.Context) ([]model.PackageRepository, error) {
			return []model.PackageRepository{{ID: 1, URL: "http://repo"}}, nil
		},
	}
	got, err := ns.listPackageRepositories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestImageBase_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"focal"`), nil
		},
	}
	b, err := ns.imageBase(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.OS != "focal" {
		t.Errorf("unexpected base: %+v", b)
	}
}

func TestImageBase_Error(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := ns.imageBase(context.Background())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestListBootImages_Error(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := ns.listBootImages(context.Background())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestListBootImages_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"focal"`), nil
		},
	}
	ns.bootResource = &mockBootResource{
		listFn: func(ctx context.Context) ([]model.BootImage, error) {
			return []model.BootImage{
				{Name: "boot/focal/amd64", ArchitectureStatusMap: map[string]string{"amd64": "true"}},
			}, nil
		},
	}
	ns.bootSource = &mockBootSource{
		listFn: func(ctx context.Context) ([]entity.BootSource, error) {
			return []entity.BootSource{{ID: 1, URL: "src1"}}, nil
		},
	}
	ns.bootSourceSelection = &mockBootSourceSelection{
		listFn: func(ctx context.Context, id int) ([]entity.BootSourceSelection, error) {
			return []entity.BootSourceSelection{
				{Release: "focal"},
			}, nil
		},
	}
	got, err := ns.listBootImages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected boot images")
	}
	if !got[0].Default {
		t.Error("expected default to be true")
	}
	if !strings.Contains(got[0].Source, "src1") {
		t.Errorf("unexpected source: %s", got[0].Source)
	}
}

func TestGetConfiguration_Success(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"ntp1 ntp2"`), nil
		},
	}
	ns.packageRepository = &mockPackageRepository{
		listFn: func(ctx context.Context) ([]model.PackageRepository, error) {
			return []model.PackageRepository{{ID: 1, URL: "http://repo"}}, nil
		},
	}
	ns.bootResource = &mockBootResource{
		listFn: func(ctx context.Context) ([]model.BootImage, error) {
			return []model.BootImage{{DistroSeries: "focal", Name: "Ubuntu"}}, nil
		},
	}
	ns.bootSource = &mockBootSource{
		listFn: func(ctx context.Context) ([]entity.BootSource, error) {
			return []entity.BootSource{
				{ID: 1, URL: "http://boot-source"},
			}, nil
		},
	}
	ns.bootSourceSelection = &mockBootSourceSelection{
		listFn: func(ctx context.Context, id int) ([]entity.BootSourceSelection, error) {
			return []entity.BootSourceSelection{
				{Release: "focal"},
			}, nil
		},
	}
	got, err := ns.GetConfiguration(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.NTPServers, []string{"ntp1", "ntp2"}) {
		t.Errorf("unexpected NTPServers: %v", got.NTPServers)
	}
	if len(got.PackageRepositories) != 1 || got.PackageRepositories[0].ID != 1 {
		t.Errorf("unexpected PackageRepositories: %v", got.PackageRepositories)
	}
	if len(got.BootImages) != 1 || got.BootImages[0].DistroSeries != "focal" {
		t.Errorf("unexpected BootImages: %v", got.BootImages)
	}
}

func TestGetConfiguration_NTPError(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return nil, errors.New("fail ntp")
		},
	}
	_, err := ns.GetConfiguration(context.Background())
	if err == nil || !strings.Contains(err.Error(), "fail ntp") {
		t.Errorf("expected ntp error, got %v", err)
	}
}

func TestGetConfiguration_PackageRepoError(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"ntp1"`), nil
		},
	}
	ns.packageRepository = &mockPackageRepository{
		listFn: func(ctx context.Context) ([]model.PackageRepository, error) {
			return nil, errors.New("fail repo")
		},
	}
	_, err := ns.GetConfiguration(context.Background())
	if err == nil || !strings.Contains(err.Error(), "fail repo") {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestGetConfiguration_BootImagesError(t *testing.T) {
	ns := newTestNexusServiceConfiguration()
	ns.server = &mockServer{
		getFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`"ntp1"`), nil
		},
	}
	ns.packageRepository = &mockPackageRepository{
		listFn: func(ctx context.Context) ([]model.PackageRepository, error) {
			return []model.PackageRepository{}, nil
		},
	}
	ns.bootResource = &mockBootResource{
		listFn: func(ctx context.Context) ([]model.BootImage, error) {
			return nil, errors.New("fail boot")
		},
	}
	_, err := ns.GetConfiguration(context.Background())
	if err == nil || !strings.Contains(err.Error(), "fail boot") {
		t.Errorf("expected boot error, got %v", err)
	}
}
