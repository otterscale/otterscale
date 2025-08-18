package app

import (
	"context"
	"fmt"
	"slices"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"

	pb "github.com/openhdc/otterscale/api/application/v1"
	"github.com/openhdc/otterscale/api/application/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type ApplicationService struct {
	pbconnect.UnimplementedApplicationServiceHandler

	uc *core.ApplicationUseCase
}

func NewApplicationService(uc *core.ApplicationUseCase) *ApplicationService {
	return &ApplicationService{uc: uc}
}

var _ pbconnect.ApplicationServiceHandler = (*ApplicationService)(nil)

func (s *ApplicationService) ListApplications(ctx context.Context, req *connect.Request[pb.ListApplicationsRequest]) (*connect.Response[pb.ListApplicationsResponse], error) {
	apps, err := s.uc.ListApplications(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	publicAddress, err := s.uc.GetPublicAddress(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListApplicationsResponse{}
	resp.SetApplications(toProtoApplications(apps, publicAddress))
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, req *connect.Request[pb.GetApplicationRequest]) (*connect.Response[pb.Application], error) {
	app, err := s.uc.GetApplication(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	metadata, err := s.uc.GetChartMetadataFromApplication(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), app)
	if err != nil {
		return nil, err
	}
	app.ChartMetadata = metadata
	publicAddress, err := s.uc.GetPublicAddress(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := toProtoApplication(app, publicAddress)
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) ListReleases(ctx context.Context, req *connect.Request[pb.ListReleasesRequest]) (*connect.Response[pb.ListReleasesResponse], error) {
	releases, err := s.uc.ListReleases(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListReleasesResponse{}
	resp.SetReleases(toProtoReleases(releases))
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) CreateRelease(ctx context.Context, req *connect.Request[pb.CreateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	release, err := s.uc.CreateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml(), req.Msg.GetValuesMap())
	if err != nil {
		return nil, err
	}
	resp := toProtoRelease(release)
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) UpdateRelease(ctx context.Context, req *connect.Request[pb.UpdateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	release, err := s.uc.UpdateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	resp := toProtoRelease(release)
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) DeleteRelease(ctx context.Context, req *connect.Request[pb.DeleteReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) RollbackRelease(ctx context.Context, req *connect.Request[pb.RollbackReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.RollbackRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) ListCharts(ctx context.Context, req *connect.Request[pb.ListChartsRequest]) (*connect.Response[pb.ListChartsResponse], error) {
	charts, err := s.uc.ListCharts(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListChartsResponse{}
	resp.SetCharts(toProtoCharts(charts))
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) GetChart(ctx context.Context, req *connect.Request[pb.GetChartRequest]) (*connect.Response[pb.Application_Chart], error) {
	ch, err := s.uc.GetChart(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	metadata := &chart.Metadata{}
	if len(ch.Versions) > 0 {
		metadata = ch.Versions[0].Metadata
	}
	resp := toProtoChart(metadata, ch.Versions...)
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) GetChartMetadata(ctx context.Context, req *connect.Request[pb.GetChartMetadataRequest]) (*connect.Response[pb.Application_Chart_Metadata], error) {
	metadata, err := s.uc.GetChartMetadata(ctx, req.Msg.GetChartRef())
	if err != nil {
		return nil, err
	}
	resp := toProtoChartMetadata(metadata)
	return connect.NewResponse(resp), nil
}

func (s *ApplicationService) ListStorageClasses(ctx context.Context, req *connect.Request[pb.ListStorageClassesRequest]) (*connect.Response[pb.ListStorageClassesResponse], error) {
	storageClasses, err := s.uc.ListStorageClasses(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStorageClassesResponse{}
	resp.SetStorageClasses(toProtoStorageClasses(storageClasses))
	return connect.NewResponse(resp), nil
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
	if a.ChartMetadata != nil {
		ret.SetMetadata(toProtoChartMetadata(a.ChartMetadata))
	}
	return ret
}

func toProtoContainers(cs []corev1.Container) []*pb.Application_Container {
	ret := []*pb.Application_Container{}
	for i := range cs {
		ret = append(ret, toProtoContainer(&cs[i]))
	}
	return ret
}

func toProtoContainer(c *corev1.Container) *pb.Application_Container {
	ret := &pb.Application_Container{}
	ret.SetImageName(c.Image)
	ret.SetImagePullPolicy(string(c.ImagePullPolicy))
	return ret
}

func toProtoServices(ss []corev1.Service) []*pb.Application_Service {
	ret := []*pb.Application_Service{}
	for i := range ss {
		ret = append(ret, toProtoService(&ss[i]))
	}
	return ret
}

func toProtoService(s *corev1.Service) *pb.Application_Service {
	ret := &pb.Application_Service{}
	ret.SetName(s.Name)
	ret.SetType(string(s.Spec.Type))
	ret.SetClusterIp(s.Spec.ClusterIP)
	ret.SetPorts(toProtoServicePorts(s.Spec.Ports))
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	return ret
}

func toProtoServicePorts(ps []corev1.ServicePort) []*pb.Application_Service_Port {
	ret := []*pb.Application_Service_Port{}
	for i := range ps {
		ret = append(ret, toProtoServicePort(&ps[i]))
	}
	return ret
}

func toProtoServicePort(p *corev1.ServicePort) *pb.Application_Service_Port {
	ret := &pb.Application_Service_Port{}
	ret.SetPort(p.Port)
	ret.SetNodePort(p.NodePort)
	ret.SetProtocol(string(p.Protocol))
	ret.SetTargetPort(p.TargetPort.String())
	return ret
}

func toProtoPods(ps []corev1.Pod) []*pb.Application_Pod {
	ret := []*pb.Application_Pod{}
	for i := range ps {
		ret = append(ret, toProtoPod(&ps[i]))
	}
	return ret
}

func toProtoPod(p *corev1.Pod) *pb.Application_Pod {
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

func toProtoLastCondition(c *corev1.PodCondition) *pb.Application_Condition {
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
	ret.SetStatus(string(s.PersistentVolumeClaim.Status.Phase))
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
	ret.SetRevision(int32(r.Version)) //nolint:gosec
	ret.SetChartName(r.Chart.Name())
	ret.SetVersion(toProtoChartVersion(&repo.ChartVersion{
		Metadata: &chart.Metadata{
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

func toProtoChart(cmd *chart.Metadata, vs ...*repo.ChartVersion) *pb.Application_Chart {
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

func toProtoChartMaintainers(ms []*chart.Maintainer) []*pb.Application_Chart_Maintainer {
	ret := []*pb.Application_Chart_Maintainer{}
	for i := range ms {
		ret = append(ret, toProtoChartMaintainer(ms[i]))
	}
	return ret
}

func toProtoChartMaintainer(m *chart.Maintainer) *pb.Application_Chart_Maintainer {
	ret := &pb.Application_Chart_Maintainer{}
	ret.SetName(m.Name)
	ret.SetEmail(m.Email)
	ret.SetUrl(m.URL)
	return ret
}

func toProtoChartDependencies(ds []*chart.Dependency) []*pb.Application_Chart_Dependency {
	ret := []*pb.Application_Chart_Dependency{}
	for i := range ds {
		ret = append(ret, toProtoChartDependency(ds[i]))
	}
	return ret
}

func toProtoChartDependency(d *chart.Dependency) *pb.Application_Chart_Dependency {
	ret := &pb.Application_Chart_Dependency{}
	ret.SetName(d.Name)
	ret.SetVersion(d.Version)
	ret.SetCondition(d.Condition)
	return ret
}

func toProtoChartVersions(vs ...*repo.ChartVersion) []*pb.Application_Chart_Version {
	ret := []*pb.Application_Chart_Version{}
	for _, v := range vs {
		ret = append(ret, toProtoChartVersion(v))
	}
	return ret
}

func toProtoChartVersion(v *repo.ChartVersion) *pb.Application_Chart_Version {
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

func toProtoChartMetadata(md *core.ChartMetadata) *pb.Application_Chart_Metadata {
	ret := &pb.Application_Chart_Metadata{}
	ret.SetValuesYaml(md.ValuesYAML)
	ret.SetReadmeMd(md.ReadmeMD)
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

func containerStatusesReadyString(statuses []corev1.ContainerStatus) string {
	ready := 0
	for i := range statuses {
		if statuses[i].Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, len(statuses))
}

func containerStatusesRestartString(statuses []corev1.ContainerStatus) string {
	restart := int32(0)
	var lastTerminatedAt metav1.Time
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

func countHealthies(pods []corev1.Pod) int32 {
	phases := []corev1.PodPhase{corev1.PodRunning, corev1.PodSucceeded}
	count := int32(0)
	for i := range pods {
		if slices.Contains(phases, pods[i].Status.Phase) {
			count++
		}
	}
	return count
}

func accessModesToStrings(modes []corev1.PersistentVolumeAccessMode) []string {
	ret := make([]string, len(modes))
	for i := range modes {
		ret[i] = string(modes[i])
	}
	return ret
}
