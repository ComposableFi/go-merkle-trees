package merkle

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/ComposableFi/go-merkle-trees/helpers"
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
	sortLeavesByIndex(p.leaves)
	var leafIndices []uint32
	for _, l := range p.leaves {
		leafIndices = append(leafIndices, l.Index)
	}
	proofIndicesLayers := proofIndeciesByLayers(leafIndices, p.totalLeavesCount)
	var proofLayers [][]Leaf
	proofCopy := make([][]byte, len(p.proofHashes))
	copy(proofCopy, p.proofHashes)
	for _, proofIndices := range proofIndicesLayers {
		var proofHashes [][]byte
		for i := 0; i < len(proofIndices); i++ {
			proofHashes = append(proofHashes, proofCopy[0])
			proofCopy = proofCopy[1:]
		}
		m := MapIndiceAndLeaves(proofIndices, proofHashes)
		proofLayers = append(proofLayers, m)
	}

	if len(proofLayers) > 0 {
		firstLayer := proofLayers[0]
		firstLayer = append(firstLayer, p.leaves...)
		sortLeavesByIndex(firstLayer)
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
// bottom to top, as a vector of lower hex strings.
func (p Proof) ProofHashesHex() []string {
	var hexList []string
	for _, p := range p.proofHashes {
		hex := hex.EncodeToString(p)
		hexList = append(hexList, hex)
	}
	return hexList
}

// proofIndeciesByLayers returns the proof indices by layers
func proofIndeciesByLayers(sortedLeafIndices []uint32, leavsCount uint32) [][]uint32 {
	depth := treeDepth(leavsCount)
	unevenLayers := unevenLayers(leavsCount)
	var proofIndices [][]uint32
	for layerIndex := uint32(0); layerIndex < depth; layerIndex++ {
		siblingIndices := helpers.SiblingIndecies(sortedLeafIndices)
		leavesCount := unevenLayers[layerIndex]
		layerLastNodeIndex := sortedLeafIndices[len(sortedLeafIndices)-1]
		if layerLastNodeIndex == uint32(leavesCount)-1 {
			_, siblingIndices = helpers.PopFromUint32Queue(siblingIndices)
		}

		proofNodesIndices := helpers.Difference(siblingIndices, sortedLeafIndices)
		proofIndices = append(proofIndices, proofNodesIndices)
		sortedLeafIndices = helpers.ParentIndecies(sortedLeafIndices)
	}
	return proofIndices

}

// unevenLayers returns map of indices that are not even
func unevenLayers(treeLeavesCount uint32) map[uint32]uint32 {
	depth := treeDepth(treeLeavesCount)
	unevenLayers := make(map[uint32]uint32)
	for i := uint32(0); i < depth; i++ {
		unevenLayer := treeLeavesCount%2 != 0
		if unevenLayer {
			unevenLayers[i] = treeLeavesCount
		}
		treeLeavesCount = uint32(math.Ceil(float64(treeLeavesCount) / 2))
	}
	return unevenLayers
}
