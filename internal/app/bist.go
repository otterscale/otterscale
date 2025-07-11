package app

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/bist/v1"
	"github.com/openhdc/otterscale/api/bist/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type BISTService struct {
	pbconnect.UnimplementedBISTServiceHandler

	uc  *core.BISTUseCase
	suc *core.StorageUseCase
}

func NewBISTService(uc *core.BISTUseCase, suc *core.StorageUseCase) *BISTService {
	return &BISTService{uc: uc, suc: suc}
}

var _ pbconnect.BISTServiceHandler = (*BISTService)(nil)

func (s *BISTService) ListTestResults(ctx context.Context, req *connect.Request[pb.ListTestResultsRequest]) (*connect.Response[pb.ListTestResultsResponse], error) {
	results, err := s.uc.ListResults(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListTestResultsResponse{}
	resp.SetTestResults(toProtoTestResults(results))
	return connect.NewResponse(resp), nil
}

func (s *BISTService) CreateTestResult(ctx context.Context, req *connect.Request[pb.CreateTestResultRequest]) (*connect.Response[pb.TestResult], error) {
	target := req.Msg.GetType()
	if target == pb.TestResult_UNSPECIFIED {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid test result type"))
	}

	if target == pb.TestResult_BLOCK {
		// [TODO] if exist
		if !s.suc.IsPoolExist(ctx, req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), "otterscale_pool") {
			_, err := s.suc.CreatePool(ctx, req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), "otterscale_pool", "replicated", false, 1, 0, 0, []string{"rbd"})
			if err != nil {
				return nil, err
			}
			_, err = s.suc.CreateImage(ctx, req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), "otterscale_pool", "otterscale_image", 4194304, 4194304, 1, 1073741824, true, true, true, true, true)
			if err != nil {
				return nil, err
			}
		}
	}

	result, err := s.uc.CreateResult(ctx, strings.ToLower(target.String()), req.Msg.GetName(), req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), toCoreFIO(req.Msg.GetFio()), toCoreWarp(req.Msg.GetWarp()))
	if err != nil {
		return nil, err
	}
	resp := toProtoTestResult(result)
	return connect.NewResponse(resp), nil
}

