package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

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
	var lemmas [][]byte
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

	leaves := MapIndiceAndLeaves(indices, mtree.Nodes)
	SortIndicesAndLeavesByIndex(leaves)

	return Proof{
		Leaves: leaves,
		Lemmas: lemmas,
		Merge:  mtree.Merge,
	}, nil

}

// GetRoot returns the root value of merkle tree
func (mtree Tree) GetRoot() interface{} {
	if len(mtree.Nodes) == 0 {
		return 0
	}
	return mtree.Nodes[0]
}
