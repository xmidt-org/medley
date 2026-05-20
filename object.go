// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"iter"
)

// Object represents a value that can be hashed. Typically, this will be an
// identifier of some sort representing an arbitrary value.
type Object interface {
	~[]byte | ~string
}

// Append appends an Object to a slice of bytes.
func Append[O Object](dst []byte, v O) []byte {
	return append(dst, v...)
}

// Objecter is a strategy for producing a hashable Object from an arbitrary value.
type Objecter[O Object, V any] func(V) O

// Objectify transforms a sequence of arbitrary values into a sequence of (Object, Value)
// tuples. No deduplication is done, so if there are duplicate values they are preserved
// as is.
func Objectify[O Object, V any](f Objecter[O, V], values iter.Seq[V]) iter.Seq2[O, V] {
	return func(yield func(O, V) bool) {
		for value := range values {
			obj := f(value)
			if !yield(obj, value) {
				return
			}
		}
	}
}

// Stringify handles the common case where a sequence of strings is used as their own
// hashable objects. This method is a handy convenience when hashing things like host names.
func Stringify[V ~string](values iter.Seq[V]) iter.Seq2[string, V] {
	return Objectify(
		func(v V) string { return string(v) },
		values,
	)
}
