package helpers

// SiblingIndex returns index of a sibling element
func SiblingIndex(idx uint32) uint32 {
	if isLeftIndex(idx) {
		// Right sibling index
		return idx + 1
	}
	// Left sibling index
	return idx - 1
}

func isLeftIndex(idx uint32) bool {
	return idx%2 == 0
}

// SiblingIndecies returns indecirs of sibling elements
func SiblingIndecies(idxs []uint32) []uint32 {
	var siblings []uint32
	for _, i := range idxs {
		siblings = append(siblings, SiblingIndex(i))
	}
	return siblings
}

// ParentIndex returns index of a parent element
func ParentIndex(idx uint32) uint32 {
	if isLeftIndex(idx) {
		return idx / 2
	}
	return SiblingIndex(idx) / 2
}

// ParentIndecies returns indecirs of parent elements
func ParentIndecies(idxs []uint32) []uint32 {
	var parents []uint32
	for _, i := range idxs {
		parents = append(parents, ParentIndex(i))
	}
	parents = removeDuplicateIndex(parents)
	return parents
}

// removeDuplicateIndex removes all duplicate values from uint32 slice of indices
func removeDuplicateIndex(strSlice []uint32) []uint32 {
	allKeys := make(map[uint32]bool)
	list := []uint32{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
