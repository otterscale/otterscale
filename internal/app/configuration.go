package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/configuration/v1"
	"github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/configuration"
	"github.com/otterscale/otterscale/internal/core/configuration/bist"
)

type ConfigurationService struct {
	pbconnect.UnimplementedConfigurationServiceHandler

	configuration *configuration.UseCase
	bist          *bist.UseCase
}

func NewConfigurationService(configuration *configuration.UseCase, bist *bist.UseCase) *ConfigurationService {
	return &ConfigurationService{
		configuration: configuration,
		bist:          bist,
	}
}

var _ pbconnect.ConfigurationServiceHandler = (*ConfigurationService)(nil)

func (s *ConfigurationService) GetConfiguration(ctx context.Context, _ *pb.GetConfigurationRequest) (*pb.Configuration, error) {
	config, err := s.configuration.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	resp := toProtoConfiguration(config)
	return resp, nil
}

func (s *ConfigurationService) UpdateNTPServer(ctx context.Context, req *pb.UpdateNTPServerRequest) (*pb.Configuration_NTPServer, error) {
	ntpServers, err := s.configuration.UpdateNTPServer(ctx, req.GetAddresses())
	if err != nil {
		return nil, err
	}

	resp := toProtoNTPServer(ntpServers)
	return resp, nil
}

func (s *ConfigurationService) UpdatePackageRepository(ctx context.Context, req *pb.UpdatePackageRepositoryRequest) (*pb.Configuration_PackageRepository, error) {
	repo, err := s.configuration.UpdatePackageRepository(ctx, int(req.GetId()), req.GetUrl())
	if err != nil {
		return nil, err
	}

	resp := toProtoPackageRepository(repo)
	return resp, nil
}

func (s *ConfigurationService) UpdateHelmRepository(_ context.Context, req *pb.UpdateHelmRepositoryRequest) (*pb.Configuration_HelmRepository, error) {
	helmRepository, err := s.configuration.UpdateHelmRepository(req.GetUrls())
	if err != nil {
		return nil, err
	}

	resp := toProtoHelmRepository(helmRepository)
	return resp, nil
}

func (s *ConfigurationService) CreateBootImage(ctx context.Context, req *pb.CreateBootImageRequest) (*pb.Configuration_BootImage, error) {
	image, err := s.configuration.CreateBootImage(ctx, req.GetDistroSeries(), req.GetArchitectures())
	if err != nil {
		return nil, err
	}

	resp := toProtoBootImage(image)
	return resp, nil
}

func (s *ConfigurationService) UpdateBootImage(ctx context.Context, req *pb.UpdateBootImageRequest) (*pb.Configuration_BootImage, error) {
	image, err := s.configuration.UpdateBootImage(ctx, int(req.GetId()), req.GetDistroSeries(), req.GetArchitectures())
	if err != nil {
		return nil, err
	}

	resp := toProtoBootImage(image)
	return resp, nil
}

