package charm

import (
	"context"
	"slices"
	"time"
)

type Charm struct {
	ID              string     `json:"id"`
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	Result          Result     `json:"result"`
	DefaultArtifact Artifact   `json:"default-release"`
	Artifacts       []Artifact `json:"channel-map"`
}

type Base struct {
	Architecture string `json:"architecture"`
	Channel      string `json:"channel"`
	Name         string `json:"name"`
}

type Channel struct {
	Base       Base      `json:"base"`
	Name       string    `json:"name"`
	ReleasedAt time.Time `json:"released-at"`
	Risk       string    `json:"risk"`
	Track      string    `json:"track"`
}

type Revision struct {
	Bases     []Base    `json:"bases"`
	CreatedAt time.Time `json:"created-at"`
	Revision  int       `json:"revision"`
	Version   string    `json:"version"`
}

type Artifact struct {
	Channel  Channel  `json:"channel"`
	Revision Revision `json:"revision"`
}

type ResultCategory struct {
	Featured bool   `json:"featured"`
	Name     string `json:"name"`
}

type ResultMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ResultPublisher struct {
	DisplayName string `json:"display-name"`
	ID          string `json:"id"`
	Username    string `json:"username"`
	Validation  string `json:"validation"`
}

type Result struct {
	BugsURL      string           `json:"bugs-url"`
	Categories   []ResultCategory `json:"categories"`
	DeployableOn []string         `json:"deployable-on"`
	Description  string           `json:"description"`
	License      string           `json:"license"`
	Media        []ResultMedia    `json:"media"`
	Publisher    *ResultPublisher `json:"publisher"`
	StoreURL     string           `json:"store-url"`
	StoreURLOld  string           `json:"store-url-old"`
	Summary      string           `json:"summary"`
	Title        string           `json:"title"`
	Unlisted     bool             `json:"unlisted"`
	Website      string           `json:"website"`
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type CharmRepo interface {
	List(ctx context.Context) ([]Charm, error)
	Get(ctx context.Context, name string) (*Charm, error)
	ListArtifacts(ctx context.Context, name string) ([]Artifact, error)
}

type UseCase struct {
	charm CharmRepo
}

func NewUseCase(charm CharmRepo) *UseCase {
	return &UseCase{
		charm: charm,
	}
}

func (uc *UseCase) ListCharms(ctx context.Context) ([]Charm, error) {
	charms, err := uc.charm.List(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(charms, func(charm Charm) bool {
		return slices.Contains(charm.Result.DeployableOn, "kubernetes") || charm.Type != "charm"
	}), nil
}

func (uc *UseCase) GetCharm(ctx context.Context, name string) (*Charm, error) {
	return uc.charm.Get(ctx, name)
}

func (uc *UseCase) ListArtifacts(ctx context.Context, name string) ([]Artifact, error) {
	return uc.charm.ListArtifacts(ctx, name)
}
