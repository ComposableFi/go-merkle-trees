package merkle

// PopFromLeafQueue pops last element in a leaf hash slice
func PopFromLeafQueue(slice [][]Leaf) ([]Leaf, [][]Leaf) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// PopFromPartialtree pops last element in a partial tree slice
func PopFromPartialtree(slice []PartialTree) (PartialTree, []PartialTree) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// MapIndiceAndLeaves maps the indices and leaves of a tree
func MapIndiceAndLeaves(indices []uint32, leaves [][]byte) (result []Leaf) {
	for i, idx := range indices {
		leaf := leaves[i]
		result = append(result, Leaf{Index: idx, Hash: leaf})
	}
	return result
}
