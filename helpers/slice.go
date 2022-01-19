package helpers

// PopFromUint32Queue pops front element of uint32 slice
func PopFromUint32Queue(slice []uint32) (uint32, []uint32) {
	popElem, newSlice := slice[len(slice)-1], slice[:len(slice)-1]
	return popElem, newSlice

}

// PopFromInterfaceQueue pops front element of interface slice
func PopFromInterfaceQueue(slice [][]byte) ([]byte, [][]byte) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice

}

// PopTwoElementsFromEndInterfaceQueue pops last two element of uint32 slice
func PopTwoElementsFromEndInterfaceQueue(slice [][]byte) ([]byte, []byte, [][]byte) {
	sliceLen := len(slice)
	lastElem, beforeLastElem, newSlice := slice[sliceLen-1], slice[sliceLen-2], slice[:sliceLen-2]
	return beforeLastElem, lastElem, newSlice
}

// Uint32SliceContains checks if a slice contains specific uint32 number
func Uint32SliceContains(slice []uint32, num uint32) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

// Difference finds the elements of first slice that are not present in the second slice
func Difference(slice1 []uint32, slice2 []uint32) []uint32 {
	var diff []uint32

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
