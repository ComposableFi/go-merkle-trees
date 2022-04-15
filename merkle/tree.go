package merkle

import (
	"encoding/hex"
	"math"
	"sort"

	"github.com/ComposableFi/go-merkle-trees/types"
)

// FromLeaves clones the leaves and builds the tree from them
func (t Tree) FromLeaves(leaves [][]byte) (Tree, error) {
	t.append(leaves)
	err := t.commit()
	if err != nil {
		return Tree{}, err
	}
	return t, nil
}

// Root returns the tree root - the top hash of the tree. Used in the inclusion proof verification.
func (t *Tree) Root() []byte {
	layers := t.layerLeaves()
	if len(layers) > 0 {
		lastLayer := layers[len(layers)-1]
		firstItem := lastLayer[0]
		return firstItem.Hash
	}
	return []byte{}
}

// RootHex returns a hex encoded string instead of
func (t *Tree) RootHex() string {
	root := t.Root()
	return hex.EncodeToString(root)
}

// helperNodes returns helper nodes required to build a partial tree for the given indices
// to be able to extract a root from it. Useful in constructing Merkle proofs
func (t *Tree) helperNodes(leafIndices []uint64) [][]byte {
	var helperNodes [][]byte
	for _, layer := range t.helperNodeLeaves(leafIndices) {
		for _, li := range layer {
			helperNodes = append(helperNodes, li.Hash)
		}
	}
	return helperNodes
}

// helperNodeLeaves gets all helper nodes required to build a partial merkle tree for the given indices,
// cloning all required hashes into the resulting vector.
func (t *Tree) helperNodeLeaves(leafIndeceis []uint64) [][]types.Leaf {
	var helperNodes [][]types.Leaf
	for _, treeLayer := range t.layerLeaves() {
		siblings := siblingIndecies(leafIndeceis)
		helperIndices := sliceDifference(siblings, leafIndeceis)

		var helpersLayer []types.Leaf
		for _, idx := range helperIndices {
			leaf, found := leafAtIndex(treeLayer, idx)
			if found {
				helpersLayer = append(helpersLayer, leaf)
			}
		}

		helperNodes = append(helperNodes, helpersLayer)

		leafIndeceis = parentIndecies(leafIndeceis)
	}
	return helperNodes
}

// Proof Returns the Merkle proof required to prove the inclusion of items in a data set.
func (t *Tree) Proof(leafIndices []uint64) Proof {
	leavesLen := t.leavesLen()
	leaves := t.leaves()
	var proofLeaves []types.Leaf

	for _, index := range leafIndices {
		for _, leaf := range leaves {
			if leaf.Index == index {
				proofLeaves = append(proofLeaves, leaf)
				break
			}
		}
	}
	return NewProof(proofLeaves, t.helperNodes(leafIndices), leavesLen, t.hasher)
}

// insert inserts a new types.Leaf. Please note it won't modify the root just yet; For the changes
// to be applied to the root, [`MerkleTree::commit`] method should be called first. To get the
// root of the new tree without applying the changes, you can use
func (t *Tree) insert(leaf []byte) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaf)
}

// append appends leaves to the tree. Behaves similarly to [`MerkleTree::insert`], but for a list of
// items. Takes ownership of the elements.
func (t *Tree) append(leaves [][]byte) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaves...)
}

// commit commits the changes made by [`MerkleTree::insert`] and [`MerkleTree::append`]
// and modifies the root.
// Commits are saved to the history, so the tree can be rolled back to any previous commit
func (t *Tree) commit() error {
	diff, err := t.uncommittedDiff()
	if err != nil {
		return err
	}
	if len(diff.layers) > 0 {
		t.history = append(t.history, diff)
		t.currentWorkingTree.mergeUnverified(diff)
		t.UncommittedLeaves = [][]byte{}
	}
	return nil
}

// rollback rolls back one commit and reverts the tree to the previous state.
// Removes the most recent commit from the history.
func (t *Tree) rollback() {
	_, t.history = popFromPartialtree(t.history)
	t.currentWorkingTree.clear()
	for _, commit := range t.history {
		t.currentWorkingTree.mergeUnverified(commit)
	}
}

// uncommittedRoot calculates the root of the uncommitted changes as if they were committed.
// Will return the same hash as root of merkle tree after commit
func (t *Tree) uncommittedRoot() ([]byte, error) {
	shadowTree, err := t.uncommittedDiff()
	if err != nil {
		return []byte{}, err
	}
	return shadowTree.Root(), nil
}

