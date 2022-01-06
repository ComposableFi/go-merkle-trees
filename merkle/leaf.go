package merkle

// PPopFromLeafHashQueue pops first front element in a leaf hash slice
func PopFromLeafHashQueue(slice [][]leafIndexAndHash) ([]leafIndexAndHash, [][]leafIndexAndHash) {
	popElem, newSlice := slice[0], slice[1:]
	return popElem, newSlice
}
