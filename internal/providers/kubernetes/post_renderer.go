package kubernetes

import (
	"bytes"
	"maps"

	"sigs.k8s.io/kustomize/kyaml/kio"
)

type postRenderer struct {
	extraLabels      map[string]string
	extraAnnotations map[string]string
}

func newPostRenderer(extraLabels, extraAnnotations map[string]string) *postRenderer {
	return &postRenderer{
		extraLabels:      extraLabels,
		extraAnnotations: extraAnnotations,
	}
}

func (p *postRenderer) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	if len(p.extraLabels) == 0 && len(p.extraAnnotations) == 0 {
		return renderedManifests, nil
	}

	nodes, err := kio.FromBytes(renderedManifests.Bytes())
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		// labels
		labels := node.GetLabels()
		maps.Copy(labels, p.extraLabels)
		if err := node.SetLabels(labels); err != nil {
			return nil, err
		}

		// annotations
		annotations := node.GetAnnotations()
		maps.Copy(annotations, p.extraAnnotations)
		if err := node.SetAnnotations(annotations); err != nil {
			return nil, err
		}
	}

	str, err := kio.StringAll(nodes)
	if err != nil {
		return nil, err
	}

	return bytes.NewBufferString(str), nil
}
