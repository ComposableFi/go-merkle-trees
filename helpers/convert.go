package helpers

import (
	"bytes"
	"encoding/gob"
)

// GetInterfaceBytes returns bytes of an interface
func GetInterfaceBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
