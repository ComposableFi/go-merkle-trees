package merkle

// siblingIndecies returns indecies of sibling elements
func siblingIndecies(idxs []uint64) []uint64 {

	indicesCount := len(idxs)
	siblings := make([]uint64, indicesCount)

	// append all sibling indices to result
	for i := 0; i < indicesCount; i++ {
		siblings[i] = siblingIndex(idxs[i])
	}

	return siblings
}

// parentIndecies returns indecies of parent elements
func parentIndecies(idxs []uint64) []uint64 {

	var parents []uint64

	// loop through all indices
	for i := 0; i < len(idxs); i++ {

		// get parent index
		parentIndex := parentIndex(idxs[i])

		// checl if it is duplicate
		isDuplicate := false
		for j := 0; j < len(parents); j++ {
			if parentIndex == parents[j] {
				isDuplicate = true
			}
		}

		// appent to result if it is not duplicated
		if !isDuplicate {
			parents = append(parents, parentIndex)
		}

	}

	return parents
}

// siblingIndex returns index of a sibling element
// ex. idx = 100 then the bitwise operation idx ^ 1 returns 01100101
// which is the binary representation of 101 and 101 ^ 1 return 01100100
// which is the binary representation of 100
func siblingIndex(idx uint64) uint64 {
	return idx ^ 1
}

// isEvenIndex returns true if the index is even
// ex. idx = 100 then the bitwise operation idx&1 returns 0
// 1   => 00000001
// 100 => 01100100
// 101 => 01100101
// so the 100&1 returns 0
func isEvenIndex(idx uint64) bool {
	return idx&1 == 0
}

// getLeftIndex returns the left node index using bitwise operation
// the left index is multiply by 2
func getLeftIndex(idx int) int {
	return idx << 1
}

// getRightIndex returns right node index, this is next index the left index
func getRightIndex(idx int) int {
	return getLeftIndex(idx) + 1
}

// parentIndex returns index of a parent element
// the half of sibling index is the parent node index
func parentIndex(idx uint64) uint64 {
	return siblingIndex(idx) / 2
}

// extractNewIndicesFromSiblings finds the sibling indices which is not present in leaf indices
func extractNewIndicesFromSiblings(siblingIndices []uint64, leafIndices []uint64) []uint64 {
	return sliceDifferences(siblingIndices, leafIndices)
}
