package cluster

type UseCase struct {
	namespace NamespaceRepo
	node      NodeRepo
}

func NewUseCase(namespace NamespaceRepo, node NodeRepo) *UseCase {
	return &UseCase{
		namespace: namespace,
		node:      node,
	}
}
