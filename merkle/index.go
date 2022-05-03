package merkle

import (
	"math"
)

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
	var lastParentSeend uint64 = math.MaxUint64
	for i := 0; i < indicesCount; i++ {
		parentIndex := parentIndex(idxs[i])
		if parentIndex == lastParentSeend {
			continue
		}
		parents = append(parents, parentIndex)
		lastParentSeend = parentIndex
	}
	return parents
}

// siblingIndex returns index of a sibling element
func siblingIndex(idx uint64) uint64 {
	return idx ^ 1
}

func isEvenIndex(idx uint64) bool {
	return idx%2 == 0
}

func getLeftIndex(idx int) int {
	return idx * 2
}

func getRightIndex(idx int) int {
	return getLeftIndex(idx) + 1
}

// parentIndex returns index of a parent element
func parentIndex(idx uint64) uint64 {
	return (idx ^ 1) / 2
}
