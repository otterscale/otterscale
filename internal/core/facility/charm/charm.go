package charm

import (
	"context"
	"slices"
	"time"
)

type Charm struct {
	ID              string          `json:"id"`
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Result          CharmResult     `json:"result"`
	DefaultArtifact CharmArtifact   `json:"default-release"`
	Artifacts       []CharmArtifact `json:"channel-map"`
}

type CharmBase struct {
	Architecture string `json:"architecture"`
	Channel      string `json:"channel"`
	Name         string `json:"name"`
}

type CharmChannel struct {
	Base       CharmBase `json:"base"`
	Name       string    `json:"name"`
	ReleasedAt time.Time `json:"released-at"`
	Risk       string    `json:"risk"`
	Track      string    `json:"track"`
}

type CharmRevision struct {
	Bases     []CharmBase `json:"bases"`
	CreatedAt time.Time   `json:"created-at"`
	Revision  int         `json:"revision"`
	Version   string      `json:"version"`
}

type CharmArtifact struct {
	Channel  CharmChannel  `json:"channel"`
	Revision CharmRevision `json:"revision"`
}

type CharmResultCategory struct {
	Featured bool   `json:"featured"`
	Name     string `json:"name"`
}

type CharmResultMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type CharmResultPublisher struct {
	DisplayName string `json:"display-name"`
	ID          string `json:"id"`
	Username    string `json:"username"`
	Validation  string `json:"validation"`
}

type CharmResult struct {
	BugsURL      string                `json:"bugs-url"`
	Categories   []CharmResultCategory `json:"categories"`
	DeployableOn []string              `json:"deployable-on"`
	Description  string                `json:"description"`
	License      string                `json:"license"`
	Media        []CharmResultMedia    `json:"media"`
	Publisher    *CharmResultPublisher `json:"publisher"`
	StoreURL     string                `json:"store-url"`
	StoreURLOld  string                `json:"store-url-old"`
	Summary      string                `json:"summary"`
	Title        string                `json:"title"`
	Unlisted     bool                  `json:"unlisted"`
	Website      string                `json:"website"`
}

type CharmRepo interface {
	List(ctx context.Context) ([]Charm, error)
	Get(ctx context.Context, name string) (*Charm, error)
	ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error)
}

type CharmUseCase struct {
	charm CharmRepo
}

func NewCharmUseCase(charm CharmRepo) *CharmUseCase {
	return &CharmUseCase{
		charm: charm,
	}
}

func (uc *CharmUseCase) ListCharms(ctx context.Context) ([]Charm, error) {
	charms, err := uc.charm.List(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(charms, func(charm Charm) bool {
		return slices.Contains(charm.Result.DeployableOn, "kubernetes") || charm.Type != "charm"
	}), nil
}

func (uc *CharmUseCase) GetCharm(ctx context.Context, name string) (*Charm, error) {
	return uc.charm.Get(ctx, name)
}

func (uc *CharmUseCase) ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error) {
	return uc.charm.ListArtifacts(ctx, name)
}
