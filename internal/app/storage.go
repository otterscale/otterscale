package app

import (
	"context"
	"fmt"
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
	pools, err := s.storage.ListPools(ctx, req.GetScope(), toPoolApplication(req.GetApplication()))
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
		toPoolType(req.GetType()),
		req.GetEcOverwrites(),
		req.GetReplicatedSize(),
		req.GetQuotaBytes(),
		req.GetQuotaObjects(),
		toPoolApplications(req.GetApplications()),
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
	var (
		sizePtr *file.Bytes
		uidPtr  *uint32
		gidPtr  *uint32
		modePtr *file.UnixMode
		poolPtr *string
		nsPtr   *bool
	)

	if req.HasSize() {
		v := file.Bytes(req.GetSize())
		sizePtr = &v
	}
	if req.HasUid() {
		v := req.GetUid()
		uidPtr = &v
	}
	if req.HasGid() {
		v := req.GetGid()
		gidPtr = &v
	}
	if req.HasMode() {
		v := file.UnixMode(req.GetMode())
		modePtr = &v
	}
	if req.HasPoolLayout() {
		v := req.GetPoolLayout()
		poolPtr = &v
	}
	if req.HasIsNamespaceIsolated() {
		v := req.GetIsNamespaceIsolated()
		nsPtr = &v
	}

	subvolume, err := s.file.CreateSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetSubvolumeName(), sizePtr, uidPtr, gidPtr, modePtr, poolPtr, nsPtr)
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) GetSubvolume(ctx context.Context, req *pb.GetSubvolumeRequest) (*pb.Subvolume, error) {
	subvolume, err := s.file.GetSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetSubvolumeName())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) UpdateSubvolume(ctx context.Context, req *pb.UpdateSubvolumeRequest) (*pb.Subvolume, error) {
	var noShrinkPtr *bool
	if req.HasNoShrink() {
		v := req.GetNoShrink()
		noShrinkPtr = &v
	}

	subvolume, err := s.file.UpdateSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetSubvolumeName(), file.Bytes(req.GetNewSize()), noShrinkPtr)
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolume(subvolume)
	return resp, nil
}

func (s *StorageService) DeleteSubvolume(ctx context.Context, req *pb.DeleteSubvolumeRequest) (*emptypb.Empty, error) {
	var isForcePtr *bool
	if req.HasIsForce() {
		v := req.GetIsForce()
		isForcePtr = &v
	}
	if err := s.file.DeleteSubvolume(ctx, req.GetScope(), req.GetVolumeName(), req.GetGroupName(), req.GetSubvolumeName(), isForcePtr); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *StorageService) ListNFSExports(ctx context.Context, req *pb.ListNFSExportsRequest) (*pb.ListNFSExportsResponse, error) {
	nfsexports, err := s.file.ListNFSExports(ctx, req.GetScope(), req.GetClusterId())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListNFSExportsResponse{}
	resp.SetNfsExports(toProtoNfsExports(nfsexports))
	return resp, nil
}

func (s *StorageService) ApplyNFSExport(ctx context.Context, req *pb.ApplyNFSExportRequest) (*pb.NfsExport, error) {
	filesystem, group, subvolume := toCephFsTarget(req.GetTarget())

	spec := toNFSExportSpecFromApply(req)

	e, err := s.file.ApplyNFSExport(ctx, req.GetScope(), req.GetClusterId(), filesystem, group, subvolume, spec)
	if err != nil {
		return nil, err
	}

	return toProtoNfsExport(e), nil
}

func (s *StorageService) GetNFSExport(ctx context.Context, req *pb.GetNFSExportRequest) (*pb.NfsExport, error) {
	nfsexport, err := s.file.GetNFSExport(ctx, req.GetScope(), req.GetClusterId(), req.GetPseudoPath())
	if err != nil {
		return nil, err
	}

	resp := toProtoNfsExport(nfsexport)
	return resp, nil
}