func (s *ConfigurationService) SetDefaultBootImage(ctx context.Context, req *pb.SetDefaultBootImageRequest) (*emptypb.Empty, error) {
	if err := s.configuration.SetDefaultBootImage(ctx, req.GetDistroSeries()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ConfigurationService) ImportBootImages(ctx context.Context, _ *pb.ImportBootImagesRequest) (*emptypb.Empty, error) {
	if err := s.configuration.ImportBootImages(ctx); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ConfigurationService) IsImportingBootImages(ctx context.Context, _ *pb.IsImportingBootImagesRequest) (*pb.IsImportingBootImagesResponse, error) {
	isImporting, err := s.configuration.IsImportingBootImages(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.IsImportingBootImagesResponse{}
	resp.SetImporting(isImporting)
	return resp, nil
}

func (s *ConfigurationService) ListBootImageSelections(_ context.Context, _ *pb.ListBootImageSelectionsRequest) (*pb.ListBootImageSelectionsResponse, error) {
	selections, err := s.configuration.ListBootImageSelections()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListBootImageSelectionsResponse{}
	resp.SetBootImageSelections(toProtoBootImageSelections(selections))
	return resp, nil
}

func (s *ConfigurationService) ListTestResults(ctx context.Context, _ *pb.ListTestResultsRequest) (*pb.ListTestResultsResponse, error) {
	results, err := s.bist.ListResults(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListTestResultsResponse{}
	resp.SetTestResults(toProtoTestResults(results))
	return resp, nil
}

func (s *ConfigurationService) CreateTestResult(ctx context.Context, req *pb.CreateTestResultRequest) (*pb.TestResult, error) {
	switch req.WhichKind() {
	case pb.CreateTestResultRequest_Fio_case:
		fio := req.GetFio()

		result, err := s.bist.CreateFIOResult(ctx, req.GetName(), req.GetCreatedBy(), toCoreFIOTarget(fio.GetCephBlockDevice(), fio.GetNetworkFileSystem()), toCoreFIOInput(fio.GetInput()))
		if err != nil {
			return nil, err
		}

		resp := toProtoTestResult(result)
		return resp, nil

	case pb.CreateTestResultRequest_Warp_case:
		warp := req.GetWarp()

		result, err := s.bist.CreateWarpResult(ctx, req.GetName(), req.GetCreatedBy(), toCoreWarpTarget(warp.GetInternalObjectService(), warp.GetExternalObjectService()), toCoreWarpInput(warp.GetInput()))
		if err != nil {
			return nil, err
		}
		resp := toProtoTestResult(result)
		return resp, nil
	}

	return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("kind is empty"))
}

func (s *ConfigurationService) DeleteTestResult(ctx context.Context, req *pb.DeleteTestResultRequest) (*emptypb.Empty, error) {
	if err := s.bist.DeleteResult(ctx, req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ConfigurationService) ListInternalObjectServices(ctx context.Context, req *pb.ListInternalObjectServicesRequest) (*pb.ListInternalObjectServicesResponse, error) {
	services, err := s.bist.ListInternalObjectServices(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListInternalObjectServicesResponse{}
	resp.SetInternalObjectServices(toProtoInternalObjectServices(services))
	return resp, nil
}

func toCoreFIOTarget(c *pb.CephBlockDevice, n *pb.NetworkFileSystem) bist.FIOTarget {
	ret := bist.FIOTarget{}

	if c != nil {
		ret.Ceph = &bist.FIOTargetCeph{
			Scope: c.GetScope(),
		}
	}

	if n != nil {
		ret.NFS = &bist.FIOTargetNFS{
			Host: n.GetHost(),
			Path: n.GetPath(),
		}
	}

	return ret
}

func toCoreWarpTarget(i *pb.InternalObjectService, e *pb.ExternalObjectService) bist.WarpTarget {
	ret := bist.WarpTarget{}

	if i != nil {
		ret.Internal = &bist.WarpTargetInternal{
			Type:  strings.ToLower(i.GetType().String()),
			Scope: i.GetScope(),
			Host:  i.GetHost(),
		}
	}

	if e != nil {
		ret.External = &bist.WarpTargetExternal{
			Host:      e.GetHost(),
			AccessKey: e.GetAccessKey(),
			SecretKey: e.GetSecretKey(),
		}
	}

	return ret
}

func toCoreFIOInput(p *pb.FIO_Input) *bist.FIOInput {
	return &bist.FIOInput{
		AccessMode: strings.ToLower(strings.ReplaceAll(p.GetAccessMode().String(), "_", "")),
		JobCount:   p.GetJobCount(),
		RunTime:    p.GetRunTimeSeconds(),
		BlockSize:  p.GetBlockSizeBytes(),
		FileSize:   p.GetFileSizeBytes(),
		IODepth:    p.GetIoDepth(),
	}
}

func toCoreWarpInput(p *pb.Warp_Input) *bist.WarpInput {
	return &bist.WarpInput{
		Operation:   strings.ToLower(p.GetOperation().String()),
		Duration:    p.GetDurationSeconds(),
		ObjectSize:  p.GetObjectSizeBytes(),
		ObjectCount: p.GetObjectCount(),
	}
}

func toProtoNTPServer(addresses []string) *pb.Configuration_NTPServer {
	ret := &pb.Configuration_NTPServer{}
	ret.SetAddresses(addresses)
	return ret
}

func toProtoPackageRepositories(prs []configuration.PackageRepository) []*pb.Configuration_PackageRepository {
	ret := []*pb.Configuration_PackageRepository{}

	for i := range prs {
		ret = append(ret, toProtoPackageRepository(&prs[i]))
	}

	return ret
}

func toProtoPackageRepository(pr *configuration.PackageRepository) *pb.Configuration_PackageRepository {
	ret := &pb.Configuration_PackageRepository{}
	ret.SetId(int64(pr.ID))
	ret.SetName(pr.Name)
	ret.SetUrl(pr.URL)
	ret.SetEnabled(pr.Enabled)
	return ret
}

func toProtoBootImages(bis []configuration.BootImage) []*pb.Configuration_BootImage {
	ret := []*pb.Configuration_BootImage{}

	for i := range bis {
		ret = append(ret, toProtoBootImage(&bis[i]))
	}

	return ret
}

func toProtoBootImage(bi *configuration.BootImage) *pb.Configuration_BootImage {
	ret := &pb.Configuration_BootImage{}
	ret.SetSource(bi.Source)
	ret.SetDistroSeries(bi.DistroSeries)
	ret.SetName(bi.Name)
	ret.SetId(int64(bi.ID))
	ret.SetArchitectures(bi.Architectures)
	ret.SetArchitectureStatusMap(bi.ArchitectureStatusMap)
	ret.SetDefault(bi.Default)
	return ret
}

func toProtoBootImageSelections(biss []configuration.BootImageSelection) []*pb.Configuration_BootImageSelection {
	ret := []*pb.Configuration_BootImageSelection{}

	for i := range biss {
		ret = append(ret, toProtoBootImageSelection(&biss[i]))
	}

	return ret
}

func toProtoBootImageSelection(bis *configuration.BootImageSelection) *pb.Configuration_BootImageSelection {
	ret := &pb.Configuration_BootImageSelection{}
	ret.SetDistroSeries(bis.DistroSeries)
	ret.SetName(bis.DisplayName)
	ret.SetArchitectures(bis.Architectures)
	return ret
}

func toProtoHelmRepository(urls []string) *pb.Configuration_HelmRepository {
	ret := &pb.Configuration_HelmRepository{}
	ret.SetUrls(urls)
	return ret
}

func toProtoConfiguration(c *configuration.Configuration) *pb.Configuration {
	ret := &pb.Configuration{}
	ret.SetNtpServer(toProtoNTPServer(c.NTPServers))
	ret.SetPackageRepositories(toProtoPackageRepositories(c.PackageRepositories))
	ret.SetBootImages(toProtoBootImages(c.BootImages))
	ret.SetHelmRepository(toProtoHelmRepository(c.HelmRepositorys))
	return ret
}

func toProtoTestResults(rs []bist.Result) []*pb.TestResult {
	ret := []*pb.TestResult{}

	for i := range rs {
		ret = append(ret, toProtoTestResult(&rs[i]))
	}

	return ret
}

func toProtoTestResult(r *bist.Result) *pb.TestResult {
	ret := &pb.TestResult{}
	ret.SetUid(r.UID)
	ret.SetName(r.Name)
	ret.SetStatus(toProtoTestResultStatus(r.Status))
	ret.SetCreatedBy(r.CreatedBy)
	ret.SetStartedAt(timestamppb.New(r.StartTime))
	ret.SetCompletedAt(timestamppb.New(r.CompletionTime))

	if r.FIO != nil {
		ret.SetFio(toProtoFIO(r.FIO))
	}

	if r.Warp != nil {
		ret.SetWarp(toProtoWarp(r.Warp))
	}

	return ret
}

func toProtoFIO(f *bist.FIO) *pb.FIO {
	ret := &pb.FIO{}

	if f.Target.Ceph != nil {
		ret.SetCephBlockDevice(toProtoCephBlockDevice(f.Target.Ceph))
	}

	if f.Target.NFS != nil {
		ret.SetNetworkFileSystem(toProtoNetworkFileSystem(f.Target.NFS))
	}

	if f.Input != nil {
		ret.SetInput(toProtoFIOInput(f.Input))
	}

	if f.Output != nil {
		ret.SetOutput(toProtoFIOOutput(f.Output))
	}

	return ret
}

func toProtoCephBlockDevice(f *bist.FIOTargetCeph) *pb.CephBlockDevice {
	ret := &pb.CephBlockDevice{}
	ret.SetScope(f.Scope)
	return ret
}

func toProtoNetworkFileSystem(f *bist.FIOTargetNFS) *pb.NetworkFileSystem {
	ret := &pb.NetworkFileSystem{}
	ret.SetHost(f.Host)
	ret.SetPath(f.Path)
	return ret
}

func toProtoFIOInput(f *bist.FIOInput) *pb.FIO_Input {
	ret := &pb.FIO_Input{}
	ret.SetAccessMode(toProtoFIOInputAccessMode(f.AccessMode))
	ret.SetJobCount(f.JobCount)
	ret.SetRunTimeSeconds(f.RunTime)
	ret.SetBlockSizeBytes(f.BlockSize)
	ret.SetFileSizeBytes(f.FileSize)
	ret.SetIoDepth(f.IODepth)
	return ret
}

func toProtoFIOOutput(f *bist.FIOOutput) *pb.FIO_Output {
	ret := &pb.FIO_Output{}

	if f.Read != nil {
		ret.SetRead(toProtoFIOOutputThroughput(f.Read))
	}

	if f.Write != nil {
		ret.SetWrite(toProtoFIOOutputThroughput(f.Write))
	}

	if f.Trim != nil {
		ret.SetTrim(toProtoFIOOutputThroughput(f.Trim))
	}

	return ret
}

func toProtoFIOOutputThroughput(f *bist.FIOThroughput) *pb.FIO_Output_Throughput {
	latency := &pb.FIO_Output_Throughput_Latency{}
	latency.SetMinNanoseconds(f.Latency.Min)
	latency.SetMaxNanoseconds(f.Latency.Max)
	latency.SetMeanNanoseconds(f.Latency.Mean)

	ret := &pb.FIO_Output_Throughput{}
	ret.SetIoBytes(f.IOBytes)
	ret.SetBandwidthBytes(f.Bandwidth)
	ret.SetIoPerSecond(f.IOPS)
	ret.SetTotalIos(f.TotalIOs)
	ret.SetLatency(latency)
	return ret
}

func toProtoWarp(f *bist.Warp) *pb.Warp {
	ret := &pb.Warp{}

	if f.Target.Internal != nil {
		ret.SetInternalObjectService(toProtoInternalObjectService(f.Target.Internal))
	}

	if f.Target.External != nil {
		ret.SetExternalObjectService(toProtoExternalObjectService(f.Target.External))
	}

	if f.Input != nil {
		ret.SetInput(toProtoWarpInput(f.Input))
	}

	if f.Output != nil {
		ret.SetOutput(toProtoWarpOutput(f.Output))
	}

	return ret
}

func toProtoInternalObjectServices(ss []bist.WarpTargetInternal) []*pb.InternalObjectService {
	ret := []*pb.InternalObjectService{}

	for i := range ss {
		ret = append(ret, toProtoInternalObjectService(&ss[i]))
	}

	return ret
}

func toProtoInternalObjectService(s *bist.WarpTargetInternal) *pb.InternalObjectService {
	ret := &pb.InternalObjectService{}
	ret.SetType(toProtoInternalObjectServiceType(s.Type))
	ret.SetScope(s.Scope)
	ret.SetHost(s.Host)
	return ret
}

func toProtoExternalObjectService(f *bist.WarpTargetExternal) *pb.ExternalObjectService {
	ret := &pb.ExternalObjectService{}
	ret.SetHost(f.Host)
	ret.SetAccessKey(f.AccessKey)
	ret.SetSecretKey(f.SecretKey)
	return ret
}

func toProtoWarpInput(w *bist.WarpInput) *pb.Warp_Input {
	ret := &pb.Warp_Input{}
	ret.SetOperation(toProtoWarpInputOperation(w.Operation))
	ret.SetDurationSeconds(w.Duration)
	ret.SetObjectSizeBytes(w.ObjectSize)
	ret.SetObjectCount(w.ObjectCount)
	return ret
}

func toProtoWarpOutput(f *bist.WarpOutput) *pb.Warp_Output {
	ret := &pb.Warp_Output{}

	for i := range f.Operations {
		if f.Operations[i].Type == http.MethodGet {
			ret.SetGet(toProtoWarpOutputThroughput(&f.Operations[i]))
		}

		if f.Operations[i].Type == http.MethodPut {
			ret.SetPut(toProtoWarpOutputThroughput(&f.Operations[i]))
		}

		if f.Operations[i].Type == http.MethodDelete {
			ret.SetDelete(toProtoWarpOutputThroughput(&f.Operations[i]))
		}
	}

	return ret
}

func toProtoWarpOutputThroughput(f *bist.WarpOperation) *pb.Warp_Output_Throughput {
	bps := &pb.Warp_Output_Throughput_Metrics{}
	bps.SetFastestPerSecond(f.Throughput.Metrics.FastestBPS)
	bps.SetMedianPerSecond(f.Throughput.Metrics.MedianBPS)
	bps.SetSlowestPerSecond(f.Throughput.Metrics.SlowestBPS)

	ops := &pb.Warp_Output_Throughput_Metrics{}
	ops.SetFastestPerSecond(f.Throughput.Metrics.FastestOPS)
	ops.SetMedianPerSecond(f.Throughput.Metrics.MedianOPS)
	ops.SetSlowestPerSecond(f.Throughput.Metrics.SlowestOPS)

	ret := &pb.Warp_Output_Throughput{}
	ret.SetTotalBytes(f.Throughput.TotalBytes)
	ret.SetTotalObjects(f.Throughput.TotalObjects)
	ret.SetTotalOperations(f.Throughput.TotalOperations)
	ret.SetBytes(bps)
	ret.SetObjects(ops)
	return ret
}

func toProtoTestResultStatus(s string) pb.TestResult_Status {
	v, ok := pb.TestResult_Status_value[strings.ToUpper(s)]
	if ok {
		return pb.TestResult_Status(v)
	}
	return pb.TestResult_RUNNING
}

func toProtoFIOInputAccessMode(s string) pb.FIO_Input_AccessMode {
	switch s {
	case "read":
		return pb.FIO_Input_READ

	case "write":
		return pb.FIO_Input_WRITE

	case "trim":
		return pb.FIO_Input_TRIM

	case "readwrite":
		return pb.FIO_Input_READ_WRITE

	case "trimwrite":
		return pb.FIO_Input_TRIM_WRITE

	case "randread":
		return pb.FIO_Input_RAND_READ

	case "randwrite":
		return pb.FIO_Input_RAND_WRITE

	case "randtrim":
		return pb.FIO_Input_RAND_TRIM

	case "randrw":
		return pb.FIO_Input_RAND_RW

	case "randtrimwrite":
		return pb.FIO_Input_RAND_TRIM_WRITE
	}

	return pb.FIO_Input_READ
}

func toProtoWarpInputOperation(s string) pb.Warp_Input_Operation {
	switch s {
	case "get":
		return pb.Warp_Input_GET

	case "put":
		return pb.Warp_Input_PUT

	case "delete":
		return pb.Warp_Input_DELETE

	case "list":
		return pb.Warp_Input_LIST

	case "stat":
		return pb.Warp_Input_STAT

	case "mixed":
		return pb.Warp_Input_MIXED
	}

	return pb.Warp_Input_GET
}

func toProtoInternalObjectServiceType(s string) pb.InternalObjectService_Type {
	v, ok := pb.InternalObjectService_Type_value[strings.ToUpper(s)]
	if ok {
		return pb.InternalObjectService_Type(v)
	}
	return pb.InternalObjectService_UNSPECIFIED
}
