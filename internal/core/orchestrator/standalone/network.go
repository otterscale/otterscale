package standalone

import (
	"context"
	"errors"
	"net"
	"slices"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/network"
)

func (uc *StandaloneUseCase) reserveIP(ctx context.Context, machine *machine.Machine, comment string) (ip net.IP, releaseFunc func() error, err error) {
	links := machine.BootInterface.Links
	if len(links) == 0 {
		return nil, nil, connect.NewError(connect.CodeInvalidArgument, errors.New("machine has no network links"))
	}

	subnet := &links[0].Subnet

	ip, err = uc.getFreeIP(ctx, subnet)
	if err != nil {
		return nil, nil, err
	}

	ipRange, err := uc.ipRange.Create(ctx, subnet.ID, ip.String(), ip.String(), comment)
	if err != nil {
		return nil, nil, err
	}

	return ip, func() error { return uc.ipRange.Delete(ctx, ipRange.ID) }, nil
}

func (uc *StandaloneUseCase) getFreeIP(ctx context.Context, subnet *network.Subnet) (net.IP, error) {
	skip := []uint32{}

	used, err := uc.getUsedIPs(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}

	skip = append(skip, used...)

	reserved, err := uc.getReservedIPs(ctx, subnet.CIDR)
	if err != nil {
		return nil, err
	}

	skip = append(skip, reserved...)

	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
	if err != nil {
		return nil, err
	}

	ip := ipToUint32(ipNet.IP)
	mask := ipToUint32(net.IP(ipNet.Mask))
	network := ip & mask
	broadcast := network | ^mask

	next := false // get next to prevent time gap

	for i := network + 1; i < broadcast; i++ {
		if slices.Contains(skip, i) {
			continue
		}

		if next {
			return uint32ToIP(i), nil
		}

		next = true
	}

	return nil, connect.NewError(connect.CodeResourceExhausted, errors.New("no free IP found"))
}

func (uc *StandaloneUseCase) getUsedIPs(ctx context.Context, subnetID int) ([]uint32, error) {
	ips, err := uc.subnet.GetIPAddresses(ctx, subnetID)
	if err != nil {
		return nil, err
	}

	record := []uint32{}

	for i := range ips {
		record = append(record, ipToUint32(ips[i].IP))
	}

	return record, nil
}

func (uc *StandaloneUseCase) getReservedIPs(ctx context.Context, cidr string) ([]uint32, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ipRanges, err := uc.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}

	record := []uint32{}

	for i := range ipRanges {
		if ipNet.Contains(ipRanges[i].StartIP) && ipNet.Contains(ipRanges[i].EndIP) {
			start := ipToUint32(ipRanges[i].StartIP)
			end := ipToUint32(ipRanges[i].EndIP)

			for i := start; i <= end; i++ {
				record = append(record, i)
			}
		}
	}

	return record, nil
}

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IP{
		byte(n >> 24), //nolint:mnd // shift
		byte(n >> 16), //nolint:mnd // shift
		byte(n >> 8),  //nolint:mnd // shift
		byte(n),
	}
}
