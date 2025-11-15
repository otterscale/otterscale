package config

type UseCase struct {
	configMap ConfigMapRepo
	secret    SecretRepo
}

func NewUseCase(configMap ConfigMapRepo, secret SecretRepo) *UseCase {
	return &UseCase{
		configMap: configMap,
		secret:    secret,
	}
}
