package merkle

import (
	"errors"

	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/ComposableFi/go-merkle-trees/helpers"
	"github.com/ComposableFi/go-merkle-trees/types"
)

// build is a wrapper for buildTree
func (pt *PartialTree) build(partialLayers [][]types.Leaf, depth uint64) (PartialTree, error) {
	layers, err := pt.buildTree(partialLayers, depth)
	if err != nil {
		return PartialTree{}, err
	}
	return PartialTree{layers: layers}, nil
}

// buildTree is a general algorithm for building a partial tree. It can be used to extract root
// from merkle proof, or if a complete set of leaves provided as a first argument and no
// helper indices given, will construct the whole tree.
func (pt *PartialTree) buildTree(partialLayers [][]types.Leaf, fullTreeDepth uint64) ([][]types.Leaf, error) {
	reversedLayers := reverseLayers(partialLayers)
	var currentLayer []types.Leaf
	var partialTree [][]types.Leaf
	for i := uint64(0); i < fullTreeDepth; i++ {

		if len(reversedLayers) > 0 {
			var nodes []types.Leaf
			nodes, reversedLayers = PopFromLeafQueue(reversedLayers)
			currentLayer = append(currentLayer, nodes...)
		}

		sortLeavesByIndex(currentLayer)

		partialTree = append(partialTree, currentLayer)

		var indices []uint64
		var nodes [][]byte
		for i := 0; i < len(currentLayer); i++ {
			indices = append(indices, currentLayer[i].Index)
			nodes = append(nodes, currentLayer[i].Hash)
		}
		// freeup for next round
		currentLayer = make([]types.Leaf, 0)

		parentLayerIndices := helpers.ParentIndecies(indices)

		for i := 0; i < len(parentLayerIndices); i++ {
			parnetNodeIndex := parentLayerIndices[i]
			leftIndex := i * 2
			if len(nodes) > leftIndex {
				leftHash := nodes[leftIndex]
				rightIndex := i*2 + 1

				var hash, rightHash []byte
				if len(nodes) > rightIndex {
					rightHash = nodes[rightIndex]
				} else {
					rightHash = nil
				}
				var err error
				hash, err = hasher.MergeAndHash(pt.hasher, leftHash, rightHash)
				if err != nil {
					return [][]types.Leaf{}, err
				}

				currentLayer = append(currentLayer, types.Leaf{
					Index: parnetNodeIndex,
					Hash:  hash,
				})
			} else {
				return [][]types.Leaf{}, errors.New("not enough helper nodes")
			}
		}

	}
	if len(currentLayer) > 0 {
		partialTree = append(partialTree, currentLayer)
	}
	return partialTree, nil
}

// Root returns the root of the tree
func (pt *PartialTree) Root() []byte {
	if len(pt.layers) > 0 {
		lastLayer := pt.layers[len(pt.layers)-1]
		firstItem := lastLayer[0]
		return firstItem.Hash
	}
	return nil
}

// contains checks if a node index is present in a layer
func (pt *PartialTree) contains(layerIndex, nodeIndex uint64) bool {
	layer, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {
		for _, l := range layer {
			if nodeIndex == l.Index {
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
func (pt *PartialTree) mergeUnverified(other PartialTree) {
	depthDifference := len(other.layers) - len(pt.layers)
	var combinedTreeSize uint64
	if depthDifference > 0 {
		combinedTreeSize = uint64(len(other.layers))
	} else {
		combinedTreeSize = uint64(len(pt.layers))
	}

	for layerIndex := uint64(0); layerIndex < combinedTreeSize; layerIndex++ {
		var combinedLayer, filteredLayer []types.Leaf

		selfLayer, ok := layerAtIndex(pt.layers, uint64(layerIndex))
		if ok {
			for _, node := range selfLayer {
				if !other.contains(layerIndex, node.Index) {
					filteredLayer = append(filteredLayer, node)
				}
			}
			combinedLayer = append(combinedLayer, filteredLayer...)

		}

		otherLayer, ok := layerAtIndex(other.layers, layerIndex)
		if ok {
			combinedLayer = append(combinedLayer, otherLayer...)
		}

		sortLeavesByIndex(otherLayer)

		pt.upsertLayer(layerIndex, combinedLayer)

	}
}

// upsertLayer replaces layer at a given index with a new layer. Used during tree merge
func (pt *PartialTree) upsertLayer(layerIndex uint64, newLayer []types.Leaf) {
	_, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {
		pt.layers[layerIndex] = newLayer
	} else {
		pt.layers = append(pt.layers, newLayer)
	}

}

// layerNodes returns all hashes of all layers
func (pt *PartialTree) layerNodes() [][][]byte {
	var allHashes [][][]byte
	for _, l := range pt.getLayers() {
		var layerHashes [][]byte
		for _, h := range l {
			layerHashes = append(layerHashes, h.Hash)
		}
		allHashes = append(allHashes, layerHashes)
	}
	return allHashes
}

// getLayers returns partial tree layers
func (pt *PartialTree) getLayers() [][]types.Leaf {
	return pt.layers
}

// clear clears all elements in the tree
func (pt *PartialTree) clear() {
	pt.layers = [][]types.Leaf{}
}

// reverseLayers reverses a slice of types.Leaf slice
func reverseLayers(s [][]types.Leaf) [][]types.Leaf {
	a := make([][]types.Leaf, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}
