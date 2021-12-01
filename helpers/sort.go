package helpers

import "sort"

// SortUint64Slice sorts slice of uint64 numbers
func SortUint64Slice(slice []uint64) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
}

// ReverseSortUint64Slice sorts slice of uint64 numbers reversely
func ReverseSortUint64Slice(slice []uint64) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] > slice[j] })
}
