package merkle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ComposableFi/merkle-go/merkle"
)

type MergeInt32 struct{}

func (m MergeInt32) Merge(left, right interface{}) interface{} {
	var r, l int
	r = right.(int)
	l = left.(int)
	merged := r - l
	return merged
}

func TestBuildEmpty(t *testing.T) {
	var leaves []interface{}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, int(0), int(len(tree.Nodes)))
	require.Equal(t, int(0), tree.GetRoot())
}

func TestBuildOne(t *testing.T) {
	var leaves = []interface{}{1}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []interface{}{1}, tree.Nodes)
}

func TestBuildTwo(t *testing.T) {
	var leaves = []interface{}{1, 2}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []interface{}{1, 1, 2}, tree.Nodes)
}

func TestBuildFive(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []interface{}{1, 1, 2, 2, 3, 5, 7, 11, 13}, tree.Nodes)
}

func TestBuildRootDirectly(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	root := mtree.BuildRoot(leaves)
	require.Equal(t, int(1), root)
}

func TestBuiltRootIsSameAsTreeRoot(t *testing.T) {
	var leaves []interface{}
	var start int
	var end = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i)
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	tree := mtree.BuildTree(leaves)
	root := mtree.BuildRoot(leaves)
	require.Equal(t, tree.GetRoot(), root)
}

func TestVerifyRetrieveLeaves(t *testing.T) {
	var leaves = []interface{}{2, 3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeInt32{},
	}
	proof, err := mtree.BuildTreeAndProof(leaves, []uint32{0, 3})
	require.NoError(t, err)
	retrievedLeaves, err := proof.RetriveLeaves(leaves)
	require.NoError(t, err)

	require.Equal(t, []interface{}{2, 7}, retrievedLeaves)

	retreivedRoot, err := proof.GetRoot(retrievedLeaves)
	require.NoError(t, err)
	mroot := mtree.BuildRoot(leaves)
	require.Equal(t, mroot, retreivedRoot)

	proof.Indices = []uint32{}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.Error(t, err)
	require.Equal(t, []interface{}{}, retrievedLeaves)

	proof.Indices = []uint32{4}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)

	proof.Indices = []uint32{11}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)
}
