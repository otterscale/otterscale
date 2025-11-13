package config

type ConfigUseCase struct {
	configMap ConfigMapRepo
	secret    SecretRepo
}

func NewConfigUseCase(configMap ConfigMapRepo, secret SecretRepo) *ConfigUseCase {
	return &ConfigUseCase{
		configMap: configMap,
		secret:    secret,
	}
}
