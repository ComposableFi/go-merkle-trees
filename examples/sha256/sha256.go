package main

import (
	"crypto/sha256"
)

type Sha256Hasher struct{}

func (hr Sha256Hasher) Hash(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
