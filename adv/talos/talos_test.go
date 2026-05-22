// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package talos_test

import (
	"bytes"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/go-adv/adv/internal/constants"
	"github.com/siderolabs/go-adv/adv/talos"
)

func TestMarshalUnmarshal(t *testing.T) {
	a, err := talos.NewADV(bytes.NewReader(make([]byte, talos.Size)))
	assert.Error(t, err)
	require.NotNil(t, a)

	const (
		val1 = "value1"
		val2 = "value2"
		val3 = "value3"
	)

	assert.True(t, a.SetTag(constants.Reserved1, val1))
	assert.True(t, a.SetTag(constants.Reserved2, val2))
	assert.True(t, a.SetTag(constants.Reserved3, val3))

	b, err := a.Bytes()
	require.NoError(t, err)
	assert.Len(t, b, talos.Size)

	// test recoverable corruption
	for _, c := range []struct {
		zeroOut [][2]int
	}{
		{},
		{
			zeroOut: [][2]int{
				{0, 2},
			},
		},
		{
			zeroOut: [][2]int{
				{30, 1000},
			},
		},
		{
			zeroOut: [][2]int{
				{8, 4},
				{40, 2},
			},
		},
		{
			zeroOut: [][2]int{
				{0, talos.Length},
			},
		},
	} {
		corrupted := slices.Clone(b)

		for _, z := range c.zeroOut {
			copy(corrupted[z[0]:z[0]+z[1]], make([]byte, z[1]))
		}

		a, err = talos.NewADV(bytes.NewReader(corrupted))
		require.NoError(t, err)
		require.NotNil(t, a)

		val, ok := a.ReadTag(constants.Reserved1)
		assert.True(t, ok)
		assert.Equal(t, val1, val)

		val, ok = a.ReadTag(constants.Reserved2)
		assert.True(t, ok)
		assert.Equal(t, val2, val)

		val, ok = a.ReadTag(constants.Reserved3)
		assert.True(t, ok)
		assert.Equal(t, val3, val)

		tags := a.ListTags()
		slices.Sort(tags)
		assert.Equal(t, []uint8{constants.Reserved1, constants.Reserved2, constants.Reserved3}, tags)
	}
}
