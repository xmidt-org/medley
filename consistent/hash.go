package consistent

import (
	"sync"

	"github.com/xmidt-org/medley"
)

const (
	// DefaultVnodes is the default number of virtual nodes per node
	DefaultVnodes = 211
)

// Config represents the available configuration options for a consistent hash.
type Config struct {
	// Algorithm is the name of the algorithm to use.  If this field
	// is unset, medley.DefaultAlgorithm() is used.
	Algorithm string `json:"algorithm"`

	// Vnodes is the number of virtual nodes for each node.  If this field
	// is unset or nonpositive, DefaultVnodes is used.
	Vnodes int `json:"vnodes"`

	// Extensions is an optional set of algorithms beyond this package's builtins.
	// The Algorithm field can refer to a key within this map.
	Extensions map[string]medley.Algorithm `json:"-"`
}

// Hash represents a consistent hash.  This type is backward compatible with
// https://github.com/billhathaway/consistentHash.
//
// A Hash is safe for concurrent use.
type Hash struct {
	alg    medley.Algorithm
	vnodes int

	// updateLock ensures that only (1) update can happen at a time.
	// updates (e.g. add, rehash) do some precomputation prior to updating
	// the set of nodes, and that precomputation shouldn't stop read operations
	// from continuing.
	updateLock sync.Mutex

	// assigner is only used by code holding the update lock
	assigner *assigner

	// nodeLock only protects the internal hash storage
	nodeLock sync.RWMutex
	nodes    medley.NodeSet
	ring     ring
}

// New constructs a consistent Hash from configuration
func New(cfg Config) (h *Hash, err error) {
	var alg medley.Algorithm
	alg, err = medley.FindAlgorithm(cfg.Algorithm, cfg.Extensions)
	if err == nil {
		if cfg.Vnodes < 1 {
			cfg.Vnodes = DefaultVnodes
		}

		h = &Hash{
			alg:      alg,
			vnodes:   cfg.Vnodes,
			assigner: newAssigner(alg),
		}
	}

	return
}

// Len returns the count of nodes associated with this hash.  This is
// not the size of the internal storage of the hash.  Len()*Vnodes()
// would return the total size of the internal hash ring.
func (h *Hash) Len() (l int) {
	h.nodeLock.RLock()
	l = h.nodes.Len()
	h.nodeLock.RUnlock()

	return
}

// Algorithm returns the medley algorithm associated with this hash
func (h *Hash) Algorithm() medley.Algorithm {
	return h.alg
}

// Vnodes returns the number of virtual nodes added for each node
func (h *Hash) Vnodes() int {
	return h.vnodes
}

// Get obtains the closest Node associated with a Key.  The hash ring is
// walked clockwise to find the nearest node.
func (h *Hash) Get(k medley.Key) (n medley.Node, err error) {
	h.nodeLock.RLock()
	defer h.nodeLock.RUnlock()

	if h.ring.Len() > 0 {
		hasher := h.alg()
		_, err = k.WriteTo(hasher)

		if err == nil {
			n, err = h.ring.closest(hasher.Sum64())
		}
	} else {
		err = ErrEmpty
	}

	return
}

// Add inserts nodes and their corresponding vnodes into this hash.
// Any nodes already present are left intact.
//
// This method reorders the nodes slice in-place using NodeSet.Filter.
func (h *Hash) Add(nodes []medley.Node) (added int) {
	h.updateLock.Lock()
	defer h.updateLock.Unlock()

	// precomputation only needs the read lock
	h.nodeLock.RLock()
	_, notIn := h.nodes.Filter(nodes)
	h.nodeLock.RUnlock()

	added = len(notIn)
	if added == 0 {
		// optimization: if nothing to do, return
		return
	}

	// now acquire the write lock
	h.nodeLock.Lock()
	defer h.nodeLock.Unlock()

	h.ring.grow(h.vnodes * added)
	for _, n := range notIn {
		h.assigner.reset(n)
		for r := 0; r < h.vnodes; r++ {
			h.ring.add(n, h.assigner.next())
		}
	}

	h.nodes.AddAll(notIn...)
	h.ring.sort()
	return
}

// Remove deletes nodes and their vnodes from this hash.
//
// This method reorders the nodes slice in-place using NodeSet.Filter.
func (h *Hash) Remove(nodes []medley.Node) (removed int) {
	h.updateLock.Lock()
	defer h.updateLock.Unlock()

	// precomputation only requires the read lock
	h.nodeLock.RLock()
	in, _ := h.nodes.Filter(nodes)
	h.nodeLock.RUnlock()

	removed = len(in)
	if removed == 0 {
		// optimization: if nothing to do, return
		return
	}

	toRemove := medley.NewNodeSet(in...)

	// now acquire the write lock to modify the set of nodes
	h.nodeLock.Lock()
	defer h.nodeLock.Unlock()

	h.nodes.RemoveAll(in...)
	h.ring.removeIf(toRemove.Has)
	h.ring.sort()
	return
}

// Rehash restructures this hash so that the given nodes are the
// only ones in the ring.  Any nodes currently part of this hash that
// are not in the set of nodes passed to this method are removed.
// Nodes passed to this method that are not in this hash are added.
//
// Unlike Add and Remove, this method does not reorder the nodes slice.
//
// The separate counts of nodes added and removed are returned.
func (h *Hash) Rehash(nodes []medley.Node) (added, removed int) {
	h.updateLock.Lock()
	defer h.updateLock.Unlock()

	// precomputation only requires the read lock
	h.nodeLock.RLock()
	var (
		// rehash will become our new node set
		rehash   = medley.NewNodeSet(nodes...)
		toAdd    medley.NodeSet
		toRemove medley.NodeSet
	)

	for n := range h.nodes {
		if !rehash.Has(n) {
			toRemove.Add(n)
		}
	}

	for n := range rehash {
		if !h.nodes.Has(n) {
			toAdd.Add(n)
		}
	}

	h.nodeLock.RUnlock()

	added = toAdd.Len()
	removed = toRemove.Len()
	if removed == 0 && added == 0 {
		// optimization: if nothing to do, return
		return
	}

	// now acquire the write lock to modify the set of nodes
	h.nodeLock.Lock()
	defer h.nodeLock.Unlock()

	if removed > 0 {
		h.ring.removeIf(toRemove.Has)
	}

	if added > 0 {
		h.ring.grow(h.vnodes * added)
		for n := range toAdd {
			h.assigner.reset(n)
			for r := 0; r < h.vnodes; r++ {
				h.ring.add(n, h.assigner.next())
			}
		}
	}

	h.nodes = rehash
	h.ring.sort()
	return
}