func (s *BISTService) DeleteTestResult(ctx context.Context, req *connect.Request[pb.DeleteTestResultRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteResult(ctx, req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *BISTService) ListS3S(ctx context.Context, req *connect.Request[pb.ListS3SRequest]) (*connect.Response[pb.ListS3SResponse], error) {
	s3s, err := s.uc.ListS3s(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListS3SResponse{}
	resp.SetS3S(toProtoS3s(s3s))
	return connect.NewResponse(resp), nil
}

func toCoreFIO(p *pb.TestResult_FIO) *core.BISTFIO {
	if p == nil {
		return nil
	}
	return &core.BISTFIO{
		AccessMode:  strings.ToLower(strings.ReplaceAll(p.GetAccessMode().String(), "_", "")),
		NFSEndpoint: p.GetNfsEndpoint(),
		NFSPath:     p.GetNfsPath(),
		JobCount:    p.GetJobCount(),
		RunTime:     p.GetRunTime(),
		BlockSize:   p.GetBlockSize(),
		FileSize:    p.GetFileSize(),
		IODepth:     p.GetIoDepth(),
	}
}

func toCoreWarp(p *pb.TestResult_Warp) *core.BISTWarp {
	if p == nil {
		return nil
	}
	return &core.BISTWarp{
		Operation:  strings.ToLower(p.GetOperation().String()),
		Endpoint:   p.GetEndpoint(),
		AccessKey:  p.GetAccessKey(),
		SecretKey:  p.GetSecretKey(),
		Duration:   p.GetDuration(),
		ObjectSize: p.GetObjectSize(),
		ObjectNum:  p.GetObjectNum(),
	}
}

func toProtoTestResults(rs []core.BISTResult) []*pb.TestResult {
	ret := []*pb.TestResult{}
	for i := range rs {
		ret = append(ret, toProtoTestResult(&rs[i]))
	}
	return ret
}

func toProtoTestResult(r *core.BISTResult) *pb.TestResult {
	target := toProtoTestResultType(r.Type)
	ret := &pb.TestResult{}
	ret.SetType(target)
	ret.SetName(r.Name)
	ret.SetStatus(r.Status)
	if r.StartTime != nil {
		ret.SetStartTime(r.StartTime.String())
	}
	if r.CompleteTime != nil {
		ret.SetCompleteTime(r.CompleteTime.String())
	}
	if target == pb.TestResult_S3 {
		ret.SetWarp(toProtoTestResultWarp(r.Warp))
		if r.CompleteTime != nil {
			ret.SetWarpResult(toProtoWarpResult(r.WarpResult))
		}
	} else {
		ret.SetFio(toProtoTestResultFIO(r.FIO))
		if r.CompleteTime != nil {
			ret.SetFioResult(toProtoFIOResult(r.FIOResult))
		}
	}
	return ret
}

func toProtoTestResultFIO(f *core.BISTFIO) *pb.TestResult_FIO {
	ret := &pb.TestResult_FIO{}
	ret.SetAccessMode(toProtoTestResultFIOAccessMode(f.AccessMode))
	ret.SetScopeUuid(f.ScopeUUID)
	ret.SetFacilityName(f.FacilityName)
	ret.SetNfsEndpoint(f.NFSEndpoint)
	ret.SetNfsPath(f.NFSPath)
	ret.SetJobCount(f.JobCount)
	ret.SetRunTime(f.RunTime)
	ret.SetBlockSize(f.BlockSize)
	ret.SetFileSize(f.FileSize)
	ret.SetIoDepth(f.IODepth)
	return ret
}

func toProtoTestResultWarp(w *core.BISTWarp) *pb.TestResult_Warp {
	ret := &pb.TestResult_Warp{}
	ret.SetOperation(toProtoTestResultWarpOperation(w.Operation))
	ret.SetEndpoint(w.Endpoint)
	ret.SetAccessKey(w.AccessKey)
	ret.SetSecretKey(w.SecretKey)
	ret.SetDuration(w.Duration)
	ret.SetObjectSize(w.ObjectSize)
	ret.SetObjectNum(w.ObjectNum)
	return ret
}

func toProtoWarpResult(wr *core.WarpResult) *pb.TestResult_WarpResult {
	ret := &pb.TestResult_WarpResult{}
	ret.SetWarpOps(toProtoWarpResultOps(wr.WarpOperations))
	return ret
}

func toProtoWarpResultOps(wos []core.WarpOperation) []*pb.TestResult_WarpResult_WarpOps {
	ret := []*pb.TestResult_WarpResult_WarpOps{}
	for i := range wos {
		ret = append(ret, toProtoWarpResultOp(&wos[i]))
	}
	return ret
}

func toProtoWarpResultOp(wo *core.WarpOperation) *pb.TestResult_WarpResult_WarpOps {
	ret := &pb.TestResult_WarpResult_WarpOps{}
	ret.SetType(wo.Type)
	ret.SetOps(float64(wo.Throughput.Ops))
	ret.SetBytes(float64(wo.Throughput.Bytes))
	ret.SetStartTime(wo.Throughput.StartTime.String())
	ret.SetEndTime(wo.Throughput.EndTime.String())
	ret.SetFastestBps(float64(wo.Throughput.Segmented.FastestBps))
	ret.SetFastestOps(float64(wo.Throughput.Segmented.FastestOps))
	ret.SetMedianBps(float64(wo.Throughput.Segmented.MedianBps))
	ret.SetMedianOps(float64(wo.Throughput.Segmented.MedianOps))
	ret.SetSlowestBps(float64(wo.Throughput.Segmented.SlowestBps))
	ret.SetSlowestOps(float64(wo.Throughput.Segmented.SlowestOps))
	return ret
}

func toProtoFIOResult(fr *core.FIOResult) *pb.TestResult_FIOResult {
	ret := &pb.TestResult_FIOResult{}

	ret.SetRead(toProtoFIOReadResult(fr))
	ret.SetWrite(toProtoFIOWriteResult(fr))
	ret.SetTrim(toProtoFIOTrimResult(fr))

	return ret
}

func toProtoFIOReadResult(fr *core.FIOResult) *pb.TestResult_FIOResult_FIOValues {
	ret := &pb.TestResult_FIOResult_FIOValues{}
	ret.SetBw(uint64(fr.Jobs[0].Read.Bw))
	ret.SetIoBytes(uint64(fr.Jobs[0].Read.IoBytes))
	ret.SetIops(uint64(fr.Jobs[0].Read.Iops))
	ret.SetRuntime(uint64(fr.Jobs[0].Read.Runtime))
	ret.SetTotalIos(uint64(fr.Jobs[0].Read.TotalIos))
	ret.SetLatency(toProtoFIOLatency(&fr.Jobs[0].Read.LatNs))

	return ret
}

func toProtoFIOWriteResult(fr *core.FIOResult) *pb.TestResult_FIOResult_FIOValues {
	ret := &pb.TestResult_FIOResult_FIOValues{}
	ret.SetBw(uint64(fr.Jobs[0].Write.Bw))
	ret.SetIoBytes(uint64(fr.Jobs[0].Write.IoBytes))
	ret.SetIops(uint64(fr.Jobs[0].Write.Iops))
	ret.SetRuntime(uint64(fr.Jobs[0].Write.Runtime))
	ret.SetTotalIos(uint64(fr.Jobs[0].Write.TotalIos))
	ret.SetLatency(toProtoFIOLatency(&fr.Jobs[0].Write.LatNs))
	return ret
}

func toProtoFIOTrimResult(fr *core.FIOResult) *pb.TestResult_FIOResult_FIOValues {
	ret := &pb.TestResult_FIOResult_FIOValues{}
	ret.SetBw(uint64(fr.Jobs[0].Trim.Bw))
	ret.SetIoBytes(uint64(fr.Jobs[0].Trim.IoBytes))
	ret.SetIops(uint64(fr.Jobs[0].Trim.Iops))
	ret.SetRuntime(uint64(fr.Jobs[0].Trim.Runtime))
	ret.SetTotalIos(uint64(fr.Jobs[0].Trim.TotalIos))
	ret.SetLatency(toProtoFIOLatency(&fr.Jobs[0].Trim.LatNs))
	return ret
}

func toProtoFIOLatency(l *core.LatNs) *pb.TestResult_FIOResult_FIOValues_Latency {
	ret := &pb.TestResult_FIOResult_FIOValues_Latency{}
	ret.SetMin(uint64(l.Min))
	ret.SetMax(uint64(l.Max))
	ret.SetMean(float64(l.Mean))
	return ret
}

func toProtoS3s(ss []core.BISTS3) []*pb.S3 {
	ret := []*pb.S3{}
	for i := range ss {
		ret = append(ret, toProtoS3(&ss[i]))
	}
	return ret
}

func toProtoS3(s *core.BISTS3) *pb.S3 {
	ret := &pb.S3{}
	ret.SetType(toProtoS3Type(s.Type))
	ret.SetName(s.Name)
	ret.SetEndpoint(s.Endpoint)
	return ret
}

func toProtoTestResultType(s string) pb.TestResult_Type {
	v, ok := pb.TestResult_Type_value[strings.ToUpper(s)]
	if ok {
		return pb.TestResult_Type(v)
	}
	return pb.TestResult_UNSPECIFIED
}

func toProtoTestResultFIOAccessMode(s string) pb.TestResult_FIO_AccessMode {
	switch s {
	case "read":
		return pb.TestResult_FIO_READ
	case "write":
		return pb.TestResult_FIO_WRITE
	case "trim":
		return pb.TestResult_FIO_TRIM
	case "readwrite":
		return pb.TestResult_FIO_READ_WRITE
	case "trimwrite":
		return pb.TestResult_FIO_TRIM_WRITE
	case "randread":
		return pb.TestResult_FIO_RAND_READ
	case "randwrite":
		return pb.TestResult_FIO_RAND_WRITE
	case "randtrim":
		return pb.TestResult_FIO_RAND_TRIM
	case "randrw":
		return pb.TestResult_FIO_RAND_RW
	case "randtrimwrite":
		return pb.TestResult_FIO_RAND_TRIM_WRITE
	}
	return pb.TestResult_FIO_READ
}

func toProtoTestResultWarpOperation(s string) pb.TestResult_Warp_Operation {
	switch s {
	case "get":
		return pb.TestResult_Warp_GET
	case "put":
		return pb.TestResult_Warp_PUT
	case "delete":
		return pb.TestResult_Warp_DELETE
	case "list":
		return pb.TestResult_Warp_LIST
	case "stat":
		return pb.TestResult_Warp_STAT
	case "mixed":
		return pb.TestResult_Warp_MIXED
	}
	return pb.TestResult_Warp_GET
}

func toProtoS3Type(s string) pb.S3_Type {
	switch s {
	case "ceph":
		return pb.S3_CEPH_RADOS_GATEWAY
	case "minio":
		return pb.S3_MINIO
	}
	return pb.S3_UNSPECIFIED
}
