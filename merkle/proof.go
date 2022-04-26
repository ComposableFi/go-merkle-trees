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
	treeDepth := treeDepth(p.totalLeavesCount)
	sortLeavesAscending(p.leaves)
	leafIndices := make([]uint64, len(p.leaves))
	for i := 0; i < len(p.leaves); i++ {
		leafIndices[i] = p.leaves[i].Index
	}
	proofIndicesLayers := proofIndeciesByLayers(leafIndices, p.totalLeavesCount)
	proofLayersCount := len(proofIndicesLayers)
	proofLayers := make(Layers, proofLayersCount)
	proofCopy := make([][]byte, len(p.proofHashes))
	copy(proofCopy, p.proofHashes)
	for i := 0; i < proofLayersCount; i++ {
		proofIndices := proofIndicesLayers[i]
		proofIndicesCount := len(proofIndices)
		proofHashes := make([][]byte, proofIndicesCount)
		for j := 0; j < proofIndicesCount; j++ {
			proofHashes[j] = proofCopy[0]
			proofCopy = proofCopy[1:]
		}
		proofLayers[i] = mapIndiceToLeaves(proofIndices, proofHashes)
	}

	if len(proofLayers) > 0 {
		firstLayer := proofLayers[0]
		firstLayer = append(firstLayer, p.leaves...)
		sortLeavesAscending(firstLayer)
		proofLayers[0] = firstLayer
	} else {
		proofLayers = append(proofLayers, p.leaves)
	}
	partialTree := NewPartialTree(p.hasher)
	PartialTree, err := partialTree.build(proofLayers, treeDepth)
	if err != nil {
		return []byte{}, err
	}
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

// proofIndeciesByLayers returns the proof indices by layers
func proofIndeciesByLayers(sortedLeafIndices []uint64, leavsCount uint64) [][]uint64 {
	depth := treeDepth(leavsCount)
	unevenLayers := unevenLayers(leavsCount)
	var proofIndices [][]uint64
	for layerIndex := uint64(0); layerIndex < depth; layerIndex++ {
		siblingIndices := siblingIndecies(sortedLeafIndices)
		leavesCount := unevenLayers[layerIndex]
		layerLastNodeIndex := sortedLeafIndices[len(sortedLeafIndices)-1]
		if layerLastNodeIndex == uint64(leavesCount)-1 {
			siblingIndices = popFromIndexQueue(siblingIndices)
		}

		proofNodesIndices := sliceDifference(siblingIndices, sortedLeafIndices)
		proofIndices = append(proofIndices, proofNodesIndices)
		sortedLeafIndices = parentIndecies(sortedLeafIndices)
	}
	return proofIndices
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

// mapIndiceToLeaves maps the indices and leaves of a tree
func mapIndiceToLeaves(indices []uint64, leavesHashes [][]byte) (result Leaves) {
	indicesLen := len(indices)
	result = make(Leaves, indicesLen)
	for i := 0; i < indicesLen; i++ {
		result[i] = types.Leaf{Index: indices[i], Hash: leavesHashes[i]}
	}
	return result
}
