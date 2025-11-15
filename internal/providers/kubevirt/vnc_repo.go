package kubevirt

import (
	"net/url"
	"strconv"

	kvcorev1 "github.com/otterscale/kubevirt-client-go/kubevirt/typed/core/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vnc"
)

type vncRepo struct {
	kubevirt *KubeVirt
}

func NewVNCRepo(kubevirt *KubeVirt) vnc.VNCRepo {
	return &vncRepo{
		kubevirt: kubevirt,
	}
}

var _ vnc.VNCRepo = (*vncRepo)(nil)

func (r *vncRepo) Stream(scope, namespace, name string) (kvcorev1.StreamInterface, error) {
	config, err := r.kubevirt.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	queryParams.Add("moveCursor", strconv.FormatBool(true))
	queryParams.Add("preserveSession", strconv.FormatBool(true))

	return kvcorev1.AsyncSubresourceHelper(config, "virtualmachineinstances", namespace, name, "vnc", queryParams)
}
