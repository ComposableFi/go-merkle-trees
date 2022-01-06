package merkle

import (
	"errors"

	"github.com/ComposableFi/merkle-go/helpers"
)

func (pt PartialTree) fromLeaves(leaves []Hash) (PartialTree, error) {
	var leafTuples [][]leafIndexAndHash
	for i := 0; i < len(leaves); i++ {
		leafTuples = append(leafTuples, []leafIndexAndHash{
			{
				index: uint32(i),
				hash:  leaves[i],
			},
		})
	}
	return pt.build(leafTuples, getTreeDepth(len(leaves)))
}

func (pt PartialTree) build(partialLayers [][]leafIndexAndHash, depth int) (PartialTree, error) {
	layers, err := pt.buildTree(partialLayers, depth)
	if err != nil {
		return PartialTree{}, err
	}
	return PartialTree{layers: layers}, nil
}

func (pt PartialTree) buildTree(partialLayers [][]leafIndexAndHash, fullTreeDepth int) ([][]leafIndexAndHash, error) {
	reversedLayers := reverse(partialLayers)
	var currentLayer []leafIndexAndHash
	var partialTree [][]leafIndexAndHash
	for i := 0; i < fullTreeDepth; i++ {
		var nodes []leafIndexAndHash

		if len(reversedLayers) > 0 {
			nodes, reversedLayers = PopFromLeafHashQueue(reversedLayers)
			currentLayer = append(currentLayer, nodes...)
		}

		sortLeafAndHashByIndex(currentLayer)

		partialTree = append(partialTree, currentLayer)

		var indices []uint32
		for _, l := range currentLayer {
			indices = append(indices, l.index)
		}
		// freeup for next round
		currentLayer = []leafIndexAndHash{}

		parentLayerIndices := helpers.GetParentIndecies(indices)

		for i := 0; i < len(parentLayerIndices); i++ {
			parnetNodeIndex := parentLayerIndices[i]
			if leftNode, err := getLeafAndHashAtIndex(nodes, uint32(i*2)); err != nil {
				rightNode, _ := getLeafAndHashAtIndex(nodes, uint32(i*2+1))
				hash, err := pt.hasher.ConcatAndHash(leftNode.hash, rightNode.hash)
				if err != nil {
					return [][]leafIndexAndHash{}, err
				}
				currentLayer = append(currentLayer, leafIndexAndHash{
					index: parnetNodeIndex,
					hash:  hash,
				})
			} else {
				return [][]leafIndexAndHash{}, errors.New("not enought helper nodes")
			}
		}

	}
	partialTree = append(partialTree, currentLayer)
	return partialTree, nil
}

func (pt PartialTree) depth() int {
	return len(pt.layers) - 1
}

func (pt PartialTree) getRoot() Hash {
	lastLayer := pt.layers[len(pt.layers)-1]
	firstItem := lastLayer[0]
	return firstItem.hash
}

func (pt PartialTree) contains(layerIndex, nodeIndex uint32) bool {
	layer, ok := getLayerAtIndex(pt.layers, layerIndex)
	if ok {
		for _, l := range layer {
			if nodeIndex == l.index {
				return true
			}
		}
	}
	return false
}

/// Consumes other partial tree into itself, replacing any conflicting nodes with nodes from
/// `other` in the process. Doesn't rehash the nodes, so the integrity of the result is
/// not verified. It gives an advantage in speed, but should be used only if the integrity of
/// the tree can't be broken, for example, it is used in the `.commit` method of the
/// `MerkleTree`, since both partial trees are essentially constructed in place and there's
/// no need to verify integrity of the result.
func (pt PartialTree) mergeUnverified(other PartialTree) {
	depthDifference := len(other.layers) - len(pt.layers)
	var combinedTreeSize uint32
	if depthDifference > 0 {
		combinedTreeSize = uint32(len(other.layers))
	} else {
		combinedTreeSize = uint32(len(pt.layers))
	}

	for layerIndex := uint32(0); layerIndex < combinedTreeSize; layerIndex++ {
		var combinedLayer, filteredLayer []leafIndexAndHash

		selfLayer, ok := getLayerAtIndex(pt.layers, uint32(layerIndex))
		if ok {
			for _, node := range selfLayer {
				if !other.contains(layerIndex, node.index) {
					filteredLayer = append(filteredLayer, node)
				}
			}
			combinedLayer = append(combinedLayer, filteredLayer...)

		}

		otherLayer, ok := getLayerAtIndex(other.layers, layerIndex)
		if ok {
			combinedLayer = append(combinedLayer, otherLayer...)
		}

		sortLeafAndHashByIndex(otherLayer)

		pt.upsertLayer(layerIndex, combinedLayer)

	}
}

func (pt PartialTree) layerNodes() [][]Hash {
	var allHashes [][]Hash
	for _, l := range pt.getLayers() {
		var layerHashes []Hash
		for _, h := range l {
			layerHashes = append(layerHashes, h.hash)
		}
		allHashes = append(allHashes, layerHashes)
	}
	return allHashes
}

func (pt PartialTree) getLayers() [][]leafIndexAndHash {
	return pt.layers
}

func (pt PartialTree) clear() {
	pt.layers = [][]leafIndexAndHash{}
}

func (pt PartialTree) upsertLayer(layerIndex uint32, newLayer []leafIndexAndHash) {
	layer, ok := getLayerAtIndex(pt.layers, layerIndex)
	if ok {
		layer = []leafIndexAndHash{}
		layer = append(layer, newLayer...)
	} else {
		pt.layers = append(pt.layers, newLayer)
	}

}

func reverse(s [][]leafIndexAndHash) [][]leafIndexAndHash {
	a := make([][]leafIndexAndHash, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}
