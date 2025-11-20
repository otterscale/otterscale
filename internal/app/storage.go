package app

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/storage/v1"
	"github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/storage"
	"github.com/otterscale/otterscale/internal/core/storage/block"
	"github.com/otterscale/otterscale/internal/core/storage/file"
	"github.com/otterscale/otterscale/internal/core/storage/object"
	"github.com/otterscale/otterscale/internal/core/storage/smb"
)

type StorageService struct {
	pbconnect.UnimplementedStorageServiceHandler

	storage *storage.UseCase
	block   *block.UseCase
	file    *file.UseCase
	object  *object.UseCase
	smb     *smb.UseCase
}

func NewStorageService(storage *storage.UseCase, block *block.UseCase, file *file.UseCase, object *object.UseCase, smb *smb.UseCase) *StorageService {
	return &StorageService{
		storage: storage,
		block:   block,
		file:    file,
		object:  object,
		smb:     smb,
	}
}

var _ pbconnect.StorageServiceHandler = (*StorageService)(nil)

func (s *StorageService) ListMonitors(ctx context.Context, req *pb.ListMonitorsRequest) (*pb.ListMonitorsResponse, error) {
	mons, err := s.storage.ListMonitors(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListMonitorsResponse{}
	resp.SetMonitors(toProtoMonitors(mons))
	return resp, nil
}

func (s *StorageService) ListObjectStorageDaemons(ctx context.Context, req *pb.ListObjectStorageDaemonsRequest) (*pb.ListObjectStorageDaemonsResponse, error) {
	osds, err := s.storage.ListObjectStorageDaemons(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListObjectStorageDaemonsResponse{}
	resp.SetObjectStorageDaemons(toProtoObjectStorageDaemons(osds))
	return resp, nil
}

func (s *StorageService) DoSMART(ctx context.Context, req *pb.DoSMARTRequest) (*pb.DoSMARTResponse, error) {
	outputs, err := s.storage.DoSMART(ctx, req.GetScope(), req.GetOsdName())
	if err != nil {
		return nil, err
	}

	resp := &pb.DoSMARTResponse{}
	resp.SetDeviceOutputMap(toProtoDeviceOutputMap(outputs))
	return resp, nil
}

func (s *StorageService) ListPools(ctx context.Context, req *pb.ListPoolsRequest) (*pb.ListPoolsResponse, error) {
	pools, err := s.storage.ListPools(ctx, req.GetScope(), req.GetApplication())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListPoolsResponse{}
	resp.SetPools(toProtoPools(pools))
	return resp, nil
}

func (s *StorageService) CreatePool(ctx context.Context, req *pb.CreatePoolRequest) (*pb.Pool, error) {
	pool, err := s.storage.CreatePool(ctx,
		req.GetScope(),
		req.GetPoolName(),
		strings.ToLower(req.GetPoolType().String()),
		req.GetEcOverwrites(),
		req.GetReplicatedSize(),
		req.GetQuotaBytes(),
		req.GetQuotaObjects(),
		req.GetApplications(),
	)
	if err != nil {
		return nil, err
	}

	resp := toProtoPool(pool)
	return resp, nil
}

func (s *StorageService) UpdatePool(ctx context.Context, req *pb.UpdatePoolRequest) (*pb.Pool, error) {
	pool, err := s.storage.UpdatePool(ctx, req.GetScope(), req.GetPoolName(), req.GetQuotaBytes(), req.GetQuotaObjects())
	if err != nil {
		return nil, err
	}

	resp := toProtoPool(pool)
	return resp, nil
}

func (s *StorageService) DeletePool(ctx context.Context, req *pb.DeletePoolRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeletePool(ctx, req.GetScope(), req.GetPoolName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	imgs, err := s.block.ListImages(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListImagesResponse{}
	resp.SetImages(toProtoImages(imgs))
	return resp, nil
}

func (s *StorageService) CreateImage(ctx context.Context, req *pb.CreateImageRequest) (*pb.Image, error) {
	img, err := s.block.CreateImage(ctx,
		req.GetScope(),
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
	img, err := s.block.UpdateImage(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}

	resp := toProtoImage(img)
	return resp, nil
}

func (s *StorageService) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*emptypb.Empty, error) {
	if err := s.block.DeleteImage(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateImageSnapshot(ctx context.Context, req *pb.CreateImageSnapshotRequest) (*pb.Image_Snapshot, error) {
	snap, err := s.block.CreateImageSnapshot(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}

	resp := toProtoImageSnapshot(snap)
	return resp, nil
}

func (s *StorageService) DeleteImageSnapshot(ctx context.Context, req *pb.DeleteImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.block.DeleteImageSnapshot(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) RollbackImageSnapshot(ctx context.Context, req *pb.RollbackImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.block.RollbackImageSnapshot(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ProtectImageSnapshot(ctx context.Context, req *pb.ProtectImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.block.ProtectImageSnapshot(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) UnprotectImageSnapshot(ctx context.Context, req *pb.UnprotectImageSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.block.UnprotectImageSnapshot(ctx, req.GetScope(), req.GetPoolName(), req.GetImageName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListVolumes(ctx context.Context, req *pb.ListVolumesRequest) (*pb.ListVolumesResponse, error) {
	volumes, err := s.file.ListVolumes(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListVolumesResponse{}
	resp.SetVolumes(toProtoVolumes(volumes))
	return resp, nil
}

func (s *StorageService) ListSubvolumes(ctx context.Context, req *pb.ListSubvolumesRequest) (*pb.ListSubvolumesResponse, error) {
	subvolumes, err := s.file.ListSubvolumes(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListSubvolumesResponse{}
	resp.SetSubvolumes(toProtoSubvolumes(subvolumes))
	return resp, nil
}

func (s *StorageService) CreateSubvolume(ctx context.Context, req *pb.CreateSubvolumeRequest) (*pb.Subvolume, error) {
	subvolume, err := s.file.CreateSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetQuotaBytes(), req.GetExport())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) UpdateSubvolume(ctx context.Context, req *pb.UpdateSubvolumeRequest) (*pb.Subvolume, error) {
	subvolume, err := s.file.UpdateSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) DeleteSubvolume(ctx context.Context, req *pb.DeleteSubvolumeRequest) (*emptypb.Empty, error) {
	if err := s.file.DeleteSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) GrantSubvolumeExportAccess(ctx context.Context, req *pb.GrantSubvolumeExportAccessRequest) (*emptypb.Empty, error) {
	if err := s.file.GrantSubvolumeClient(ctx, req.GetScope(), req.GetSubvolumeName(), req.GetClientIp()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) RevokeSubvolumeExportAccess(ctx context.Context, req *pb.RevokeSubvolumeExportAccessRequest) (*emptypb.Empty, error) {
	if err := s.file.RevokeSubvolumeClient(ctx, req.GetScope(), req.GetSubvolumeName(), req.GetClientIp()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateSubvolumeSnapshot(ctx context.Context, req *pb.CreateSubvolumeSnapshotRequest) (*pb.Subvolume_Snapshot, error) {
	snapshot, err := s.file.CreateSubvolumeSnapshot(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolumeSnapshot(snapshot)
	return resp, nil
}

func (s *StorageService) DeleteSubvolumeSnapshot(ctx context.Context, req *pb.DeleteSubvolumeSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.file.DeleteSubvolumeSnapshot(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListSubvolumeGroups(ctx context.Context, req *pb.ListSubvolumeGroupsRequest) (*pb.ListSubvolumeGroupsResponse, error) {
	groups, err := s.file.ListSubvolumeGroups(ctx, req.GetScope(), req.GetVolumeName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListSubvolumeGroupsResponse{}
	resp.SetSubvolumeGroups(toProtoSubvolumeGroups(groups))
	return resp, nil
}

func (s *StorageService) CreateSubvolumeGroup(ctx context.Context, req *pb.CreateSubvolumeGroupRequest) (*pb.SubvolumeGroup, error) {
	group, err := s.file.CreateSubvolumeGroup(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolumeGroup(group)
	return resp, nil
}

func (s *StorageService) UpdateSubvolumeGroup(ctx context.Context, req *pb.UpdateSubvolumeGroupRequest) (*pb.SubvolumeGroup, error) {
	group, err := s.file.UpdateSubvolumeGroup(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetQuotaBytes())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolumeGroup(group)
	return resp, nil
}

func (s *StorageService) DeleteSubvolumeGroup(ctx context.Context, req *pb.DeleteSubvolumeGroupRequest) (*emptypb.Empty, error) {
	if err := s.file.DeleteSubvolumeGroup(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListBuckets(ctx context.Context, req *pb.ListBucketsRequest) (*pb.ListBucketsResponse, error) {
	buckets, err := s.object.ListBuckets(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListBucketsResponse{}
	resp.SetBuckets(toProtoBuckets(buckets))
	return resp, nil
}

func (s *StorageService) CreateBucket(ctx context.Context, req *pb.CreateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.object.CreateBucket(ctx, req.GetScope(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), toACL(req.GetAcl().String()))
	if err != nil {
		return nil, err
	}

	resp := toProtoBucket(bucket)
	return resp, nil
}

func (s *StorageService) UpdateBucket(ctx context.Context, req *pb.UpdateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.object.UpdateBucket(ctx, req.GetScope(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), toACL(req.GetAcl().String()))
	if err != nil {
		return nil, err
	}

	resp := toProtoBucket(bucket)
	return resp, nil
}

func (s *StorageService) DeleteBucket(ctx context.Context, req *pb.DeleteBucketRequest) (*emptypb.Empty, error) {
	if err := s.object.DeleteBucket(ctx, req.GetScope(), req.GetBucketName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.object.ListUsers(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListUsersResponse{}
	resp.SetUsers(toProtoUsers(users))
	return resp, nil
}

func (s *StorageService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.object.CreateUser(ctx, req.GetScope(), req.GetUserId(), req.GetUserName(), req.GetSuspended())
	if err != nil {
		return nil, err
	}

	resp := toProtoUser(user)
	return resp, nil
}

func (s *StorageService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := s.object.UpdateUser(ctx, req.GetScope(), req.GetUserId(), req.GetUserName(), req.GetSuspended())
	if err != nil {
		return nil, err
	}

	resp := toProtoUser(user)
	return resp, nil
}

func (s *StorageService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.object.DeleteUser(ctx, req.GetScope(), req.GetUserId()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) CreateUserKey(ctx context.Context, req *pb.CreateUserKeyRequest) (*pb.User_Key, error) {
	key, err := s.object.CreateUserKey(ctx, req.GetScope(), req.GetUserId())
	if err != nil {
		return nil, err
	}

	resp := toProtoUserKey(key)
	return resp, nil
}

func (s *StorageService) DeleteUserKey(ctx context.Context, req *pb.DeleteUserKeyRequest) (*emptypb.Empty, error) {
	if err := s.object.DeleteUserKey(ctx, req.GetScope(), req.GetUserId(), req.GetAccessKey()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListSMBShares(ctx context.Context, req *pb.ListSMBSharesRequest) (*pb.ListSMBSharesResponse, error) {
	shares, hostname, err := s.smb.ListSMBShares(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListSMBSharesResponse{}
	resp.SetSmbShares(toProtoSMBShares(shares, hostname))
	return resp, nil
}

func (s *StorageService) CreateSMBShare(ctx context.Context, req *pb.CreateSMBShareRequest) (*pb.SMBShare, error) {
	var (
		mapToGuest   string
		securityMode string
		localUsers   []smb.User
		realm        string
		joinSource   *smb.User
	)

	commonConfig := req.GetCommonConfig()
	if commonConfig != nil {
		mapToGuest = toMapToGuest(commonConfig.GetMapToGuest())
	}

	securityConfig := req.GetSecurityConfig()
	if securityConfig != nil {
		securityMode = toSecurityMode(securityConfig.GetMode())
		localUsers = toUsers(securityConfig.GetLocalUsers())
		realm = securityConfig.GetRealm()
		joinSource = toUser(securityConfig.GetJoinSource())
	}

	share, hostname, err := s.smb.CreateSMBShare(ctx,
		req.GetScope(),
		req.GetName(),
		req.GetSizeBytes(),
		req.GetPort(),
		req.GetBrowsable(),
		req.GetReadOnly(),
		req.GetGuestOk(),
		req.GetValidUsers(),
		mapToGuest,
		securityMode,
		localUsers,
		realm,
		joinSource)
	if err != nil {
		return nil, err
	}

	resp := toProtoSMBShare(share, hostname)
	return resp, nil
}

func (s *StorageService) UpdateSMBShare(ctx context.Context, req *pb.UpdateSMBShareRequest) (*pb.SMBShare, error) {
	var (
		mapToGuest   string
		securityMode string
		localUsers   []smb.User
		realm        string
		joinSource   *smb.User
	)

	commonConfig := req.GetCommonConfig()
	if commonConfig != nil {
		mapToGuest = toMapToGuest(commonConfig.GetMapToGuest())
	}

	securityConfig := req.GetSecurityConfig()
	if securityConfig != nil {
		securityMode = toSecurityMode(securityConfig.GetMode())
		localUsers = toUsers(securityConfig.GetLocalUsers())
		realm = securityConfig.GetRealm()
		joinSource = toUser(securityConfig.GetJoinSource())
	}

	share, hostname, err := s.smb.UpdateSMBShare(ctx,
		req.GetScope(),
		req.GetName(),
		req.GetSizeBytes(),
		req.GetPort(),
		req.GetBrowsable(),
		req.GetReadOnly(),
		req.GetGuestOk(),
		req.GetValidUsers(),
		mapToGuest,
		securityMode,
		localUsers,
		realm,
		joinSource)
	if err != nil {
		return nil, err
	}

	resp := toProtoSMBShare(share, hostname)
	return resp, nil
}

func (s *StorageService) ValidateSMBUser(ctx context.Context, req *pb.ValidateSMBUserRequest) (*pb.ValidateSMBUserResponse, error) {
	result, err := s.smb.ValidateSMBUser(ctx,
		req.GetRealm(),
		req.GetUsername(),
		req.GetPassword(),
		req.GetSearchUsername(),
		req.GetTls())
	if err != nil {
		return nil, err
	}

	resp := &pb.ValidateSMBUserResponse{}
	resp.SetValid(result.Valid)
	resp.SetEntityType(toProtoEntityType(result.EntityType))
	resp.SetMessage(result.Message)
	return resp, nil
}

func toProtoEntityType(et int) pb.ValidateSMBUserResponse_EntityType {
	switch et {
	case smb.EntityTypeUser:
		return pb.ValidateSMBUserResponse_USER
	case smb.EntityTypeGroup:
		return pb.ValidateSMBUserResponse_GROUP
	default:
		return pb.ValidateSMBUserResponse_UNKNOWN
	}
}

func toACL(str string) object.BucketCannedACL {
	acl := object.BucketCannedACL(strings.ToLower(strings.Join(strings.Split(str, "_"), "-")))

	if slices.Contains(acl.Values(), acl) {
		return acl
	}

	return object.BucketCannedACLPrivate
}

func toMapToGuest(s pb.SMBShare_CommonConfig_MapToGuest) string {
	ret := strings.ToLower(s.String())
	ret = strings.ReplaceAll(ret, "_", " ")
	return ret
}

func toSecurityMode(s pb.SMBShare_SecurityConfig_Mode) string {
	ret := strings.ToLower(s.String())
	ret = strings.ReplaceAll(ret, "_", "-")
	return ret
}

func toUser(u *pb.SMBShare_SecurityConfig_User) *smb.User {
	return &smb.User{
		Username: u.GetUsername(),
		Password: u.GetPassword(),
	}
}

func toUsers(us []*pb.SMBShare_SecurityConfig_User) []smb.User {
	ret := []smb.User{}

	for _, u := range us {
		ret = append(ret, *toUser(u))
	}

	return ret
}

func toProtoStorageMachine(m *machine.Machine) *pb.Machine {
	ret := &pb.Machine{}
	ret.SetId(m.SystemID)
	ret.SetHostname(m.Hostname)
	return ret
}

func toProtoMonitors(ms []storage.Monitor) []*pb.Monitor {
	ret := []*pb.Monitor{}

	for i := range ms {
		ret = append(ret, toProtoMonitor(&ms[i]))
	}

	return ret
}

func toProtoMonitor(m *storage.Monitor) *pb.Monitor {
	ret := &pb.Monitor{}
	ret.SetLeader(m.Leader)
	ret.SetName(m.Name)
	ret.SetRank(m.Rank)
	ret.SetPublicAddress(m.PublicAddress)

	if m.Machine != nil {
		ret.SetMachine(toProtoStorageMachine(m.Machine))
	}

	return ret
}

func toProtoObjectStorageDaemons(os []storage.ObjectStorageDaemon) []*pb.ObjectStorageDaemon {
	ret := []*pb.ObjectStorageDaemon{}

	for i := range os {
		ret = append(ret, toProtoObjectStorageDaemon(&os[i]))
	}

	return ret
}

func toProtoObjectStorageDaemon(o *storage.ObjectStorageDaemon) *pb.ObjectStorageDaemon {
	ret := &pb.ObjectStorageDaemon{}
	ret.SetId(o.ID)
	ret.SetName(o.Name)
	ret.SetUp(o.Up)
	ret.SetIn(o.In)
	ret.SetExists(o.Exists)
	ret.SetDeviceClass(o.DeviceClass)
	ret.SetSizeBytes(o.Size)
	ret.SetUsedBytes(o.Used)
	ret.SetPlacementGroupCount(o.PGCount)

	if o.Machine != nil {
		ret.SetMachine(toProtoStorageMachine(o.Machine))
	}

	return ret
}

func toProtoDeviceOutputMap(m map[string][]string) map[string]*pb.DoSMARTResponse_Output {
	ret := map[string]*pb.DoSMARTResponse_Output{}

	for device, lines := range m {
		output := &pb.DoSMARTResponse_Output{}
		output.SetLines(lines)
		ret[device] = output
	}

	return ret
}

func toProtoPools(ps []storage.Pool) []*pb.Pool {
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

func toProtoPool(p *storage.Pool) *pb.Pool {
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

func toProtoImages(is []block.Image) []*pb.Image {
	ret := []*pb.Image{}

	for i := range is {
		ret = append(ret, toProtoImage(&is[i]))
	}

	return ret
}

func toProtoImage(i *block.Image) *pb.Image {
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

func toProtoImageSnapshots(ss []block.ImageSnapshot) []*pb.Image_Snapshot {
	ret := []*pb.Image_Snapshot{}

	for i := range ss {
		ret = append(ret, toProtoImageSnapshot(&ss[i]))
	}

	return ret
}

func toProtoImageSnapshot(s *block.ImageSnapshot) *pb.Image_Snapshot {
	ret := &pb.Image_Snapshot{}
	ret.SetName(s.Name)
	ret.SetProtected(s.Protected)
	ret.SetQuotaBytes(s.Quota)
	ret.SetUsedBytes(s.Used)
	return ret
}

func toProtoVolumes(vs []file.Volume) []*pb.Volume {
	ret := []*pb.Volume{}

	for i := range vs {
		ret = append(ret, toProtoVolume(&vs[i]))
	}

	return ret
}

func toProtoVolume(v *file.Volume) *pb.Volume {
	ret := &pb.Volume{}
	ret.SetId(v.ID)
	ret.SetName(v.Name)
	ret.SetCreatedAt(timestamppb.New(v.CreatedAt))
	return ret
}

func toProtoSubvolumes(ss []file.Subvolume) []*pb.Subvolume {
	ret := []*pb.Subvolume{}

	for i := range ss {
		ret = append(ret, toProtoSubvolume(&ss[i]))
	}

	return ret
}

func toProtoSubvolume(s *file.Subvolume) *pb.Subvolume {
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

func toProtoSubvolumeExport(e *file.SubvolumeExport) *pb.Subvolume_Export {
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

func toProtoSubvolumeSnapshots(ss []file.SubvolumeSnapshot) []*pb.Subvolume_Snapshot {
	ret := []*pb.Subvolume_Snapshot{}

	for i := range ss {
		ret = append(ret, toProtoSubvolumeSnapshot(&ss[i]))
	}

	return ret
}

func toProtoSubvolumeSnapshot(s *file.SubvolumeSnapshot) *pb.Subvolume_Snapshot {
	ret := &pb.Subvolume_Snapshot{}
	ret.SetName(s.Name)
	ret.SetHasPendingClones(s.HasPendingClones)
	ret.SetCreatedAt(timestamppb.New(s.CreatedAt))
	return ret
}

func toProtoSubvolumeGroups(ss []file.SubvolumeGroup) []*pb.SubvolumeGroup {
	ret := []*pb.SubvolumeGroup{}

	for i := range ss {
		ret = append(ret, toProtoSubvolumeGroup(&ss[i]))
	}

	return ret
}

func toProtoSubvolumeGroup(s *file.SubvolumeGroup) *pb.SubvolumeGroup {
	ret := &pb.SubvolumeGroup{}
	ret.SetName(s.Name)
	ret.SetMode(s.Mode)
	ret.SetPoolName(s.PoolName)
	ret.SetQuotaBytes(s.Quota)
	ret.SetUsedBytes(s.Used)
	ret.SetCreatedAt(timestamppb.New(s.CreatedAt))
	return ret
}

func toProtoBuckets(bs []object.BucketData) []*pb.Bucket {
	ret := []*pb.Bucket{}

	for i := range bs {
		ret = append(ret, toProtoBucket(&bs[i]))
	}

	return ret
}

func toProtoBucket(b *object.BucketData) *pb.Bucket {
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

func toProtoBucketGrants(gs []object.Grant) []*pb.Bucket_Grant {
	ret := []*pb.Bucket_Grant{}

	for i := range gs {
		ret = append(ret, toProtoBucketGrant(&gs[i]))
	}

	return ret
}

func toProtoBucketGrant(g *object.Grant) *pb.Bucket_Grant {
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

func toProtoUsers(us []object.User) []*pb.User {
	ret := []*pb.User{}

	for i := range us {
		ret = append(ret, toProtoUser(&us[i]))
	}

	return ret
}

func toProtoUser(u *object.User) *pb.User {
	ret := &pb.User{}
	ret.SetId(u.ID)
	ret.SetName(u.DisplayName)
	ret.SetSuspended(u.Suspended != nil && *u.Suspended == 1)
	ret.SetKeys(toProtoUserKeys(u.Keys))
	return ret
}

func toProtoUserKeys(uks []object.UserKey) []*pb.User_Key {
	ret := []*pb.User_Key{}

	for i := range uks {
		ret = append(ret, toProtoUserKey(&uks[i]))
	}

	return ret
}

func toProtoUserKey(uk *object.UserKey) *pb.User_Key {
	ret := &pb.User_Key{}
	ret.SetAccessKey(uk.AccessKey)
	ret.SetSecretKey(uk.SecretKey)
	return ret
}

func toProtoSMBShareCommonConfigMapToGuest(mapToGuest string) pb.SMBShare_CommonConfig_MapToGuest {
	if strings.EqualFold(mapToGuest, "bad user") {
		return pb.SMBShare_CommonConfig_BAD_USER
	} else if strings.EqualFold(mapToGuest, "bad password") {
		return pb.SMBShare_CommonConfig_BAD_PASSWORD
	}
	return pb.SMBShare_CommonConfig_NEVER
}

func toProtoSMBShareSecurityConfigMode(mode string) pb.SMBShare_SecurityConfig_Mode {
	if strings.EqualFold(mode, "active-directory") {
		return pb.SMBShare_SecurityConfig_ACTIVE_DIRECTORY
	}
	return pb.SMBShare_SecurityConfig_USER
}

func toProtoSMBShareSecurityConfigUsers(us []smb.User) []*pb.SMBShare_SecurityConfig_User {
	ret := []*pb.SMBShare_SecurityConfig_User{}

	for i := range us {
		ret = append(ret, toProtoSMBShareSecurityConfigUser(&us[i]))
	}

	return ret
}

// ignore password in proto for security reason
func toProtoSMBShareSecurityConfigUser(u *smb.User) *pb.SMBShare_SecurityConfig_User {
	ret := &pb.SMBShare_SecurityConfig_User{}
	ret.SetUsername(u.Username)
	return ret
}

func toProtoSMBShareCommonConfig(cc *smb.CommonConfig) *pb.SMBShare_CommonConfig {
	mapToGuest := ""

	globalConfig := cc.Spec.CustomGlobalConfig
	if globalConfig != nil && globalConfig.Configs != nil {
		mapToGuest = globalConfig.Configs[smb.MapToGuestKey]
	}

	ret := &pb.SMBShare_CommonConfig{}
	ret.SetMapToGuest(toProtoSMBShareCommonConfigMapToGuest(mapToGuest))
	return ret
}

func toProtoSMBShareSecurityConfig(sc *smb.SecurityConfig, localUsers []smb.User, joinSource *smb.User) *pb.SMBShare_SecurityConfig {
	ret := &pb.SMBShare_SecurityConfig{}
	ret.SetMode(toProtoSMBShareSecurityConfigMode(sc.Spec.Mode))
	ret.SetLocalUsers(toProtoSMBShareSecurityConfigUsers(localUsers))
	ret.SetRealm(sc.Spec.Realm)

	if joinSource != nil && joinSource.Username != "" {
		ret.SetJoinSource(toProtoSMBShareSecurityConfigUser(joinSource))
	}

	return ret
}

func toProtoSMBShares(sds []smb.ShareData, hostname string) []*pb.SMBShare {
	ret := []*pb.SMBShare{}

	for i := range sds {
		ret = append(ret, toProtoSMBShare(&sds[i], hostname))
	}

	return ret
}

func toProtoSMBShare(sd *smb.ShareData, hostname string) *pb.SMBShare {
	ret := &pb.SMBShare{}

	share := sd.Share
	if share != nil {
		ret.SetName(share.Name)

		service := sd.Service
		if service != nil {
			ports := service.Spec.Ports
			if len(ports) > 0 {
				port := ports[0].NodePort
				uri := fmt.Sprintf("smb://%s:%d/%s", hostname, port, sd.Share.Name)

				ret.SetUri(uri)
			}
		}

		deployment := sd.Deployment
		if deployment != nil {
			replicas := deployment.Spec.Replicas
			if replicas != nil {
				ret.SetReplicas(*replicas)
			}
		}

		ret.SetHealthies(countHealthies(sd.Pods))

		pvc := share.Spec.Storage.Pvc
		if pvc != nil {
			spec := pvc.Spec
			if spec != nil {
				sizeBytes, _ := spec.Resources.Requests.Storage().AsInt64()
				ret.SetSizeBytes(uint64(sizeBytes)) //nolint:gosec // ignore
			}
		}

		ret.SetBrowsable(share.Spec.Browseable)
		ret.SetReadOnly(share.Spec.ReadOnly)

		config := share.Spec.CustomShareConfig
		if config != nil {
			configs := config.Configs
			if configs != nil {
				if val, ok := configs[smb.GuestOKkey]; ok {
					ret.SetGuestOk(strings.EqualFold(val, "yes") || strings.EqualFold(val, "true") || strings.EqualFold(val, "1"))
				}

				if val, ok := configs[smb.ValidUsersKey]; ok {
					users := strings.Split(val, " ")
					for i := range users {
						users[i] = strings.TrimSpace(users[i])
					}
					ret.SetValidUsers(users)
				}
			}
		}
	}

	commonConfig := sd.CommonConfig
	if commonConfig != nil {
		ret.SetCommonConfig(toProtoSMBShareCommonConfig(commonConfig))
	}

	securityConfig := sd.SecurityConfig
	if securityConfig != nil {
		ret.SetSecurityConfig(toProtoSMBShareSecurityConfig(securityConfig, sd.LocalUsers, sd.JoinSource))
	}

	return ret
}
