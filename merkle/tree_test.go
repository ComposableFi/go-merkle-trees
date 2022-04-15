package merkle

import (
	"testing"

	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree(t *testing.T) {
	NewTree(hasher.Sha256Hasher{})
}

func TestFromLeaves(t *testing.T) {
	leaves, err := sampleHashes()
	require.NoError(t, err)

	mtree := NewTree(hasher.Sha256Hasher{})
	mtree, err = mtree.FromLeaves(leaves)
	require.NoError(t, err)
	require.Equal(t, [][]byte{}, mtree.UncommittedLeaves)
}

func TestRoot(t *testing.T) {
	leaves, err := sampleHashes()
	require.NoError(t, err)

	mtree := NewTree(hasher.Sha256Hasher{})
	mtree, err = mtree.FromLeaves(leaves)
	require.NoError(t, err)

	indicesToProve := []uint64{0, 1}
	proof := mtree.Proof(indicesToProve)
	root := mtree.Root()

	verified, err := proof.Verify(root)
	require.NoError(t, err)
	require.True(t, verified)
}

func TestProof(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e", "f"}
	var leaves [][]byte
	for i := 0; i < len(values); i++ {
		h, _ := hasher.Sha256Hasher{}.Hash([]byte(values[i]))
		leaves = append(leaves, h)
	}
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree, err := mtree.FromLeaves(leaves)
	require.NoError(t, err)

	indicesToProve := []uint64{3, 4}
	proof := mtree.Proof(indicesToProve)
	root := mtree.Root()

	verified, err := proof.Verify(root)
	require.NoError(t, err)
	require.True(t, verified)
}

func TestCorrectTreeRoot(t *testing.T) {
	testData := setupTestData()

	merkleTree, err := NewTree(hasher.Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	rootHex := merkleTree.RootHex()
	require.Equal(t, testData.expectedRootHex, rootHex)
}

func TestCorrectTreeDepth(t *testing.T) {
	testData := setupTestData()

	merkleTree, err := NewTree(hasher.Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	depth := merkleTree.depth()
	require.Equal(t, 3, depth)
}

func TestCorrectProofRoot(t *testing.T) {
	testData := setupTestData()
	indicesToProve := []uint64{3, 4}
	expectedProofHashes := []string{
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
		"252f10c83610ebca1a059c0bae8255eba2f95be4d1d7bcfa89d7248a82d9f111",
		"e5a01fee14e0ed5c48714f22180f25ad8365b53f9779f79dc4a3d7e93963f94a",
	}
	merkleTree, err := NewTree(hasher.Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	proof := merkleTree.Proof(indicesToProve)
	require.Equal(t, expectedProofHashes, proof.ProofHashesHex())
}

func TestGetCorrectRootAfterCommit(t *testing.T) {
	testData := setupTestData()
	expectedRoot := testData.expectedRootHex
	leafHashes := testData.leafHashes

	merkleTree, err := NewTree(hasher.Sha256Hasher{}).FromLeaves([][]byte{})
	require.NoError(t, err)
	merkleTree2, err := NewTree(hasher.Sha256Hasher{}).FromLeaves(leafHashes)
	require.NoError(t, err)

	merkleTree.append(leafHashes)

	root, err := merkleTree.uncommittedRootHex()
	require.NoError(t, err)

	require.Equal(t, expectedRoot, merkleTree2.RootHex())
	require.Equal(t, expectedRoot, root)

	expectedRoot = "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034"
	leaf, err := hasher.Sha256Hasher{}.Hash([]byte("g"))
	require.NoError(t, err)
	merkleTree.insert(leaf)

	uncommittedRoot, err := merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, expectedRoot, uncommittedRoot)

	require.Equal(t, []byte{}, merkleTree.Root())

	merkleTree.commit()

	hashOfH, _ := hasher.Sha256Hasher{}.Hash([]byte("h"))
	hashOfK, _ := hasher.Sha256Hasher{}.Hash([]byte("k"))
	newLeaves := [][]byte{hashOfH, hashOfK}

	merkleTree.append(newLeaves)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())
	uncommittedRootHex, err := merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", uncommittedRootHex)

	merkleTree.commit()

	leaves := merkleTree.baseLeaves()
	reconstructedTree, err := merkleTree.FromLeaves(leaves)
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", reconstructedTree.RootHex())

}

func TestChangeTheResultWenCalledTwice(t *testing.T) {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	var leaves [][]byte
	for _, v := range leafValues {
		h, _ := hasher.Sha256Hasher{}.Hash([]byte(v))
		leaves = append(leaves, h)
	}

	merkleTree := NewTree(hasher.Sha256Hasher{})

	// Appending leaves to the tree without committing
	merkleTree.append(leaves)

	require.Equal(t, []byte{}, merkleTree.Root())

	uncommittedRootHex, err := merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", uncommittedRootHex)

	merkleTree.commit()

	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.RootHex())

	uncommittedRootHex, err = merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "", uncommittedRootHex)

	gHash, _ := hasher.Sha256Hasher{}.Hash([]byte("g"))
	merkleTree.insert(gHash)

	uncommittedRootHex, err = merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", uncommittedRootHex)

	merkleTree.commit()

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())

	hashOfH, _ := hasher.Sha256Hasher{}.Hash([]byte("h"))
	hashOfK, _ := hasher.Sha256Hasher{}.Hash([]byte("k"))
	merkleTree.append([][]byte{hashOfH, hashOfK})

	merkleTree.commit()
	merkleTree.commit()

	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", merkleTree.RootHex())

	merkleTree.rollback()
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())

	merkleTree.rollback()
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.RootHex())
}

