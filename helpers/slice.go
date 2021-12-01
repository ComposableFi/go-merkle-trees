package helpers

// PopFromUint64Queue pops front element of uint64 slice
func PopFromUint64Queue(slice []uint64) (uint64, []uint64) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice

}

// PopFromEndUint64Queue pops last element of uint64 slice
func PopFromEndUint64Queue(slice []uint64) (uint64, []uint64) {
	sliceLen := len(slice)
	popElem, newSlice := slice[sliceLen-1], slice[:sliceLen-1]
	return popElem, newSlice
}

// PopTwoElementsFromEndUint64Queue pops last two element of uint64 slice
func PopTwoElementsFromEndUint64Queue(slice []uint64) (uint64, uint64, []uint64) {
	sliceLen := len(slice)
	lastElem, beforeLastElem, newSlice := slice[sliceLen-1], slice[sliceLen-2], slice[:sliceLen-2]
	return beforeLastElem, lastElem, newSlice
}

// Uint64SliceContains checks if a slice contains specific uint64 number
func Uint64SliceContains(slice []uint64, num uint64) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}
