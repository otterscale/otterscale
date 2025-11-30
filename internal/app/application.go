package app

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/apimachinery/pkg/util/duration"

	pb "github.com/otterscale/otterscale/api/application/v1"
	"github.com/otterscale/otterscale/api/application/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/registry/chart"
)

type ApplicationService struct {
	pbconnect.UnimplementedApplicationServiceHandler

	cluster    *cluster.UseCase
	chart      *chart.UseCase
	config     *config.UseCase
	persistent *persistent.UseCase
	release    *release.UseCase
	service    *service.UseCase
	workload   *workload.UseCase
}

func NewApplicationService(cluster *cluster.UseCase, chart *chart.UseCase, config *config.UseCase, persistent *persistent.UseCase, release *release.UseCase, service *service.UseCase, workload *workload.UseCase) *ApplicationService {
	return &ApplicationService{
		cluster:    cluster,
		chart:      chart,
		config:     config,
		persistent: persistent,
		release:    release,
		service:    service,
		workload:   workload,
	}
}

var _ pbconnect.ApplicationServiceHandler = (*ApplicationService)(nil)

func (s *ApplicationService) ListApplications(ctx context.Context, req *pb.ListApplicationsRequest) (*pb.ListApplicationsResponse, error) {
	apps, hostname, err := s.workload.ListApplications(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListApplicationsResponse{}
	resp.SetApplications(toProtoApplications(apps))
	resp.SetHostname(hostname)
	return resp, nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *pb.GetApplicationRequest) (*pb.Application, error) {
	app, err := s.workload.GetApplication(ctx, req.GetScope(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	// chartFile, err := s.release.GetChartFileFromApplication(ctx, req.GetScope(), req.GetNamespace(), app.Labels, app.Annotations)
	// if err != nil {
	// 	return nil, err
	// }

	// app.ChartFile = chartFile

	resp := toProtoApplication(app)
	return resp, nil
}

func (s *ApplicationService) RestartApplication(ctx context.Context, req *pb.RestartApplicationRequest) (*emptypb.Empty, error) {
	if err := s.workload.RestartApplication(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetType()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ScaleApplication(ctx context.Context, req *pb.ScaleApplicationRequest) (*emptypb.Empty, error) {
	if err := s.workload.ScaleApplication(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetType(), req.GetReplicas()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) DeleteApplicationPod(ctx context.Context, req *pb.DeleteApplicationPodRequest) (*emptypb.Empty, error) {
	err := s.workload.DeletePod(ctx, req.GetScope(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) WatchLogs(ctx context.Context, req *pb.WatchLogsRequest, stream *connect.ServerStream[pb.WatchLogsResponse]) error {
	logs, err := s.workload.StreamLogs(ctx, req.GetScope(), req.GetNamespace(), req.GetPodName(), req.GetContainerName(), req.GetDuration().AsDuration())
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
	if err := s.workload.WriteToTTYSession(req.GetSessionId(), req.GetStdin()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ExecuteTTY(ctx context.Context, req *pb.ExecuteTTYRequest, stream *connect.ServerStream[pb.ExecuteTTYResponse]) error {
	// create session pipes
	sessionID, err := s.workload.CreateTTYSession()
	if err != nil {
		return err
	}
	defer func() {
		_ = s.workload.CleanupTTYSession(sessionID)
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
		_ = s.workload.ExecuteTTY(ctx, sessionID, req.GetScope(), req.GetNamespace(), req.GetPodName(), req.GetContainerName(), req.GetCommand(), stdOutChan)
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
	releases, err := s.release.ListReleases(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListReleasesResponse{}
	resp.SetReleases(toProtoReleases(releases))
	return resp, nil
}

func (s *ApplicationService) CreateRelease(ctx context.Context, req *pb.CreateReleaseRequest) (*pb.Release, error) {
	release, err := s.release.CreateRelease(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDryRun(), req.GetChartRef(), req.GetValuesYaml(), req.GetValuesMap())
	if err != nil {
		return nil, err
	}

	resp := toProtoRelease(release)
	return resp, nil
}

func (s *ApplicationService) UpdateRelease(ctx context.Context, req *pb.UpdateReleaseRequest) (*pb.Release, error) {
	release, err := s.release.UpdateRelease(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDryRun(), req.GetChartRef(), req.GetValuesYaml())
	if err != nil {
		return nil, err
	}

	resp := toProtoRelease(release)
	return resp, nil
}

func (s *ApplicationService) DeleteRelease(ctx context.Context, req *pb.DeleteReleaseRequest) (*emptypb.Empty, error) {
	if err := s.release.DeleteRelease(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDryRun()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) RollbackRelease(ctx context.Context, req *pb.RollbackReleaseRequest) (*emptypb.Empty, error) {
	if err := s.release.RollbackRelease(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDryRun()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ApplicationService) ListConfigMaps(ctx context.Context, req *pb.ListConfigMapsRequest) (*pb.ListConfigMapsResponse, error) {
	configMaps, err := s.config.ListConfigMaps(ctx, req.GetScope(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListConfigMapsResponse{}
	resp.SetConfigMaps(toProtoConfigMaps(configMaps))
	return resp, nil
}

func (s *ApplicationService) ListSecrets(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	secrets, err := s.config.ListSecrets(ctx, req.GetScope(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListSecretsResponse{}
	resp.SetSecrets(toProtoSecrets(secrets))
	return resp, nil
}

func (s *ApplicationService) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {
	namespaces, err := s.cluster.ListNamespaces(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListNamespacesResponse{}
	resp.SetNamespaces(toProtoNamespaces(namespaces))
	return resp, nil
}

func (s *ApplicationService) ListStorageClasses(ctx context.Context, req *pb.ListStorageClassesRequest) (*pb.ListStorageClassesResponse, error) {
	storageClasses, err := s.persistent.ListStorageClasses(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListStorageClassesResponse{}
	resp.SetStorageClasses(toProtoStorageClasses(storageClasses))
	return resp, nil
}

func toProtoApplications(as []workload.Application) []*pb.Application {
	ret := []*pb.Application{}

	for i := range as {
		ret = append(ret, toProtoApplication(&as[i]))
	}

	return ret
}

func toProtoApplication(a *workload.Application) *pb.Application {
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
	ret.SetPersistentVolumeClaims(toProtoPersistentVolumeClaims(a.Persistents))
	ret.SetCreatedAt(timestamppb.New(a.CreatedAt))

	return ret
}

func toProtoContainers(cs []workload.Container) []*pb.Application_Container {
	ret := []*pb.Application_Container{}

	for i := range cs {
		ret = append(ret, toProtoContainer(&cs[i]))
	}

	return ret
}

func toProtoContainer(c *workload.Container) *pb.Application_Container {
	ret := &pb.Application_Container{}
	ret.SetImageName(c.Image)
	ret.SetImagePullPolicy(string(c.ImagePullPolicy))
	return ret
}

func toProtoServices(ss []service.Service) []*pb.Application_Service {
	ret := []*pb.Application_Service{}

	for i := range ss {
		ret = append(ret, toProtoService(&ss[i]))
	}

	return ret
}

func toProtoService(s *service.Service) *pb.Application_Service {
	ret := &pb.Application_Service{}
	ret.SetName(s.Name)
	ret.SetType(string(s.Spec.Type))
	ret.SetClusterIp(s.Spec.ClusterIP)
	ret.SetPorts(toProtoServicePorts(s.Spec.Ports))
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	return ret
}

func toProtoServicePorts(ps []service.Port) []*pb.Application_Service_Port {
	ret := []*pb.Application_Service_Port{}

	for i := range ps {
		ret = append(ret, toProtoServicePort(&ps[i]))
	}

	return ret
}

func toProtoServicePort(p *service.Port) *pb.Application_Service_Port {
	ret := &pb.Application_Service_Port{}
	ret.SetPort(p.Port)
	ret.SetNodePort(p.NodePort)
	ret.SetName(p.Name)
	ret.SetProtocol(string(p.Protocol))
	ret.SetTargetPort(p.TargetPort.String())
	return ret
}

func toProtoPods(ps []workload.Pod) []*pb.Application_Pod {
	ret := []*pb.Application_Pod{}

	for i := range ps {
		ret = append(ret, toProtoPod(&ps[i]))
	}

	return ret
}

func toProtoPod(p *workload.Pod) *pb.Application_Pod {
	ret := &pb.Application_Pod{}
	ret.SetName(p.Name)
	ret.SetPhase(string(p.Status.Phase))
	ret.SetReady(containerStatusesReadyString(p.Status.ContainerStatuses))
	ret.SetRestarts(containerStatusesRestartString(p.Status.ContainerStatuses))

	conditions := p.Status.Conditions

	for i := range conditions {
		if conditions[i].Status == workload.ConditionTrue {
			ret.SetLastCondition(toProtoApplicationCondition(&conditions[i]))

			break
		}
	}

	ret.SetCreatedAt(timestamppb.New(p.CreationTimestamp.Time))

	return ret
}

func toProtoApplicationCondition(c *workload.PodCondition) *pb.Application_Condition {
	ret := &pb.Application_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))

	probedAt := c.LastProbeTime.Time

	if !probedAt.IsZero() {
		ret.SetProbedAt(timestamppb.New(probedAt))
	}

	transitionedAt := c.LastTransitionTime.Time

	if !transitionedAt.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(transitionedAt))
	}

	return ret
}

func toProtoPersistentVolumeClaims(ss []persistent.Persistent) []*pb.Application_PersistentVolumeClaim {
	ret := []*pb.Application_PersistentVolumeClaim{}

	for i := range ss {
		ret = append(ret, toProtoPersistentVolumeClaim(&ss[i]))
	}

	return ret
}

func toProtoPersistentVolumeClaim(s *persistent.Persistent) *pb.Application_PersistentVolumeClaim {
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

func toProtoReleases(rs []release.Release) []*pb.Release {
	ret := []*pb.Release{}

	for i := range rs {
		ret = append(ret, toProtoRelease(&rs[i]))
	}

	return ret
}

func toProtoRelease(r *release.Release) *pb.Release {
	ret := &pb.Release{}
	ret.SetNamespace(r.Namespace)
	ret.SetName(r.Name)
	ret.SetRevision(int32(r.Version)) //nolint:gosec // ignore

	chart := r.Chart
	if chart != nil {
		ret.SetChart(toProtoChart(chart.Metadata))
	}

	return ret
}

func toProtoConfigMaps(cms []config.ConfigMap) []*pb.ConfigMap {
	ret := []*pb.ConfigMap{}

	for i := range cms {
		ret = append(ret, toProtoConfigMap(&cms[i]))
	}

	return ret
}

func toProtoConfigMap(cm *config.ConfigMap) *pb.ConfigMap {
	ret := &pb.ConfigMap{}
	ret.SetName(cm.Name)
	ret.SetNamespace(cm.Namespace)
	ret.SetLabels(cm.Labels)
	ret.SetData(cm.Data)
	ret.SetBinaryData(cm.BinaryData)

	if cm.Immutable != nil {
		ret.SetImmutable(*cm.Immutable)
	}

	ret.SetCreatedAt(timestamppb.New(cm.CreationTimestamp.Time))

	return ret
}

func toProtoSecrets(ss []config.Secret) []*pb.Secret {
	ret := []*pb.Secret{}

	for i := range ss {
		ret = append(ret, toProtoSecret(&ss[i]))
	}

	return ret
}

func toProtoSecret(s *config.Secret) *pb.Secret {
	ret := &pb.Secret{}
	ret.SetName(s.Name)
	ret.SetNamespace(s.Namespace)
	ret.SetLabels(s.Labels)
	ret.SetData(s.Data)
	ret.SetType(string(s.Type))
	ret.SetStringData(s.StringData)

	if s.Immutable != nil {
		ret.SetImmutable(*s.Immutable)
	}

	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))

	return ret
}

func toProtoNamespaces(ns []cluster.Namespace) []*pb.Namespace {
	ret := []*pb.Namespace{}

	for i := range ns {
		ret = append(ret, toProtoNamespace(&ns[i]))
	}

	return ret
}

func toProtoNamespace(n *cluster.Namespace) *pb.Namespace {
	ret := &pb.Namespace{}
	ret.SetName(n.Name)
	ret.SetLabels(n.Labels)
	ret.SetCreatedAt(timestamppb.New(n.CreationTimestamp.Time))
	return ret
}

func toProtoStorageClasses(scs []persistent.StorageClass) []*pb.StorageClass {
	ret := []*pb.StorageClass{}

	for i := range scs {
		ret = append(ret, toProtoStorageClass(&scs[i]))
	}

	return ret
}

func toProtoStorageClass(sc *persistent.StorageClass) *pb.StorageClass {
	reclaimPolicy := ""
	volumeBindingMode := ""

	if v := sc.ReclaimPolicy; v != nil {
		reclaimPolicy = string(*v)
	}

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

func containerStatusesReadyString(statuses []workload.ContainerStatus) string {
	ready := 0

	for i := range statuses {
		if statuses[i].Ready {
			ready++
		}
	}

	return fmt.Sprintf("%d/%d", ready, len(statuses))
}

func containerStatusesRestartString(statuses []workload.ContainerStatus) string {
	restart := int32(0)
	var lastTerminatedAt workload.Time

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

func countHealthies(pods []workload.Pod) int32 {
	count := int32(0)

	for i := range pods {
		if pods[i].Status.Phase != workload.PodPhaseRunning {
			continue
		}

		count++
	}

	return count
}

func accessModesToStrings(modes []persistent.PersistentVolumeAccessMode) []string {
	ret := make([]string, len(modes))

	for i := range modes {
		ret[i] = string(modes[i])
	}

	return ret
}
