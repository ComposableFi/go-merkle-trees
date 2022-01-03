package main

import (
	"crypto/sha256"
	"encoding/hex"
)

type MergeByteArray struct{}

func (m MergeByteArray) Merge(left, right []byte) []byte {
	h, _ := HashNodes(Node{
		Left:  left,
		Right: right,
	})
	return h
}

type Node struct {
	Left  []byte
	Right []byte
}

func HashNodes(node Node) ([]byte, error) {
	return CalculateHash(append(node.Left[:], node.Right[:]...))
}

func CalculateHash(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func HashToStr(h interface{}) string {
	return hex.EncodeToString(h.([]byte))
}
