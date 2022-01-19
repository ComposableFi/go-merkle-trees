package merkle

import (
	"encoding/hex"
	"math"
	"sort"

	"github.com/ComposableFi/merkle-go/helpers"
)

/// FromLeaves clones the leaves and builds the tree from them
func (t Tree) FromLeaves(leaves []Hash) (Tree, error) {
	t.Append(leaves)
	err := t.Commit()
	if err != nil {
		return Tree{}, err
	}
	return t, nil
}

// Root returns the tree root - the top hash of the tree. Used in the inclusion proof verification.
func (t *Tree) Root() Hash {
	layers := t.layerLeaves()
	if len(layers) > 0 {
		lastLayer := layers[len(layers)-1]
		firstItem := lastLayer[0]
		return firstItem.Hash
	}
	return Hash{}
}

// RootHex returns a hex encoded string instead of
func (t *Tree) RootHex() string {
	root := t.Root()
	return hex.EncodeToString([]byte(root))
}

// HelperNodes returns helper nodes required to build a partial tree for the given indices
// to be able to extract a root from it. Useful in constructing Merkle proofs
func (t *Tree) HelperNodes(leafIndices []uint32) []Hash {
	var helperNodes []Hash
	for _, layer := range t.HelperNodeLeaves(leafIndices) {
		for _, li := range layer {
			helperNodes = append(helperNodes, li.Hash)
		}
	}
	return helperNodes
}

// HelperNodeLeaves gets all helper nodes required to build a partial merkle tree for the given indices,
// cloning all required hashes into the resulting vector.
func (t *Tree) HelperNodeLeaves(leafIndeceis []uint32) [][]Leaf {
	var helperNodes [][]Leaf
	for _, treeLayer := range t.layerLeaves() {
		siblings := helpers.SiblingIndecies(leafIndeceis)
		helperIndices := helpers.Difference(siblings, leafIndeceis)

		var helpersLayer []Leaf
		for _, idx := range helperIndices {
			leaf, found := leafAtIndex(treeLayer, idx)
			if found {
				helpersLayer = append(helpersLayer, leaf)
			}
		}

		helperNodes = append(helperNodes, helpersLayer)

		leafIndeceis = helpers.ParentIndecies(leafIndeceis)
	}
	return helperNodes
}

// Proof Returns the Merkle proof required to prove the inclusion of items in a data set.
func (t *Tree) Proof(leafIndices []uint32) Proof {
	leavesLen := t.LeavesLen()
	leaves := t.leaves()
	var proofLeaves []Leaf

	for _, index := range leafIndices {
		for _, leaf := range leaves {
			if leaf.Index == index {
				proofLeaves = append(proofLeaves, leaf)
				break
			}
		}
	}
	return NewProof(proofLeaves, t.HelperNodes(leafIndices), leavesLen, t.hasher)
}

// Insert inserts a new leaf. Please note it won't modify the root just yet; For the changes
// to be applied to the root, [`MerkleTree::commit`] method should be called first. To get the
// root of the new tree without applying the changes, you can use
func (t *Tree) Insert(leaf Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaf)
}

// Append appends leaves to the tree. Behaves similarly to [`MerkleTree::insert`], but for a list of
// items. Takes ownership of the elements.
func (t *Tree) Append(leaves []Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaves...)
}

// Commit commits the changes made by [`MerkleTree::insert`] and [`MerkleTree::append`]
// and modifies the root.
// Commits are saved to the history, so the tree can be rolled back to any previous commit
func (t *Tree) Commit() error {
	diff, err := t.uncommittedDiff()
	if err != nil {
		return err
	}
	if len(diff.layers) > 0 {
		t.history = append(t.history, diff)
		t.currentWorkingTree.mergeUnverified(diff)
		t.UncommittedLeaves = []Hash{}
	}
	return nil
}

// Rollback rolls back one commit and reverts the tree to the previous state.
// Removes the most recent commit from the history.
func (t *Tree) Rollback() {
	_, t.history = PopFromPartialtree(t.history)
	t.currentWorkingTree.clear()
	for _, commit := range t.history {
		t.currentWorkingTree.mergeUnverified(commit)
	}
}

// uncommittedRoot calculates the root of the uncommitted changes as if they were committed.
// Will return the same hash as root of merkle tree after commit
func (t *Tree) uncommittedRoot() (Hash, error) {
	shadowTree, err := t.uncommittedDiff()
	if err != nil {
		return Hash{}, err
	}
	return shadowTree.Root(), nil
}

// UncommittedRootHex calculates the root of the uncommitted changes as if they were committed. Serializes
// the result as a hex string.
func (t *Tree) UncommittedRootHex() (string, error) {
	root, err := t.uncommittedRoot()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), nil
}

// Depth returns the tree depth. A tree depth is how many layers there is between the
// leaves and the root
func (t *Tree) Depth() int {
	return len(t.layerLeaves()) - 1
}

/// BaseLeaves returns a copy of the tree leaves - the base level of the tree.
func (t *Tree) BaseLeaves() []Hash {
	layers := t.layers()
	if len(layers) > 0 {
		return []Hash{}
	}
	return layers[0]
}

/// LeavesLen returns the number of leaves in the tree.
func (t *Tree) LeavesLen() uint32 {
	leaves := t.leaves()
	return uint32(len(leaves))
}

// leaves returns leaves of the current working tree
func (t *Tree) leaves() []Leaf {
	if len(t.layerLeaves()) > 0 {
		return t.layerLeaves()[0]
	}
	return []Leaf{}
}

// layers returns the whole tree, where the first layer is leaves and
// consequent layers are nodes.
func (t *Tree) layers() [][]Hash {
	return t.currentWorkingTree.layerNodes()
}

// layerLeaves returns leaves of the current working tree
func (t *Tree) layerLeaves() [][]Leaf {
	return t.currentWorkingTree.layers
}

/// uncommittedDiff creates a diff from a changes that weren't committed to the main tree yet. Can be used
/// to get uncommitted root or can be merged with the main tree
func (t *Tree) uncommittedDiff() (PartialTree, error) {
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, nil
	}
	commitedLeavesCount := t.LeavesLen()
	var shadowIndecies []uint32
	for i, _ := range t.UncommittedLeaves {
		shadowIndecies = append(shadowIndecies, commitedLeavesCount+uint32(i))
	}
	var shadowNodeLeaves []Leaf
	for i := 0; i < len(shadowIndecies); i++ {
		x := Leaf{Index: shadowIndecies[i], Hash: t.UncommittedLeaves[i]}
		shadowNodeLeaves = append(shadowNodeLeaves, x)
	}
	partialTreeLeaves := t.HelperNodeLeaves(shadowIndecies)
	leavesInNewTree := t.LeavesLen() + uint32(len(t.UncommittedLeaves))
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

// leafAtIndex returns leaf object at the index
func leafAtIndex(leavesAndHash []Leaf, index uint32) (Leaf, bool) {
	for _, l := range leavesAndHash {
		if l.Index == index {
			return l, true
		}
	}
	return Leaf{}, false
}

// layerAtIndex returns layer object at the index
func layerAtIndex(layers [][]Leaf, index uint32) ([]Leaf, bool) {
	if len(layers) > int(index) {
		return layers[index], true
	}
	return []Leaf{}, false
}

// sortLeavesByIndex sorts leaves by their index
func sortLeavesByIndex(li []Leaf) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index < li[j].Index })

}

// treeDepth returns the depth of a tree
func treeDepth(leaves_count uint32) uint32 {
	if leaves_count == 1 {
		return 1
	} else {
		return uint32(math.Ceil(math.Log2(float64(leaves_count))))
	}
}
