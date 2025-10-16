package app

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"k8s.io/apimachinery/pkg/util/duration"

	pb "github.com/otterscale/otterscale/api/application/v1"
	"github.com/otterscale/otterscale/api/application/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type ApplicationService struct {
	pbconnect.UnimplementedApplicationServiceHandler

	chart      *core.ChartUseCase
	release    *core.ReleaseUseCase
	kubernetes *core.KubernetesUseCase
}

func NewApplicationService(chart *core.ChartUseCase, release *core.ReleaseUseCase, kubernetes *core.KubernetesUseCase) *ApplicationService {
	return &ApplicationService{
		chart:      chart,
		release:    release,
		kubernetes: kubernetes,
	}
}

var _ pbconnect.ApplicationServiceHandler = (*ApplicationService)(nil)

func (s *ApplicationService) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {
	namespaces, err := s.kubernetes.ListNamespaces(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListNamespacesResponse{}
	resp.SetNamespaces(toProtoNamespaces(namespaces))
	return resp, nil
}

func (s *ApplicationService) ListApplications(ctx context.Context, req *pb.ListApplicationsRequest) (*pb.ListApplicationsResponse, error) {
	apps, err := s.kubernetes.ListApplications(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	publicAddress, err := s.kubernetes.GetPublicAddress(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListApplicationsResponse{}
	resp.SetApplications(toProtoApplications(apps, publicAddress))
	return resp, nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *pb.GetApplicationRequest) (*pb.Application, error) {
	app, err := s.kubernetes.GetApplication(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}
	chartFile, err := s.chart.GetChartFileFromApplication(ctx, req.GetScope(), req.GetFacility(), app)
	if err != nil {
		return nil, err
	}
	app.ChartFile = chartFile
	publicAddress, err := s.kubernetes.GetPublicAddress(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := toProtoApplication(app, publicAddress)
	return resp, nil
}

func (s *ApplicationService) RestartApplication(ctx context.Context, req *pb.RestartApplicationRequest) (*emptypb.Empty, error) {
	err := s.kubernetes.RestartApplication(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetType())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ScaleApplication(ctx context.Context, req *pb.ScaleApplicationRequest) (*emptypb.Empty, error) {
	err := s.kubernetes.ScaleApplication(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetType(), req.GetReplicas())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) DeleteApplicationPod(ctx context.Context, req *pb.DeleteApplicationPodRequest) (*emptypb.Empty, error) {
	err := s.kubernetes.DeletePod(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) WatchLogs(ctx context.Context, req *pb.WatchLogsRequest, stream *connect.ServerStream[pb.WatchLogsResponse]) error {
	logs, err := s.kubernetes.StreamLogs(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetPodName(), req.GetContainerName())
	if err != nil {
		return err
	}
	defer logs.Close()

	// read logs from the stream and send to client
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		resp := &pb.WatchLogsResponse{}
		resp.SetLog(scanner.Text())
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log stream: %w", err)
	}
	return nil
}

func (s *ApplicationService) WriteTTY(_ context.Context, req *pb.WriteTTYRequest) (*emptypb.Empty, error) {
	if err := s.kubernetes.WriteToTTYSession(req.GetSessionId(), req.GetStdin()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ExecuteTTY(ctx context.Context, req *pb.ExecuteTTYRequest, stream *connect.ServerStream[pb.ExecuteTTYResponse]) error {
	// create session pipes
	sessionID, err := s.kubernetes.CreateTTYSession()
	if err != nil {
		return err
	}
	defer func() {
		_ = s.kubernetes.CleanupTTYSession(sessionID)
	}()

	// send session id to client
	resp := &pb.ExecuteTTYResponse{}
	resp.SetSessionId(sessionID)
	if err := stream.Send(resp); err != nil {
		return err
	}

	// create stdout channel
	stdOutChan := make(chan []byte)
	go func() {
		_ = s.kubernetes.ExecuteTTY(ctx, sessionID, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetPodName(), req.GetContainerName(), req.GetCommand(), stdOutChan)
	}()

	// send stdout to client
	for stdOut := range stdOutChan {
		resp := &pb.ExecuteTTYResponse{}
		resp.SetStdout(stdOut)
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *ApplicationService) ListReleases(ctx context.Context, req *pb.ListReleasesRequest) (*pb.ListReleasesResponse, error) {
	releases, err := s.release.ListReleases(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListReleasesResponse{}
	resp.SetReleases(toProtoReleases(releases))
	return resp, nil
}

func (s *ApplicationService) CreateRelease(ctx context.Context, req *pb.CreateReleaseRequest) (*pb.Application_Release, error) {
	release, err := s.release.CreateRelease(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetDryRun(), req.GetChartRef(), req.GetValuesYaml(), req.GetValuesMap())
	if err != nil {
		return nil, err
	}
	resp := toProtoRelease(release)
	return resp, nil
}

func (s *ApplicationService) UpdateRelease(ctx context.Context, req *pb.UpdateReleaseRequest) (*pb.Application_Release, error) {
	release, err := s.release.UpdateRelease(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetDryRun(), req.GetChartRef(), req.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	resp := toProtoRelease(release)
	return resp, nil
}

func (s *ApplicationService) DeleteRelease(ctx context.Context, req *pb.DeleteReleaseRequest) (*emptypb.Empty, error) {
	if err := s.release.DeleteRelease(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetDryRun()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) RollbackRelease(ctx context.Context, req *pb.RollbackReleaseRequest) (*emptypb.Empty, error) {
	if err := s.release.RollbackRelease(ctx, req.GetScope(), req.GetFacility(), req.GetNamespace(), req.GetName(), req.GetDryRun()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ListCharts(ctx context.Context, _ *pb.ListChartsRequest) (*pb.ListChartsResponse, error) {
	charts, err := s.chart.ListCharts(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListChartsResponse{}
	resp.SetCharts(toProtoCharts(charts))
	return resp, nil
}

func (s *ApplicationService) GetChart(ctx context.Context, req *pb.GetChartRequest) (*pb.Application_Chart, error) {
	ch, err := s.chart.GetChart(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	metadata := &core.ChartMetadata{}
	if len(ch.Versions) > 0 {
		metadata = ch.Versions[0].Metadata
	}
	resp := toProtoChart(metadata, ch.Versions...)
	return resp, nil
}

func (s *ApplicationService) GetChartMetadata(_ context.Context, req *pb.GetChartMetadataRequest) (*pb.Application_Chart_Metadata, error) {
	file, err := s.chart.GetChartFile(req.GetChartRef())
	if err != nil {
		return nil, err
	}
	resp := toProtoChartMetadata(file)
	return resp, nil
}

func (s *ApplicationService) ListStorageClasses(ctx context.Context, req *pb.ListStorageClassesRequest) (*pb.ListStorageClassesResponse, error) {
	storageClasses, err := s.kubernetes.ListStorageClasses(ctx, req.GetScope(), req.GetFacility())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStorageClassesResponse{}
	resp.SetStorageClasses(toProtoStorageClasses(storageClasses))
	return resp, nil
}

func (s *ApplicationService) UploadChart(ctx context.Context, req *pb.UploadChartRequest) (*emptypb.Empty, error) {
	err := s.chart.UploadChart(ctx, req.GetChartContent())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func toProtoNamespaces(ns []core.Namespace) []*pb.Namespace {
	ret := []*pb.Namespace{}
	for i := range ns {
		ret = append(ret, toProtoNamespace(&ns[i]))
	}
	return ret
}

func toProtoNamespace(n *core.Namespace) *pb.Namespace {
	ret := &pb.Namespace{}
	ret.SetName(n.Name)
	ret.SetLabels(n.Labels)
	ret.SetCreatedAt(timestamppb.New(n.CreationTimestamp.Time))
	return ret
}

func toProtoApplications(as []core.Application, publicAddress string) []*pb.Application {
	ret := []*pb.Application{}
	for i := range as {
		ret = append(ret, toProtoApplication(&as[i], publicAddress))
	}
	return ret
}

func toProtoApplication(a *core.Application, publicAddress string) *pb.Application {
	replicas := int32(0)
	if a.Replicas != nil {
		replicas = *a.Replicas
	}
	ret := &pb.Application{}
	ret.SetType(a.Type)
	ret.SetName(a.Name)
	ret.SetNamespace(a.Namespace)
	ret.SetLabels(a.Labels)
	ret.SetReplicas(replicas)
	ret.SetHealthies(countHealthies(a.Pods))
	ret.SetContainers(toProtoContainers(a.Containers))
	ret.SetServices(toProtoServices(a.Services))
	ret.SetPods(toProtoPods(a.Pods))
	ret.SetPersistentVolumeClaims(toProtoPersistentVolumeClaims(a.Storages))
	ret.SetCreatedAt(timestamppb.New(a.ObjectMeta.CreationTimestamp.Time))
	ret.SetPublicAddress(publicAddress)
	if a.ChartFile != nil {
		ret.SetMetadata(toProtoChartMetadata(a.ChartFile))
	}
	return ret
}

func toProtoContainers(cs []core.Container) []*pb.Application_Container {
	ret := []*pb.Application_Container{}
	for i := range cs {
		ret = append(ret, toProtoContainer(&cs[i]))
	}
	return ret
}

func toProtoContainer(c *core.Container) *pb.Application_Container {
	ret := &pb.Application_Container{}
	ret.SetImageName(c.Image)
	ret.SetImagePullPolicy(string(c.ImagePullPolicy))
	return ret
}

func toProtoServices(ss []core.Service) []*pb.Application_Service {
	ret := []*pb.Application_Service{}
	for i := range ss {
		ret = append(ret, toProtoService(&ss[i]))
	}
	return ret
}

func toProtoService(s *core.Service) *pb.Application_Service {
	ret := &pb.Application_Service{}
	ret.SetName(s.Name)
	ret.SetType(string(s.Spec.Type))
	ret.SetClusterIp(s.Spec.ClusterIP)
	ret.SetPorts(toProtoServicePorts(s.Spec.Ports))
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	return ret
}

func toProtoServicePorts(ps []core.ServicePort) []*pb.Application_Service_Port {
	ret := []*pb.Application_Service_Port{}
	for i := range ps {
		ret = append(ret, toProtoServicePort(&ps[i]))
	}
	return ret
}

func toProtoServicePort(p *core.ServicePort) *pb.Application_Service_Port {
	ret := &pb.Application_Service_Port{}
	ret.SetPort(p.Port)
	ret.SetNodePort(p.NodePort)
	ret.SetName(p.Name)
	ret.SetProtocol(string(p.Protocol))
	ret.SetTargetPort(p.TargetPort.String())
	return ret
}

func toProtoPods(ps []core.Pod) []*pb.Application_Pod {
	ret := []*pb.Application_Pod{}
	for i := range ps {
		ret = append(ret, toProtoPod(&ps[i]))
	}
	return ret
}

func toProtoPod(p *core.Pod) *pb.Application_Pod {
	ret := &pb.Application_Pod{}
	ret.SetName(p.Name)
	ret.SetPhase(string(p.Status.Phase))
	ret.SetReady(containerStatusesReadyString(p.Status.ContainerStatuses))
	ret.SetRestarts(containerStatusesRestartString(p.Status.ContainerStatuses))
	if len(p.Status.Conditions) > 0 {
		index := len(p.Status.Conditions) - 1
		ret.SetLastCondition(toProtoLastCondition(&p.Status.Conditions[index]))
	}
	ret.SetCreatedAt(timestamppb.New(p.CreationTimestamp.Time))
	return ret
}

func toProtoLastCondition(c *core.PodCondition) *pb.Application_Condition {
	ret := &pb.Application_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))
	ret.SetProbedAt(timestamppb.New(c.LastProbeTime.Time))
	ret.SetTransitionedAt(timestamppb.New(c.LastTransitionTime.Time))
	return ret
}

func toProtoPersistentVolumeClaims(ss []core.Storage) []*pb.Application_PersistentVolumeClaim {
	ret := []*pb.Application_PersistentVolumeClaim{}
	for i := range ss {
		ret = append(ret, toProtoPersistentVolumeClaim(&ss[i]))
	}
	return ret
}

func toProtoPersistentVolumeClaim(s *core.Storage) *pb.Application_PersistentVolumeClaim {
	ret := &pb.Application_PersistentVolumeClaim{}
	ret.SetName(s.PersistentVolumeClaim.Name)
	ret.SetStatus(string(s.Status.Phase))
	ret.SetAccessModes(accessModesToStrings(s.Spec.AccessModes))
	ret.SetCapacity(s.Spec.Resources.Requests.Storage().String())
	if s.StorageClass != nil {
		ret.SetStorageClass(toProtoStorageClass(s.StorageClass))
	}
	ret.SetCreatedAt(timestamppb.New(s.PersistentVolumeClaim.CreationTimestamp.Time))
	return ret
}

func toProtoReleases(rs []core.Release) []*pb.Application_Release {
	ret := []*pb.Application_Release{}
	for i := range rs {
		ret = append(ret, toProtoRelease(&rs[i]))
	}
	return ret
}

func toProtoRelease(r *core.Release) *pb.Application_Release {
	ret := &pb.Application_Release{}
	ret.SetNamespace(r.Namespace)
	ret.SetName(r.Name)
	ret.SetRevision(int32(r.Version)) //nolint:gosec // ignore
	ret.SetChartName(r.Chart.Name())
	ret.SetVersion(toProtoChartVersion(&core.ChartVersion{
		Metadata: &core.ChartMetadata{
			Version:    r.Chart.Metadata.Version,
			AppVersion: r.Chart.Metadata.AppVersion,
		},
	}))
	return ret
}

func toProtoCharts(cs []core.Chart) []*pb.Application_Chart {
	ret := []*pb.Application_Chart{}
	for i := range cs {
		if len(cs[i].Versions) > 0 {
			ret = append(ret, toProtoChart(cs[i].Versions[0].Metadata, cs[i].Versions[0])) // latest only
		}
	}
	return ret
}

func toProtoChart(cmd *core.ChartMetadata, vs ...*core.ChartVersion) *pb.Application_Chart {
	ret := &pb.Application_Chart{}
	ret.SetName(cmd.Name)
	ret.SetIcon(cmd.Icon)
	ret.SetDescription(cmd.Description)
	ret.SetDeprecated(cmd.Deprecated)
	ret.SetTags(cmd.Tags)
	ret.SetKeywords(cmd.Keywords)
	ret.SetLicense(getChartLicense(cmd.Annotations))
	ret.SetVerified(false) // TODO: custom
	ret.SetHome(cmd.Home)
	ret.SetSources(cmd.Sources)
	ret.SetMaintainers(toProtoChartMaintainers(cmd.Maintainers))
	ret.SetDependencies(toProtoChartDependencies(cmd.Dependencies))
	ret.SetVersions(toProtoChartVersions(vs...))
	return ret
}

func toProtoChartMaintainers(ms []*core.ChartMaintainer) []*pb.Application_Chart_Maintainer {
	ret := []*pb.Application_Chart_Maintainer{}
	for i := range ms {
		ret = append(ret, toProtoChartMaintainer(ms[i]))
	}
	return ret
}

func toProtoChartMaintainer(m *core.ChartMaintainer) *pb.Application_Chart_Maintainer {
	ret := &pb.Application_Chart_Maintainer{}
	ret.SetName(m.Name)
	ret.SetEmail(m.Email)
	ret.SetUrl(m.URL)
	return ret
}

func toProtoChartDependencies(ds []*core.ChartDependency) []*pb.Application_Chart_Dependency {
	ret := []*pb.Application_Chart_Dependency{}
	for i := range ds {
		ret = append(ret, toProtoChartDependency(ds[i]))
	}
	return ret
}

func toProtoChartDependency(d *core.ChartDependency) *pb.Application_Chart_Dependency {
	ret := &pb.Application_Chart_Dependency{}
	ret.SetName(d.Name)
	ret.SetVersion(d.Version)
	ret.SetCondition(d.Condition)
	return ret
}

func toProtoChartVersions(vs ...*core.ChartVersion) []*pb.Application_Chart_Version {
	ret := []*pb.Application_Chart_Version{}
	for _, v := range vs {
		ret = append(ret, toProtoChartVersion(v))
	}
	return ret
}

func toProtoChartVersion(v *core.ChartVersion) *pb.Application_Chart_Version {
	ret := &pb.Application_Chart_Version{}
	ret.SetChartVersion(v.Version)
	ret.SetApplicationVersion(v.AppVersion)
	if len(v.URLs) > 0 {
		ret.SetChartRef(v.URLs[0])
	}
	return ret
}

func toProtoStorageClasses(scs []core.StorageClass) []*pb.StorageClass {
	ret := []*pb.StorageClass{}
	for i := range scs {
		ret = append(ret, toProtoStorageClass(&scs[i]))
	}
	return ret
}

func toProtoStorageClass(sc *core.StorageClass) *pb.StorageClass {
	reclaimPolicy := ""
	if v := sc.ReclaimPolicy; v != nil {
		reclaimPolicy = string(*v)
	}
	volumeBindingMode := ""
	if v := sc.VolumeBindingMode; v != nil {
		volumeBindingMode = string(*v)
	}
	ret := &pb.StorageClass{}
	ret.SetName(sc.Name)
	ret.SetProvisioner(sc.Provisioner)
	ret.SetReclaimPolicy(reclaimPolicy)
	ret.SetVolumeBindingMode(volumeBindingMode)
	ret.SetParameters(sc.Parameters)
	ret.SetCreatedAt(timestamppb.New(sc.CreationTimestamp.Time))
	return ret
}

func toProtoChartMetadata(md *core.ChartFile) *pb.Application_Chart_Metadata {
	ret := &pb.Application_Chart_Metadata{}
	ret.SetValuesYaml(md.ValuesYAML)
	ret.SetReadmeMd(md.ReadmeMarkdown)
	ret.SetCustomization(toProtoChartCustomization(md.Customization))
	return ret
}

func toProtoChartCustomization(c map[string]any) *pb.Application_Chart_Customization {
	ret := &pb.Application_Chart_Customization{}
	values, err := structpb.NewStruct(c)
	if err == nil {
		ret.SetValues(values)
	}
	return ret
}

func getChartLicense(m map[string]string) string {
	keys := []string{"license", "licenses"}
	for _, key := range keys {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return ""
}

func containerStatusesReadyString(statuses []core.ContainerStatus) string {
	ready := 0
	for i := range statuses {
		if statuses[i].Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, len(statuses))
}

func containerStatusesRestartString(statuses []core.ContainerStatus) string {
	restart := int32(0)
	var lastTerminatedAt core.KubernetesTime
	for i := range statuses {
		restart += statuses[i].RestartCount
		if statuses[i].LastTerminationState.Terminated != nil {
			lastTerminatedAt = statuses[i].LastTerminationState.Terminated.FinishedAt
		}
	}
	if lastTerminatedAt.IsZero() {
		return fmt.Sprintf("%d", restart)
	}
	return fmt.Sprintf("%d (%s ago)", restart, duration.HumanDuration(time.Since(lastTerminatedAt.Time)))
}

func countHealthies(pods []core.Pod) int32 {
	count := int32(0)
	for i := range pods {
		if core.IsPodHealthy(&pods[i]) {
			count++
		}
	}
	return count
}

func accessModesToStrings(modes []core.PersistentVolumeAccessMode) []string {
	ret := make([]string, len(modes))
	for i := range modes {
		ret[i] = string(modes[i])
	}
	return ret
}
