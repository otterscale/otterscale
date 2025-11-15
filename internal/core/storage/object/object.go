package object

type UseCase struct {
	bucket BucketRepo
	user   UserRepo
}

func NewUseCase(bucket BucketRepo, user UserRepo) *UseCase {
	return &UseCase{
		bucket: bucket,
		user:   user,
	}
}
