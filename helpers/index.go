package helpers

// GetSibling returns index of a sibling element
func GetSibling(idx uint64) uint64 {
	if idx == 0 {
		return 0
	}
	return ((idx + 1) ^ 1) - 1
}

// GetParent returns index of a parent element
func GetParent(idx uint64) uint64 {
	if idx == 0 {
		return 0
	}
	return (idx - 1) >> 1
}

// IsLeft checks if an index is a left index
func IsLeft(idx uint64) bool {
	return idx&1 == 1
}
