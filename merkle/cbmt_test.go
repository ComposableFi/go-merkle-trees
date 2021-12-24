package merkle_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ComposableFi/merkle-go/merkle"
)

type MergeInt32 struct{}

func (m MergeInt32) Merge(left, right []byte) []byte {
	var r, l int32
	r = b2i(right)
	l = b2i(left)
	merged := r - l
	return i2b(merged)
}

func b2i(b []byte) int32 {
	i := int32(binary.LittleEndian.Uint64(b))
	return i
}

func i2b(i int32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func TestBuildEmpty(t *testing.T) {
	var leaves [][]byte

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, int(0), int(len(tree.Nodes)))
	require.Equal(t, []byte{0}, tree.GetRoot())
}

func TestBuildOne(t *testing.T) {
	var leaves = [][]byte{i2b(1)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1)}, tree.Nodes)
}

func TestBuildTwo(t *testing.T) {
	var leaves = [][]byte{i2b(1), i2b(2)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1), i2b(1), i2b(2)}, tree.Nodes)
}

func TestBuildFive(t *testing.T) {
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1), i2b(1), i2b(2), i2b(2), i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}, tree.Nodes)
}

func TestBuildRootDirectly(t *testing.T) {
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, i2b(1), root)
}

func TestBuiltRootIsSameAsTreeRoot(t *testing.T) {
	var leaves [][]byte
	var start int
	var end = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i2b(int32(i)))
	}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, tree.GetRoot(), root)
}

func TestVerifyRetrieveLeaves(t *testing.T) {
	var leaves = [][]byte{i2b(2), i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	proof, err := cbmt.BuildMerkleProof(leaves, []uint32{0, 3})
	require.NoError(t, err)
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)

	require.Equal(t, [][]byte{i2b(2), i2b(7)}, retrievedLeaves)

	retreivedRoot, err := proof.CalculateRootHash()
	require.NoError(t, err)
	mroot := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, mroot, retreivedRoot)

	proof.Leaves = []merkle.LeafData{}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.Error(t, err)
	require.Equal(t, [][]byte{}, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{{Index: 0, Leaf: i2b(4)}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{{Index: 0, Leaf: i2b(11)}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)
}

func TestRebuildProof(t *testing.T) {
	var leaves = [][]byte{i2b(2), i2b(3), i2b(5), i2b(7), i2b(11)}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := tree.GetRoot()

	//build proof
	proof, err := tree.BuildProof([]uint32{0, 3})
	require.NoError(t, err)
	lemmas := proof.Proofs
	leafDataList := proof.Leaves

	rebuildProof := merkle.Proof{
		Leaves: leafDataList,
		Proofs: lemmas,
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
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13), i2b(17)}
	leafIndecies := []uint32{0, 5}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}

	//build proof
	proof, err := cbmt.BuildMerkleProof(leaves, leafIndecies)
	require.NoError(t, err)
	require.Equal(t, [][]byte{i2b(13), i2b(5), i2b(4)}, proof.Proofs)
	root, err := proof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, i2b(2), root)

	leaves = [][]byte{i2b(2)}
	proof, err = cbmt.BuildMerkleProof(leaves, []uint32{0})
	require.NoError(t, err)
	require.Equal(t, 0, len(proof.Proofs))
	root, err = proof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, i2b(2), root)
}

func TestTreeRootIsTheSameAsProofRoot(t *testing.T) {
	var leaves [][]byte
	var leafIndices []uint32
	var start uint32 = 2
	var end uint32 = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i2b(int32(i)))
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
