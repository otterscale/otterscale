package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
)

func (a *KubeApp) ListReleases(ctx context.Context, req *connect.Request[v1.ListReleasesRequest]) (*connect.Response[v1.ListReleasesResponse], error) {
	rels, err := a.svc.ListReleases(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), "")
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		fmt.Println(rel.Name, rel.Namespace)
	}
	res := &v1.ListReleasesResponse{}
	return connect.NewResponse(res), nil
}
