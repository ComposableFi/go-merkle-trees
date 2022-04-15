package merkle

import (
	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/ComposableFi/go-merkle-trees/types"
)

type TestData struct {
	leafValues      []string
	expectedRootHex string
	leafHashes      [][]byte
}
type ProofTestCases struct {
	merkleTree Tree
	cases      []MerkleProofTestCase
}
type MerkleProofTestCase struct {
	LeafIndicesToProve []uint64
	Leaves             []types.Leaf
}

func setupTestData() TestData {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	expectedRootHex := "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2"
	var leafHashes [][]byte
	for i := 0; i < len(leafValues); i++ {
		h, _ := hasher.Sha256Hasher{}.Hash([]byte(leafValues[i]))
		leafHashes = append(leafHashes, h)
	}
	return TestData{
		leafValues:      leafValues,
		leafHashes:      leafHashes,
		expectedRootHex: expectedRootHex,
	}
}

func setupProofTestCases() ([]ProofTestCases, error) {
	maxCase := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "k", "l", "m", "o", "p", "r", "s",
	}
	var merkleProofCases []ProofTestCases
	for i := 0; i < len(maxCase); i++ {
		var leavesHashes [][]byte
		var Leaves []types.Leaf
		for j := 0; j < i+1; j++ {
			h, _ := hasher.Sha256Hasher{}.Hash([]byte(maxCase[j]))
			leavesHashes = append(leavesHashes, h)
			Leaves = append(Leaves, types.Leaf{Index: uint64(j), Hash: h})
		}
		possibleProofElementCombinations := combinations(Leaves)

		var cases []MerkleProofTestCase
		for _, proofElements := range possibleProofElementCombinations {
			var indices []uint64
			for _, proofElement := range proofElements {
				indices = append(indices, proofElement.Index)
				// leaves2 = append(leaves2, proofElement.Hash)

			}
			cases = append(cases, MerkleProofTestCase{LeafIndicesToProve: indices, Leaves: proofElements})
		}
		merkleTree, err := NewTree(hasher.Sha256Hasher{}).FromLeaves(leavesHashes)
		if err != nil {
			return []ProofTestCases{}, err
		}

		c := ProofTestCases{
			merkleTree: merkleTree,
			cases:      cases,
		}
		merkleProofCases = append(merkleProofCases, c)
	}
	return merkleProofCases, nil
}

func combinations(leaves []types.Leaf) [][]types.Leaf {
	return combine([]types.Leaf{}, leaves, [][]types.Leaf{})
}

func combine(active []types.Leaf, rest []types.Leaf, combinations [][]types.Leaf) [][]types.Leaf {
	if len(rest) == 0 {
		if len(active) == 0 {
			return combinations
		}
		combinations = append(combinations, active)
		return combinations
	}
	next := make([]types.Leaf, len(active))
	copy(next, active)

	if len(rest) > 0 {
		next = append(next, rest[0])
	}
	combinations = combine(next, rest[1:], combinations)
	combinations = combine(active, rest[1:], combinations)
	return combinations

}
