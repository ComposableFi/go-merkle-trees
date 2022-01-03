package merkle

import (
	"bytes"
	"log"
	"sort"
)

// MapIndiceAndLeaves maps the indices and leaves of a tree
func MapIndiceAndLeaves(indices []uint32, leaves [][]byte) (result []LeafData) {
	for _, idx := range indices {
		leaf := leaves[idx]
		result = append(result, LeafData{Index: idx, Leaf: leaf})
	}
	return result
}

func SortIndicesAndLeavesByLeafData(li []LeafData) {
	sort.Slice(li, func(i, j int) bool {
		switch bytes.Compare(li[i].Leaf, li[j].Leaf) {
		case -1:
			return true
		case 0, 1:
			return false
		default:
			log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
			return false
		}
	})
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
