// Package merkle is responsible for creating the tree and verify the proof
package merkle

import (
	"encoding/hex"
	"math"
	"sort"

	"github.com/ComposableFi/go-merkle-trees/types"
)

// FromLeaves clones the leaves and builds the tree from them
func (t Tree) FromLeaves(leaves [][]byte) (Tree, error) {

	// populate initial tree leaves
	t.append(leaves)

	// create tree
	err := t.commit()
	if err != nil {
		return Tree{}, err
	}

	// return tree
	return t, nil
}

// Root returns the tree root - the top hash of the tree. Used in the inclusion proof verification.
func (t *Tree) Root() []byte {

	layers := t.layers()
	if len(layers) > 0 {

		// get first node of last layer
		lastLayer := layers[len(layers)-1]
		firstNode := lastLayer[0]

		// return first node hash as root
		return firstNode.Hash
	}

	return []byte{}
}

// RootHex returns a hex encoded string instead of
func (t *Tree) RootHex() string {

	// get root
	root := t.Root()

	// convert to hex string and return
	return hex.EncodeToString(root)
}

// currentLayersWithSiblingsHashes returns sibling leaves required to build a partial tree for the given indices
// to be able to extract a root from it. Useful in constructing Merkle proofs
func (t *Tree) currentLayersWithSiblingsHashes(leafIndices []uint64) [][]byte {

	// loop through all layers and siblings hashes
	var newSiblingAndExistingHashes [][]byte
	for _, layer := range t.currentLayersWithSiblings(leafIndices) {
		for _, li := range layer {
			newSiblingAndExistingHashes = append(newSiblingAndExistingHashes, li.Hash)
		}
	}

	return newSiblingAndExistingHashes
}

// currentLayersWithSiblings gets all sibling layers required to build a partial merkle tree for the given indices,
// cloning all required hashes into the resulting slice.
func (t *Tree) currentLayersWithSiblings(leafIndices []uint64) Layers {

	var layersNodesWithSiblings Layers
	for _, layer := range t.layers() {
		// get siblings of leaf indices and extract newly created indices
		siblings := siblingIndecies(leafIndices)

		// detect newly extracted siblings
		newSiblingIndices := extractNewIndicesFromSiblings(siblings, leafIndices)

		// get the exisitng leaves in the layer with sibling indecies
		var existingLeavesInTree Leaves
		for i := 0; i < len(newSiblingIndices); i++ {

			leafIndex := newSiblingIndices[i]
			leaf, found := leafAtIndex(layer, leafIndex)
			if found {

				// append new sibling index
				existingLeavesInTree = append(existingLeavesInTree, leaf)
			}
		}

		// append to result
		layersNodesWithSiblings = append(layersNodesWithSiblings, existingLeavesInTree)

		// go one level up in the leafInfices
		leafIndices = parentIndecies(leafIndices)
	}

	return layersNodesWithSiblings
}

// Proof Returns the Merkle proof required to prove the inclusion of items in a data set.
func (t *Tree) Proof(proofIndices []uint64) Proof {
	leavesLen := t.leavesLen()
	leaves := t.leaves()

	// make proof leaves from proof indices
	var proofLeaves Leaves
	for i := 0; i < len(proofIndices); i++ {
		for j := 0; j < len(leaves); j++ {
			leaf := leaves[j]
			if leaf.Index == proofIndices[i] {
				proofLeaves = append(proofLeaves, leaf)
				break
			}
		}
	}

	// get all hashes of leaves and their siblings
	siblingProofHashes := t.currentLayersWithSiblingsHashes(proofIndices)

	// create new proof object using proof leaves and hashes
	return NewProof(proofLeaves, siblingProofHashes, leavesLen, t.hasher)
}

// insert inserts a new types.Leaf. Please note it won't modify the root just yet; For the changes
// to be applied to the root, commit method should be called first. To get the
// root of the new tree without applying the changes, you can use
func (t *Tree) insert(leaf []byte) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaf)
}

// append appends leaves to the tree.
func (t *Tree) append(leaves [][]byte) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaves...)
}

