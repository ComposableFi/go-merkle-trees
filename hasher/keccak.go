package hasher

import "github.com/ethereum/go-ethereum/crypto"

// Keccak256Hasher is hasher type for the keccack
type Keccak256Hasher struct{}

// Hash generates keccak hash from bytes
func (hr Keccak256Hasher) Hash(b []byte) ([]byte, error) {
	h := crypto.Keccak256Hash(b)
	return h.Bytes(), nil
}
