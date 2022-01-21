package main

import "github.com/ethereum/go-ethereum/crypto"

type Keccak256Hasher struct{}

func (hr Keccak256Hasher) Hash(b []byte) ([]byte, error) {
	h := crypto.Keccak256Hash(b)
	return h.Bytes(), nil
}
