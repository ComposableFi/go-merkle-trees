package merkle

import (
	"errors"

	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/ComposableFi/go-merkle-trees/types"
)

// build is a wrapper for buildTree
func (pt *PartialTree) build(partialLayers Layers, depth uint64) (PartialTree, error) {
	layers, err := pt.buildTree(partialLayers, depth)
	if err != nil {
		return PartialTree{}, err
	}
	return PartialTree{layers: layers}, nil
}

// buildTree is a general algorithm for building a partial tree. It can be used to extract root
// from merkle proof, or if a complete set of leaves provided as a first argument and no
// helper indices given, will construct the whole tree.
func (pt *PartialTree) buildTree(partialLayers Layers, fullTreeDepth uint64) (Layers, error) {
	reversedLayers := reverseLayers(partialLayers)
	var currentLayer Leaves
	var partialTree Layers
	for i := uint64(0); i < fullTreeDepth; i++ {

		if len(reversedLayers) > 0 {
			var nodes Leaves
			nodes, reversedLayers = popLayer(reversedLayers)
			currentLayer = append(currentLayer, nodes...)
		}

		sortLeavesAscending(currentLayer)

		partialTree = append(partialTree, currentLayer)

		indices, hashes := extractIndicesAndHashes(currentLayer)

		// freeup for next round
		currentLayer = make(Leaves, 0)

		parentIndices := parentIndecies(indices)

		for i := 0; i < len(parentIndices); i++ {
			parnetNodeIndex := parentIndices[i]
			leftIndex := getLeftIndex(i)
			if len(hashes) > leftIndex {
				rightIndex := getRightIndex(i)

				leftHash := hashes[leftIndex]
				var rightHash []byte
				if len(hashes) > rightIndex {
					rightHash = hashes[rightIndex]
				} else {
					rightHash = nil
				}

				hash, err := hasher.MergeAndHash(pt.hasher, leftHash, rightHash)
				if err != nil {
					return Layers{}, err
				}

				currentLayer = append(currentLayer, types.Leaf{
					Index: parnetNodeIndex,
					Hash:  hash,
				})
			} else {
				return Layers{}, errors.New("not enough helper nodes")
			}
		}

	}
	if len(currentLayer) > 0 {
		partialTree = append(partialTree, currentLayer)
	}
	return partialTree, nil
}

// Root returns the root of the tree, it is the first item hash of the last layer
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
	layerLeaves, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {
		for i := 0; i < len(layerLeaves); i++ {
			l := layerLeaves[i]
			if nodeIndex == l.Index {
				return true
			}
		}

	}
	return false
}

// mergeUnverified gets other partial tree into itself, replacing any conflicting nodes with nodes from
// `other` in the process. Doesn't rehash the nodes, so the integrity of the result is
// not verified. It gives an advantage in speed, but should be used only if the integrity of
// the tree can't be broken, for example, it is used in the `.commit` method of the
// `MerkleTree`, since both partial trees are essentially constructed in place and there's
// no need to verify integrity of the result.
func (pt *PartialTree) mergeUnverified(other PartialTree) {
	depthDifference := len(other.layers) - len(pt.layers)
	var combinedTreeSize uint64
	if depthDifference > 0 {
		combinedTreeSize = uint64(len(other.layers))
	} else {
		combinedTreeSize = uint64(len(pt.layers))
	}

	for layerIndex := uint64(0); layerIndex < combinedTreeSize; layerIndex++ {
		var combinedLayer, filteredLayer Leaves

		selfLayer, ok := layerAtIndex(pt.layers, layerIndex)
		if ok {
			for i := 0; i < len(selfLayer); i++ {
				node := selfLayer[i]
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

		sortLeavesAscending(otherLayer)

		pt.upsertLayer(layerIndex, combinedLayer)

	}
}

// upsertLayer replaces layer at a given index with a new layer. Used during tree merge
func (pt *PartialTree) upsertLayer(layerIndex uint64, newLayer Leaves) {
	_, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {
		pt.layers[layerIndex] = newLayer
	} else {
		pt.layers = append(pt.layers, newLayer)
	}

}

// layerNodesHashes returns all hashes of all layers
func (pt *PartialTree) layerNodesHashes() [][][]byte {
	layers := pt.getLayers()
	layersCount := len(layers)
	allHashes := make([][][]byte, layersCount)
	for i := 0; i < layersCount; i++ {
		l := layers[i]
		leavesCount := len(l)
		layerHashes := make([][]byte, leavesCount)
		for j := 0; j < leavesCount; j++ {
			layerHashes[j] = l[j].Hash
		}
		allHashes[i] = layerHashes
	}
	return allHashes
}

// getLayers returns partial tree layers
func (pt *PartialTree) getLayers() Layers {
	return pt.layers
}

// reverseLayers reverses a slice of types.Leaf slice
func reverseLayers(s Layers) Layers {
	a := make(Layers, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

// popLayer pops last element in the layers
func popLayer(slice Layers) (Leaves, Layers) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

func extractIndicesAndHashes(leaves Leaves) ([]uint64, [][]byte) {
	leavesLen := len(leaves)
	indices := make([]uint64, leavesLen)
	hashes := make([][]byte, leavesLen)
	for i := 0; i < leavesLen; i++ {
		l := leaves[i]
		indices[i], hashes[i] = l.Index, l.Hash
	}
	return indices, hashes
}