func TestRollbackPreviousCommit(t *testing.T) {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	var leaves [][]byte
	for _, v := range leafValues {
		h, _ := hasher.Sha256Hasher{}.Hash([]byte(v))
		leaves = append(leaves, h)
	}

	merkleTree := NewTree(hasher.Sha256Hasher{})
	merkleTree.append(leaves)

	require.Equal(t, []byte{}, merkleTree.Root())

	merkleTree.commit()

	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.RootHex())

	gHash, _ := hasher.Sha256Hasher{}.Hash([]byte("g"))
	merkleTree.insert(gHash)

	uncommittedRootHex, err := merkleTree.uncommittedRootHex()
	require.NoError(t, err)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", uncommittedRootHex)

	merkleTree.commit()

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())

	hashOfH, _ := hasher.Sha256Hasher{}.Hash([]byte("h"))
	hashOfK, _ := hasher.Sha256Hasher{}.Hash([]byte("k"))
	merkleTree.append([][]byte{hashOfH, hashOfK})

	uncommittedRootHex, err = merkleTree.uncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", uncommittedRootHex)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())

	merkleTree.commit()

	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", merkleTree.RootHex())

	merkleTree.rollback()
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.RootHex())

	merkleTree.rollback()
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.RootHex())
}

func sampleHashes() ([][]byte, error) {
	aHash, err := hasher.Sha256Hasher{}.Hash([]byte("a"))
	if err != nil {
		return nil, err
	}
	bHash, err := hasher.Sha256Hasher{}.Hash([]byte("b"))
	if err != nil {
		return nil, err
	}
	cHash, err := hasher.Sha256Hasher{}.Hash([]byte("c"))
	if err != nil {
		return nil, err
	}
	return [][]byte{aHash, bHash, cHash}, nil
}

// Benchmarking
func BenchmarkSha256Hash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hasher.Sha256Hasher{}.Hash([]byte("a"))
	}
}

func BenchmarkFromLeaves(b *testing.B) {
	leaves, _ := sampleHashes()
	for n := 0; n < b.N; n++ {
		mtree := NewTree(hasher.Sha256Hasher{})
		mtree.FromLeaves(leaves)
	}
}

func BenchmarkShadowIndices(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	for n := 0; n < b.N; n++ {
		mtree.getShadowIndecies()
	}
}

func BenchmarkShadowLeaves(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	shadowIndices := mtree.getShadowIndecies()
	for n := 0; n < b.N; n++ {
		mtree.getShadowLeaves(shadowIndices)
	}
}

func BenchmarkSiblingIndices(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	shadowIndices := mtree.getShadowIndecies()
	for n := 0; n < b.N; n++ {
		siblingIndecies(shadowIndices)
	}
}

func BenchmarkParentIndices(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	shadowIndices := mtree.getShadowIndecies()
	for n := 0; n < b.N; n++ {
		parentIndecies(shadowIndices)
	}
}

func BenchmarkHelperNodeLeaves(b *testing.B) {
	leaves, _ := sampleHashes()
	mtree := NewTree(hasher.Sha256Hasher{})
	mtree.append(leaves)
	shadowIndices := mtree.getShadowIndecies()
	for n := 0; n < b.N; n++ {
		mtree.helperNodeLeaves(shadowIndices)
	}
}
