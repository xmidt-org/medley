// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package internal

import "unsafe"

// UnsafeBytes uses the unsafe package to obtain a string's bytes.
// Callers must not mutate the returned byte slice.
//
// If v is empty, a non-nil, empty slice is returned.
func UnsafeBytes(v string) []byte {
	// The unsafe package has some odd quirks when dealing with 0-length strings,
	// so only use unsafe for non-empty strings.
	if len(v) > 0 {
		return unsafe.Slice(unsafe.StringData(v), len(v))
	}

	return []byte{}
}
