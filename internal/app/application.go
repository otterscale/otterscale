package app

import (
	"context"
	"fmt"
	"slices"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListApplications(ctx context.Context, req *connect.Request[pb.ListApplicationsRequest]) (*connect.Response[pb.ListApplicationsResponse], error) {
	as, err := a.svc.ListApplications(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	res := &pb.ListApplicationsResponse{}
	res.SetApplications(toProtoApplications(as))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetApplication(ctx context.Context, req *connect.Request[pb.GetApplicationRequest]) (*connect.Response[pb.Application], error) {
	app, err := a.svc.GetApplication(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoApplication(app)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListReleases(ctx context.Context, req *connect.Request[pb.ListReleasesRequest]) (*connect.Response[pb.ListReleasesResponse], error) {
	rs, err := a.svc.ListReleases(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListReleasesResponse{}
	res.SetReleases(toProtoReleases(rs))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateRelease(ctx context.Context, req *connect.Request[pb.CreateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	r, err := a.svc.CreateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	res := toProtoRelease(r)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateRelease(ctx context.Context, req *connect.Request[pb.UpdateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	r, err := a.svc.UpdateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	res := toProtoRelease(r)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteRelease(ctx context.Context, req *connect.Request[pb.DeleteReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) RollbackRelease(ctx context.Context, req *connect.Request[pb.RollbackReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.RollbackRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListCharts(ctx context.Context, req *connect.Request[pb.ListChartsRequest]) (*connect.Response[pb.ListChartsResponse], error) {
	cs, err := a.svc.ListCharts(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListChartsResponse{}
	res.SetCharts(toProtoCharts(cs))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetChart(ctx context.Context, req *connect.Request[pb.GetChartRequest]) (*connect.Response[pb.Application_Release_Chart], error) {
	c, err := a.svc.GetChart(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	md := &chart.Metadata{}
	if len(c.Versions) > 0 {
		md = c.Versions[0].Metadata
	}
	res := toProtoChart(md, c.Versions...)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetChartMetadata(ctx context.Context, req *connect.Request[pb.GetChartMetadataRequest]) (*connect.Response[pb.Application_Release_Chart_Metadata], error) {
	md, err := a.svc.GetChartMetadata(ctx, req.Msg.GetChartRef())
	if err != nil {
		return nil, err
	}
	res := toProtoChartMetadata(md)
	return connect.NewResponse(res), nil
}

func toProtoApplications(as []model.Application) []*pb.Application {
	ret := []*pb.Application{}
	for i := range as {
		ret = append(ret, toProtoApplication(&as[i]))
	}
	return ret
}

func toProtoApplication(a *model.Application) *pb.Application {
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
	ret.SetPersistentVolumeClaims(toProtoPersistentVolumeClaims(a.PersistentVolumeClaims))
	ret.SetCreatedAt(timestamppb.New(a.ObjectMeta.CreationTimestamp.Time))
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
	ret.SetLastCondition(toProtoLastCondition(p.Status.Conditions))
	ret.SetCreatedAt(timestamppb.New(p.CreationTimestamp.Time))
	return ret
}

func toProtoLastCondition(cs []corev1.PodCondition) *pb.Application_Condition {
	i := len(cs) - 1
	ret := &pb.Application_Condition{}
	ret.SetType(string(cs[i].Type))
	ret.SetStatus(string(cs[i].Status))
	ret.SetReason((cs[i].Reason))
	ret.SetMessage((cs[i].Message))
	ret.SetProbedAt(timestamppb.New(cs[i].LastProbeTime.Time))
	ret.SetTransitionedAt(timestamppb.New(cs[i].LastTransitionTime.Time))
	return ret
}

func toProtoPersistentVolumeClaims(ps []model.PersistentVolumeClaim) []*pb.Application_PersistentVolumeClaim {
	ret := []*pb.Application_PersistentVolumeClaim{}
	for i := range ps {
		ret = append(ret, toProtoPersistentVolumeClaim(&ps[i]))
	}
	return ret
}

func toProtoPersistentVolumeClaim(p *model.PersistentVolumeClaim) *pb.Application_PersistentVolumeClaim {
	ret := &pb.Application_PersistentVolumeClaim{}
	ret.SetName(p.PersistentVolumeClaim.Name)
	ret.SetStatus(string(p.PersistentVolumeClaim.Status.Phase))
	ret.SetAccessModes(accessModesToStrings(p.Spec.AccessModes))
	ret.SetCapacity(p.Spec.Resources.Requests.Storage().String())
	if p.StorageClass != nil {
		ret.SetStorageClass(toProtoStorageClass(p.StorageClass))
	}
	ret.SetCreatedAt(timestamppb.New(p.PersistentVolumeClaim.CreationTimestamp.Time))
	return ret
}

func toProtoStorageClass(sc *storagev1.StorageClass) *pb.Application_PersistentVolumeClaim_StorageClass {
	reclaimPolicy := ""
	if v := sc.ReclaimPolicy; v != nil {
		reclaimPolicy = string(*v)
	}
	volumeBindingMode := ""
	if v := sc.VolumeBindingMode; v != nil {
		volumeBindingMode = string(*v)
	}
	ret := &pb.Application_PersistentVolumeClaim_StorageClass{}
	ret.SetName(sc.Name)
	ret.SetProvisioner(sc.Provisioner)
	ret.SetReclaimPolicy(reclaimPolicy)
	ret.SetVolumeBindingMode(volumeBindingMode)
	ret.SetParameters(sc.Parameters)
	ret.SetCreatedAt(timestamppb.New(sc.CreationTimestamp.Time))
	return ret
}

func toProtoReleases(rs []model.Release) []*pb.Application_Release {
	ret := []*pb.Application_Release{}
	for i := range rs {
		ret = append(ret, toProtoRelease(&rs[i]))
	}
	return ret
}

func toProtoRelease(r *model.Release) *pb.Application_Release {
	ret := &pb.Application_Release{}
	ret.SetScopeName(r.ScopeName)
	ret.SetScopeUuid(r.ScopeUUID)
	ret.SetFacilityName(r.FacilityName)
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

func toProtoCharts(cs []model.Chart) []*pb.Application_Release_Chart {
	ret := []*pb.Application_Release_Chart{}
	for i := range cs {
		if len(cs[i].Versions) > 0 {
			ret = append(ret, toProtoChart(cs[i].Versions[0].Metadata, cs[i].Versions[0])) // latest only
		}
	}
	return ret
}

func toProtoChart(cmd *chart.Metadata, vs ...*repo.ChartVersion) *pb.Application_Release_Chart {
	ret := &pb.Application_Release_Chart{}
	ret.SetName(cmd.Name)
	ret.SetIcon(cmd.Icon)
	ret.SetDescription(cmd.Description)
	ret.SetDeprecated(cmd.Deprecated)
	ret.SetTags(cmd.Tags)
	ret.SetKeywords(cmd.Keywords)
	ret.SetLicense(getChartLicense(cmd.Annotations))
	ret.SetHome(cmd.Home)
	ret.SetSources(cmd.Sources)
	ret.SetMaintainers(toProtoChartMaintainers(cmd.Maintainers))
	ret.SetDependencies(toProtoChartDependencies(cmd.Dependencies))
	ret.SetVersions(toProtoChartVersions(vs...))
	return ret
}

func toProtoChartMaintainers(ms []*chart.Maintainer) []*pb.Application_Release_Chart_Maintainer {
	ret := []*pb.Application_Release_Chart_Maintainer{}
	for i := range ms {
		ret = append(ret, toProtoChartMaintainer(ms[i]))
	}
	return ret
}

func toProtoChartMaintainer(m *chart.Maintainer) *pb.Application_Release_Chart_Maintainer {
	ret := &pb.Application_Release_Chart_Maintainer{}
	ret.SetName(m.Name)
	ret.SetEmail(m.Email)
	ret.SetUrl(m.URL)
	return ret
}

func toProtoChartDependencies(ds []*chart.Dependency) []*pb.Application_Release_Chart_Dependency {
	ret := []*pb.Application_Release_Chart_Dependency{}
	for i := range ds {
		ret = append(ret, toProtoChartDependency(ds[i]))
	}
	return ret
}

func toProtoChartDependency(d *chart.Dependency) *pb.Application_Release_Chart_Dependency {
	ret := &pb.Application_Release_Chart_Dependency{}
	ret.SetName(d.Name)
	ret.SetVersion(d.Version)
	ret.SetCondition(d.Condition)
	return ret
}

func toProtoChartVersions(vs ...*repo.ChartVersion) []*pb.Application_Release_Chart_Version {
	ret := []*pb.Application_Release_Chart_Version{}
	for _, v := range vs {
		ret = append(ret, toProtoChartVersion(v))
	}
	return ret
}

func toProtoChartVersion(v *repo.ChartVersion) *pb.Application_Release_Chart_Version {
	ret := &pb.Application_Release_Chart_Version{}
	ret.SetChartVersion(v.Version)
	ret.SetApplicationVersion(v.AppVersion)
	if len(v.URLs) > 0 {
		ret.SetChartRef(v.URLs[0])
	}
	return ret
}

func toProtoChartMetadata(md *model.ChartMetadata) *pb.Application_Release_Chart_Metadata {
	ret := &pb.Application_Release_Chart_Metadata{}
	ret.SetValuesYaml(md.ValuesYAML)
	ret.SetReadmeMd(md.ReadmeMD)
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
