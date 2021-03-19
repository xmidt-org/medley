package medley

import (
	"io"
)

// NilNode is the canonicalized nil value for nodes
var NilNode Node = Node("")

// Node is a string type that is assigned hash values in a ring.
// A node can be any string, but is most often a URL or host name.
type Node string

// WriteTo fulfills the io.WriterTo interface and allows nodes to
// write their own hash bytes to a writer.
func (n Node) WriteTo(w io.Writer) (int64, error) {
	c, err := io.WriteString(w, string(n))
	return int64(c), err
}

// NodeSet is a set of nodes.  The zero value of this type is
// ready to use.
//
// The typical use case for a NodeSet is to track which nodes are
// part of a hash.  Without a set, the hash would likely need to do
// binary or even linear searches to determine if a node was already hashed.
type NodeSet map[Node]bool

// NewNodeSet constructs a NodeSet from a slice of nodes
func NewNodeSet(nodes ...Node) NodeSet {
	var ns NodeSet
	for _, n := range nodes {
		ns.Add(n)
	}

	return ns
}

// Len returns the count of nodes in this set
func (ns NodeSet) Len() int {
	return len(ns)
}

// Has tests if the given node is present in this set
func (ns NodeSet) Has(n Node) bool {
	return ns[n]
}

// Add inserts the given node into this set.  This set
// is initialized as needed.  This method returns true
// to indicate that the node was added, false to indicate
// that the node was already present.
func (ns *NodeSet) Add(n Node) (added bool) {
	if *ns != nil {
		added = !(*ns)[n]
	} else {
		*ns = make(NodeSet)
		added = true
	}

	if added {
		(*ns)[n] = true
	}

	return
}

// AddAll adds each of a sequence of nodes.  Any nodes already in
// this set are not modified.  This method returns the count of
// nodes actually added.
//
// As with Add, this set is initialized as needed.
func (ns *NodeSet) AddAll(nodes ...Node) (count int) {
	if *ns == nil && len(nodes) > 0 {
		*ns = make(NodeSet)
	}

	for _, n := range nodes {
		if !(*ns)[n] {
			count++
			(*ns)[n] = true
		}
	}

	return
}

// Remove deletes a node from this set.  This method returns
// true if the node was deleted, false to indicate the node
// did not exist in this set.
func (ns *NodeSet) Remove(n Node) (removed bool) {
	if *ns != nil {
		removed = (*ns)[n]
		delete(*ns, n)
	}

	return
}

// RemoveAll deletes each of a sequence of nodes.  Any node not
// present in this set is ignored.  This method returns the count
// of nodes actually deleted.
func (ns *NodeSet) RemoveAll(nodes ...Node) (count int) {
	if *ns != nil {
		for _, n := range nodes {
			if (*ns)[n] {
				count++
				delete(*ns, n)
			}
		}
	}

	return
}

// Filter examines each node in a slice to determine if it is present in this
// set.  The slice is rearranged in-place so that nodes that are in this set
// are contiguous in the first portion of the slice, while all the nodes not in
// this set follow.  The two slices that are returned point into the given slice
// and contain the nodes in and not in this set, respectively.
func (ns NodeSet) Filter(nodes []Node) (in, notIn []Node) {
	i := 0
	j := len(nodes) - 1
	for i <= j {
		if ns[nodes[i]] {
			i++
			continue
		}

		nodes[i], nodes[j] = nodes[j], nodes[i]
		j--
	}

	in = nodes[0:i]
	notIn = nodes[i:]
	return
}
