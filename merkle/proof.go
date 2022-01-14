package merkle

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/ComposableFi/merkle-go/helpers"
)

// func (p Proof) fromBytes(bytes []byte) (PartialTree, error) {
// 	return p.deserialize(bytes)
// }

// func (p Proof) deserialize(bytes []byte) (PartialTree, error) {
// 	return p.serializer.Deserialize(bytes)
// }

func (p Proof) Verify(root Hash, leafTuples []Leaf, totalLeavesCount int) (bool, error) {
	extractedRoot, err := p.GetRoot(leafTuples, int(totalLeavesCount))
	if err != nil {
		return false, err
	}
	return bytes.Equal(extractedRoot, root), nil
}

func (p Proof) GetRoot(leafTuples []Leaf, totalLeavesCount int) (Hash, error) {
	treeDepth := getTreeDepth(totalLeavesCount)
	sortLeavesByIndex(leafTuples)
	var leafIndices []uint32
	for _, l := range leafTuples {
		leafIndices = append(leafIndices, l.Index)
	}
	proofIndicesLayers := proofIndeciesByLayers(leafIndices, totalLeavesCount)
	var proofLayers [][]Leaf
	proofCopy := make([]Hash, len(p.proofHashes))
	copy(proofCopy, p.proofHashes)
	for _, proofIndices := range proofIndicesLayers {
		var proofHashes []Hash
		for i := 0; i < len(proofIndices); i++ {
			proofHashes = append(proofHashes, proofCopy[0])
			proofCopy = proofCopy[1:]
		}
		m := MapIndiceAndLeaves(proofIndices, proofHashes)
		proofLayers = append(proofLayers, m)
	}

	if len(proofLayers) > 0 {
		firstLayer := proofLayers[0]
		firstLayer = append(firstLayer, leafTuples...)
		sortLeavesByIndex(firstLayer)
		proofLayers[0] = firstLayer

	} else {
		proofLayers = append(proofLayers, leafTuples)
	}
	partialTree := NewPartialTree(p.hasher)
	PartialTree, err := partialTree.build(proofLayers, treeDepth)
	if err != nil {
		return Hash{}, err
	}
	return PartialTree.GetRoot(), err
}

func (p Proof) GetRootHex(leafTuples []Leaf, totalLeavesCount int) (string, error) {
	root, err := p.GetRoot(leafTuples, totalLeavesCount)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), nil
}

func (p Proof) ProofHashes() []Hash {
	return p.proofHashes
}

func (p Proof) ProofHashesHex() []string {
	var hexList []string
	for _, p := range p.proofHashes {
		hex := hex.EncodeToString(p)
		hexList = append(hexList, hex)
	}
	return hexList
}

func proofIndeciesByLayers(sortedLeafIndices []uint32, leavsCount int) [][]uint32 {
	depth := getTreeDepth(leavsCount)
	unevenLayers := unevenLayers(leavsCount)
	var proofIndices [][]uint32
	for layerIndex := 0; layerIndex < depth; layerIndex++ {
		siblingIndices := helpers.GetSiblingIndecies(sortedLeafIndices)
		leavesCount := unevenLayers[layerIndex]
		layerLastNodeIndex := sortedLeafIndices[len(sortedLeafIndices)-1]
		if layerLastNodeIndex == uint32(leavesCount)-1 {
			_, siblingIndices = helpers.PopFromUint32Queue(siblingIndices)
		}

		proofNodesIndices := helpers.Difference(siblingIndices, sortedLeafIndices)
		proofIndices = append(proofIndices, proofNodesIndices)
		sortedLeafIndices = helpers.GetParentIndecies(sortedLeafIndices)
	}
	return proofIndices

}

func unevenLayers(treeLeavesCount int) map[int]int {
	depth := getTreeDepth(treeLeavesCount)
	unevenLayers := make(map[int]int)
	for i := 0; i < depth; i++ {
		unevenLayer := treeLeavesCount%2 != 0
		if unevenLayer {
			unevenLayers[i] = treeLeavesCount
		}
		treeLeavesCount = int(math.Ceil(float64(treeLeavesCount) / 2))
	}
	return unevenLayers
}
