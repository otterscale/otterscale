package kube

import (
	"context"

	oscore "github.com/openhdc/otterscale/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	clonev1 "kubevirt.io/api/clone/v1beta1"
	virtv1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
)

type virtVM struct {
	kubevirt *kubevirt
	kube     *Kube
}

func NewVirtVM(kube *Kube, kubevirt *kubevirt) oscore.KubeVirtVMRepo {
	return &virtVM{
		kubevirt: kubevirt,
	}
}

var _ oscore.KubeVirtVMRepo = (*virtVM)(nil)

func (r *virtVM) CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) (*oscore.VirtualMachine, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	vm := &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		vm.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachine(namespace).Create(ctx, vm, opts)
}

func (r *virtVM) GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachine, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	virtClient.CdiClient().Discovery()
	opts := metav1.GetOptions{}
	return virtClient.VirtualMachine(namespace).Get(ctx, name, opts)
}

func (r *virtVM) ListVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) ([]oscore.VirtualMachine, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	vms, err := virtClient.VirtualMachine(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return vms.Items, nil
}

func (r *virtVM) UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) (*oscore.VirtualMachine, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	vm := &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		vm.Spec = *spec
	}

	opts := metav1.UpdateOptions{}
	return virtClient.VirtualMachine(namespace).Update(ctx, vm, opts)
}

func (r *virtVM) DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachine(namespace).Delete(ctx, name, opts)
}

func (r *virtVM) MigrateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.MigrateOptions{}
	return virtClient.VirtualMachine(namespace).Migrate(ctx, name, opts)
}

func (r *virtVM) StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.StartOptions{}
	return virtClient.VirtualMachine(namespace).Start(ctx, name, opts)
}

func (r *virtVM) RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.RestartOptions{}
	return virtClient.VirtualMachine(namespace).Restart(ctx, name, opts)
}

func (r *virtVM) StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.StopOptions{}
	return virtClient.VirtualMachine(namespace).Stop(ctx, name, opts)
}

// VirtualMachine Clone
func (r *virtVM) CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMCloneSpec) (*oscore.VirtualMachineClone, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	clone := &clonev1.VirtualMachineClone{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		clone.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachineClone(namespace).Create(ctx, clone, opts)
}

func (r *virtVM) GetVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineClone, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return virtClient.VirtualMachineClone(namespace).Get(ctx, name, opts)
}

func (r *virtVM) ListVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) ([]oscore.VirtualMachineClone, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	clones, err := virtClient.VirtualMachineClone(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return clones.Items, nil
}

func (r *virtVM) DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachineClone(namespace).Delete(ctx, name, opts)
}

// VirtualMachine snapshot
func (r *virtVM) CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSnapshotSpec) (*oscore.VirtualMachineSnapshot, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	snapshot := &snapshotv1.VirtualMachineSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		snapshot.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachineSnapshot(namespace).Create(ctx, snapshot, opts)
}

func (r *virtVM) GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineSnapshot, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return virtClient.VirtualMachineSnapshot(namespace).Get(ctx, name, opts)
}

func (r *virtVM) ListVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSnapshotSpec) ([]oscore.VirtualMachineSnapshot, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	snapshots, err := virtClient.VirtualMachineSnapshot(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return snapshots.Items, err
}

func (r *virtVM) DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSnapshotSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachineSnapshot(namespace).Delete(ctx, name, opts)
}

// VirtualMachineRestore
func (r *virtVM) CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMRestoreSpec) (*oscore.VirtualMachineRestore, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	vm := &snapshotv1.VirtualMachineRestore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		vm.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachineRestore(namespace).Create(ctx, vm, opts)
}

func (r *virtVM) GetVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineRestore, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return virtClient.VirtualMachineRestore(namespace).Get(ctx, name, opts)
}

func (r *virtVM) ListVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMSnapshotSpec) ([]oscore.VirtualMachineRestore, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	restores, err := virtClient.VirtualMachineRestore(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return restores.Items, err
}

func (r *virtVM) DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMRestoreSpec) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachineRestore(namespace).Delete(ctx, name, opts)
}

func (r *virtVM) GetVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineInstance, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}
	virtClient.CdiClient().Discovery()
	opts := metav1.GetOptions{}
	return virtClient.VirtualMachineInstance(namespace).Get(ctx, name, opts)
}

func (r *virtVM) ListVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) ([]oscore.VirtualMachineInstance, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	vms, err := virtClient.VirtualMachineInstance(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return vms.Items, nil
}

func (r *virtVM) UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMISpec) (*oscore.VirtualMachineInstance, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	vmi := &virtv1.VirtualMachineInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		vmi.Spec = *spec
	}

	opts := metav1.UpdateOptions{}
	return virtClient.VirtualMachineInstance(namespace).Update(ctx, vmi, opts)
}

func (r *virtVM) DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return virtClient.VirtualMachineInstance(namespace).Delete(ctx, name, opts)
}

func (r *virtVM) MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *oscore.VMIMigrateSpec) (*oscore.VirtualMachineInstanceMigration, error) {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return nil, err
	}

	vmiMigration := &virtv1.VirtualMachineInstanceMigration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if spec != nil {
		vmiMigration.Spec = *spec
	}

	opts := metav1.CreateOptions{}
	return virtClient.VirtualMachineInstanceMigration(namespace).Create(ctx, vmiMigration, opts)
}

func (r *virtVM) PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.PauseOptions{}
	return virtClient.VirtualMachineInstance(namespace).Pause(ctx, name, opts)
}

func (r *virtVM) UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error {
	virtClient, err := r.kubevirt.virtClient(config)
	if err != nil {
		return err
	}

	opts := &virtv1.UnpauseOptions{}
	return virtClient.VirtualMachineInstance(namespace).Unpause(ctx, name, opts)
}
