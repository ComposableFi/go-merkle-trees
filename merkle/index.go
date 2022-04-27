package merkle

// siblingIndecies returns indecies of sibling elements
func siblingIndecies(idxs []uint64) []uint64 {
	indicesCount := len(idxs)
	siblings := make([]uint64, indicesCount)
	for i := 0; i < indicesCount; i++ {
		siblings[i] = siblingIndex(idxs[i])
	}
	return siblings
}

// parentIndecies returns indecies of parent elements
func parentIndecies(idxs []uint64) []uint64 {
	indicesCount := len(idxs)
	var parents []uint64
	for i := 0; i < indicesCount; i++ {
		parentIndex := parentIndex(idxs[i])
		isDuplicate := false
		for j := 0; j < len(parents); j++ {
			if parentIndex == parents[j] {
				isDuplicate = true
			}
		}
		if !isDuplicate {
			parents = append(parents, parentIndex)
		}

	}
	return parents
}

// siblingIndex returns index of a sibling element
func siblingIndex(idx uint64) uint64 {
	return idx ^ 1
}

// isEvenIndex returns true if the index is even
// ex. idx = 100 then the bitwise operation idx&1 return 0
func isEvenIndex(idx uint64) bool {
	return idx&1 == 0
}

// getLeftIndex returns the left node index using bitwise operation
// the left index is multiply by 2
func getLeftIndex(idx int) int {
	return idx << 1
}

// getRightIndex returns right node index, this is next to the left index
func getRightIndex(idx int) int {
	return getLeftIndex(idx) + 1
}

// parentIndex returns index of a parent element
func parentIndex(idx uint64) uint64 {
	return (idx ^ 1) / 2
}

// extractNewIndicesFromSiblings finds the sibling indices which is not present in leaf indices
func extractNewIndicesFromSiblings(siblingIndices []uint64, leafIndices []uint64) []uint64 {
	return sliceDifferences(siblingIndices, leafIndices)
}
