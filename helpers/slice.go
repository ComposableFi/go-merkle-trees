package helpers

// PopFromUint32Queue pops front element of uint32 slice
func PopFromUint32Queue(slice []uint32) (uint32, []uint32) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice

}

// PopFromInterfaceQueue pops front element of interface slice
func PopFromInterfaceQueue(slice []interface{}) (interface{}, []interface{}) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice

}

// PopFromEndUint32Queue pops last element of uint32 slice
func PopFromEndUint32Queue(slice []uint32) (uint32, []uint32) {
	sliceLen := len(slice)
	popElem, newSlice := slice[sliceLen-1], slice[:sliceLen-1]
	return popElem, newSlice
}

// PopTwoElementsFromEndUint32Queue pops last two element of uint32 slice
func PopTwoElementsFromEndUint32Queue(slice []uint32) (uint32, uint32, []uint32) {
	sliceLen := len(slice)
	lastElem, beforeLastElem, newSlice := slice[sliceLen-1], slice[sliceLen-2], slice[:sliceLen-2]
	return beforeLastElem, lastElem, newSlice
}

// PopTwoElementsFromEndInterfaceQueue pops last two element of uint32 slice
func PopTwoElementsFromEndInterfaceQueue(slice []interface{}) (interface{}, interface{}, []interface{}) {
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
