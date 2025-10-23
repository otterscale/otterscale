package core

import (
	"context"
	"errors"
	"net"
	"slices"
	"strconv"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity"
)

func (uc *OrchestratorUseCase) reserveIP(ctx context.Context, machineID, comment string) (net.IP, error) {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}
	links := machine.BootInterface.Links
	if len(links) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("machine has no network links"))
	}
	subnet := &links[0].Subnet
	ip, err := uc.getFreeIP(ctx, subnet)
	if err != nil {
		return nil, err
	}
	if _, err := uc.createIPRange(ctx, subnet.ID, ip.String(), comment); err != nil {
		return nil, err
	}
	return ip, nil
}

func (uc *OrchestratorUseCase) getFreeIP(ctx context.Context, subnet *entity.Subnet) (net.IP, error) {
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

func (uc *OrchestratorUseCase) getUsedIPs(ctx context.Context, subnetID int) ([]uint32, error) {
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

func (uc *OrchestratorUseCase) getReservedIPs(ctx context.Context, cidr string) ([]uint32, error) {
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

func (uc *OrchestratorUseCase) createIPRange(ctx context.Context, subnetID int, ip, comment string) (*IPRange, error) {
	params := &entity.IPRangeParams{
		Type:    "reserved",
		Subnet:  strconv.Itoa(subnetID),
		StartIP: ip,
		EndIP:   ip,
		Comment: comment,
	}
	return uc.ipRange.Create(ctx, params)
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
