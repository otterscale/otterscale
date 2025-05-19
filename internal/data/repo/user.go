package repo

// import (
// 	"context"

// 	"github.com/openhdc/otterscale/internal/data/repo/ent"
// 	"github.com/openhdc/otterscale/internal/domain/model"
// 	"github.com/openhdc/otterscale/internal/domain/service"
// )

// type userRepo struct {
// 	repo *Repo
// }

// func NewUser(repo *Repo) service.UserRepo {
// 	return &userRepo{
// 		repo: repo,
// 	}
// }

// var _ service.UserRepo = (*userRepo)(nil)

// func (r *userRepo) Get(ctx context.Context, id string) (*model.User, error) {
// 	e, err := r.repo.User.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return toUser(e), nil
// }

// func toUser(e *ent.User) *model.User {
// 	return &model.User{}
// }
