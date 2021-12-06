package merkle

import (
	"sort"
)

// MapIndiceAndLeaves maps the indices and leaves of a tree
func MapIndiceAndLeaves(indices []uint32, leaves []interface{}) (result []LeafIndex) {
	for i, idx := range indices {
		leaf := leaves[i]
		result = append(result, LeafIndex{Index: idx, Leaf: leaf})
	}
	return result
}

// SortIndicesAndLeavesByIndexReversely sorts the leaf index slice reversely by index
func SortIndicesAndLeavesByIndexReversely(li []LeafIndex) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index > li[j].Index })
}

// PopFromLeafIndexQueue pops first front element in a leaf index slice
func PopFromLeafIndexQueue(slice []LeafIndex) (LeafIndex, []LeafIndex) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice
}
