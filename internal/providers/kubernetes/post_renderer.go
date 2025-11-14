package kubernetes

import (
	"bytes"
	"maps"

	"sigs.k8s.io/kustomize/kyaml/kio"
)

type postRenderer struct {
	ExtraLabels      map[string]string
	ExtraAnnotations map[string]string
}

func newPostRenderer(extraLabels, extraAnnotations map[string]string) *postRenderer {
	return &postRenderer{
		ExtraLabels:      extraLabels,
		ExtraAnnotations: extraAnnotations,
	}
}

func (p *postRenderer) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	if len(p.ExtraLabels) == 0 && len(p.ExtraAnnotations) == 0 {
		return renderedManifests, nil
	}

	nodes, err := kio.FromBytes(renderedManifests.Bytes())
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		// labels
		labels := node.GetLabels()
		maps.Copy(labels, p.ExtraLabels)
		if err := node.SetLabels(labels); err != nil {
			return nil, err
		}

		// annotations
		annotations := node.GetAnnotations()
		maps.Copy(annotations, p.ExtraAnnotations)
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