// uncommittedRootHex calculates the root of the uncommitted changes as if they were committed. Serializes
// the result as a hex string.
func (t *Tree) uncommittedRootHex() (string, error) {
	root, err := t.uncommittedRoot()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), nil
}

// depth returns the tree depth. A tree depth is how many layers there is between the
// leaves and the root
func (t *Tree) depth() int {
	return len(t.layerLeaves()) - 1
}

// baseLeaves returns a copy of the tree leaves - the base level of the tree.
func (t *Tree) baseLeaves() [][]byte {
	layers := t.layers()
	if len(layers) > 0 {
		return [][]byte{}
	}
	return layers[0]
}

// leavesLen returns the number of leaves in the tree.
func (t *Tree) leavesLen() uint64 {
	leaves := t.leaves()
	return uint64(len(leaves))
}

// leaves returns leaves of the current working tree
func (t *Tree) leaves() []types.Leaf {
	if len(t.layerLeaves()) > 0 {
		return t.layerLeaves()[0]
	}
	return []types.Leaf{}
}

// layers returns the whole tree, where the first layer is leaves and
// consequent layers are nodes.
func (t *Tree) layers() [][][]byte {
	return t.currentWorkingTree.layerNodes()
}

// layerLeaves returns leaves of the current working tree
func (t *Tree) layerLeaves() [][]types.Leaf {
	return t.currentWorkingTree.layers
}

/// uncommittedDiff creates a diff from a changes that weren't committed to the main tree yet. Can be used
/// to get uncommitted root or can be merged with the main tree
func (t *Tree) uncommittedDiff() (PartialTree, error) {
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, nil
	}

	shadowIndecies := t.getShadowIndecies()
	shadowNodeLeaves := t.getShadowLeaves(shadowIndecies)

	partialTreeLeaves := t.helperNodeLeaves(shadowIndecies)
	leavesInNewTree := t.leavesLen() + uint64(len(t.UncommittedLeaves))
	uncommittedTreeDepth := treeDepth(leavesInNewTree)
	if len(partialTreeLeaves) == 0 {
		partialTreeLeaves = append(partialTreeLeaves, shadowNodeLeaves)
	} else {
		firstLayer := partialTreeLeaves[0]
		firstLayer = append(firstLayer, shadowNodeLeaves...)
		sortLeavesByIndex(firstLayer)
		partialTreeLeaves[0] = firstLayer
	}
	tree := NewPartialTree(t.hasher)
	return tree.build(partialTreeLeaves, uncommittedTreeDepth)
}

// getShadowIndecies returns shadow indices of the uncommited leaves
func (t *Tree) getShadowIndecies() []uint64 {
	if len(t.UncommittedLeaves) == 0 {
		return []uint64{}
	}

	commitedLeavesCount := t.leavesLen()
	var shadowIndecies []uint64
	for i := 0; i < len(t.UncommittedLeaves); i++ {
		shadowIndecies = append(shadowIndecies, commitedLeavesCount+uint64(i))
	}
	return shadowIndecies
}

// getShadowLeaves returns shadow leaves of the uncommited leaves
func (t *Tree) getShadowLeaves(shadowIndecies []uint64) []types.Leaf {
	var shadowNodeLeaves []types.Leaf
	for i := 0; i < len(shadowIndecies); i++ {
		leaf := types.Leaf{Index: shadowIndecies[i], Hash: t.UncommittedLeaves[i]}
		shadowNodeLeaves = append(shadowNodeLeaves, leaf)
	}
	return shadowNodeLeaves
}

// leafAtIndex returns types.Leaf object at the index
func leafAtIndex(leavesAndHash []types.Leaf, index uint64) (types.Leaf, bool) {
	for i := 0; i < len(leavesAndHash); i++ {
		l := leavesAndHash[i]
		if l.Index == index {
			return l, true
		}
	}
	return types.Leaf{}, false
}

// layerAtIndex returns layer object at the index
func layerAtIndex(layers [][]types.Leaf, index uint64) ([]types.Leaf, bool) {
	if len(layers) > int(index) {
		return layers[index], true
	}
	return []types.Leaf{}, false
}

// sortLeavesByIndex sorts leaves by their index
func sortLeavesByIndex(li []types.Leaf) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index < li[j].Index })
}

// treeDepth returns the depth of a tree
func treeDepth(leavesCount uint64) uint64 {
	if leavesCount == 1 {
		return 1
	}
	return uint64(math.Ceil(math.Log2(float64(leavesCount))))
}
