package containerimage

import (
	"context"
	"fmt"
	"time"

	chart "github.com/otterscale/otterscale/internal/core/application/chart"
)

type ContainerImage struct {
	Namespace string
	Name      string
	Ref       string
	Size      int64
	Tag       string 
	CreateAt  time.Time
}

type ContainerImageRepo interface {
	ListRepos(ctx context.Context, scope, namespace, endpoint string) ([]string, error)
	ListTags(ctx context.Context, scope, namespace, endpoint, name string) ([]string, error)
	GetSize(ctx context.Context, scope, namespace, endpoint, name, tag string) (int64, error)
	GetCreateAt(ctx context.Context, scope, namespace, endpoint, name, tag string) (time.Time, error)
	UploadImage(ctx context.Context, scope, namespace, endpoint, name, tag string, imageTar []byte) error
	DeleteImage(ctx context.Context, scope, namespace, endpoint, name, tag string) error
}

type UseCase struct {
	containerimage    ContainerImageRepo
	chart    *chart.UseCase
}

func NewUseCase(containerimage ContainerImageRepo, chartUC *chart.UseCase) *UseCase {
	return &UseCase{
		containerimage:    containerimage,
		chart:    chartUC,
	}
}

func (uc *UseCase) ListContainerImages(ctx context.Context, scope, namespace, endpoint string) ([]ContainerImage, error) {
	var (
		images []ContainerImage
	)

	repos, err := uc.containerimage.ListRepos(ctx, scope ,namespace, endpoint)
	if err != nil {
		return images, err
	}

	for _, name := range repos { 
        tags, err := uc.containerimage.ListTags(ctx, scope, namespace, endpoint, name) 
        if err != nil {
            return images, err
        }

        for _, tag := range tags { 
            size, err := uc.containerimage.GetSize(ctx, scope, namespace, endpoint, name, tag)
            if err != nil {
                return images, err
            }

            createdAt, err := uc.containerimage.GetCreateAt(ctx, scope, namespace, endpoint, name, tag)
            if err != nil {
                return images, err
            }

			repoPath := name
            if namespace != "" {
                repoPath = fmt.Sprintf("%s/%s", namespace, name)
            }

            ref := fmt.Sprintf("%s/%s:%s", endpoint, repoPath, tag)

            images = append(images, ContainerImage{
                Namespace: namespace,
                Name:      name,
                Ref:       ref,
                Size:      size,
                Tag:       tag,
                CreateAt:  createdAt,
            })
        }
    }

	return images, nil
}

func (uc *UseCase) UploadContainerImage(ctx context.Context, scope, namespace, endpoint, name, tag string, imageTar []byte) error {
	if len(imageTar) == 0 {
        return fmt.Errorf("image content is empty")
    }

    if err := uc.containerimage.UploadImage(ctx, scope, namespace, endpoint, name, tag, imageTar); err != nil {
        return err
    }

    return nil
}

func (uc *UseCase) DeleteContainerImage(ctx context.Context, scope, namespace, endpoint, name, tag string) error {
	if err := uc.containerimage.DeleteImage(ctx, scope, namespace, endpoint, name, tag); err != nil {
        return err
    }

    return nil
}

