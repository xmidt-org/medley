// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import "github.com/xmidt-org/medley"

// node is a single hash ring node for a service.
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
