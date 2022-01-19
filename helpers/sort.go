package helpers

import (
	"bytes"
	"log"
	"sort"
)

// SortUint32Slice sorts slice of uint32 numbers
func SortUint32Slice(slice []uint32) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
}

// ReverseSortUint32Slice sorts slice of uint32 numbers reversely
func ReverseSortUint32Slice(slice []uint32) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] > slice[j] })
}

// ReverseSortByteArraySlice sorts slice of []byte reversely
func ReverseSortByteArraySlice(slice [][]byte) {
	sort.Slice(slice, func(i, j int) bool {
		// bytes package already implements Comparable for []byte.
		switch bytes.Compare(slice[i], slice[j]) {
		case -1:
			return false
		case 0, 1:
			return true
		default:
			log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
			return false
		}
	})
}
