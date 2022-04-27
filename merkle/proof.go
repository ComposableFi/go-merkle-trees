package merkle

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/ComposableFi/go-merkle-trees/types"
)

// Verify uses proof to verify that a given set of elements is contained in the original data
// set the proof was made for.
func (p Proof) Verify(root []byte) (bool, error) {
	extractedRoot, err := p.Root()
	if err != nil {
		return false, err
	}
	return bytes.Equal(extractedRoot, root), nil
}

// Root calculates Merkle root based on provided leaves and proof hashes. Used inside the
// Verify method, but sometimes can be used on its own.
func (p Proof) Root() ([]byte, error) {

	sortLeavesAscending(p.leaves)

	// extract proof leaves indices
	leafIndices := make([]uint64, len(p.leaves))
	for i := 0; i < len(p.leaves); i++ {
		leafIndices[i] = p.leaves[i].Index
	}

	proofLayers := p.proofLayers(leafIndices)

	if len(proofLayers) > 0 {
		// set the first layer as proof
		firstLayer := proofLayers[0]
		firstLayer = append(firstLayer, p.leaves...)
		sortLeavesAscending(firstLayer)
		proofLayers[0] = firstLayer
	} else {
		proofLayers = append(proofLayers, p.leaves)
	}

	// build the partial tree from proof leaves
	treeDepth := treeDepth(p.totalLeavesCount)
	partialTree := NewPartialTree(p.hasher)
	PartialTree, err := partialTree.build(proofLayers, treeDepth)
	if err != nil {
		return []byte{}, err
	}

	// return root of partial tree
	return PartialTree.Root(), err
}

// RootHex calculates the root and serializes it into a hex string.
func (p Proof) RootHex() (string, error) {
	root, err := p.Root()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), nil
}

// ProofHashes returns all hashes from the proof, sorted from the left to right,
// bottom to top.
func (p Proof) ProofHashes() [][]byte {
	return p.proofHashes
}

// ProofHashesHex returns all hashes from the proof, sorted from the left to right,
// bottom to top, as a slice of lower hex strings.
func (p Proof) ProofHashesHex() []string {
	hashesLen := len(p.proofHashes)
	hexList := make([]string, hashesLen)
	for i := 0; i < hashesLen; i++ {
		hexList[i] = hex.EncodeToString(p.proofHashes[i])
	}
	return hexList
}

// proofLayers returns the proof indices by layers
func (p Proof) proofLayers(leafIndices []uint64) Layers {
	depth := treeDepth(p.totalLeavesCount)
	unevenLayers := unevenLayers(p.totalLeavesCount)
	proofLayers := make(Layers, depth)

	// copied proof index
	lastProofIndex := 0

	// loop through depth of tree and update proof indices
	for layerIndex := uint64(0); layerIndex < depth; layerIndex++ {
		// get siblings without last event index
		siblingIndices := popLastEvenIndexFromSiblings(leafIndices, unevenLayers[layerIndex])

		// append proof indices inot the result
		proofNodesIndices := extractNewIndicesFromSiblings(siblingIndices, leafIndices)

		// set the proof leaves from proof hashes
		proofIndicesCount := len(proofNodesIndices)
		proofLeaves := make(Leaves, proofIndicesCount)
		for j := 0; j < proofIndicesCount; j++ {
			proofLeaves[j] = types.Leaf{
				Index: proofNodesIndices[j],
				Hash:  p.proofHashes[lastProofIndex],
			}
			lastProofIndex++
		}

		// use proof indices and hash to set the layer leaves
		proofLayers[layerIndex] = proofLeaves

		// go one level up in leaves
		leafIndices = parentIndecies(leafIndices)
	}
	return proofLayers
}

// popLastEvenIndex removes last uneven index from siblings
func popLastEvenIndexFromSiblings(leafIndices []uint64, unevenLeavesCount uint64) []uint64 {
	siblingIndices := siblingIndecies(leafIndices)
	layerLastNodeIndex := leafIndices[len(leafIndices)-1]
	if layerLastNodeIndex == unevenLeavesCount-1 {
		siblingIndices = popFromIndexQueue(siblingIndices)
	}
	return siblingIndices
}

// unevenLayers returns map of indices that are not even
func unevenLayers(treeLeavesCount uint64) map[uint64]uint64 {
	depth := treeDepth(treeLeavesCount)
	unevenLayers := make(map[uint64]uint64)
	for i := uint64(0); i < depth; i++ {
		if !isEvenIndex(treeLeavesCount) {
			unevenLayers[i] = treeLeavesCount
		}
		treeLeavesCount = uint64(math.Ceil(float64(treeLeavesCount) / 2))
	}
	return unevenLayers
}
