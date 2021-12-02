package helpers

// GetSibling returns index of a sibling element
func GetSibling(idx uint32) uint32 {
	if idx == 0 {
		return 0
	}
	return ((idx + 1) ^ 1) - 1
}

// GetParent returns index of a parent element
func GetParent(idx uint32) uint32 {
	if idx == 0 {
		return 0
	}
	return (idx - 1) >> 1
}

// IsLeft checks if an index is a left index
func IsLeft(idx uint32) bool {
	return idx&1 == 1
}
