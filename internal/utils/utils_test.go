package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceContains(t *testing.T) {
	// strings
	stringSlice := []string{"first", "second", "third"}
	require.True(t, SliceContains(stringSlice, "second"))
	require.False(t, SliceContains(stringSlice, "fourth"))

	// floats
	floatSlice := []float64{1.1, 1.2, 1.3}
	require.True(t, SliceContains(floatSlice, 1.2))
	require.False(t, SliceContains(floatSlice, 1.4))

	// structs
	type temp struct {
		test int
	}
	structSlice := []temp{{test: 1}, {test: 2}, {test: 3}}
	require.True(t, SliceContains(structSlice, temp{test: 2}))
	require.False(t, SliceContains(structSlice, temp{test: 4}))

	// pointers
	ptrSlice := []*temp{&structSlice[0], &structSlice[1], &structSlice[2]}
	otherPtr := &temp{test: 1}
	require.True(t, SliceContains(ptrSlice, &structSlice[0]))
	require.False(t, SliceContains(ptrSlice, otherPtr))
}

func TestHashString(t *testing.T) {
	require.True(t, len(HashString("first")) > 0)

	// confirm that same string hashes the same
	require.Equal(t, HashString("first"), HashString("first"))

	// confirm that different strings hash differently.
	// this is an extreme over-simplification and does not
	// address hash collisions at all.
	require.NotEqual(t, HashString("first"), HashString("second"))
}
