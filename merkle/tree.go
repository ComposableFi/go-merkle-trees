package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

// BuildRoot builds the merkle root from a merkle tree
func (mtree Tree) BuildRoot(leaves []interface{}) interface{} {
	if len(leaves) == 0 {
		return 0
	}
	var queue []interface{}

	for len(leaves) >= 2 {
		var leaf1, leaf2 interface{}
		leaf1, leaf2, leaves = helpers.PopTwoElementsFromEndInterfaceQueue(leaves)

		m := mtree.Merge.Merge(leaf1, leaf2)
		queue = append(queue, m)
	}
	if len(leaves) == 1 {
		leaf := leaves[0]
		queue = append([]interface{}{leaf}, queue...)
	}

	for len(queue) > 1 {
		var right, left interface{}
		right, queue = helpers.PopFromInterfaceQueue(queue)
		left, queue = helpers.PopFromInterfaceQueue(queue)
		queue = append(queue, mtree.Merge.Merge(left, right))
	}

	res, _ := helpers.PopFromInterfaceQueue(queue)
	return res
}

// BuildTree constructs merkle tree into the merkle
func (mtree Tree) BuildTree(leaves []interface{}) Tree {
	leaveLen := len(leaves)

	if leaveLen > 0 {

		nodes := make([]interface{}, leaveLen-1)
		nodes = append(nodes, leaves...)

		for i := leaveLen - 2; 0 <= i; i-- {
			left := nodes[(i<<1)+1]
			right := nodes[(i<<1)+2]
			merged := mtree.Merge.Merge(left, right)
			nodes[i] = merged
		}
		return Tree{
			Nodes: nodes,
			Merge: mtree.Merge,
		}
	}
	return Tree{Merge: mtree.Merge}
}

// BuildProof builds the merkle proof by leaf index slice
func (mtree Tree) BuildProof(leafIndices []uint32) (Proof, error) {
	if len(leafIndices) == 0 || len(mtree.Nodes) == 0 {
		return Proof{Merge: mtree.Merge}, errors.New("empty nodes or indices not allowed")
	}

	leavesCount := uint32((len(mtree.Nodes) >> 1) + 1)
	var indices []uint32

	for _, leafIdx := range leafIndices {
		indices = append(indices, leavesCount+leafIdx-1)
	}
	helpers.ReverseSortUint32Slice(indices)

	if indices[0] >= (leavesCount<<1)-1 {
		return Proof{Merge: mtree.Merge}, errors.New("first element of indices is not valid")
	}
	var lemmas []interface{}
	queue := append([]uint32{}, indices...)

	for {

		if len(queue) == 0 {
			break
		}

		var idx uint32
		idx, queue = helpers.PopFromUint32Queue(queue)

		if idx == 0 {
			if len(queue) != 0 {
				return Proof{Merge: mtree.Merge}, errors.New("queue is not empty")
			}
			break
		}
		sibling := helpers.GetSibling(idx)

		if len(queue) > 0 && sibling == queue[0] {
			_, queue = helpers.PopFromUint32Queue(queue)
		} else {
			lemmas = append(lemmas, mtree.Nodes[sibling])
		}
		parent := helpers.GetParent(idx)
		if parent != 0 {
			queue = append(queue, parent)
		}
	}

	helpers.SortUint32Slice(indices)

	return Proof{
		Indices: indices,
		Lemmas:  lemmas,
		Merge:   mtree.Merge,
	}, nil

}

// BuildTreeAndProof builds merkle proof from the tree and leaves
func (mtree Tree) BuildTreeAndProof(leaves []interface{}, leafIndices []uint32) (Proof, error) {
	return mtree.BuildTree(leaves).BuildProof(leafIndices)
}

// GetRoot returns the root value of merkle tree
func (mtree Tree) GetRoot() interface{} {
	if len(mtree.Nodes) == 0 {
		return 0
	}
	return mtree.Nodes[0]
}
