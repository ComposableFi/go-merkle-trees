package merkle

import (
	"encoding/hex"
	"errors"
	"math"
	"sort"

	"github.com/ComposableFi/merkle-go/helpers"
)

func (t Tree) FromLeaves(leaves []Hash) Tree {
	t.append(leaves)
	t.commit()
	return t
}

func (t Tree) getRoot() Hash {
	layers := t.layerTuples()
	lastLayer := layers[len(layers)-1]
	firstItem := lastLayer[0]
	return firstItem.hash
}

func (t Tree) getRootHex() string {
	root := t.getRoot()
	return hex.EncodeToString([]byte(root))
}
func (t Tree) HelperNodes(leafIndices []uint32) []Hash {
	var helperNodes []Hash
	for _, layer := range t.HelperNodeTuples(leafIndices) {
		for _, li := range layer {
			helperNodes = append(helperNodes, li.hash)
		}
	}
	return helperNodes
}
func (t Tree) HelperNodeTuples(leafIndeceis []uint32) [][]leafIndexAndHash {
	var helpersLayer []leafIndexAndHash
	var helperNodes [][]leafIndexAndHash
	for _, treeLayer := range t.layerTuples() {
		siblings := helpers.GetSiblingIndecies(leafIndeceis)
		helperIndices := helpers.Difference(siblings, leafIndeceis)

		for _, idx := range helperIndices {
			i, _ := getLeafAndHashAtIndex(treeLayer, idx)
			helpersLayer[idx] = i
		}
		helperNodes = append(helperNodes, helpersLayer)
		leafIndeceis = helpers.GetParentIndecies(leafIndeceis)
	}
	return helperNodes
}

func (t Tree) insert(leaf Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaf)
}

func (t Tree) append(leaves []Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaves...)
}

func (t Tree) commit() error {
	diff, err := t.uncommittedDiff()
	if err != nil {
		return err
	}
	t.history = append(t.history, diff)
	t.currentWorkingTree.mergeUnverified(diff)
	t.UncommittedLeaves = []Hash{}
	return nil
}

func (t Tree) uncommittedRoot() (Hash, error) {
	shadowTree, err := t.uncommittedDiff()
	if err != nil {
		return Hash{}, err
	}
	return shadowTree.getRoot(), nil
}

func (t Tree) uncommittedRootHex() (string, error) {
	root, err := t.uncommittedRoot()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), err
}

func (t Tree) abortCommitted() {
	t.UncommittedLeaves = make([]Hash, 0)
}

func (t Tree) depth() int {
	return len(t.layerTuples()) - 1
}

func (t Tree) leaves() []Hash {
	layers := t.layers()
	if len(layers) > 0 {
		return []Hash{}
	}
	return layers[0]
}

func (t Tree) leavesLen() int {
	leaves := t.leavesTuples()
	return len(leaves)
}

func (t Tree) leavesTuples() []leafIndexAndHash {
	return t.layerTuples()[0]
}

func (t Tree) layers() [][]Hash {
	return t.currentWorkingTree.layerNodes()
}

func (t Tree) layerTuples() [][]leafIndexAndHash {
	return t.currentWorkingTree.layers
}

func (t Tree) uncommittedDiff() (PartialTree, error) {
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, errors.New("leaves can not be empty!")
	}
	commitedLeavesCount := t.leavesLen()
	var shadowIndecies []uint32
	for i, _ := range t.UncommittedLeaves {
		shadowIndecies = append(shadowIndecies, uint32(commitedLeavesCount+i))
	}
	var shadowNodeTuples []leafIndexAndHash
	for _, idx := range shadowIndecies {
		x := leafIndexAndHash{index: idx, hash: t.UncommittedLeaves[idx]}
		shadowNodeTuples = append(shadowNodeTuples, x)
	}
	partialTreeTuples := t.HelperNodeTuples(shadowIndecies)
	leavesInNewTree := t.leavesLen() + len(t.UncommittedLeaves)
	uncommittedTreeDepth := getTreeDepth(leavesInNewTree)
	if len(partialTreeTuples) == 0 {
		partialTreeTuples = append(partialTreeTuples, shadowNodeTuples)
	} else {
		firstLayer := partialTreeTuples[0]
		firstLayer = append(firstLayer, shadowNodeTuples...)
		sortLeafAndHashByIndex(firstLayer)
	}
	return NewPartialTree(t.hasher).build(partialTreeTuples, uncommittedTreeDepth)
}

func getLeafAndHashAtIndex(leavesAndHash []leafIndexAndHash, index uint32) (leafIndexAndHash, error) {
	for _, l := range leavesAndHash {
		if l.index == index {
			return l, nil
		}
	}
	return leafIndexAndHash{}, errors.New("leaf not found")
}

func getLayerAtIndex(layers [][]leafIndexAndHash, index uint32) ([]leafIndexAndHash, bool) {
	if len(layers) > int(index) {
		return layers[index], true
	}
	return []leafIndexAndHash{}, false
}

func sortLeafAndHashByIndex(li []leafIndexAndHash) {
	sort.Slice(li, func(i, j int) bool { return li[i].index < li[j].index })

}
func getTreeDepth(leaves_count int) int {
	if leaves_count == 1 {
		return 1
	} else {
		return int(math.Ceil(math.Log2(float64(leaves_count))))
	}
}