// commit commits the changes made by insert and append
// and modifies the root.
func (t *Tree) commit() error {

	// get difference committed and not committed tree layers
	diff, err := t.uncommittedDiff()
	if err != nil {
		return err
	}

	// if there is new layers update the tree
	if len(diff.layers) > 0 {

		// merge existing and newly created partial tree
		t.currentWorkingTree.mergeUnverifiedLayers(diff)

		// free up the uncommitted leaves after storing the tree
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

	// get uncommited root
	root, err := t.uncommittedRoot()
	if err != nil {
		return "", err
	}

	// convert to hex string and return
	return hex.EncodeToString(root), nil
}

// depth returns the tree depth. A tree depth is how many layers there is between the
// leaves and the root
func (t *Tree) depth() int {
	return len(t.layers()) - 1
}

// baseLeaves returns a copy of the tree leaves - the base level of the tree.
func (t *Tree) baseLeaves() [][]byte {

	// get all hashes of leaves of all layersLeavesHashes
	layersLeavesHashes := t.layersNodesHashes()

	// if leaves are available
	if len(layersLeavesHashes) > 0 {
		return layersLeavesHashes[0]
	}

	return [][]byte{}
}

// leavesLen returns the number of leaves in the tree.
func (t *Tree) leavesLen() uint64 {
	leaves := t.leaves()
	return uint64(len(leaves))
}

// leaves returns leaves of the first layer that has the complete tree
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

	// if there is no uncommitted leaves, there is no more partial
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, nil
	}

	// get uncommitted partial layer
	partialTreeLayers, uncommittedTreeDepth := t.uncommitedPartialTreeLayers()

	// build partial tree and return
	tree := NewPartialTree(t.hasher)
	return tree.build(partialTreeLayers, uncommittedTreeDepth)
}

// uncommitedPartialTreeLayers calculates reserved indices and leaves then returns uncommitted partial tree layers
func (t *Tree) uncommitedPartialTreeLayers() (Layers, uint64) {

	// reserve indices for uncommitted leaves
	reservedIndecies := t.getUncommitedReservedIndecies()

	// extract uncommitted leaves from uncommitted indices
	reservedLeaves := t.getUncommitedReservedLeaves(reservedIndecies)

	// update layers with new siblings of each layer
	partialTreeLayers := t.currentLayersWithSiblings(reservedIndecies)

	// upsert partial layer by uncommitted reseved nodes
	partialTreeLayers = upsertUncommitedReservedLayers(partialTreeLayers, reservedLeaves)

	// calculate new tree depth
	leavesInNewTree := t.leavesLen() + uint64(len(t.UncommittedLeaves))
	uncommittedTreeDepth := treeDepth(leavesInNewTree)

	return partialTreeLayers, uncommittedTreeDepth
}

// getUncommitedReservedIndecies returns uncommitted reserved indices of the uncommitted leaves
func (t *Tree) getUncommitedReservedIndecies() []uint64 {

	// if there are no uncommitted leaves there is nothing to reserve
	if len(t.UncommittedLeaves) == 0 {
		return []uint64{}
	}

	commitedLeavesCount := t.leavesLen()
	unCommitedLeavesCount := len(t.UncommittedLeaves)

	// populate uncommitted indices according to the last committed leaves indices
	reservedIndecies := make([]uint64, unCommitedLeavesCount)
	for i := 0; i < unCommitedLeavesCount; i++ {
		reservedIndecies[i] = commitedLeavesCount + uint64(i)
	}

	return reservedIndecies
}

// getUncommitedReservedLeaves returns uncommitted reserved leaves of the uncommitted leaves
func (t *Tree) getUncommitedReservedLeaves(reservedIndecies []uint64) Leaves {

	// read uncommitted leaves hashes and set into reserved leaves by indices
	indicesCount := len(reservedIndecies)
	reservedLeaves := make(Leaves, indicesCount)
	for i := 0; i < indicesCount; i++ {
		reservedLeaves[i] = types.Leaf{Index: reservedIndecies[i], Hash: t.UncommittedLeaves[i]}
	}

	return reservedLeaves
}

// leafAtIndex returns leaf object by index
func leafAtIndex(leaves Leaves, index uint64) (types.Leaf, bool) {

	// loop through leaves and return leaf at certain index
	for i := 0; i < len(leaves); i++ {
		l := leaves[i]
		if l.Index == index {
			return l, true
		}
	}

	return types.Leaf{}, false
}

// layerAtIndex returns layer object by intex
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

	// math formula for tree depth
	// logarithm of the number of leaf nodes in the tree
	// https://en.wikipedia.org/wiki/Merkle_tree
	depth := math.Log2(float64(leavesCount))

	// round and return
	return uint64(math.Ceil(depth))
}

func upsertUncommitedReservedLayers(partialTreeLayers Layers, reservedNodeLeaves Leaves) Layers {

	if len(partialTreeLayers) == 0 {
		// no partial layers available yet, so we ned to create one
		partialTreeLayers = append(partialTreeLayers, reservedNodeLeaves)
	} else {
		// get first layer and append leaves
		firstLayer := partialTreeLayers[0]
		firstLayer = append(firstLayer, reservedNodeLeaves...)

		// sort leaves after new nodes addition
		sortLeavesAscending(firstLayer)

		// update the first layer
		partialTreeLayers[0] = firstLayer
	}

	return partialTreeLayers
}
