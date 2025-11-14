package vmi

import (
	"context"
)

type VirtualMachineInstanceMigrationRepo interface {
	Migrate(ctx context.Context, scope, namespace, name, hostname string) error
}

func (uc *VirtualMachineInstanceUseCase) MigrateInstance(ctx context.Context, scope, namespace, name, hostname string) error {
	return uc.virtualMachineInstanceMigration.Migrate(ctx, scope, namespace, name, hostname)
}
