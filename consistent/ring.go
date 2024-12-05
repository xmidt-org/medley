package consistent

import (
	"sort"

	"github.com/xmidt-org/medley"
)

// node is a single hash node for a service.
type node[S medley.Service] struct {
	token   uint64
	service S
}

// nodes is the Ring's primary storage.
type nodes[S medley.Service] []*node[S]

func (ns nodes[S]) Len() int {
	return len(ns)
}

func (ns nodes[S]) Less(i, j int) bool {
	return ns[i].token < ns[j].token
}

func (ns nodes[S]) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

// sort sorts this set of nodes by token.
func (ns nodes[S]) sort() {
	sort.Sort(ns)
}

// serviceNodes holds nodes calculated for each service.
type serviceNodes[S medley.Service] map[S]nodes[S]

// Len returns the count of services in this map.
func (sn serviceNodes[S]) Len() int {
	return len(sn)
}

// addAll adds precomputed nodes to this map, and also appends each
// service's nodes to slice of all nodes. the newly appended slice of
// all nodes is returned.
func (sn serviceNodes[S]) addAll(orig serviceNodes[S], all nodes[S]) nodes[S] {
	for svc, snodes := range orig {
		sn[svc] = snodes
		all = append(all, snodes...)
	}

	return all
}

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
	hasher   hasher[S]
	services serviceNodes[S]
	nodes    nodes[S]
}

// Find performs a hash on the given object and returns the nearest
// service. If this ring is empty, this method returns medley.ErrNoServices.
func (r *Ring[S]) Find(object string) (svc S, err error) {
	if len(r.nodes) > 0 {
		node := r.nearest(
			r.hasher.hashString(object),
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
		currentServices = make(serviceNodes[S])
		newServices     = make(medley.ServiceSet[S])
	)

	for _, svc := range services {
		if nodes, exists := current.services[svc]; exists {
			currentServices[svc] = nodes
		} else {
			newServices.Add(svc)
		}
	}

	updated = (newServices.Len() != 0 || currentServices.Len() != current.services.Len())
	if updated {
		newServiceCount := currentServices.Len() + newServices.Len()
		next = &Ring[S]{
			hasher:   current.hasher,
			services: make(serviceNodes[S], newServiceCount),
			nodes:    make(nodes[S], 0, current.hasher.totalCount(newServiceCount)),
		}

		// copy over the cached, previously computed nodes
		next.nodes = next.services.addAll(currentServices, next.nodes)

		// compute the necessary new nodes
		for svc := range newServices {
			snodes := next.hasher.serviceNodes(svc)
			next.services[svc] = snodes
			next.nodes = append(next.nodes, snodes...)
		}

		next.nodes.sort()
	} else {
		next = current
	}

	return
}
