package hasher

import (
	"crypto/sha256"
)

// Sha256Hasher is hasher type for the sha256
type Sha256Hasher struct{}

// Hash generates sha256 hash from bytes
func (hr Sha256Hasher) Hash(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
