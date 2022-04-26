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

func isEvenIndex(idx uint64) bool {
	return idx&1 != 1
}

func getLeftIndex(idx int) int {
	return idx << 1
}

func getRightIndex(idx int) int {
	return getLeftIndex(idx) + 1
}

// parentIndex returns index of a parent element
func parentIndex(idx uint64) uint64 {
	return (idx ^ 1) / 2
}
