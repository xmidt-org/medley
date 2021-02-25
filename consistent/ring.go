package consistent

import (
	"errors"
	"sort"

	"github.com/xmidt-org/medley"
)

var (
	// ErrEmpty indicates that a search operation was attempted over an empty hash
	ErrEmpty = errors.New("no hash nodes defined")
)

// hashEntry is a tuple of hash information about a Node.  Hash entries
// are typically assigned points on the unit circle by placing them into
// a ring and sorting the ring.
type hashEntry struct {
	Node  medley.Node
	Value uint64
}

// ring is a circle of hash entries, ordered by hash value and node comparison.
// A ring implements sort.Interface.
type ring []hashEntry

// Len returns the number of hash entries in this ring
func (r ring) Len() int {
	return len(r)
}

// Less compares two hash entries in this ring.  The hash Values are compared
// first.  If there is a collision, i.e. the two entries have the same Value,
// then the Nodes are compared.
func (r ring) Less(i, j int) bool {
	ith, jth := r[i], r[j]
	if ith.Value < jth.Value {
		return true
	} else if ith.Value == jth.Value {
		return ith.Node < jth.Node
	}

	return false
}

// Swap swaps two hash entries in this ring
func (r ring) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// grow ensures that ring has enough capacity for the given additional
// hash entries.  Use this method prior to adding multiple hash entries
// to reduce the number of allocations.
//
// This method uses a similar algorithm to strings.Builder.  If more capacity
// is needed, this method allocates even more capacity to keep the total number
// of memory allocations as small as possible.
func (r *ring) grow(n int) {
	if cap(*r)-len(*r) < n {
		larger := make([]hashEntry, len(*r), 2*cap(*r)+n)
		copy(larger, *r)
		*r = larger
	}
}

// add appends a new hash entry.  This method does not re-sort the ring.
func (r *ring) add(n medley.Node, v uint64) {
	*r = append(*r, hashEntry{Node: n, Value: v})
}

// removeIf deletes all hash entries whose node matches the given predicate.
// This method does not re-sort the ring.
//
// medley.NodeSet.Has can be passed to this method directly.
func (r *ring) removeIf(p func(medley.Node) bool) {
	end := len(*r) - 1
	var pos int
	for pos <= end {
		if p((*r)[pos].Node) {
			(*r)[pos], (*r)[end] = (*r)[end], (*r)[pos]
			(*r)[end] = hashEntry{}
			end--
		} else {
			pos++
		}
	}

	*r = (*r)[0 : end+1]
}

// closest returns the Node closest to the given hash value by walking clockwise
// from the given value v to the next node, wrapping around as necessary.
//
// This method returns an error if this ring is empty.
//
// See: https://pkg.go.dev/sort#Search
func (r ring) closest(v uint64) (n medley.Node, err error) {
	if len(r) > 0 {
		p := sort.Search(len(r),
			func(i int) bool {
				return r[i].Value >= v
			},
		)

		if p < len(r) {
			n = r[p].Node
		} else {
			n = r[0].Node
		}
	} else {
		err = ErrEmpty
	}

	return
}

// sort reorders the hash entries in this ring so that Search can do
// a binary search.  This method should be called after any sequence of
// operations that could modify the contents of the ring.
//
// This is not a stable sort, though in practice that doesn't matter due
// to the order property defined for HashEntry.
//
// See: https://pkg.go.dev/sort#Sort
func (r ring) sort() {
	sort.Sort(r)
}
