package merkle

import "github.com/ComposableFi/go-merkle-trees/types"

// PopFromLeafQueue pops last element in a types.Leaf hash slice
func PopFromLeafQueue(slice [][]types.Leaf) ([]types.Leaf, [][]types.Leaf) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// PopFromPartialtree pops last element in a partial tree slice
func PopFromPartialtree(slice []PartialTree) (PartialTree, []PartialTree) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// MapIndiceAndLeaves maps the indices and leaves of a tree
func MapIndiceAndLeaves(indices []uint64, leaves [][]byte) (result []types.Leaf) {
	for i, idx := range indices {
		leaf := leaves[i]
		result = append(result, types.Leaf{Index: idx, Hash: leaf})
	}
	return result
}
