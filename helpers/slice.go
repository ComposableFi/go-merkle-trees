package helpers

// PopFromUint32Queue pops front element of uint64 slice
func PopFromUint32Queue(slice []uint64) (uint64, []uint64) {
	popElem, newSlice := slice[len(slice)-1], slice[:len(slice)-1]
	return popElem, newSlice

}

// Uint32SliceContains checks if a slice contains specific uint64 number
func Uint32SliceContains(slice []uint64, num uint64) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == num {
			return true
		}
	}
	return false
}

// Difference finds the elements of first slice that are not present in the second slice
func Difference(slice1 []uint64, slice2 []uint64) []uint64 {
	var diff []uint64

	for i := 0; i < len(slice1); i++ {
		found := false
		for j := 0; j < len(slice2); j++ {
			if slice1[i] == slice2[j] {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, slice1[i])
		}
	}

	return diff
}
