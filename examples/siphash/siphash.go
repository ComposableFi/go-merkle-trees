package main

import (
	sh "github.com/dchest/siphash"
)

type MergeUint64 struct{}

func (m MergeUint64) Merge(left, right interface{}) interface{} {
	h := HashNodes(Node{
		Left:  left.(uint64),
		Right: right.(uint64),
	})
	return h
}

type Node struct {
	Left  uint64
	Right uint64
}

func HashNodes(node Node) uint64 {
	sum64 := sh.Hash(node.Left, node.Right, []byte{})
	return sum64
}
