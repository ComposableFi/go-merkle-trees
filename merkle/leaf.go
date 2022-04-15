package merkle

import "github.com/ComposableFi/go-merkle-trees/types"

// popFromLeafQueue pops last element in a types.Leaf hash slice
func popFromLeafQueue(slice [][]types.Leaf) ([]types.Leaf, [][]types.Leaf) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// popFromPartialtree pops last element in a partial tree slice
func popFromPartialtree(slice []PartialTree) (PartialTree, []PartialTree) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// mapIndiceAndLeaves maps the indices and leaves of a tree
func mapIndiceAndLeaves(indices []uint64, leaves [][]byte) (result []types.Leaf) {
	for i, idx := range indices {
		leaf := leaves[i]
		result = append(result, types.Leaf{Index: idx, Hash: leaf})
	}
	return result
}
