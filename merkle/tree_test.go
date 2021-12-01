package merkle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ComposableFi/merkle-go/merkle"
)

type MergeUint64 struct{}

func (m MergeUint64) Merge(left, right uint64) uint64 {
	// lbs := make([]byte, 8)
	// binary.BigEndian.PutUint64(lbs, left)
	// rbs := make([]byte, 8)
	// binary.BigEndian.PutUint64(rbs, right)
	// h := hash.HashSha1(hash.Node{
	// 	Left:  lbs,
	// 	Right: rbs,
	// })
	// return binary.LittleEndian.Uint64(h[:])
	// h := hash.HashNodes(hash.Node{
	// 	Left:  left,
	// 	Right: right,
	// })

	// return h
	y := float64(right) - float64(left)
	return uint64(y)
}

func TestBuildEmpty(t *testing.T) {
	var leaves []uint64

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, uint64(0), uint64(len(tree.Nodes)))
	require.Equal(t, uint64(0), tree.GetRoot())
}

func TestBuildOne(t *testing.T) {
	var leaves = []uint64{1}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []uint64{1}, tree.Nodes)
}

func TestBuildTwo(t *testing.T) {
	var leaves = []uint64{1, 2}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []uint64{1, 1, 2}, tree.Nodes)
}

func TestBuildFive(t *testing.T) {
	var leaves = []uint64{3, 5, 7, 11, 13}

	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	require.Equal(t, []uint64{1, 1, 2, 2, 3, 5, 7, 11, 13}, tree.Nodes)
}

func TestBuildRootDirectly(t *testing.T) {
	var leaves = []uint64{3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	root := mtree.BuildRoot(leaves)
	require.Equal(t, uint64(1), root)
}

func TestBuiltRootIsSameAsTreeRoot(t *testing.T) {
	var leaves []uint64
	var start uint64
	var end uint64 = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i)
	}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	tree := mtree.BuildTree(leaves)
	root := mtree.BuildRoot(leaves)
	require.Equal(t, tree.GetRoot(), root)
}

func TestVerifyRetrieveLeaves(t *testing.T) {
	var leaves = []uint64{2, 3, 5, 7, 11, 13}
	mtree := merkle.Tree{
		Nodes: leaves,
		Merge: MergeUint64{},
	}
	proof, err := mtree.BuildTreeAndProof(leaves, []uint64{0, 3})
	require.NoError(t, err)
	retrievedLeaves, err := proof.RetriveLeaves(leaves)
	require.NoError(t, err)

	require.Equal(t, []uint64{2, 7}, retrievedLeaves)

	retreivedRoot, err := proof.GetRoot(retrievedLeaves)
	require.NoError(t, err)
	mroot := mtree.BuildRoot(leaves)
	require.Equal(t, mroot, retreivedRoot)

	proof.Indices = []uint64{}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.Error(t, err)
	require.Equal(t, []uint64{}, retrievedLeaves)

	proof.Indices = []uint64{4}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)

	proof.Indices = []uint64{11}
	retrievedLeaves, err = proof.RetriveLeaves(leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)
}
