package helpers

// PopFromUint32Queue pops front element of uint64 slice
func PopFromUint32Queue(slice []uint64) (uint64, []uint64) {
	popElem, newSlice := slice[len(slice)-1], slice[:len(slice)-1]
	return popElem, newSlice

}

// PopFromInterfaceQueue pops front element of interface slice
func PopFromInterfaceQueue(slice [][]byte) ([]byte, [][]byte) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice

}

// PopTwoElementsFromEndInterfaceQueue pops last two element of uint64 slice
func PopTwoElementsFromEndInterfaceQueue(slice [][]byte) ([]byte, []byte, [][]byte) {
	sliceLen := len(slice)
	lastElem, beforeLastElem, newSlice := slice[sliceLen-1], slice[sliceLen-2], slice[:sliceLen-2]
	return beforeLastElem, lastElem, newSlice
}

// Uint32SliceContains checks if a slice contains specific uint64 number
func Uint32SliceContains(slice []uint64, num uint64) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

// Difference finds the elements of first slice that are not present in the second slice
func Difference(slice1 []uint64, slice2 []uint64) []uint64 {
	var diff []uint64

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}
