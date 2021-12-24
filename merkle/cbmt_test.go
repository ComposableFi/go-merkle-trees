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

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, int(0), int(len(tree.Nodes)))
	require.Equal(t, int(0), tree.GetRoot())
}

func TestBuildOne(t *testing.T) {
	var leaves = []interface{}{1}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, []interface{}{1}, tree.Nodes)
}

func TestBuildTwo(t *testing.T) {
	var leaves = []interface{}{1, 2}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, []interface{}{1, 1, 2}, tree.Nodes)
}

func TestBuildFive(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, []interface{}{1, 1, 2, 2, 3, 5, 7, 11, 13}, tree.Nodes)
}

func TestBuildRootDirectly(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, int(1), root)
}

func TestBuiltRootIsSameAsTreeRoot(t *testing.T) {
	var leaves []interface{}
	var start int
	var end = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i)
	}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, tree.GetRoot(), root)
}

func TestVerifyRetrieveLeaves(t *testing.T) {
	var leaves = []interface{}{2, 3, 5, 7, 11, 13}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	proof, err := cbmt.BuildMerkleProof(leaves, []uint32{0, 3})
	require.NoError(t, err)
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)

	require.Equal(t, []interface{}{2, 7}, retrievedLeaves)

	retreivedRoot, err := proof.CalculateRootHash()
	require.NoError(t, err)
	mroot := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, mroot, retreivedRoot)

	proof.Leaves = []merkle.LeafData{}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.Error(t, err)
	require.Equal(t, []interface{}{}, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{merkle.LeafData{Index: 0, Leaf: 4}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{merkle.LeafData{Index: 0, Leaf: 11}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)
}

func TestRebuildProof(t *testing.T) {
	var leaves = []interface{}{3, 5, 7, 11, 13}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := tree.GetRoot()

	//build proof
	proof, err := tree.BuildProof([]uint32{0, 3})
	require.NoError(t, err)
	lemmas := proof.Lemmas
	leafDataList := proof.Leaves

	// rebuild proof
	var neededLeaves []interface{}

	for _, v := range leafDataList {
		neededLeaves = append(neededLeaves, tree.Nodes[v.Index])
	}

	rebuildProof := merkle.Proof{
		Leaves: leafDataList,
		Lemmas: lemmas,
		Merge:  MergeInt32{},
	}

	isValid, err := rebuildProof.VerifyRootHash(root)
	require.NoError(t, err)
	require.Equal(t, true, isValid)

	rebuiltRoot, err := rebuildProof.CalculateRootHash()
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
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}

	//build proof
	proof, err := cbmt.BuildMerkleProof(leaves, leafIndecies)
	require.NoError(t, err)
	require.Equal(t, []interface{}{13, 5, 4}, proof.Lemmas)
	root, err := proof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, int(2), root)

	leaves = []interface{}{2}
	proof, err = cbmt.BuildMerkleProof(leaves, []uint32{0})
	require.NoError(t, err)
	require.Equal(t, 0, len(proof.Lemmas))
	root, err = proof.CalculateRootHash()
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
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	proof, err := cbmt.BuildMerkleProof(leaves, leafIndices)
	require.NoError(t, err)

	proofRoot, err := proof.CalculateRootHash()
	require.NoError(t, err)

	treeRoot := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, proofRoot, treeRoot)
}
