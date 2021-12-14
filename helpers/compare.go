package helpers

import "bytes"

// AreInterfacesEqual checks if two interface are equal
func AreInterfacesEqual(first, second interface{}) bool {
	switch first.(type) {
	case []byte:
		return bytes.Equal(first.([]byte), second.([]byte))
	default:
		return true //first == second
	}
}
