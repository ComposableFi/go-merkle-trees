package merkle

import "math"

const halfDivider = 2

// siblingIndecies returns indecies of sibling elements
func siblingIndecies(leafIndices []uint64) []uint64 {

	indicesCount := len(leafIndices)
	siblings := make([]uint64, indicesCount)

	// append all sibling indices to result
	for i := 0; i < indicesCount; i++ {
		siblings[i] = siblingIndex(leafIndices[i])
	}

	return siblings
}

// parentIndecies returns indecies of parent elements
func parentIndecies(leafIndices []uint64) []uint64 {
	indicesCount := len(leafIndices)
	var parents []uint64
	var lastParentSeend uint64 = math.MaxUint64
	for i := 0; i < indicesCount; i++ {
		parentIndex := parentIndex(leafIndices[i])
		if parentIndex == lastParentSeend {
			continue
		}
		parents = append(parents, parentIndex)
		lastParentSeend = parentIndex
	}
	return parents
}

// siblingIndex returns index of a sibling element
// ex. index = 100 then the bitwise operation index ^ 1 returns 01100101
// which is the binary representation of 101 and 101 ^ 1 return 01100100
// which is the binary representation of 100
func siblingIndex(index uint64) uint64 {
	return index ^ 1
}

// isEvenIndex returns true if the index is even
// ex. index = 100 then the bitwise operation index&1 returns 0
// 1   => 00000001
// 100 => 01100100
// 101 => 01100101
// so the 100&1 returns 0
func isEvenIndex(index uint64) bool {
	return index&1 == 0
}

// getLeftIndex returns the left node index using bitwise operation
// the left index is multiply by 2
func getLeftIndex(index int) int {
	return index << 1
}

// getRightIndex returns right node index, this is next index the left index
func getRightIndex(index int) int {
	return getLeftIndex(index) + 1
}

// parentIndex returns index of a parent element
// the half of sibling index is the parent node index
func parentIndex(index uint64) uint64 {
	return siblingIndex(index) / halfDivider
}

// extractNewIndicesFromSiblings finds the sibling indices which is not present in leaf indices
func extractNewIndicesFromSiblings(siblingIndices []uint64, leafIndices []uint64) []uint64 {
	return sliceDifferences(siblingIndices, leafIndices)
}
