// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"encoding/binary"
	"io"
	"iter"
	"slices"

	"github.com/xmidt-org/medley/internal"
)

// Object is a hashable sequence of bytes. Typically, an Object represents the
// hashable part of an arbitrary type. Using an Object allows the hashable bytes to be
// determined once and then reused.
type Object struct {
	b []byte
}

// Len returns the number of bytes this in this Object.
func (obj Object) Len() int {
	return len(obj.b)
}

// Append appends this object's bytes to the supplied buffer,
// and returns the resulting slice.
func (obj Object) Append(buf []byte) []byte {
	buf = slices.Grow(buf, len(obj.b))
	return append(buf, obj.b...)
}

// ToHash writes this Object's contents to the given writer.
// The writer must not return errors. This method is appropriate
// when writing to an in-memory buffer or to a hash.Hash.
func (obj Object) ToHash(dst io.Writer) {
	dst.Write(obj.b)
}

// WriteTo writes this Object's contents to the given writer.
// This method allows an Object to be used as an io.WriterTo
// and is appropriate when writing to something external, such
// as a file.
func (obj Object) WriteTo(dst io.Writer) (int64, error) {
	c, err := dst.Write(obj.b)
	return int64(c), err
}

// Bytes returns an Object which contains the given bytes. The caller
// must not mutate b.
func Bytes(b []byte) Object {
	if b == nil {
		b = []byte{}
	}

	return Object{b: b}
}

// String returns an Object with the given string's bytes. This function
// does not reallocate memory for the string's contents.
func String(v string) Object {
	return Object{
		b: internal.UnsafeBytes(v),
	}
}

// Integer creates a hashable Object for a given integral type using the given byte order.
// Only certain integral types are directly supported, but callers may convert other integral
// types as needed before calling this function.
func Integer[U uint16 | uint32 | uint64](v U, order binary.ByteOrder) (obj Object) {
	switch i := any(v).(type) {
	case uint16:
		obj.b = make([]byte, 2)
		order.PutUint16(obj.b, i)

	case uint32:
		obj.b = make([]byte, 4)
		order.PutUint32(obj.b, i)

	case uint64:
		obj.b = make([]byte, 8)
		order.PutUint64(obj.b, i)
	}

	return
}

// Objecter is a closure that can produce a hashable Object from an arbitrary value.
type Objecter[V any] func(V) Object

// Objectify transforms a sequence of values into a sequence of (Object, value) tuples.
// The given Objecter is used to produce each value's corresponding Object.
//
// This function is primarily useful when you have a sequence of values that need some
// form of hashing, and you do not want that API to have a compile-time dependence on medley.
func Objectify[V any](o Objecter[V], values iter.Seq[V]) iter.Seq2[Object, V] {
	return func(yield func(Object, V) bool) {
		for value := range values {
			if !yield(o(value), value) {
				return
			}
		}
	}
}
