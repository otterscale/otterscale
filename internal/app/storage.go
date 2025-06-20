package app

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/storage/v1"
	"github.com/openhdc/otterscale/api/storage/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type StorageService struct {
	pbconnect.UnimplementedStorageServiceHandler

	uc *core.StorageUseCase
}

func NewStorageService(uc *core.StorageUseCase) *StorageService {
	return &StorageService{uc: uc}
}

var _ pbconnect.StorageServiceHandler = (*StorageService)(nil)

func (s *StorageService) ListMONs(ctx context.Context, req *connect.Request[pb.ListMONsRequest]) (*connect.Response[pb.ListMONsResponse], error) {
	mons, err := s.uc.ListMONs(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMONsResponse{}
	resp.SetMons(toProtoMONs(mons))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListOSDs(ctx context.Context, req *connect.Request[pb.ListOSDsRequest]) (*connect.Response[pb.ListOSDsResponse], error) {
	osds, err := s.uc.ListOSDs(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListOSDsResponse{}
	resp.SetOsds(toProtoOSDs(osds))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) DoSMART(ctx context.Context, req *connect.Request[pb.DoSMARTRequest]) (*connect.Response[pb.DoSMARTResponse], error) {
	outputs, err := s.uc.DoSMART(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetOsdName())
	if err != nil {
		return nil, err
	}
	resp := &pb.DoSMARTResponse{}
	resp.SetDeviceOutputMap(toDeviceOutputMap(outputs))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListPools(ctx context.Context, req *connect.Request[pb.ListPoolsRequest]) (*connect.Response[pb.ListPoolsResponse], error) {
	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListPoolsResponse{}
	resp.SetPools(toProtoPools(pools))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) CreatePool(ctx context.Context, req *connect.Request[pb.CreatePoolRequest]) (*connect.Response[pb.Pool], error) {
	if req.Msg.GetPoolType() == pb.CreatePoolRequest_UNSPECIFIED {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid pool type"))
	}
	apps := []string{}
	for _, app := range req.Msg.GetApplications() {
		apps = append(apps, strings.ToLower(app.String()))
	}
	pool, err := s.uc.CreatePool(ctx,
		req.Msg.GetScopeUuid(),
		req.Msg.GetFacilityName(),
		req.Msg.GetPoolName(),
		strings.ToLower(req.Msg.GetPoolType().String()),
		req.Msg.GetEcOverwrites(),
		int(req.Msg.GetQuotaMaxBytes()),
		int(req.Msg.GetQuotaMaxObjects()),
		apps,
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoPool(pool)
	return connect.NewResponse(resp), nil
}

func (s *StorageService) UpdatePool(ctx context.Context, req *connect.Request[pb.UpdatePoolRequest]) (*connect.Response[pb.Pool], error) {
	pool, err := s.uc.UpdatePool(ctx,
		req.Msg.GetScopeUuid(),
		req.Msg.GetFacilityName(),
		req.Msg.GetPoolName(),
		int(req.Msg.GetQuotaMaxBytes()),
		int(req.Msg.GetQuotaMaxObjects()),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoPool(pool)
	return connect.NewResponse(resp), nil
}

func (s *StorageService) DeletePool(ctx context.Context, req *connect.Request[pb.DeletePoolRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeletePool(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetPoolName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListImages(ctx context.Context, req *connect.Request[pb.ListImagesRequest]) (*connect.Response[pb.ListImagesResponse], error) {
	imgs, err := s.uc.ListImages(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListImagesResponse{}
	resp.SetImages(toProtoImages(imgs))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListVolumes(ctx context.Context, req *connect.Request[pb.ListVolumesRequest]) (*connect.Response[pb.ListVolumesResponse], error) {
	volumes, err := s.uc.ListVolumes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVolumesResponse{}
	resp.SetVolumes(toProtoVolumes(volumes))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListSubvolumes(ctx context.Context, req *connect.Request[pb.ListSubvolumesRequest]) (*connect.Response[pb.ListSubvolumesResponse], error) {
	subvolumes, err := s.uc.ListSubvolumes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetVolumeName(), req.Msg.GetGroupName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSubvolumesResponse{}
	resp.SetSubvolumes(toProtoSubvolumes(subvolumes))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListSubvolumeGroups(ctx context.Context, req *connect.Request[pb.ListSubvolumeGroupsRequest]) (*connect.Response[pb.ListSubvolumeGroupsResponse], error) {
	groups, err := s.uc.ListSubvolumeGroups(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetVolumeName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSubvolumeGroupsResponse{}
	resp.SetSubvolumeGroups(toProtoSubvolumeGroups(groups))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) ListBuckets(ctx context.Context, req *connect.Request[pb.ListBucketsRequest]) (*connect.Response[pb.ListBucketsResponse], error) {
	buckets, err := s.uc.ListBuckets(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListBucketsResponse{}
	resp.SetBuckets(toProtoBuckets(buckets))
	return connect.NewResponse(resp), nil
}

// func (s *StorageService) ListRoles(ctx context.Context, req *connect.Request[pb.ListRolesRequest]) (*connect.Response[pb.ListRolesResponse], error) {
// 	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListPoolsResponse{}
// 	resp.SetPools(toProtoPools(pools))
// 	return connect.NewResponse(resp), nil
// }

func (s *StorageService) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	users, err := s.uc.ListUsers(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListUsersResponse{}
	resp.SetUsers(toProtoUsers(users))
	return connect.NewResponse(resp), nil
}

// func (s *StorageService) ListAccessKeys(ctx context.Context, req *connect.Request[pb.ListAccessKeysRequest]) (*connect.Response[pb.ListAccessKeysResponse], error) {
// 	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListPoolsResponse{}
// 	resp.SetPools(toProtoPools(pools))
// 	return connect.NewResponse(resp), nil
// }

// func (s *StorageService) ListSnapshots(ctx context.Context, req *connect.Request[pb.ListSnapshotsRequest]) (*connect.Response[pb.ListSnapshotsResponse], error) {
// 	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListPoolsResponse{}
// 	resp.SetPools(toProtoPools(pools))
// 	return connect.NewResponse(resp), nil
// }

// func (s *StorageService) ListSnapshotSchedules(ctx context.Context, req *connect.Request[pb.ListSnapshotSchedulesRequest]) (*connect.Response[pb.ListSnapshotSchedulesResponse], error) {
// 	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListSnapshotSchedulesResponse{}
// 	resp.SetSnapshotSchedules(toProtoSnapshotSchedules(pools))
// 	return connect.NewResponse(resp), nil
// }

func toProtoMONs(ms []core.MON) []*pb.MON {
	ret := []*pb.MON{}
	for i := range ms {
		ret = append(ret, toProtoMON(&ms[i]))
	}
	return ret
}

func toProtoMON(m *core.MON) *pb.MON {
	ret := &pb.MON{}
	ret.SetName(m.Name)
	return ret
}

func toProtoOSDs(os []core.OSD) []*pb.OSD {
	ret := []*pb.OSD{}
	for i := range os {
		ret = append(ret, toProtoOSD(&os[i]))
	}
	return ret
}

func toProtoOSD(o *core.OSD) *pb.OSD {
	ret := &pb.OSD{}
	ret.SetName(o.Name)
	return ret
}

func toDeviceOutputMap(m map[string][]string) map[string]*pb.DoSMARTResponse_Output {
	ret := map[string]*pb.DoSMARTResponse_Output{}
	for device, lines := range m {
		output := &pb.DoSMARTResponse_Output{}
		output.SetLines(lines)
		ret[device] = output
	}
	return ret
}

func toProtoPools(ps []core.Pool) []*pb.Pool {
	ret := []*pb.Pool{}
	for i := range ps {
		ret = append(ret, toProtoPool(&ps[i]))
	}
	return ret
}

func toProtoPool(p *core.Pool) *pb.Pool {
	ret := &pb.Pool{}
	ret.SetName(p.Name)
	return ret
}

func toProtoImages(is []core.RBDImage) []*pb.Image {
	ret := []*pb.Image{}
	for i := range is {
		ret = append(ret, toProtoImage(&is[i]))
	}
	return ret
}

func toProtoImage(i *core.RBDImage) *pb.Image {
	ret := &pb.Image{}
	ret.SetName(i.Name)
	return ret
}

func toProtoVolumes(vs []core.Volume) []*pb.Volume {
	ret := []*pb.Volume{}
	for i := range vs {
		ret = append(ret, toProtoVolume(&vs[i]))
	}
	return ret
}

func toProtoVolume(v *core.Volume) *pb.Volume {
	ret := &pb.Volume{}
	ret.SetName(v.Name)
	return ret
}

func toProtoSubvolumes(ss []core.Subvolume) []*pb.Subvolume {
	ret := []*pb.Subvolume{}
	for i := range ss {
		ret = append(ret, toProtoSubvolume(&ss[i]))
	}
	return ret
}

func toProtoSubvolume(s *core.Subvolume) *pb.Subvolume {
	ret := &pb.Subvolume{}
	ret.SetName(s.Name)
	return ret
}

func toProtoSubvolumeGroups(ss []core.SubvolumeGroup) []*pb.SubvolumeGroup {
	ret := []*pb.SubvolumeGroup{}
	for i := range ss {
		ret = append(ret, toProtoSubvolumeGroup(&ss[i]))
	}
	return ret
}

func toProtoSubvolumeGroup(s *core.SubvolumeGroup) *pb.SubvolumeGroup {
	ret := &pb.SubvolumeGroup{}
	ret.SetName(s.Name)
	return ret
}

func toProtoBuckets(bs []core.RGWBucket) []*pb.Bucket {
	ret := []*pb.Bucket{}
	for i := range bs {
		ret = append(ret, toProtoBucket(&bs[i]))
	}
	return ret
}

func toProtoBucket(b *core.RGWBucket) *pb.Bucket {
	ret := &pb.Bucket{}
	ret.SetName(b.Bucket)
	return ret
}

func toProtoRoles(rs []core.RGWRole) []*pb.Role {
	ret := []*pb.Role{}
	for i := range rs {
		ret = append(ret, toProtoRole(&rs[i]))
	}
	return ret
}

func toProtoRole(r *core.RGWRole) *pb.Role {
	ret := &pb.Role{}
	ret.SetName(r.Name)
	return ret
}

func toProtoUsers(us []core.RGWUser) []*pb.User {
	ret := []*pb.User{}
	for i := range us {
		ret = append(ret, toProtoUser(&us[i]))
	}
	return ret
}

func toProtoUser(u *core.RGWUser) *pb.User {
	ret := &pb.User{}
	ret.SetName(u.Name)
	return ret
}

// func toProtoAccessKeys(as []core.RGWAccessKey) []*pb.AccessKey {
// 	ret := []*pb.AccessKey{}
// 	for i := range as {
// 		ret = append(ret, toProtoAccessKey(&as[i]))
// 	}
// 	return ret
// }

// func toProtoAccessKey(a *core.RGWAccessKey) *pb.AccessKey {
// 	ret := &pb.AccessKey{}
// 	ret.SetName(a.Name)
// 	return ret
// }

// func toProtoSnapshots(ss []core.Snapshot) []*pb.Snapshot {
// 	ret := []*pb.Snapshot{}
// 	for i := range ss {
// 		ret = append(ret, toProtoSnapshot(&ss[i]))
// 	}
// 	return ret
// }

// func toProtoSnapshot(s *core.Snapshot) *pb.Snapshot {
// 	ret := &pb.Snapshot{}
// 	ret.SetName(s.Name)
// 	return ret
// }

// func toProtoSnapshotSchedules(ss []core.SnapshotSchedule) []*pb.SnapshotSchedule {
// 	ret := []*pb.SnapshotSchedule{}
// 	for i := range ss {
// 		ret = append(ret, toProtoSnapshotSchedule(&ss[i]))
// 	}
// 	return ret
// }

// func toProtoSnapshotSchedule(s *core.SnapshotSchedule) *pb.SnapshotSchedule {
// 	ret := &pb.SnapshotSchedule{}
// 	ret.SetName(s.Name)
// 	return ret
// }
