package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

// GetRoot returns the root value of a merkle proof
func (p *Proof) GetRoot(leaves []uint64) (uint64, error) {
	if len(leaves) != len(p.Indices) || len(leaves) == 0 {
		return 0, errors.New("leaves could not be empty")
	}

	helpers.SortUint64Slice(leaves)

	pre := MapIndiceAndLeaves(p.Indices, leaves)
	SortIndicesAndLeavesByIndexReversely(pre)

	var queue []LeafIndex
	queue = append(queue, pre...)

	lemmaIdx := 0
	for {
		if len(queue) == 0 {
			break
		}

		var leafIdx LeafIndex
		leafIdx, queue = PopFromLeafIndexQueue(queue)

		if leafIdx.Index == 0 {
			if lemmaIdx <= len(p.Lemmas) && len(queue) == 0 {
				return leafIdx.Leaf, nil
			}
			return 0, errors.New("there are more unprocessed queue items")
		}

		var sibling uint64
		if len(queue) > 0 && queue[0].Index == helpers.GetSibling(leafIdx.Index) {
			var sibLeaf LeafIndex
			sibLeaf, queue = PopFromLeafIndexQueue(queue)
			sibling = sibLeaf.Leaf
		} else {
			sibling = p.Lemmas[lemmaIdx]
			lemmaIdx++
		}

		var parentNode uint64
		if helpers.IsLeft(leafIdx.Index) {
			parentNode = p.Merge.Merge(leafIdx.Leaf, sibling)
		} else {
			parentNode = p.Merge.Merge(sibling, leafIdx.Leaf)
		}
		queue = append(queue, LeafIndex{
			Index: helpers.GetParent(leafIdx.Index),
			Leaf:  parentNode,
		})
	}
	return 0, errors.New("")
}

// Verify verifies the root value against tree leaves
func (p *Proof) Verify(root uint64, leaves []uint64) (bool, error) {
	r, err := p.GetRoot(leaves)
	if err != nil {
		return false, err
	}
	if root == r {
		return true, nil
	}
	return false, nil
}

// RetriveLeaves returns the leaves of a merkle proof
func (p Proof) RetriveLeaves(leaves []uint64) ([]uint64, error) {
	if len(leaves) == 0 || len(p.Indices) == 0 {
		return []uint64{}, errors.New("leaves or indecies should not be empty")
	}

	leavesCount := uint64(len(leaves))
	var validIndicesRange []uint64
	for i := uint64(leavesCount - 1); i < uint64((leavesCount<<1)-1); i++ {
		validIndicesRange = append(validIndicesRange, i)
	}

	var allProofsInRange = true

	for _, v := range p.Indices {
		if !helpers.Uint64SliceContains(validIndicesRange, v) {
			allProofsInRange = false
			break
		}
	}

	var res []uint64
	if allProofsInRange {
		for _, v := range p.Indices {
			res = append(res, leaves[uint64(v+1-leavesCount)])
		}
	}

	return res, nil
}
