package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

// GetRoot returns the root value of a merkle proof
func (p *Proof) GetRoot() (interface{}, error) {
	if len(p.Leaves) == 0 {
		return 0, errors.New("leaves count should equal to indices and not be empty")
	}

	// TODO: write function for interace sorting helpers.SortUint32Slice(leaves)

	SortIndicesAndLeavesByIndexReversely(p.Leaves)

	var queue []LeafData
	queue = append(queue, p.Leaves...)

	lemmaIdx := 0
	for {
		if len(queue) == 0 {
			break
		}

		var leafIdx LeafData
		leafIdx, queue = PopFromLeafIndexQueue(queue)

		if leafIdx.Index == 0 {
			if lemmaIdx <= len(p.Lemmas) && len(queue) == 0 {
				return leafIdx.Leaf, nil
			}
			return 0, errors.New("there are more unprocessed queue items")
		}

		var sibling interface{}
		if len(queue) > 0 && queue[0].Index == helpers.GetSibling(leafIdx.Index) {
			var sibLeaf LeafData
			sibLeaf, queue = PopFromLeafIndexQueue(queue)
			sibling = sibLeaf.Leaf
		} else {
			sibling = p.Lemmas[lemmaIdx]
			lemmaIdx++
		}

		var parentNode interface{}
		if helpers.IsLeft(leafIdx.Index) {
			parentNode = p.Merge.Merge(leafIdx.Leaf, sibling)
		} else {
			parentNode = p.Merge.Merge(sibling, leafIdx.Leaf)
		}
		queue = append(queue, LeafData{
			Index: helpers.GetParent(leafIdx.Index),
			Leaf:  parentNode,
		})
	}
	return 0, errors.New("root not found")
}

// Verify verifies the root value against tree leaves
func (p *Proof) Verify(root interface{}) (bool, error) {
	r, err := p.GetRoot()
	if err != nil {
		return false, err
	}
	if helpers.AreInterfacesEqual(root, r) {
		return true, nil
	}
	return false, nil
}
