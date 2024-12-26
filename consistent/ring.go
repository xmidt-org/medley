package consistent

import (
	"sort"

	"github.com/xmidt-org/medley"
)

// Ring is a hash circle that distributes services randomly
// along a circle. A Ring should be created through a Builder.
//
// A Ring is a valid medley.Locator, and can be used to find services
// based on a consistent hash.
//
// The hash values for services are the same as github.com/billhathaway/consistentHash
// for backward compatibility.
//
// Rings are immutable once created. To handle an updated set of services,
// use the Update function.
type Ring[S medley.Service] struct {
	hasher hasher[S]

	// cache holds each individual service's nodes.  This is used
	// primarly to quickly rehash a ring, since we don't need to spend
	// compute computing tokens that we've already computed.
	cache medley.Map[S, nodes[S]]

	// nodes is the ring's storage
	nodes nodes[S]
}

// Find performs a hash on the given object and returns the nearest
// service. If this ring is empty, this method returns medley.ErrNoServices.
func (r *Ring[S]) Find(object []byte) (svc S, err error) {
	if len(r.nodes) > 0 {
		node := r.nearest(
			r.hasher.sum64(object),
		)

		svc = node.service
	} else {
		err = medley.ErrNoServices
	}

	return
}

// nearest returns the nearest node to the target hash value.
func (r *Ring[S]) nearest(target uint64) *node[S] {
	i := sort.Search(
		r.nodes.Len(),
		func(p int) bool {
			return r.nodes[p].token >= target
		},
	)

	if i >= r.nodes.Len() {
		i = 0
	}

	return r.nodes[i]
}

// Update checks if a set of services constitutes an update to the given Ring.
//
// If the services slice is the same as the services already hashed by the current Ring, then
// the current Ring is returned as is along with false to indicate that no update was necessary.
//
// If the services slice does represent an update, a new, distinct Ring is created with the same
// hash configuration but containing the given services. If any services in the slice were already
// hashed by the Ring, those points on the circle are copied over in order to reduce the total
// time spent hashing. This method returns true in this case, to indicate that an update was
// necessary.
//
// The current Ring is not modified by this function.
func Update[S medley.Service](current *Ring[S], services ...S) (next *Ring[S], updated bool) {
	var (
		cache                   = make(medley.Map[S, nodes[S]], len(services))
		nodes                   = make(nodes[S], 0, current.hasher.ringSize(len(services)))
		newCount, existingCount int
	)

	for update := range current.cache.Update(services...) {
		if update.Exists {
			existingCount++
			cache[update.Service] = update.Value
			nodes = append(nodes, update.Value...)
		} else {
			newCount++
			snodes := current.hasher.serviceNodes(update.Service)
			cache[update.Service] = snodes
			nodes = append(nodes, snodes...)
		}
	}

	updated = (newCount > 0 || existingCount != len(current.cache))
	if updated {
		next = &Ring[S]{
			hasher: current.hasher,
			cache:  cache,
			nodes:  nodes,
		}

		sort.Sort(next.nodes)
	} else {
		next = current
	}

	return
}
