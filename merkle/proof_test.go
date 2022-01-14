package merkle_test

import (
	"testing"

	"github.com/ComposableFi/merkle-go/merkle"
	"github.com/stretchr/testify/require"
)

func TestCorrectProofs(t *testing.T) {
	testData := setup()
	expectedRoot := testData.expectedRootHex
	leafHashes := testData.leafHashes
	indicesToProve := []uint32{3, 4}
	var leavesToProve []merkle.Hash
	for _, i := range indicesToProve {
		leavesToProve = append(leavesToProve, leafHashes[i])
	}

	merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)

	proof := merkleTree.Proof(indicesToProve)
	leafTuples := merkle.MapIndiceAndLeaves(indicesToProve, leavesToProve)

	extractedRoot, err := proof.GetRootHex(leafTuples, len(testData.leafValues))
	require.NoError(t, err)

	require.Equal(t, expectedRoot, extractedRoot)

	testCases, err := setupProofTestCases()
	require.NoError(t, err)
	for k, testCase := range testCases {
		t.Logf("Proof Case: %v", k)
		merkleTree := testCase.merkleTree
		root := merkleTree.GetRoot()
		for k2, c := range testCase.cases {
			t.Logf("Test Case: %v", k2)
			t.Logf("Indices: %v", c.LeafIndicesToProve)
			t.Logf("leafTuples: %v", c.LeafTuples)
			proof := merkleTree.Proof(c.LeafIndicesToProve)
			extractedRoot, err := proof.GetRoot(c.LeafTuples, merkleTree.GetLeavesLen())
			require.NoError(t, err)
			require.Equal(t, root, extractedRoot)
		}

	}
}
