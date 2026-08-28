package handler

// cleanObject strips metadata.managedFields (server-side apply bookkeeping)
// and the kubectl.kubernetes.io/last-applied-configuration annotation.
//
// Sanitizing is a presentation concern, so it happens here rather than in the
// domain layer. Working on map[string]any keeps the handler free of
// k8s.io/apimachinery imports.
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
