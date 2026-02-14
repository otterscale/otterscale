package chisel

import (
	"fmt"
	"hash/fnv"
)

// addressAllocator manages a pool of unique loopback addresses in the
// 127.x.x.x range. Each cluster is assigned a distinct address so
// that chisel can route reverse-tunnel traffic without port conflicts.
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

// allocate picks a unique loopback address for the given cluster by
// hashing the name and probing linearly until an unused address is
// found.
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
	return "", fmt.Errorf("exhausted loopback address space (%d hosts)", maxHosts)
}

// release returns a previously allocated host to the pool.
func (a *addressAllocator) release(host string) {
	delete(a.usedHosts, host)
}

// hashKey returns a deterministic 32-bit hash of the given key using
// FNV-1a so that the same cluster name tends to land on the same
// starting index.
func hashKey(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

// hostFromIndex maps a linear index (0 – maxHosts-1) to a unique
// loopback address in the range 127.1.1.1 – 127.254.254.254.
// Octets 0 and 255 are avoided to stay clear of network/broadcast
// conventions.
func hostFromIndex(idx uint32) string {
	a := idx / (254 * 254)
	b := (idx / 254) % 254
	c := idx % 254
	return fmt.Sprintf("127.%d.%d.%d", a+1, b+1, c+1)
}
