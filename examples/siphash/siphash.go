package main

import (
	"encoding/binary"

	sh "github.com/dchest/siphash"
)

type MergeUint64 struct{}

func (m MergeUint64) Merge(left, right []byte) []byte {
	h := HashNodes(Node{
		Left:  b2i(left),
		Right: b2i(right),
	})
	return i2b(h)
}

type Node struct {
	Left  uint64
	Right uint64
}

func HashNodes(node Node) uint64 {
	sum64 := sh.Hash(node.Left, node.Right, []byte{})
	return sum64
}

func b2i(b []byte) uint64 {
	i := binary.LittleEndian.Uint64(b)
	return i
}

func i2b(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}
