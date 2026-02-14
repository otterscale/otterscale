package handler

// cleanObject strips noisy metadata from a raw Kubernetes object map:
//   - metadata.managedFields (server-side apply bookkeeping)
//   - the kubectl.kubernetes.io/last-applied-configuration annotation
//
// This is a presentation concern: the domain layer returns raw
// Kubernetes objects and the handler sanitises them before serialising
// to protobuf. Operating on map[string]any keeps the handler layer
// free of k8s.io/apimachinery imports.
func cleanObject(obj map[string]any) {
	metadata, ok := obj["metadata"].(map[string]any)
	if !ok {
		return
	}
	delete(metadata, "managedFields")

	annotations, ok := metadata["annotations"].(map[string]any)
	if !ok || len(annotations) == 0 {
		return
	}
	delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
	if len(annotations) == 0 {
		delete(metadata, "annotations")
	}
}
