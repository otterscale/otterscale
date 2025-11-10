package app

import (
	"context"
	"errors"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/storage/v1"
	"github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type StorageService struct {
	pbconnect.UnimplementedStorageServiceHandler

	storage *core.StorageUseCase
}

func NewStorageService(storage *core.StorageUseCase) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

var _ pbconnect.StorageServiceHandler = (*StorageService)(nil)

func (s *StorageService) ListMONs(ctx context.Context, req *pb.ListMONsRequest) (*pb.ListMONsResponse, error) {
	mons, err := s.storage.ListMONs(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMONsResponse{}
	resp.SetMons(toProtoMONs(mons))
	return resp, nil
}

func (s *StorageService) ListOSDs(ctx context.Context, req *pb.ListOSDsRequest) (*pb.ListOSDsResponse, error) {
	osds, err := s.storage.ListOSDs(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListOSDsResponse{}
	resp.SetOsds(toProtoOSDs(osds))
	return resp, nil
}

func (s *StorageService) DoSMART(ctx context.Context, req *pb.DoSMARTRequest) (*pb.DoSMARTResponse, error) {
	outputs, err := s.storage.DoSMART(ctx, req.GetScope(), req.GetFacility(), req.GetOsdName())
	if err != nil {
		return nil, err
	}
	resp := &pb.DoSMARTResponse{}
	resp.SetDeviceOutputMap(toDeviceOutputMap(outputs))
	return resp, nil
}

func (s *StorageService) ListPools(ctx context.Context, req *pb.ListPoolsRequest) (*pb.ListPoolsResponse, error) {
	pools, err := s.storage.ListPools(ctx, req.GetScope(), req.GetFacility(), req.GetApplication())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListPoolsResponse{}
	resp.SetPools(toProtoPools(pools))
	return resp, nil
}

func (s *StorageService) CreatePool(ctx context.Context, req *pb.CreatePoolRequest) (*pb.Pool, error) {
	poolType := []string{"cephfs", "rbd", "rgw"}
	apps := []string{}
	for _, app := range req.GetApplications() {
		if !slices.Contains(poolType, app) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid pool type"))
		}
		apps = append(apps, app)
	}
	pool, err := s.storage.CreatePool(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetPoolName(),
		strings.ToLower(req.GetPoolType().String()),
		req.GetEcOverwrites(),
		req.GetReplicatedSize(),
		req.GetQuotaBytes(),
		req.GetQuotaObjects(),
		apps,
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoPool(pool)
	return resp, nil
}

func (s *StorageService) UpdatePool(ctx context.Context, req *pb.UpdatePoolRequest) (*pb.Pool, error) {
	pool, err := s.storage.UpdatePool(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetPoolName(),
		req.GetQuotaBytes(),
		req.GetQuotaObjects(),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoPool(pool)
	return resp, nil
}

func (s *StorageService) DeletePool(ctx context.Context, req *pb.DeletePoolRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeletePool(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	imgs, err := s.storage.ListImages(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListImagesResponse{}
	resp.SetImages(toProtoImages(imgs))
	return resp, nil
}

func (s *StorageService) CreateImage(ctx context.Context, req *pb.CreateImageRequest) (*pb.Image, error) {
	img, err := s.storage.CreateImage(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetPoolName(),
		req.GetImageName(),
		req.GetObjectSizeBytes(),
		req.GetStripeUnitBytes(),
		req.GetStripeCount(),
		req.GetQuotaBytes(),
		req.GetLayering(),
		req.GetExclusiveLock(),
		req.GetObjectMap(),
		req.GetFastDiff(),
		req.GetDeepFlatten(),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoImage(img)
	return resp, nil
}

func (s *StorageService) UpdateImage(ctx context.Context, req *pb.UpdateImageRequest) (*pb.Image, error) {
	img, err := s.storage.UpdateImage(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetPoolName(),
		req.GetImageName(),
		req.GetQuotaBytes(),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoImage(img)
	return resp, nil
}

func (s *StorageService) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteImage(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateImageSnapshot(ctx context.Context, req *pb.CreateImageSnapshotRequest) (*pb.Image_Snapshot, error) {
	snap, err := s.storage.CreateImageSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}
	resp := toProtoImageSnapshot(snap)
	return resp, nil
}

func (s *StorageService) DeleteImageSnapshot(ctx context.Context, req *pb.DeleteImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteImageSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) RollbackImageSnapshot(ctx context.Context, req *pb.RollbackImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.storage.RollbackImageSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ProtectImageSnapshot(ctx context.Context, req *pb.ProtectImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.storage.ProtectImageSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) UnprotectImageSnapshot(ctx context.Context, req *pb.UnprotectImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.storage.UnprotectImageSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListVolumes(ctx context.Context, req *pb.ListVolumesRequest) (*pb.ListVolumesResponse, error) {
	volumes, err := s.storage.ListVolumes(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVolumesResponse{}
	resp.SetVolumes(toProtoVolumes(volumes))
	return resp, nil
}

func (s *StorageService) ListSubvolumes(ctx context.Context, req *pb.ListSubvolumesRequest) (*pb.ListSubvolumesResponse, error) {
	subvolumes, err := s.storage.ListSubvolumes(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetGroupName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSubvolumesResponse{}
	resp.SetSubvolumes(toProtoSubvolumes(subvolumes))
	return resp, nil
}

func (s *StorageService) CreateSubvolume(ctx context.Context, req *pb.CreateSubvolumeRequest) (*pb.Subvolume, error) {
	subvolume, err := s.storage.CreateSubvolume(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetVolumeName(),
		req.GetSubvolumeName(),
		req.GetGroupName(),
		req.GetQuotaBytes(),
		req.GetExport(),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) UpdateSubvolume(ctx context.Context, req *pb.UpdateSubvolumeRequest) (*pb.Subvolume, error) {
	subvolume, err := s.storage.UpdateSubvolume(ctx,
		req.GetScope(),
		req.GetFacility(),
		req.GetVolumeName(),
		req.GetSubvolumeName(),
		req.GetGroupName(),
		req.GetQuotaBytes(),
	)
	if err != nil {
		return nil, err
	}
	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) DeleteSubvolume(ctx context.Context, req *pb.DeleteSubvolumeRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteSubvolume(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) GrantSubvolumeExportAccess(ctx context.Context, req *pb.GrantSubvolumeExportAccessRequest) (*emptypb.Empty, error) {
	if err := s.storage.GrantSubvolumeClient(ctx, req.GetScope(), req.GetFacility(), req.GetSubvolumeName(), req.GetClientIp()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) RevokeSubvolumeExportAccess(ctx context.Context, req *pb.RevokeSubvolumeExportAccessRequest) (*emptypb.Empty, error) {
	if err := s.storage.RevokeSubvolumeClient(ctx, req.GetScope(), req.GetFacility(), req.GetSubvolumeName(), req.GetClientIp()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateSubvolumeSnapshot(ctx context.Context, req *pb.CreateSubvolumeSnapshotRequest) (*pb.Subvolume_Snapshot, error) {
	snapshot, err := s.storage.CreateSubvolumeSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}
	resp := toProtoSubvolumeSnapshot(snapshot)
	return resp, nil
}

func (s *StorageService) DeleteSubvolumeSnapshot(ctx context.Context, req *pb.DeleteSubvolumeSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteSubvolumeSnapshot(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListSubvolumeGroups(ctx context.Context, req *pb.ListSubvolumeGroupsRequest) (*pb.ListSubvolumeGroupsResponse, error) {
	groups, err := s.storage.ListSubvolumeGroups(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListSubvolumeGroupsResponse{}
	resp.SetSubvolumeGroups(toProtoSubvolumeGroups(groups))
	return resp, nil
}

func (s *StorageService) CreateSubvolumeGroup(ctx context.Context, req *pb.CreateSubvolumeGroupRequest) (*pb.SubvolumeGroup, error) {
	group, err := s.storage.CreateSubvolumeGroup(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetGroupName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}
	resp := toProtoSubvolumeGroup(group)
	return resp, nil
}

func (s *StorageService) UpdateSubvolumeGroup(ctx context.Context, req *pb.UpdateSubvolumeGroupRequest) (*pb.SubvolumeGroup, error) {
	group, err := s.storage.UpdateSubvolumeGroup(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetGroupName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}
	resp := toProtoSubvolumeGroup(group)
	return resp, nil
}

func (s *StorageService) DeleteSubvolumeGroup(ctx context.Context, req *pb.DeleteSubvolumeGroupRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteSubvolumeGroup(ctx, req.GetScope(), req.GetFacility(), req.GetVolumeName(), req.GetGroupName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListBuckets(ctx context.Context, req *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error) {
	buckets, err := s.storage.ListBuckets(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListBucketsResponse{}
	resp.SetBuckets(toProtoBuckets(buckets))
	return resp, nil
}

func (s *StorageService) CreateBucket(ctx context.Context, req *pb.CreateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.storage.CreateBucket(ctx, req.GetScope(), req.GetFacility(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), s.ACL(req.GetAcl().String()))
	if err != nil {
		return nil, err
	}
	resp := toProtoBucket(bucket)
	return resp, nil
}

func (s *StorageService) UpdateBucket(ctx context.Context, req *pb.UpdateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.storage.UpdateBucket(ctx, req.GetScope(), req.GetFacility(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), s.ACL(req.GetAcl().String()))
	if err != nil {
		return nil, err
	}
	resp := toProtoBucket(bucket)
	return resp, nil
}

func (s *StorageService) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteBucket(ctx, req.GetScope(), req.GetFacility(), req.GetBucketName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.storage.ListUsers(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListUsersResponse{}
	resp.SetUsers(toProtoUsers(users))
	return resp, nil
}

func (s *StorageService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.storage.CreateUser(ctx, req.GetScope(), req.GetFacility(), req.GetUserId(), req.GetUserName(), req.GetSuspended())
	if err != nil {
		return nil, err
	}
	resp := toProtoUser(user)
	return resp, nil
}

func (s *StorageService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := s.storage.UpdateUser(ctx, req.GetScope(), req.GetFacility(), req.GetUserId(), req.GetUserName(), req.GetSuspended())
	if err != nil {
		return nil, err
	}
	resp := toProtoUser(user)
	return resp, nil
}

func (s *StorageService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteUser(ctx, req.GetScope(), req.GetFacility(), req.GetUserId()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateUserKey(ctx context.Context, req *pb.CreateUserKeyRequest) (*pb.User_Key, error) {
	key, err := s.storage.CreateUserKey(ctx, req.GetScope(), req.GetFacility(), req.GetUserId())
	if err != nil {
		return nil, err
	}
	resp := toProtoUserKey(key)
	return resp, nil
}

func (s *StorageService) DeleteUserKey(ctx context.Context, req *pb.DeleteUserKeyRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteUserKey(ctx, req.GetScope(), req.GetFacility(), req.GetUserId(), req.GetAccessKey()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateSMBShare(ctx context.Context, req *pb.CreateSMBShareRequest) (*pb.SMBShare, error) {
	adAuth, localAuth := toCoreSMBShareAuth(req.GetActiveDirectory(), req.GetLocalUser())

	securityMode := ""
	switch req.GetSecurityMode() {
	case pb.SMBShare_USER:
		securityMode = "user"
	case pb.SMBShare_ACTIVE_DIRECTORY:
		securityMode = "active-directory"
	}

	namespace := normalizeNamespace(req.GetNamespace())

	share, err := s.storage.CreateSMBShare(ctx, req.GetScope(), req.GetFacility(), namespace, req.GetName(), toCoreSMBShareConfig(req.GetBrowsable(), req.GetReadOnly(), req.GetGuestOk(), req.GetMapToGuest(), req.GetValidUsers()), req.GetSizeBytes(), securityMode, adAuth, localAuth)
	if err != nil {
		return nil, err
	}

	return toProtoSMBShare(share), nil
}

func (s *StorageService) ListSMBShares(ctx context.Context, req *pb.ListSMBSharesRequest) (*pb.ListSMBSharesResponse, error) {
	shares, err := s.storage.ListSMBShares(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListSMBSharesResponse{}
	resp.SetSmbShares(toProtoSMBShares(shares))
	return resp, nil
}

func (s *StorageService) UpdateSMBShare(ctx context.Context, req *pb.UpdateSMBShareRequest) (*pb.SMBShare, error) {
	adAuth, localAuth := toCoreSMBShareAuth(req.GetActiveDirectory(), req.GetLocalUser())

	namespace := normalizeNamespace(req.GetNamespace())

	share, err := s.storage.UpdateSMBShare(ctx, req.GetScope(), req.GetFacility(), namespace, req.GetName(), toCoreSMBShareConfig(req.GetBrowsable(), req.GetReadOnly(), req.GetGuestOk(), req.GetMapToGuest(), req.GetValidUsers()), req.GetSizeBytes(), localAuth, adAuth)
	if err != nil {
		return nil, err
	}

	return toProtoSMBShare(share), nil
}

func (s *StorageService) ACL(str string) core.RGWBucketCannedACL {
	acl := core.RGWBucketCannedACL(strings.ToLower(strings.Join(strings.Split(str, "_"), "-")))
	if slices.Contains(acl.Values(), acl) {
		return acl
	}
	return core.RGWBucketCannedACLPrivate
}

func toProtoStorageMachine(id, hostname string) *pb.Machine {
	if id == "" || hostname == "" {
		return nil
	}
	ret := &pb.Machine{}
	ret.SetId(id)
	ret.SetHostname(hostname)
	return ret
}

func toProtoMONs(ms []core.MON) []*pb.MON {
	ret := []*pb.MON{}
	for i := range ms {
		ret = append(ret, toProtoMON(&ms[i]))
	}
	return ret
}

func toProtoMON(m *core.MON) *pb.MON {
	ret := &pb.MON{}
	ret.SetLeader(m.Leader)
	ret.SetName(m.Name)
	ret.SetRank(m.Rank)
	ret.SetPublicAddress(m.PublicAddress)
	ret.SetMachine(toProtoStorageMachine(m.MachineID, m.MachineHostname))
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
	ret.SetId(o.ID)
	ret.SetName(o.Name)
	ret.SetUp(o.Up)
	ret.SetIn(o.In)
	ret.SetExists(o.Exists)
	ret.SetDeviceClass(o.DeviceClass)
	ret.SetSizeBytes(o.Size)
	ret.SetUsedBytes(o.Used)
	ret.SetPlacementGroupCount(o.PGCount)
	ret.SetMachine(toProtoStorageMachine(o.MachineID, o.MachineHostname))
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

func toProtoPoolType(s string) pb.PoolType {
	if pt, ok := pb.PoolType_value[strings.ToUpper(s)]; ok {
		return pb.PoolType(pt)
	}
	return pb.PoolType_UNSPECIFIED
}

func toProtoPool(p *core.Pool) *pb.Pool {
	ret := &pb.Pool{}
	ret.SetId(p.ID)
	ret.SetName(p.Name)
	ret.SetUpdating(p.Updating)
	ret.SetPoolType(toProtoPoolType(p.Type))
	ret.SetEcOverwrites(p.ECOverwrites)
	ret.SetDataChunks(p.DataChunks)
	ret.SetCodingChunks(p.CodingChunks)
	ret.SetReplicatedSize(p.ReplicatedSize)
	ret.SetQuotaBytes(p.QuotaBytes)
	ret.SetQuotaObjects(p.QuotaObjects)
	ret.SetUsedBytes(p.UsedBytes)
	ret.SetUsedObjects(p.UsedObjects)
	ret.SetMaxBytes(p.MaxBytes)
	ret.SetPlacementGroupCount(p.PlacementGroupCount)
	ret.SetPlacementGroupState(p.PlacementGroupState)
	ret.SetCreatedAt(timestamppb.New(p.CreatedAt))
	ret.SetApplications(p.Applications)
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
	ret.SetPoolName(i.PoolName)
	ret.SetObjectSizeBytes(i.ObjectSize)
	ret.SetStripeUnitBytes(i.StripeUnit)
	ret.SetStripeCount(i.StripeCount)
	ret.SetQuotaBytes(i.Quota)
	ret.SetUsedBytes(i.Used)
	ret.SetObjectCount(i.ObjectCount)
	ret.SetLayering(i.FeatureLayering)
	ret.SetExclusiveLock(i.FeatureExclusiveLock)
	ret.SetObjectMap(i.FeatureObjectMap)
	ret.SetFastDiff(i.FeatureFastDiff)
	ret.SetDeepFlatten(i.FeatureDeepFlatten)
	ret.SetCreatedAt(timestamppb.New(i.CreatedAt))
	ret.SetSnapshots(toProtoImageSnapshots(i.Snapshots))
	return ret
}

func toProtoImageSnapshots(ss []core.RBDImageSnapshot) []*pb.Image_Snapshot {
	ret := []*pb.Image_Snapshot{}
	for i := range ss {
		ret = append(ret, toProtoImageSnapshot(&ss[i]))
	}
	return ret
}

func toProtoImageSnapshot(s *core.RBDImageSnapshot) *pb.Image_Snapshot {
	ret := &pb.Image_Snapshot{}
	ret.SetName(s.Name)
	ret.SetProtected(s.Protected)
	ret.SetQuotaBytes(s.Quota)
	ret.SetUsedBytes(s.Used)
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
	ret.SetId(v.ID)
	ret.SetName(v.Name)
	ret.SetCreatedAt(timestamppb.New(v.CreatedAt))
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
	ret.SetPath(s.Path)
	ret.SetMode(s.Mode)
	ret.SetPoolName(s.PoolName)
	ret.SetQuotaBytes(s.Quota)
	ret.SetUsedBytes(s.Used)
	ret.SetCreatedAt(timestamppb.New(s.CreatedAt))
	if s.Export != nil {
		ret.SetExport(toProtoSubvolumeExport(s.Export))
	}
	ret.SetSnapshots(toProtoSubvolumeSnapshots(s.Snapshots))
	return ret
}

func toProtoSubvolumeExport(e *core.SubvolumeExport) *pb.Subvolume_Export {
	if len(e.Clients) == 0 {
		return nil
	}
	ret := &pb.Subvolume_Export{}
	ret.SetIp(e.IP)
	ret.SetPath(e.Path)
	ret.SetClients(e.Clients)
	ret.SetCommand(e.Command)
	return ret
}

func toProtoSubvolumeSnapshots(ss []core.SubvolumeSnapshot) []*pb.Subvolume_Snapshot {
	ret := []*pb.Subvolume_Snapshot{}
	for i := range ss {
		ret = append(ret, toProtoSubvolumeSnapshot(&ss[i]))
	}
	return ret
}

func toProtoSubvolumeSnapshot(s *core.SubvolumeSnapshot) *pb.Subvolume_Snapshot {
	ret := &pb.Subvolume_Snapshot{}
	ret.SetName(s.Name)
	ret.SetHasPendingClones(s.HasPendingClones)
	ret.SetCreatedAt(timestamppb.New(s.CreatedAt))
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
	ret.SetMode(s.Mode)
	ret.SetPoolName(s.PoolName)
	ret.SetQuotaBytes(s.Quota)
	ret.SetUsedBytes(s.Used)
	ret.SetCreatedAt(timestamppb.New(s.CreatedAt))
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
	ret.SetOwner(b.Owner)
	ret.SetName(b.Bucket.Bucket)
	if b.Policy != nil {
		ret.SetPolicy(*b.Policy)
	}
	ret.SetGrants(toProtoBucketGrants(b.Grants))
	usedBytes := b.Usage.RgwMain.Size
	if usedBytes != nil {
		ret.SetUsedBytes(*usedBytes)
	}
	usedObjects := b.Usage.RgwMain.NumObjects
	if usedObjects != nil {
		ret.SetUsedObjects(*usedObjects)
	}
	createdAt := b.CreationTime
	if createdAt != nil {
		ret.SetCreatedAt(timestamppb.New(*createdAt))
	}
	return ret
}

func toProtoBucketGrants(gs []core.RGWGrant) []*pb.Bucket_Grant {
	ret := []*pb.Bucket_Grant{}
	for i := range gs {
		ret = append(ret, toProtoBucketGrant(&gs[i]))
	}
	return ret
}

func toProtoBucketGrant(g *core.RGWGrant) *pb.Bucket_Grant {
	ret := &pb.Bucket_Grant{}
	ret.SetType(string(g.Grantee.Type))
	id := g.Grantee.ID
	if id != nil {
		ret.SetId(*id)
	}
	name := g.Grantee.DisplayName
	if name != nil {
		ret.SetName(*name)
	}
	uri := g.Grantee.URI
	if uri != nil {
		ret.SetUri(*uri)
	}
	ret.SetPermission(string(g.Permission))
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
	ret.SetId(u.ID)
	ret.SetName(u.DisplayName)
	ret.SetSuspended(u.Suspended != nil && *u.Suspended == 1)
	ret.SetKeys(toProtoUserKeys(u.Keys))
	return ret
}

func toProtoUserKeys(uks []core.RGWUserKey) []*pb.User_Key {
	ret := []*pb.User_Key{}
	for i := range uks {
		ret = append(ret, toProtoUserKey(&uks[i]))
	}
	return ret
}

func toProtoUserKey(uk *core.RGWUserKey) *pb.User_Key {
	ret := &pb.User_Key{}
	ret.SetAccessKey(uk.AccessKey)
	ret.SetSecretKey(uk.SecretKey)
	return ret
}

func toCoreSMBShareConfig(browsable, readOnly, guestOk bool, mapToGuest pb.SMBShare_MapToGuest, validUsers []string) *core.SMBShareConfig {
	mapToGuestStr := ""
	switch mapToGuest {
	case pb.SMBShare_NEVER:
		mapToGuestStr = "never"
	case pb.SMBShare_BAD_USER:
		mapToGuestStr = "bad user"
	case pb.SMBShare_BAD_PASSWORD:
		mapToGuestStr = "bad password"
	}

	validUsersStr := ""
	if len(validUsers) > 0 {
		validUsersStr = strings.Join(validUsers, " ")
	}

	return &core.SMBShareConfig{
		Browseable: browsable,
		ReadOnly:   readOnly,
		GuestOk:    guestOk,
		ValidUsers: validUsersStr,
		MapToGuest: mapToGuestStr,
	}
}

func toCoreSMBShareAuth(activeDirectory *pb.SMBShare_ActiveDirectory, localUser *pb.SMBShare_LocalUser) (*core.ADAuth, *core.LocalAuth) {
	var adAuth *core.ADAuth
	var localAuth *core.LocalAuth

	if activeDirectory != nil {
		adAuth = &core.ADAuth{
			Realm:        activeDirectory.GetRealm(),
			JoinUsername: activeDirectory.GetJoinSource().GetUsername(),
			JoinPassword: activeDirectory.GetJoinSource().GetPassword(),
		}
	}

	if localUser != nil {
		users := make([]core.SMBUser, 0, len(localUser.GetUsers()))
		for _, u := range localUser.GetUsers() {
			users = append(users, core.SMBUser{
				Username: u.GetUsername(),
				Password: u.GetPassword(),
			})
		}
		localAuth = &core.LocalAuth{
			Users: users,
		}
	}

	return adAuth, localAuth
}

func normalizeNamespace(namespace string) string {
	if namespace == "" {
		return "default"
	}
	return namespace
}

func toProtoSMBShares(shares []core.SMBShare) []*pb.SMBShare {
	ret := make([]*pb.SMBShare, 0, len(shares))
	for i := range shares {
		ret = append(ret, toProtoSMBShare(&shares[i]))
	}
	return ret
}

func toProtoSMBShare(share *core.SMBShare) *pb.SMBShare {
	ret := &pb.SMBShare{}
	ret.SetName(share.Name)
	ret.SetNamespace(share.Namespace)
	ret.SetStatus(share.Status)
	ret.SetSizeBytes(share.SizeBytes)
	ret.SetBrowsable(share.Browseable)
	ret.SetReadOnly(share.ReadOnly)
	ret.SetGuestOk(share.GuestOk)

	switch share.MapToGuest {
	case "never":
		ret.SetMapToGuest(pb.SMBShare_NEVER)
	case "bad user":
		ret.SetMapToGuest(pb.SMBShare_BAD_USER)
	case "bad password":
		ret.SetMapToGuest(pb.SMBShare_BAD_PASSWORD)
	}

	// Convert string SecurityMode to enum
	switch share.SecurityMode {
	case "user":
		ret.SetSecurityMode(pb.SMBShare_USER)
	case "active-directory":
		ret.SetSecurityMode(pb.SMBShare_ACTIVE_DIRECTORY)
	}

	// Set authentication information
	if share.ADAuth != nil {
		adAuth := &pb.SMBShare_ActiveDirectory{}
		adAuth.SetRealm(share.ADAuth.Realm)
		if share.ADAuth.JoinUsername != "" {
			joinSource := &pb.SMBShare_User{}
			joinSource.SetUsername(share.ADAuth.JoinUsername)
			adAuth.SetJoinSource(joinSource)
		}
		ret.SetActiveDirectory(adAuth)
	} else if share.LocalAuth != nil {
		localAuth := &pb.SMBShare_LocalUser{}
		users := make([]*pb.SMBShare_User, 0, len(share.LocalAuth.Users))
		for _, u := range share.LocalAuth.Users {
			user := &pb.SMBShare_User{}
			user.SetUsername(u.Username)
			user.SetPassword(u.Password)
			users = append(users, user)
		}
		localAuth.SetUsers(users)
		ret.SetLocalUser(localAuth)
	}

	// Set valid users
	if len(share.ValidUsers) > 0 {
		ret.SetValidUsers(share.ValidUsers)
	}

	return ret
}
