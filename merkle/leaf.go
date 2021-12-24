package merkle

import (
	"sort"
)

// MapIndiceAndLeaves maps the indices and leaves of a tree
func MapIndiceAndLeaves(indices []uint32, leaves [][]byte) (result []LeafData) {
	for i, idx := range indices {
		leaf := leaves[i]
		result = append(result, LeafData{Index: idx, Leaf: leaf})
	}
	return result
}

// SortIndicesAndLeavesByIndex sorts the leaf index slice reversely by index
func SortIndicesAndLeavesByIndex(li []LeafData) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index < li[j].Index })
}

// SortIndicesAndLeavesByIndexReversely sorts the leaf index slice reversely by index
func SortIndicesAndLeavesByIndexReversely(li []LeafData) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index > li[j].Index })
}

// PopFromLeafIndexQueue pops first front element in a leaf index slice
func PopFromLeafIndexQueue(slice []LeafData) (LeafData, []LeafData) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice
}
