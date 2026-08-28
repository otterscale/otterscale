package chisel

import (
	"fmt"
	"hash/fnv"

	"github.com/otterscale/otterscale/internal/core"
)

// addressAllocator hands each cluster a distinct 127.x.x.x address, so chisel
// can route reverse-tunnel traffic without port conflicts.
//
// All methods must be called with the parent Service's mu held.
type addressAllocator struct {
	usedHosts map[string]struct{}
}

func newAddressAllocator() *addressAllocator {
	return &addressAllocator{
		usedHosts: make(map[string]struct{}),
	}
}

// allocate hashes the cluster name, then probes linearly from there.
func (a *addressAllocator) allocate(cluster string) (string, error) {
	base := hashKey(cluster)
	for i := range uint32(maxHosts) {
		candidate := hostFromIndex((base + i) % uint32(maxHosts))
		if _, exists := a.usedHosts[candidate]; exists {
			continue
		}
		a.usedHosts[candidate] = struct{}{}
		return candidate, nil
	}
	return "", &core.DomainError{
		Code:    core.ErrorCodeResourceExhausted,
		Message: fmt.Sprintf("exhausted loopback address space (%d hosts)", maxHosts),
	}
}

func (a *addressAllocator) release(host string) {
	delete(a.usedHosts, host)
}

// hashKey is FNV-1a, so the same cluster name tends to land on the same
// starting index.
func hashKey(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

// hostFromIndex maps 0 – maxHosts-1 onto 127.1.1.1 – 127.254.254.254, avoiding
// octets 0 and 255 to stay clear of network/broadcast conventions.
func hostFromIndex(idx uint32) string {
	a := idx / (254 * 254)
	b := (idx / 254) % 254
	c := idx % 254
	return fmt.Sprintf("127.%d.%d.%d", a+1, b+1, c+1)
}
