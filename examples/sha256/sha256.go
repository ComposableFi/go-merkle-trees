package main

import (
	"crypto/sha256"

	"github.com/ComposableFi/merkle-go/merkle"
)

type Sha256Hasher struct{}

func (hr Sha256Hasher) Hash(b []byte) (merkle.Hash, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
