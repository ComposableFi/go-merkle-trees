package merkle

import (
	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/ComposableFi/go-merkle-trees/types"
)

// build is a wrapper for buildTree
func (pt *PartialTree) build(partialLayers Layers, depth uint64) (PartialTree, error) {

	// build partial tree layers
	layers, err := pt.buildTree(partialLayers, depth)
	if err != nil {
		return PartialTree{}, err
	}

	return PartialTree{layers: layers}, nil
}

// buildTree is a general algorithm for building a partial tree. It can be used to extract root
// from merkle proof, or if a complete set of leaves provided as a first argument and no
// helper indices given, will construct the whole tree.
// the layers need to be reversed because we are going to process the tree from the bottom and merge left and right nodes to get parent
func (pt *PartialTree) buildTree(partialLayers Layers, fullTreeDepth uint64) (Layers, error) {

	// reverse the layers to process backward
	reversedLayers := reverseLayers(partialLayers)

	var currentLayer Leaves
	var partialTree Layers

	// loop through all indices of full tree depth
	for i := uint64(0); i < fullTreeDepth; i++ {

		// add nodes to current layer for following process
		if len(reversedLayers) > 0 {
			var nodes Leaves
			nodes, reversedLayers = popLayer(reversedLayers)
			currentLayer = append(currentLayer, nodes...)
		}

		sortLeavesAscending(currentLayer)

		partialTree = append(partialTree, currentLayer)

		// to get siblings we need to have indices and hashes in separate slices
		indices, hashes := extractIndicesAndHashes(currentLayer)

		// freeup for next round
		currentLayer = make(Leaves, 0)

		// get parent indices to set the merged node hash
		parentIndices := parentIndecies(indices)

		// loop through parents and set the merged hash
		for i := 0; i < len(parentIndices); i++ {
			parnetNodeIndex := parentIndices[i]
			leftIndex := getLeftIndex(i)
			if len(hashes) > leftIndex {
				rightIndex := getRightIndex(i)

				// calculate left and right hash
				leftHash := hashes[leftIndex]
				var rightHash []byte
				if len(hashes) > rightIndex {
					rightHash = hashes[rightIndex]
				} else {
					rightHash = nil
				}

				// merge left and right hash and merge them
				hash, err := hasher.MergeAndHash(pt.hasher, leftHash, rightHash)
				if err != nil {
					return Layers{}, err
				}

				// append parent node to the current layer for next round
				currentLayer = append(currentLayer, types.Leaf{
					Index: parnetNodeIndex,
					Hash:  hash,
				})
			} else {
				// it means we have not enough parent indices to match hashes with
				return Layers{}, errNotEnoughParentNodes
			}
		}

	}

	// update and return partial tree after traversing the whole depth of full tree
	if len(currentLayer) > 0 {
		partialTree = append(partialTree, currentLayer)
	}

	return partialTree, nil
}

// Root returns the root of the tree, it is the first item hash of the last layer
func (pt *PartialTree) Root() []byte {

	if len(pt.layers) > 0 {
		// get the last layer
		lastLayer := pt.layers[len(pt.layers)-1]

		// get the first leaf of top layer
		firstItem := lastLayer[0]

		// return the hash of most top node as partial tree root
		return firstItem.Hash
	}

	// no root if no layers is available
	return nil
}

// contains checks if a node index is present in a layer
func (pt *PartialTree) contains(layerIndex, nodeIndex uint64) bool {

	// check if the layer exists
	layerLeaves, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {

		// check all leaves indices
		for i := 0; i < len(layerLeaves); i++ {

			// if leaves of layer have index
			if nodeIndex == layerLeaves[i].Index {
				return true
			}

		}
	}

	// layer or node in the layer not exist
	return false
}

// mergeUnverifiedLayers gets other partial tree into itself, replacing any conflicting nodes with nodes from
// `other` in the process. Doesn't rehash the nodes, so the integrity of the result is
// not verified. It gives an advantage in speed, but should be used only if the integrity of
// the tree can't be broken, for example, it is used in the `.commit` method of the
// `MerkleTree`, since both partial trees are essentially constructed in place and there's
// no need to verify integrity of the result.
func (pt *PartialTree) mergeUnverifiedLayers(other PartialTree) {

	// calculate size of combined layers of new partial tree and current tree
	depthDifference := len(other.layers) - len(pt.layers)
	var combinedTreeSize uint64
	if depthDifference > 0 {
		combinedTreeSize = uint64(len(other.layers))
	} else {
		combinedTreeSize = uint64(len(pt.layers))
	}

	// loop until we reach the combined size
	for layerIndex := uint64(0); layerIndex < combinedTreeSize; layerIndex++ {
		var combinedLayer, filteredLayer Leaves

		// populate existing layer nodes that are missing in the new partial tree
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

		// append new tree to the combined layer
		otherLayer, ok := layerAtIndex(other.layers, layerIndex)
		if ok {
			combinedLayer = append(combinedLayer, otherLayer...)
		}

		// sort combined and make it final
		sortLeavesAscending(otherLayer)

		// update or insert all of processes combined nodes into the layer
		pt.upsertLayer(layerIndex, combinedLayer)

	}

}

// upsertLayer replaces layer at a given index with a new layer. Used during tree merge
func (pt *PartialTree) upsertLayer(layerIndex uint64, newLayer Leaves) {

	// check layer existance
	_, ok := layerAtIndex(pt.layers, layerIndex)
	if ok {
		// layer exists then update
		pt.layers[layerIndex] = newLayer
	} else {
		// layer not exists then insert
		pt.layers = append(pt.layers, newLayer)
	}

}

// layerNodesHashes returns all hashes of all layers
func (pt *PartialTree) layerNodesHashes() [][][]byte {
	layers := pt.getLayers()
	layersCount := len(layers)
	allHashes := make([][][]byte, layersCount)

	// loop through all layers
	for i := 0; i < layersCount; i++ {
		l := layers[i]
		leavesCount := len(l)
		layerHashes := make([][]byte, leavesCount)

		// loop through all of nodes of this layer
		for j := 0; j < leavesCount; j++ {
			layerHashes[j] = l[j].Hash
		}

		// update the result
		allHashes[i] = layerHashes
	}

	return allHashes
}

// getLayers returns partial tree layers
func (pt *PartialTree) getLayers() Layers {
	return pt.layers
}

// reverseLayers reverses a slice of types.Leaf slice
func reverseLayers(layers Layers) Layers {

	// make cls copy to prevent modification of original layers
	cls := make(Layers, len(layers))
	copy(cls, layers)

	for i := len(cls)/2 - 1; i >= 0; i-- {
		opp := len(cls) - 1 - i
		// swap the items
		cls[i], cls[opp] = cls[opp], cls[i]
	}

	return cls
}

// popLayer pops last element in the layers
func popLayer(slice Layers) (Leaves, Layers) {
	popElem, newSlice := slice[len(slice)-1], slice[0:len(slice)-1]
	return popElem, newSlice
}

// extractIndicesAndHashes makes indices and hashes separated into two different slices
func extractIndicesAndHashes(leaves Leaves) ([]uint64, [][]byte) {

	leavesLen := len(leaves)
	indices := make([]uint64, leavesLen)
	hashes := make([][]byte, leavesLen)

	// loop through leaves and add the index and hash to different slices
	for i := 0; i < leavesLen; i++ {
		l := leaves[i]
		indices[i], hashes[i] = l.Index, l.Hash
	}

	return indices, hashes
}
