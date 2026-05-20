// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"slices"
	"sort"
)

// hashNode is a single node in a consistent hash ring.
type hashNode[V any] struct {
	token uint64
	value V
}

// hashNodes is the raw storage for a hash ring.
// It is sortable by token in ascending order.
type hashNodes[V any] []hashNode[V]

// allocateHashNodes returns an empty hashNodes with enough space preallocated
// to hold a hash ring of the given size.
func allocateHashNodes[V any](elements int) hashNodes[V] {
	return make(hashNodes[V], 0, elements)
}

// Len returns the length of this raw hashNodes slice. This method is
// required to implement sort.Interface.
func (hns hashNodes[V]) Len() int {
	return len(hns)
}

// Less compares two elements of this hashNodes slice, returning true if
// the ith element is less than the jth element. Non-nil elements are considered
// less than nil elements, so that nil elements are sorted at the end of the slice.
//
// This method is required to implement sort.Interface.
func (hns hashNodes[V]) Less(i, j int) bool {
	return hns[i].token < hns[j].token
}

// Swap swaps two nodes. This method is required to implement sort.Interface.
func (hns hashNodes[V]) Swap(i, j int) {
	hns[i], hns[j] = hns[j], hns[i]
}

// sort sorts this slice in ascending order by token.
func (hns hashNodes[V]) sort() {
	sort.Sort(hns)
}

// nearest returns the value nearest the given object. If this slice
// is empty, this method panics. Callers should check the length first
// before calling this method.
func (hns hashNodes[V]) nearest(object uint64) (value V) {
	pos := sort.Search(len(hns), func(i int) bool {
		return hns[i].token >= object
	})

	if pos < len(hns) {
		value = hns[pos].value
	} else {
		value = hns[0].value
	}

	return
}

// grow allocates enough space for n more additions, and returns the
// possibly increased slice.
func (hns hashNodes[V]) grow(n int) hashNodes[V] {
	return slices.Grow(hns, n)
}

// append adds a single node, and returns the appended slice. After this method is used,
// sort() must be called before doing searches.
func (hns hashNodes[V]) append(token uint64, value V) hashNodes[V] {
	hns = append(hns, hashNode[V]{
		token: token,
		value: value,
	})

	return hns
}

// clear zeroes out every node. This can relieve some gc pressure for large slices.
func (hns hashNodes[V]) clear() {
	clear(hns)
}
