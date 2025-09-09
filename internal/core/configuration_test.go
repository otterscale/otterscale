package core

import (
	"context"
	"errors"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/stretchr/testify/assert"
)

// Mock ServerRepo
type confMockServerRepo struct {
	values map[string]string
}

func (m *confMockServerRepo) Get(ctx context.Context, name string) (string, error) {
	if v, ok := m.values[name]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}

func (m *confMockServerRepo) Update(ctx context.Context, name, value string) error {
	m.values[name] = value
	return nil
}

// Mock ScopeRepo
type confMockScopeRepo struct{}

func (m *confMockScopeRepo) List(ctx context.Context) ([]Scope, error) {
	return []Scope{{UUID: "uuid-1"}}, nil
}

func (m *confMockScopeRepo) Create(ctx context.Context, name string) (*Scope, error) {
	return &Scope{Name: name, UUID: "uuid-2"}, nil
}

// Mock ScopeConfigRepo
type mockScopeConfigRepo struct {
	setCalled bool
}

func (m *mockScopeConfigRepo) List(ctx context.Context, uuid string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (m *mockScopeConfigRepo) Set(ctx context.Context, uuid string, config map[string]any) error {
	m.setCalled = true
	return nil
}

func (m *mockScopeConfigRepo) Unset(ctx context.Context, uuid string, keys ...string) error {
	return nil
}

// Mock BootResourceRepo
type mockBootResourceRepo struct{}

func (m *mockBootResourceRepo) List(ctx context.Context) ([]entity.BootResource, error) {
	return []entity.BootResource{
		{Name: "ubuntu/focal", Architecture: "amd64", Type: "release"},
	}, nil
}

func (m *mockBootResourceRepo) Import(ctx context.Context) error {
	return nil
}

func (m *mockBootResourceRepo) IsImporting(ctx context.Context) (bool, error) {
	return false, nil
}

// Mock BootSourceRepo
type mockBootSourceRepo struct{}

func (m *mockBootSourceRepo) List(ctx context.Context) ([]entity.BootSource, error) {
	return []entity.BootSource{
		{ID: 1, URL: "http://example.com"},
	}, nil
}

// Mock BootSourceSelectionRepo
type mockBootSourceSelectionRepo struct{}

func (m *mockBootSourceSelectionRepo) List(ctx context.Context, id int) ([]entity.BootSourceSelection, error) {
	return []entity.BootSourceSelection{
		{Release: "focal", OS: "ubuntu", Arches: []string{"amd64"}},
	}, nil
}

func (m *mockBootSourceSelectionRepo) Create(ctx context.Context, bootSourceID int, params *entity.BootSourceSelectionParams) (*entity.BootSourceSelection, error) {
	return &entity.BootSourceSelection{
		Release: params.Release,
		OS:      params.OS,
		Arches:  params.Arches,
	}, nil
}

// Mock PackageRepositoryRepo
type mockPackageRepositoryRepo struct{}

func (m *mockPackageRepositoryRepo) List(ctx context.Context) ([]PackageRepository, error) {
	return []PackageRepository{{ID: 1, Name: "main", URL: "http://archive.ubuntu.com/ubuntu"}}, nil
}

func (m *mockPackageRepositoryRepo) Update(ctx context.Context, id int, params *entity.PackageRepositoryParams) (*PackageRepository, error) {
	return &PackageRepository{ID: id, Name: "main", URL: params.URL}, nil
}

func TestConfigurationUseCase_GetConfiguration(t *testing.T) {
	server := &confMockServerRepo{values: map[string]string{
		"ntp_servers":           "0.pool.ntp.org 1.pool.ntp.org",
		"default_distro_series": "focal",
	}}
	uc := NewConfigurationUseCase(
		server,
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	cfg, err := uc.GetConfiguration(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, cfg.NTPServers, "0.pool.ntp.org")
	assert.NotEmpty(t, cfg.PackageRepositories)
	assert.NotEmpty(t, cfg.BootImages)
}

func TestConfigurationUseCase_UpdateNTPServer(t *testing.T) {
	server := &confMockServerRepo{values: map[string]string{"ntp_servers": "0.pool.ntp.org"}}
	uc := NewConfigurationUseCase(
		server,
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	ntps, err := uc.UpdateNTPServer(context.Background(), []string{"1.pool.ntp.org", "2.pool.ntp.org"})
	assert.NoError(t, err)
	assert.Contains(t, ntps, "1.pool.ntp.org")
}

func TestConfigurationUseCase_UpdatePackageRepository(t *testing.T) {
	server := &confMockServerRepo{values: map[string]string{}}
	scopeConfig := &mockScopeConfigRepo{}
	uc := NewConfigurationUseCase(
		server,
		&confMockScopeRepo{},
		scopeConfig,
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	repo, err := uc.UpdatePackageRepository(context.Background(), 1, "http://mirror.test/ubuntu", false)
	assert.NoError(t, err)
	assert.Equal(t, "http://mirror.test/ubuntu", repo.URL)
	assert.True(t, scopeConfig.setCalled)
}

func TestConfigurationUseCase_CreateBootImage(t *testing.T) {
	uc := NewConfigurationUseCase(
		&confMockServerRepo{values: map[string]string{}},
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	img, err := uc.CreateBootImage(context.Background(), "focal", []string{"amd64"})
	assert.NoError(t, err)
	assert.Equal(t, "focal", img.DistroSeries)
	assert.Equal(t, "ubuntu", img.Name)
	assert.Contains(t, img.ArchitectureStatusMap, "amd64")
}

func TestConfigurationUseCase_SetDefaultBootImage(t *testing.T) {
	server := &confMockServerRepo{values: map[string]string{}}
	uc := NewConfigurationUseCase(
		server,
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	err := uc.SetDefaultBootImage(context.Background(), "jammy")
	assert.NoError(t, err)
	assert.Equal(t, "ubuntu", server.values["default_osystem"])
	assert.Equal(t, "jammy", server.values["default_distro_series"])
	assert.Equal(t, "jammy", server.values["commissioning_distro_series"])
}

func TestConfigurationUseCase_ImportBootImages(t *testing.T) {
	uc := NewConfigurationUseCase(
		&confMockServerRepo{values: map[string]string{}},
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	err := uc.ImportBootImages(context.Background())
	assert.NoError(t, err)
}

func TestConfigurationUseCase_IsImportingBootImages(t *testing.T) {
	uc := NewConfigurationUseCase(
		&confMockServerRepo{values: map[string]string{}},
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	importing, err := uc.IsImportingBootImages(context.Background())
	assert.NoError(t, err)
	assert.False(t, importing)
}

func TestConfigurationUseCase_ListBootImageSelections(t *testing.T) {
	uc := NewConfigurationUseCase(
		&confMockServerRepo{values: map[string]string{}},
		&confMockScopeRepo{},
		&mockScopeConfigRepo{},
		&mockBootResourceRepo{},
		&mockBootSourceRepo{},
		&mockBootSourceSelectionRepo{},
		&mockPackageRepositoryRepo{},
	)
	selections, err := uc.ListBootImageSelections()
	assert.NoError(t, err)
	assert.NotEmpty(t, selections)
}
