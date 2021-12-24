package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

// BuildMerkleRoot builds the merkle root from leaves
func (cbmt CBMT) BuildMerkleRoot(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return nil
	}
	var queue [][]byte

	for len(leaves) >= 2 {
		var leaf1, leaf2 []byte
		leaf1, leaf2, leaves = helpers.PopTwoElementsFromEndInterfaceQueue(leaves)

		m := cbmt.Merge.Merge(leaf1, leaf2)
		queue = append(queue, m)
	}
	if len(leaves) == 1 {
		leaf := leaves[0]
		queue = append([][]byte{leaf}, queue...)
	}

	for len(queue) > 1 {
		var right, left []byte
		right, queue = helpers.PopFromInterfaceQueue(queue)
		left, queue = helpers.PopFromInterfaceQueue(queue)
		queue = append(queue, cbmt.Merge.Merge(left, right))
	}

	res, _ := helpers.PopFromInterfaceQueue(queue)
	return res
}

// BuildMerkleTree constructs merkle tree
func (cbmt CBMT) BuildMerkleTree(leaves [][]byte) Tree {
	leaveLen := len(leaves)

	if leaveLen > 0 {

		nodes := make([][]byte, leaveLen-1)
		nodes = append(nodes, leaves...)

		for i := leaveLen - 2; 0 <= i; i-- {
			left := nodes[(i<<1)+1]
			right := nodes[(i<<1)+2]
			merged := cbmt.Merge.Merge(left, right)
			nodes[i] = merged
		}
		return Tree{
			Nodes: nodes,
			Merge: cbmt.Merge,
		}
	}
	return Tree{Merge: cbmt.Merge}
}

// BuildMerkleProof builds merkle proof from the leaves and leaf indices
func (cbmt CBMT) BuildMerkleProof(leaves [][]byte, leafIndices []uint32) (Proof, error) {
	return cbmt.BuildMerkleTree(leaves).BuildProof(leafIndices)
}

// RetriveLeaves returns the leaves of a merkle proof
func (cbmt CBMT) RetriveLeaves(p Proof, leaves [][]byte) ([][]byte, error) {
	if len(leaves) == 0 || len(p.Leaves) == 0 {
		return [][]byte{}, errors.New("leaves or indecies should not be empty")
	}

	leavesCount := uint32(len(leaves))
	var validIndicesRange []uint32
	for i := uint32(leavesCount - 1); i < uint32((leavesCount<<1)-1); i++ {
		validIndicesRange = append(validIndicesRange, i)
	}

	var allProofsInRange = true

	for _, v := range p.Leaves {
		if !helpers.Uint32SliceContains(validIndicesRange, v.Index) {
			allProofsInRange = false
			break
		}
	}

	var res [][]byte
	if allProofsInRange {
		for _, v := range p.Leaves {
			res = append(res, leaves[v.Index+1-leavesCount])
		}
	}

	return res, nil
}
