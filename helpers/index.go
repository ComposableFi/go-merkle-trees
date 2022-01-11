package helpers

// GetSiblingIndex returns index of a sibling element
func GetSiblingIndex(idx uint32) uint32 {
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

// GetSiblingIndecies returns indecirs of sibling elements
func GetSiblingIndecies(idxs []uint32) []uint32 {
	var siblings []uint32
	for _, i := range idxs {
		siblings = append(siblings, GetSiblingIndex(i))
	}
	return siblings
}

// GetParentIndex returns index of a parent element
func GetParentIndex(idx uint32) uint32 {
	if isLeftIndex(idx) {
		return idx / 2
	}
	return GetSiblingIndex(idx) / 2
}

// GetParentIndecies returns indecirs of parent elements
func GetParentIndecies(idxs []uint32) []uint32 {
	var parents []uint32
	for _, i := range idxs {
		parents = append(parents, GetParentIndex(i))
	}
	parents = removeDuplicateIndex(parents)
	return parents
}

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
