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
	layers := t.layers()
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

// helperNodesHashes returns helper nodes required to build a partial tree for the given indices
// to be able to extract a root from it. Useful in constructing Merkle proofs
func (t *Tree) helperNodesHashes(leafIndices []uint64) [][]byte {
	var helperNodesHashes [][]byte
	for _, layer := range t.helperNodeLayers(leafIndices) {
		for _, li := range layer {
			helperNodesHashes = append(helperNodesHashes, li.Hash)
		}
	}
	return helperNodesHashes
}

// helperNodeLayers gets all helper nodes required to build a partial merkle tree for the given indices,
// cloning all required hashes into the resulting slice.
func (t *Tree) helperNodeLayers(leafIndices []uint64) Layers {
	var helperNodes Layers
	for _, treeLayer := range t.layers() {
		siblings := siblingIndecies(leafIndices)
		helperIndices := sliceDifference(siblings, leafIndices)

		var helpersLayer Leaves
		for _, idx := range helperIndices {
			leaf, found := leafAtIndex(treeLayer, idx)
			if found {
				helpersLayer = append(helpersLayer, leaf)
			}
		}

		helperNodes = append(helperNodes, helpersLayer)

		leafIndices = parentIndecies(leafIndices)
	}
	return helperNodes
}

// Proof Returns the Merkle proof required to prove the inclusion of items in a data set.
func (t *Tree) Proof(leafIndices []uint64) Proof {
	leavesLen := t.leavesLen()
	leaves := t.leaves()
	var proofLeaves Leaves

	for _, index := range leafIndices {
		for _, leaf := range leaves {
			if leaf.Index == index {
				proofLeaves = append(proofLeaves, leaf)
				break
			}
		}
	}
	return NewProof(proofLeaves, t.helperNodesHashes(leafIndices), leavesLen, t.hasher)
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
		t.currentWorkingTree.mergeUnverified(diff)
		t.UncommittedLeaves = [][]byte{}
	}
	return nil
}

// uncommittedRoot calculates the root of the uncommitted changes as if they were committed.
// Will return the same hash as root of merkle tree after commit
func (t *Tree) uncommittedRoot() ([]byte, error) {
	uncommitedTree, err := t.uncommittedDiff()
	if err != nil {
		return []byte{}, err
	}
	return uncommitedTree.Root(), nil
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
	return len(t.layers()) - 1
}

// baseLeaves returns a copy of the tree leaves - the base level of the tree.
func (t *Tree) baseLeaves() [][]byte {
	layers := t.layersNodesHashes()
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
func (t *Tree) leaves() Leaves {
	if len(t.layers()) > 0 {
		return t.layers()[0]
	}
	return Leaves{}
}

// layersNodesHashes returns the whole tree, where the first layer is leaves and
// consequent layersNodesHashes are nodes.
func (t *Tree) layersNodesHashes() [][][]byte {
	return t.currentWorkingTree.layerNodesHashes()
}

// layers returns leaves of the current working tree
func (t *Tree) layers() Layers {
	return t.currentWorkingTree.layers
}

// uncommittedDiff creates a diff from a changes that weren't committed to the main tree yet. Can be used
// to get uncommitted root or can be merged with the main tree
func (t *Tree) uncommittedDiff() (PartialTree, error) {
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, nil
	}

	partialTreeLayers, uncommittedTreeDepth := t.uncommitedPartialTreeLayers()

	tree := NewPartialTree(t.hasher)
	return tree.build(partialTreeLayers, uncommittedTreeDepth)
}

// uncommitedPartialTreeLayers calculates reserved indices and leaves then returns uncommited partial tree layers
func (t *Tree) uncommitedPartialTreeLayers() (Layers, uint64) {
	reservedIndecies := t.getUncommitedReservedIndecies()
	reservedNodeLeaves := t.getUncommitedReservedLeaves(reservedIndecies)

	partialTreeLayers := t.helperNodeLayers(reservedIndecies)
	partialTreeLayers = appendUncommitedReservedLayers(partialTreeLayers, reservedNodeLeaves)

	leavesInNewTree := t.leavesLen() + uint64(len(t.UncommittedLeaves))
	uncommittedTreeDepth := treeDepth(leavesInNewTree)
	return partialTreeLayers, uncommittedTreeDepth
}

// getUncommitedReservedIndecies returns uncommited reserved indices of the uncommited leaves
func (t *Tree) getUncommitedReservedIndecies() []uint64 {
	if len(t.UncommittedLeaves) == 0 {
		return []uint64{}
	}

	commitedLeavesCount := t.leavesLen()
	unCommitedLeavesCount := len(t.UncommittedLeaves)
	reservedIndecies := make([]uint64, unCommitedLeavesCount)
	for i := 0; i < unCommitedLeavesCount; i++ {
		reservedIndecies[i] = commitedLeavesCount + uint64(i)
	}
	return reservedIndecies
}

// getUncommitedReservedLeaves returns uncommited reserved leaves of the uncommited leaves
func (t *Tree) getUncommitedReservedLeaves(reservedIndecies []uint64) Leaves {
	indicesCount := len(reservedIndecies)
	reservedNodeLeaves := make(Leaves, indicesCount)
	for i := 0; i < indicesCount; i++ {
		reservedNodeLeaves[i] = types.Leaf{Index: reservedIndecies[i], Hash: t.UncommittedLeaves[i]}
	}
	return reservedNodeLeaves
}

// leafAtIndex returns types.Leaf object at the index
func leafAtIndex(leavesAndHash Leaves, index uint64) (types.Leaf, bool) {
	for i := 0; i < len(leavesAndHash); i++ {
		l := leavesAndHash[i]
		if l.Index == index {
			return l, true
		}
	}
	return types.Leaf{}, false
}

// layerAtIndex returns layer object at the index
func layerAtIndex(layers Layers, index uint64) (Leaves, bool) {
	if len(layers) > int(index) {
		return layers[index], true
	}
	return Leaves{}, false
}

// sortLeavesAscending sorts leaves by their index
func sortLeavesAscending(li Leaves) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index < li[j].Index })
}

// treeDepth returns the depth of a tree
func treeDepth(leavesCount uint64) uint64 {
	if leavesCount == 1 {
		return 1
	}
	return uint64(math.Ceil(math.Log2(float64(leavesCount))))
}

func appendUncommitedReservedLayers(partialTreeLayers Layers, reservedNodeLeaves Leaves) Layers {
	if len(partialTreeLayers) == 0 {
		partialTreeLayers = append(partialTreeLayers, reservedNodeLeaves)
	} else {
		firstLayer := partialTreeLayers[0]
		firstLayer = append(firstLayer, reservedNodeLeaves...)
		sortLeavesAscending(firstLayer)
		partialTreeLayers[0] = firstLayer
	}
	return partialTreeLayers
}

// popPartialtree pops last element in a partial tree slice
func popPartialtree(slice []PartialTree) (PartialTree, []PartialTree) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}
