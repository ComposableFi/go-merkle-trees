package helpers

import "sort"

// SiblingIndex returns index of a sibling element
func SiblingIndex(idx uint64) uint64 {
	if isLeftIndex(idx) {
		// Right sibling index
		return idx + 1
	}
	// Left sibling index
	return idx - 1
}

func isLeftIndex(idx uint64) bool {
	return idx%2 == 0
}

// SiblingIndecies returns indecirs of sibling elements
func SiblingIndecies(idxs []uint64) []uint64 {
	var siblings []uint64
	for i := 0; i < len(idxs); i++ {
		siblings = append(siblings, SiblingIndex(idxs[i]))
	}
	return siblings
}

// ParentIndex returns index of a parent element
func ParentIndex(idx uint64) uint64 {
	if isLeftIndex(idx) {
		return idx / 2
	}
	return SiblingIndex(idx) / 2
}

// ParentIndecies returns indecirs of parent elements
func ParentIndecies(idxs []uint64) []uint64 {
	var parents []uint64
	for i := 0; i < len(idxs); i++ {
		parents = append(parents, ParentIndex(idxs[i]))
	}
	parents = removeDuplicateIndex(parents)
	return parents
}

// removeDuplicateIndex removes all duplicate values from uint64 slice of indices
func removeDuplicateIndex(s []uint64) []uint64 {
	// if there are 0 or 1 items we return the slice itself.
	if len(s) < 2 {
		return s
	}

	// make the slice ascending sorted.
	sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })

	uniqPointer := 0

	for i := 1; i < len(s); i++ {
		// compare a current item with the item under the unique pointer.
		// if they are not the same, write the item next to the right of the unique pointer.
		if s[uniqPointer] != s[i] {
			uniqPointer++
			s[uniqPointer] = s[i]
		}
	}

	return s[:uniqPointer+1]
}
