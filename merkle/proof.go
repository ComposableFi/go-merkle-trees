package merkle

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/ComposableFi/go-merkle-trees/types"
)

// Verify uses proof to verify that a given set of elements is contained in the original data
// set the proof was made for.
func (p Proof) Verify(expectedRoot []byte) (bool, error) {

	// extract root from proof
	extractedRoot, err := p.Root()
	if err != nil {
		return false, err
	}

	// return true if extracted root is uqual to expected root
	return bytes.Equal(extractedRoot, expectedRoot), nil
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

// proofLayers returns the proof layers by indices
func (p Proof) proofLayers(leafIndices []uint64) Layers {

	depth := treeDepth(p.totalLeavesCount)
	proofLayers := make(Layers, depth)

	// get uneven layers to remove any of uneven sibling indices in following loop
	unevenLayers := unevenLayersCountMap(p.totalLeavesCount)

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

			// set the proof leaf at index
			proofLeaves[j] = types.Leaf{
				Index: proofNodesIndices[j],
				Hash:  p.proofHashes[lastProofIndex],
			}

			// set the last checked index for next round of partent loop
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

	// get siblings
	siblingIndices := siblingIndecies(leafIndices)

	// remove from siblings if last node index is equal to eneven layer index
	lastNodeIndex := leafIndices[len(leafIndices)-1]
	if lastNodeIndex == unevenLeavesCount-1 {
		siblingIndices = popFromIndexQueue(siblingIndices)
	}

	return siblingIndices
}

// unevenLayersCountMap returns map of layer indices that are not even
func unevenLayersCountMap(totalLeavesCount uint64) map[uint64]uint64 {

	// set depth to prevent modification by loop
	depth := treeDepth(totalLeavesCount)

	unevenLayers := make(map[uint64]uint64)

	// loop until reach the full depth of tree
	for i := uint64(0); i < depth; i++ {

		// if count is not even, append it to result
		if !isEvenIndex(totalLeavesCount) {
			unevenLayers[i] = totalLeavesCount
		}

		// update the index to check and make it half
		totalLeavesCount = uint64(math.Ceil(float64(totalLeavesCount) / 2))
	}

	return unevenLayers
}
