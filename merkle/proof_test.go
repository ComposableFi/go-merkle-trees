package merkle_test

import (
	"testing"

	"github.com/ComposableFi/merkle-go/merkle"
	"github.com/stretchr/testify/require"
)

func TestRebuildProof(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	root := tree.GetRoot()

	//build proof
	proof, err := tree.BuildProof([]uint32{0, 3})
	require.NoError(t, err)
	lemmas := proof.Lemmas
	indices := proof.Indices

	// rebuild proof
	var neededLeaves []interface{}

	for _, v := range indices {
		neededLeaves = append(neededLeaves, tree.Nodes[v])
	}

	rebuildProof := merkle.Proof{
		Indices: indices,
		Lemmas:  lemmas,
		Merge:   MergeInt32{},
	}

	isValid, err := rebuildProof.Verify(root, neededLeaves)
	require.NoError(t, err)
	require.Equal(t, true, isValid)

	rebuiltRoot, err := rebuildProof.GetRoot(neededLeaves)
	require.NoError(t, err)
	require.Equal(t, root, rebuiltRoot)
}

func TestBuildProof(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13, 17}
	leafIndecies := []uint32{0, 5}
	var proofLeaves []interface{}
	for _, idx := range leafIndecies {
		proofLeaves = append(proofLeaves, leaves[idx])
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}

	//build proof
	proof, err := mtree.BuildTreeAndProof(leaves, leafIndecies)
	require.NoError(t, err)
	require.Equal(t, []interface{}{13, 5, 4}, proof.Lemmas)
	root, err := proof.GetRoot(proofLeaves)
	require.NoError(t, err)
	require.Equal(t, int(2), root)

	leaves = []interface{}{2}
	mtree.Nodes = leaves
	proof, err = mtree.BuildTreeAndProof(leaves, []uint32{0})
	require.NoError(t, err)
	require.Equal(t, 0, len(proof.Lemmas))
	root, err = proof.GetRoot(leaves)
	require.NoError(t, err)
	require.Equal(t, int(2), root)
}

func TestTreeRootIsTheSameAsProofRoot(t *testing.T) {
	var leaves []interface{}
	var leafIndices []uint32
	var start uint32 = 2
	var end uint32 = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, int(i))
		leafIndices = append(leafIndices, i-start)
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	proof, err := mtree.BuildTreeAndProof(leaves, leafIndices)
	require.NoError(t, err)

	proofRoot, err := proof.GetRoot(leaves)
	require.NoError(t, err)

	treeRoot := mtree.BuildRoot(leaves)
	require.Equal(t, proofRoot, treeRoot)
}
