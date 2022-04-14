package merkle

// popFromIndexQueue pops front element of uint64 slice
func popFromIndexQueue(slice []uint64) (uint64, []uint64) {
	popElem, newSlice := slice[len(slice)-1], slice[:len(slice)-1]
	return popElem, newSlice

}

// SliceDifference finds the elements of first slice that are not present in the second slice
func SliceDifference(slice1 []uint64, slice2 []uint64) []uint64 {
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
