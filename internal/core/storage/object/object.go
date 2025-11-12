package object

type ObjectUseCase struct {
	bucket BucketRepo
	user   UserRepo
}

func NewObjectUseCase(bucket BucketRepo, user UserRepo) *ObjectUseCase {
	return &ObjectUseCase{
		bucket: bucket,
		user:   user,
	}
}
