package infra

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/service/infra/kube"
)

var ProviderSet = wire.NewSet(
	kube.NewClientset,
	kube.NewNamespace,
	kube.NewCronJob,
	kube.NewJob,
)
