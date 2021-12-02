package main

import (
	"crypto/sha256"
	"encoding/hex"
)

type MergeByteArray struct{}

func (m MergeByteArray) Merge(left, right interface{}) interface{} {
	h, _ := HashNodes(Node{
		Left:  left.([]byte),
		Right: right.([]byte),
	})
	return h
}

type Node struct {
	Left  []byte
	Right []byte
}

func HashNodes(node Node) ([]byte, error) {
	var l, r []byte
	l, err := CalculateHash(node.Left)
	if err != nil {
		panic(err)
	}
	r, err = CalculateHash(node.Right)
	if err != nil {
		panic(err)
	}
	return CalculateHash(append(l[:], r[:]...))
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
