package configuration

import (
	"context"
	"strings"
)

type ProvisionerRepo interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

func (uc *ConfigurationUseCase) UpdateNTPServer(ctx context.Context, addresses []string) ([]string, error) {
	if err := uc.provisioner.Update(ctx, "ntp_servers", strings.Join(addresses, " ")); err != nil {
		return nil, err
	}

	return uc.listNTPServers(ctx)
}

func (uc *ConfigurationUseCase) SetDefaultBootImage(ctx context.Context, distroSeries string) error {
	if err := uc.provisioner.Update(ctx, "default_osystem", "ubuntu"); err != nil {
		return err
	}

	if err := uc.provisioner.Update(ctx, "default_distro_series", distroSeries); err != nil {
		return err
	}

	return uc.provisioner.Update(ctx, "commissioning_distro_series", distroSeries)
}

func (uc *ConfigurationUseCase) listNTPServers(ctx context.Context) ([]string, error) {
	ntpServers, err := uc.provisioner.Get(ctx, "ntp_servers")
	if err != nil {
		return nil, err
	}

	return strings.Split(ntpServers, " "), nil
}