func (s *StorageService) DeleteNFSExport(ctx context.Context, req *pb.DeleteNFSExportRequest) (*emptypb.Empty, error) {
	if err := s.file.DeleteNFSExport(ctx, req.GetScope(), req.GetClusterId(), req.GetPseudoPath()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

/*func (s *StorageService) CreateSubvolumeSnapshot(ctx context.Context, req *pb.CreateSubvolumeSnapshotRequest) (*pb.Subvolume_Snapshot, error) {
	snapshot, err := s.file.CreateSubvolumeSnapshot(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubvolumeSnapshot(snapshot)
	return resp, nil
}*/

/*func (s *StorageService) DeleteSubvolumeSnapshot(ctx context.Context, req *pb.DeleteSubvolumeSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.file.DeleteSubvolumeSnapshot(ctx, req.GetScope(), req.GetVolumeName(), req.GetSubvolumeName(), req.GetGroupName(), req.GetSnapshotName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}*/

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
	buckets, uri, err := s.object.ListBuckets(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListBucketsResponse{}
	resp.SetBuckets(toProtoBuckets(buckets))
	resp.SetServiceUri(uri)
	return resp, nil
}

func (s *StorageService) CreateBucket(ctx context.Context, req *pb.CreateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.object.CreateBucket(ctx, req.GetScope(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), toACL(req.GetAcl()))
	if err != nil {
		return nil, err
	}

	resp := toProtoBucket(bucket)
	return resp, nil
}

func (s *StorageService) UpdateBucket(ctx context.Context, req *pb.UpdateBucketRequest) (*pb.Bucket, error) {
	bucket, err := s.object.UpdateBucket(ctx, req.GetScope(), req.GetBucketName(), req.GetOwner(), req.GetPolicy(), toACL(req.GetAcl()))
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
	users, uri, err := s.object.ListUsers(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListUsersResponse{}
	resp.SetUsers(toProtoUsers(users))
	resp.SetServiceUri(uri)
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

	resp := toProtoUserKey(key, false)
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
		mapToGuest   smb.MapToGuest
		securityMode smb.SecurityMode
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
		mapToGuest   smb.MapToGuest
		securityMode smb.SecurityMode
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

func (s *StorageService) ValidateSMBUser(_ context.Context, req *pb.ValidateSMBUserRequest) (*pb.ValidateSMBUserResponse, error) {
	entityType, err := s.smb.ValidateSMBUser(
		req.GetRealm(),
		req.GetUsername(),
		req.GetPassword(),
		req.GetSearchUsername(),
		req.GetTls())
	if err != nil {
		return nil, err
	}

	resp := &pb.ValidateSMBUserResponse{}
	resp.SetEntityType(toProtoEntityType(entityType))
	return resp, nil
}

func toPoolType(pt pb.Pool_Type) storage.PoolType {
	switch pt {
	case pb.Pool_TYPE_ERASURE:
		return storage.PoolTypeErasure

	case pb.Pool_TYPE_REPLICATED:
		return storage.PoolTypeReplicated

	default:
		return storage.PoolTypeUnspecified
	}
}

func toPoolApplications(pas []pb.Pool_Application) []storage.PoolApplication {
	ret := []storage.PoolApplication{}

	for _, pa := range pas {
		ret = append(ret, toPoolApplication(pa))
	}

	return ret
}

func toPoolApplication(pa pb.Pool_Application) storage.PoolApplication {
	switch pa {
	case pb.Pool_APPLICATION_BLOCK:
		return storage.PoolApplicationBlock

	case pb.Pool_APPLICATION_FILE:
		return storage.PoolApplicationFile

	case pb.Pool_APPLICATION_OBJECT:
		return storage.PoolApplicationObject

	default:
		return storage.PoolApplicationUnspecified
	}
}

func toACL(acl pb.Bucket_ACL) object.BucketCannedACL {
	switch acl {
	case pb.Bucket_ACL_PRIVATE:
		return object.BucketCannedACLPrivate

	case pb.Bucket_ACL_PUBLIC_READ:
		return object.BucketCannedACLPublicRead

	case pb.Bucket_ACL_PUBLIC_READ_WRITE:
		return object.BucketCannedACLPublicReadWrite

	case pb.Bucket_ACL_AUTHENTICATED_READ:
		return object.BucketCannedACLAuthenticatedRead

	default:
		return object.BucketCannedACLPrivate
	}
}

func toMapToGuest(mtg pb.SMBShare_CommonConfig_MapToGuest) smb.MapToGuest {
	switch mtg {
	case pb.SMBShare_CommonConfig_MAP_TO_GUEST_NEVER:
		return smb.MapToGuestNever

	case pb.SMBShare_CommonConfig_MAP_TO_GUEST_BAD_USER:
		return smb.MapToGuestBadUser

	case pb.SMBShare_CommonConfig_MAP_TO_GUEST_BAD_PASSWORD:
		return smb.MapToGuestBadPassword

	default:
		return smb.MapToGuestNever
	}
}

func toSecurityMode(sm pb.SMBShare_SecurityConfig_Mode) smb.SecurityMode {
	switch sm {
	case pb.SMBShare_SecurityConfig_MODE_USER:
		return smb.SecurityModeUser

	case pb.SMBShare_SecurityConfig_MODE_ACTIVE_DIRECTORY:
		return smb.SecurityModeActiveDirectory

	default:
		return smb.SecurityModeUser
	}
}

func toUsers(us []*pb.SMBShare_SecurityConfig_User) []smb.User {
	ret := []smb.User{}

	for _, u := range us {
		ret = append(ret, *toUser(u))
	}

	return ret
}

func toUser(u *pb.SMBShare_SecurityConfig_User) *smb.User {
	return &smb.User{
		Username: u.GetUsername(),
		Password: u.GetPassword(),
	}
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

func toProtoPoolType(pt storage.PoolType) pb.Pool_Type {
	switch pt {
	case storage.PoolTypeErasure:
		return pb.Pool_TYPE_ERASURE

	case storage.PoolTypeReplicated:
		return pb.Pool_TYPE_REPLICATED

	default:
		return pb.Pool_TYPE_UNSPECIFIED
	}
}

func toProtoPoolApplications(pas []storage.PoolApplication) []pb.Pool_Application {
	ret := []pb.Pool_Application{}

	for i := range pas {
		ret = append(ret, toProtoPoolApplication(pas[i]))
	}

	return ret
}

func toProtoPoolApplication(pa storage.PoolApplication) pb.Pool_Application {
	switch pa {
	case storage.PoolApplicationBlock:
		return pb.Pool_APPLICATION_BLOCK

	case storage.PoolApplicationFile:
		return pb.Pool_APPLICATION_FILE

	case storage.PoolApplicationObject:
		return pb.Pool_APPLICATION_OBJECT

	default:
		return pb.Pool_APPLICATION_UNSPECIFIED
	}
}

func toProtoPool(p *storage.Pool) *pb.Pool {
	ret := &pb.Pool{}
	ret.SetId(p.ID)
	ret.SetName(p.Name)
	ret.SetUpdating(p.Updating)
	ret.SetType(toProtoPoolType(p.Type))
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
	ret.SetApplications(toProtoPoolApplications(p.Applications))
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
	ret := make([]*pb.Subvolume, 0, len(ss))
	for i := range ss {
		if p := toProtoSubvolume(&ss[i]); p != nil {
			ret = append(ret, p)
		}
	}
	return ret
}

func toProtoSubvolume(s *file.Subvolume) *pb.Subvolume {
	if s == nil {
		return nil
	}
	ret := &pb.Subvolume{}
	ret.SetVolumeName(s.Key.VolumeName)
	ret.SetGroupName(s.Key.GroupName)
	ret.SetSubvolumeName(s.Key.SubvolumeName)

	info := &pb.SubvolumeInfo{}
	info.SetPath(s.Info.Path)
	info.SetState(cephStateToPB(string(s.Info.State)))

	info.SetUid(s.Info.UID)
	info.SetGid(s.Info.GID)
	info.SetMode(uint32(s.Info.Mode))

	info.SetBytesPercent(s.Info.BytesPercent)
	info.SetBytesUsed(uint64(s.Info.BytesUsed))
	if s.Info.BytesQuota != nil {
		info.SetBytesQuota(uint64(*s.Info.BytesQuota))
	}

	info.SetDataPool(s.Info.DataPool)
	info.SetPoolNamespace(s.Info.PoolNamespace)

	info.SetAtime(timestamppb.New(s.Info.Atime))
	info.SetMtime(timestamppb.New(s.Info.Mtime))
	info.SetCtime(timestamppb.New(s.Info.Ctime))
	info.SetCreatedAt(timestamppb.New(s.Info.CreatedAt))

	info.SetFeatures(toProtoSubvolumeFeatures(s.Info.Features))

	ret.SetInfo(info)
	return ret
}

func cephStateToPB(s string) pb.SubvolumeInfo_SubvolumeState {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "init":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_INIT
	case "pending":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_PENDING
	case "in_progress", "in-progress", "inprogress":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_IN_PROGRESS
	case "failed", "fail", "error":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_FAILED
	case "complete", "completed", "ready":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_COMPLETE
	case "canceled", "cancelled", "cancel":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_CANCELED
	case "snapshot_retained", "snapshot-retained":
		return pb.SubvolumeInfo_SUBVOLUME_STATE_SNAPSHOT_RETAINED
	default:
		return pb.SubvolumeInfo_SUBVOLUME_STATE_UNSPECIFIED
	}
}

func toProtoSubvolumeFeatures(fs []file.Feature) []pb.SubvolumeInfo_Feature {
	ret := make([]pb.SubvolumeInfo_Feature, 0, len(fs))
	for _, f := range fs {
		switch strings.ToLower(string(f)) {
		case "snapshot-clone":
			ret = append(ret, pb.SubvolumeInfo_Feature(pb.SubvolumeInfo_FEATURE_SNAPSHOT_CLONE))
		case "snapshot-autoprotect":
			ret = append(ret, pb.SubvolumeInfo_Feature(pb.SubvolumeInfo_FEATURE_SNAPSHOT_AUTOPROTECT))
		case "snapshot-retention":
			ret = append(ret, pb.SubvolumeInfo_Feature(pb.SubvolumeInfo_FEATURE_SNAPSHOT_RETENTION))
		default:
			ret = append(ret, pb.SubvolumeInfo_Feature(pb.SubvolumeInfo_FEATURE_UNSPECIFIED))
		}
	}
	return ret
}

/*func toProtoSubvolumeSnapshots(ss []file.SubvolumeSnapshot) []*pb.Subvolume_Snapshot {
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
}*/

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
		ret = append(ret, toProtoUserKey(&uks[i], true))
	}

	return ret
}

// ignore secret key in proto for security reason
func toProtoUserKey(uk *object.UserKey, ignore bool) *pb.User_Key {
	ret := &pb.User_Key{}
	ret.SetAccessKey(uk.AccessKey)

	if !ignore {
		ret.SetSecretKey(uk.SecretKey)
	}

	return ret
}

func toProtoSMBShareCommonConfigMapToGuest(mtg string) pb.SMBShare_CommonConfig_MapToGuest {
	switch mtg {
	case smb.MapToGuestNever.String():
		return pb.SMBShare_CommonConfig_MAP_TO_GUEST_NEVER

	case smb.MapToGuestBadUser.String():
		return pb.SMBShare_CommonConfig_MAP_TO_GUEST_BAD_USER

	case smb.MapToGuestBadPassword.String():
		return pb.SMBShare_CommonConfig_MAP_TO_GUEST_BAD_PASSWORD

	default:
		return pb.SMBShare_CommonConfig_MAP_TO_GUEST_NEVER
	}
}

func toProtoSMBShareSecurityConfigMode(sm string) pb.SMBShare_SecurityConfig_Mode {
	switch sm {
	case smb.SecurityModeUser.String():
		return pb.SMBShare_SecurityConfig_MODE_USER

	case smb.SecurityModeActiveDirectory.String():
		return pb.SMBShare_SecurityConfig_MODE_ACTIVE_DIRECTORY

	default:
		return pb.SMBShare_SecurityConfig_MODE_USER
	}
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

	if joinSource != nil {
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

func toProtoEntityType(et smb.EntityType) pb.ValidateSMBUserResponse_EntityType {
	switch et {
	case smb.EntityTypeUser:
		return pb.ValidateSMBUserResponse_ENTITY_TYPE_USER

	case smb.EntityTypeGroup:
		return pb.ValidateSMBUserResponse_ENTITY_TYPE_GROUP

	default:
		return pb.ValidateSMBUserResponse_ENTITY_TYPE_UNKNOWN
	}
}

func toCephFsTarget(t *pb.CephFsTarget) (filesystem, group, subvolume string) {
	if t == nil {
		return "", "", ""
	}
	return t.GetFileSystemName(), t.GetGroupName(), t.GetSubvolumeName()
}

func toNFSExportSpecFromApply(req *pb.ApplyNFSExportRequest) file.NFSExportSpec {
	specPB := req.GetSpec()

	var spec file.NFSExportSpec

	if specPB == nil {
		return spec
	}

	spec.PseudoPath = specPB.GetPseudoPath()

	if specPB.HasEnableSecurityLabel() {
		spec.EnableSecurityLabel = specPB.GetEnableSecurityLabel()
	}

	if specPB.HasAccessType() {
		spec.AccessType = toAccessType(specPB.GetAccessType())
	}
	if specPB.HasSquash() {
		spec.Squash = toSquash(specPB.GetSquash())
	}

	ps := make([]file.NFSProtocol, 0, len(specPB.GetProtocols()))
	for _, p := range specPB.GetProtocols() {
		if pp := toProtocol(p); pp != 0 {
			ps = append(ps, pp)
		}
	}
	spec.Protocols = ps

	ts := make([]file.NFSTransport, 0, len(specPB.GetTransports()))
	for _, t := range specPB.GetTransports() {
		if tt := toTransport(t); tt != "" {
			ts = append(ts, tt)
		}
	}
	spec.Transports = ts

	cs := make([]file.Client, 0, len(specPB.GetClients()))
	for _, r := range specPB.GetClients() {
		if r == nil {
			continue
		}
		c := file.Client{Addresses: r.GetAddresses()}
		if r.HasAccessType() {
			c.AccessType = toAccessType(r.GetAccessType())
		}
		if r.HasSquash() {
			c.Squash = toSquash(r.GetSquash())
		}
		cs = append(cs, c)
	}
	spec.Clients = cs

	st := make([]file.Sectype, 0, len(specPB.GetSectypes()))
	for _, x := range specPB.GetSectypes() {
		if ss := toSecType(x); ss != "" {
			st = append(st, ss)
		}
	}
	spec.Sectypes = st

	return spec
}

func toAccessType(v pb.NfsAccessType) file.AccessType {
	switch v {
	case pb.NfsAccessType_NFS_ACCESS_TYPE_RW:
		return file.AccessTypeRW
	case pb.NfsAccessType_NFS_ACCESS_TYPE_RO:
		return file.AccessTypeRO
	case pb.NfsAccessType_NFS_ACCESS_TYPE_NONE:
		return file.AccessTypeNone
	default:
		return file.AccessTypeUnspecified
	}
}

func toSquash(v pb.NfsSquash) file.Squash {
	switch v {
	case pb.NfsSquash_NFS_SQUASH_NONE:
		return file.SquashNone
	case pb.NfsSquash_NFS_SQUASH_ROOT:
		return file.SquashRoot
	case pb.NfsSquash_NFS_SQUASH_ALL:
		return file.SquashAll
	case pb.NfsSquash_NFS_SQUASH_ROOTID:
		return file.SquashRootID
	default:
		return file.SquashUnspecified
	}
}

func toProtocol(v pb.NfsProtocol) file.NFSProtocol {
	switch v {
	case pb.NfsProtocol_NFS_PROTOCOL_V3:
		return file.NFSProtocolV3
	case pb.NfsProtocol_NFS_PROTOCOL_V4:
		return file.NFSProtocolV4
	default:
		return 0
	}
}

func toTransport(v pb.NfsTransport) file.NFSTransport {
	switch v {
	case pb.NfsTransport_NFS_TRANSPORT_TCP:
		return file.NFSTransportTCP
	case pb.NfsTransport_NFS_TRANSPORT_UDP:
		return file.NFSTransportUDP
	default:
		return ""
	}
}

func toSecType(v pb.NfsSecType) file.Sectype {
	switch v {
	case pb.NfsSecType_NFS_SECTYPE_SYS:
		return file.SysSec
	case pb.NfsSecType_NFS_SECTYPE_NONE:
		return file.NoneSec
	case pb.NfsSecType_NFS_SECTYPE_KRB5:
		return file.Krb5Sec
	case pb.NfsSecType_NFS_SECTYPE_KRB5I:
		return file.Krb5iSec
	case pb.NfsSecType_NFS_SECTYPE_KRB5P:
		return file.Krb5pSec
	default:
		return ""
	}
}

func toProtoNfsExports(es []file.NFSExport) []*pb.NfsExport {
	out := make([]*pb.NfsExport, 0, len(es))
	for i := range es {
		e := es[i]
		out = append(out, toProtoNfsExport(&e))
	}
	return out
}

func toProtoNfsExport(e *file.NFSExport) *pb.NfsExport {
	if e == nil {
		return nil
	}

	out := &pb.NfsExport{}
	out.SetClusterId(e.ClusterID)

	cephfs := &pb.CephFsTarget{}
	cephfs.SetFileSystemName(e.FileSystemName)
	cephfs.SetGroupName(e.GroupName)
	cephfs.SetSubvolumeName(e.SubvolumeName)
	out.SetCephfs(cephfs)

	info := &pb.NfsExportInfo{}
	info.SetEnableSecurityLabel(e.EnableSecurityLabel)
	info.SetPath(e.Path)
	info.SetProtocols(toProtoNfsProtocols(e.Protocols))
	info.SetPseudoPath(e.PseudoPath)
	info.SetAccessType(toProtoNfsAccessType(e.AccessType))
	info.SetSquash(toProtoNfsSquash(e.Squash))
	info.SetTransports(toProtoNfsTransports(e.Transports))
	info.SetClients(toProtoNfsClientRules(e.Clients))
	info.SetSectypes(toProtoNfsSecTypes(e.Sectypes))
	out.SetInfo(info)

	return out
}

func toProtoNfsAccessType(v file.AccessType) pb.NfsAccessType {
	switch v {
	case file.AccessTypeRW:
		return pb.NfsAccessType_NFS_ACCESS_TYPE_RW
	case file.AccessTypeRO:
		return pb.NfsAccessType_NFS_ACCESS_TYPE_RO
	case file.AccessTypeNone:
		return pb.NfsAccessType_NFS_ACCESS_TYPE_NONE
	default:
		return pb.NfsAccessType_NFS_ACCESS_TYPE_UNSPECIFIED
	}
}

func toProtoNfsSquash(v file.Squash) pb.NfsSquash {
	switch v {
	case file.SquashNone:
		return pb.NfsSquash_NFS_SQUASH_NONE
	case file.SquashRoot:
		return pb.NfsSquash_NFS_SQUASH_ROOT
	case file.SquashAll:
		return pb.NfsSquash_NFS_SQUASH_ALL
	case file.SquashRootID:
		return pb.NfsSquash_NFS_SQUASH_ROOTID
	default:
		return pb.NfsSquash_NFS_SQUASH_UNSPECIFIED
	}
}

func toProtoNfsProtocols(v []file.NFSProtocol) []pb.NfsProtocol {
	out := make([]pb.NfsProtocol, 0, len(v))
	for _, x := range v {
		p := toProtoNfsProtocol(x)
		if p != pb.NfsProtocol_NFS_PROTOCOL_UNSPECIFIED {
			out = append(out, p)
		}
	}
	return out
}

func toProtoNfsTransports(v []file.NFSTransport) []pb.NfsTransport {
	out := make([]pb.NfsTransport, 0, len(v))
	for _, x := range v {
		t := toProtoNfsTransport(x)
		if t != pb.NfsTransport_NFS_TRANSPORT_UNSPECIFIED {
			out = append(out, t)
		}
	}
	return out
}

func toProtoNfsSecTypes(v []file.Sectype) []pb.NfsSecType {
	out := make([]pb.NfsSecType, 0, len(v))
	for _, x := range v {
		s := toProtoNfsSecType(x)
		if s != pb.NfsSecType_NFS_SECTYPE_UNSPECIFIED {
			out = append(out, s)
		}
	}
	return out
}

func toProtoNfsProtocol(v file.NFSProtocol) pb.NfsProtocol {
	switch v {
	case file.NFSProtocolV3:
		return pb.NfsProtocol_NFS_PROTOCOL_V3
	case file.NFSProtocolV4:
		return pb.NfsProtocol_NFS_PROTOCOL_V4
	default:
		return pb.NfsProtocol_NFS_PROTOCOL_UNSPECIFIED
	}
}

func toProtoNfsTransport(v file.NFSTransport) pb.NfsTransport {
	switch v {
	case file.NFSTransportTCP:
		return pb.NfsTransport_NFS_TRANSPORT_TCP
	case file.NFSTransportUDP:
		return pb.NfsTransport_NFS_TRANSPORT_UDP
	default:
		return pb.NfsTransport_NFS_TRANSPORT_UNSPECIFIED
	}
}

func toProtoNfsSecType(v file.Sectype) pb.NfsSecType {
	switch v {
	case file.SysSec:
		return pb.NfsSecType_NFS_SECTYPE_SYS
	case file.NoneSec:
		return pb.NfsSecType_NFS_SECTYPE_NONE
	case file.Krb5Sec:
		return pb.NfsSecType_NFS_SECTYPE_KRB5
	case file.Krb5iSec:
		return pb.NfsSecType_NFS_SECTYPE_KRB5I
	case file.Krb5pSec:
		return pb.NfsSecType_NFS_SECTYPE_KRB5P
	default:
		return pb.NfsSecType_NFS_SECTYPE_UNSPECIFIED
	}
}

func toProtoNfsClientRules(v []file.Client) []*pb.NfsClientRule {
	out := make([]*pb.NfsClientRule, 0, len(v))
	for i := range v {
		r := v[i]
		x := &pb.NfsClientRule{}
		x.SetAddresses(r.Addresses)
		x.SetAccessType(toProtoNfsAccessType(r.AccessType))
		x.SetSquash(toProtoNfsSquash(r.Squash))
		out = append(out, x)
	}
	return out
}
