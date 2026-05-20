// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeBytes(t *testing.T) {
	t.Run("Uninitialized", func(t *testing.T) {
		var test string
		assert.Empty(t, UnsafeBytes(test))
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Empty(t, UnsafeBytes(""))
	})

	t.Run("NonEmpty", func(t *testing.T) {
		assert.Equal(
			t,
			[]byte("this is a test"),
			UnsafeBytes("this is a test"),
		)
	})
}
