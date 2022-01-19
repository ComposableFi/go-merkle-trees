package merkle_test

import (
	"crypto/sha256"

	"github.com/ComposableFi/go-merkle-trees/merkle"
	"github.com/ethereum/go-ethereum/crypto"
)

type Sha256Hasher struct{}

func (hr Sha256Hasher) Hash(b []byte) (merkle.Hash, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

type Keccak256Hasher struct{}

func (hr Keccak256Hasher) Hash(b []byte) (merkle.Hash, error) {
	h := crypto.Keccak256Hash(b)
	return h.Bytes(), nil
}

type TestData struct {
	leafValues      []string
	expectedRootHex string
	leafHashes      []merkle.Hash
}
type ProofTestCases struct {
	merkleTree merkle.Tree
	cases      []MerkleProofTestCase
}
type MerkleProofTestCase struct {
	LeafIndicesToProve []uint32
	Leaves             []merkle.Leaf
}

func setupTestData() TestData {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	expectedRootHex := "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2"
	var leafHashes []merkle.Hash
	for i := 0; i < len(leafValues); i++ {
		h, _ := Sha256Hasher{}.Hash([]byte(leafValues[i]))
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
		var leavesHashes []merkle.Hash
		var Leaves []merkle.Leaf
		for j := 0; j < i+1; j++ {
			h, _ := Sha256Hasher{}.Hash([]byte(maxCase[j]))
			leavesHashes = append(leavesHashes, h)
			Leaves = append(Leaves, merkle.Leaf{Index: uint32(j), Hash: h})
		}
		possibleProofElementCombinations := combinations(Leaves)

		var cases []MerkleProofTestCase
		for _, proofElements := range possibleProofElementCombinations {
			var indices []uint32
			// var leaves2 []merkle.Hash
			for _, proofElement := range proofElements {
				indices = append(indices, proofElement.Index)
				// leaves2 = append(leaves2, proofElement.Hash)

			}
			cases = append(cases, MerkleProofTestCase{LeafIndicesToProve: indices, Leaves: proofElements})
		}
		merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(leavesHashes)
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

func combinations(leaves []merkle.Leaf) [][]merkle.Leaf {
	return combine([]merkle.Leaf{}, leaves, [][]merkle.Leaf{})
}

func combine(active []merkle.Leaf, rest []merkle.Leaf, combinations [][]merkle.Leaf) [][]merkle.Leaf {
	if len(rest) == 0 {
		if len(active) == 0 {
			return combinations
		} else {
			combinations = append(combinations, active)
			return combinations
		}
	} else {
		next := make([]merkle.Leaf, len(active))
		copy(next, active)

		if len(rest) > 0 {
			next = append(next, rest[0])
		}
		combinations := combine(next, rest[1:], combinations)
		combinations = combine(active, rest[1:], combinations)
		return combinations
	}
}
