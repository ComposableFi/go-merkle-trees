package helpers

import "sort"

// SortUint32Slice sorts slice of uint32 numbers
func SortUint32Slice(slice []uint32) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
}

// ReverseSortUint32Slice sorts slice of uint32 numbers reversely
func ReverseSortUint32Slice(slice []uint32) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] > slice[j] })
}
