package merkle

import (
	"testing"

	"github.com/ComposableFi/go-merkle-trees/hasher"
)

func BenchmarkBuildPartialTree(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	partialTreeLayers, uncommittedTreeDepth := mtree.uncommitedPartialTreeLayers()
	ptree := NewPartialTree(mtree.hasher)
	for n := 0; n < b.N; n++ {
		ptree.build(partialTreeLayers, uncommittedTreeDepth)
	}
}

func BenchmarkReverseLayers(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	partialTreeLayers, _ := mtree.uncommitedPartialTreeLayers()
	for n := 0; n < b.N; n++ {
		reverseLayers(partialTreeLayers)
	}
}

func BenchmarkLayerNodesHashes(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	partialTreeLayers, uncommittedTreeDepth := mtree.uncommitedPartialTreeLayers()
	ptree := NewPartialTree(mtree.hasher)
	ptree.build(partialTreeLayers, uncommittedTreeDepth)
	for n := 0; n < b.N; n++ {
		ptree.layerNodesHashes()
	}
}
