package cluster

type ClusterUseCase struct {
	namespace NamespaceRepo
	node      NodeRepo
}

func NewClusterUseCase(namespace NamespaceRepo, node NodeRepo) *ClusterUseCase {
	return &ClusterUseCase{
		namespace: namespace,
		node:      node,
	}
}
