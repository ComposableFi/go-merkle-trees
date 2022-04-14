package helpers_test

import (
	"testing"

	"github.com/ComposableFi/go-merkle-trees/helpers"
)

func BenchmarkDifference(b *testing.B) {
	s1 := []uint64{11, 15, 84, 88888888, 999999999999999}
	s2 := []uint64{11, 15, 1333, 7777777777, 999999999999999}
	for n := 0; n < b.N; n++ {
		helpers.Difference(s1, s2)
	}
}

func BenchmarkUint32SliceContains(b *testing.B) {
	s := []uint64{11, 15, 84, 88888888, 999999999999999}
	for n := 0; n < b.N; n++ {
		helpers.Uint32SliceContains(s, 15)
	}
}
