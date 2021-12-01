package merkle_test

import (
	"testing"

	"github.com/ComposableFi/merkle-go/merkle"
	"github.com/stretchr/testify/require"
)

func TestRebuildProof(t *testing.T) {
	var leaves = []uint64{3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	root := tree.GetRoot()

	//build proof
	proof, err := tree.BuildProof([]uint64{0, 3})
	require.NoError(t, err)
	lemmas := proof.Lemmas
	indices := proof.Indices

	// rebuild proof
	var neededLeaves []uint64

	for _, v := range indices {
		neededLeaves = append(neededLeaves, tree.Nodes[v])
	}

	rebuildProof := merkle.Proof{
		Indices: indices,
		Lemmas:  lemmas,
		Merge:   MergeUint64{},
	}

	isValid, err := rebuildProof.Verify(root, neededLeaves)
	require.NoError(t, err)
	require.Equal(t, true, isValid)

	rebuiltRoot, err := rebuildProof.GetRoot(neededLeaves)
	require.NoError(t, err)
	require.Equal(t, root, rebuiltRoot)
}

func TestBuildProof(t *testing.T) {
	var leaves = []uint64{3, 5, 7, 11, 13, 17}
	leafIndecies := []uint64{0, 5}
	var proofLeaves []uint64
	for _, idx := range leafIndecies {
		proofLeaves = append(proofLeaves, leaves[idx])
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}

	//build proof
	proof, err := mtree.BuildTreeAndProof(leaves, leafIndecies)
	require.NoError(t, err)
	require.Equal(t, []uint64{13, 5, 4}, proof.Lemmas)
	root, err := proof.GetRoot(proofLeaves)
	require.NoError(t, err)
	require.Equal(t, uint64(2), root)

	leaves = []uint64{2}
	mtree.Nodes = leaves
	proof, err = mtree.BuildTreeAndProof(leaves, []uint64{0})
	require.NoError(t, err)
	require.Equal(t, 0, len(proof.Lemmas))
	root, err = proof.GetRoot(leaves)
	require.NoError(t, err)
	require.Equal(t, uint64(2), root)
}

func TestTreeRootIsTheSameAsProofRoot(t *testing.T) {
	var leaves, leafIndices []uint64
	var start uint64 = 2
	var end uint64 = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i)
		leafIndices = append(leafIndices, i-start)
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	proof, err := mtree.BuildTreeAndProof(leaves, leafIndices)
	require.NoError(t, err)

	proofRoot, err := proof.GetRoot(leaves)
	require.NoError(t, err)

	treeRoot := mtree.BuildRoot(leaves)
	require.Equal(t, proofRoot, treeRoot)
}
