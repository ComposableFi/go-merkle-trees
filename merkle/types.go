package merkle

type Hash []byte

// Tree is a Merkle Tree that is well suited for both basic and advanced usage.
//
// Basic features include the creation and verification of Merkle proofs from a set of leaves.
// This is often done in various cryptocurrencies.
//
// Advanced features include being able to make transactional changes to a tree with being able to
// roll back to any previously committed state of the tree. This scenario is similar to Git and
// can be found in databases and file systems.
type Tree struct {
	currentWorkingTree PartialTree
	history            []PartialTree
	UncommittedLeaves  []Hash
	hasher             Hasher
}

// NewTree creates a new instance of merkle tree. requires a hash algorithm to be specified.
func NewTree(hasher Hasher) Tree {
	return Tree{
		currentWorkingTree: NewPartialTree(hasher),
		history:            []PartialTree{},
		UncommittedLeaves:  []Hash{},
		hasher:             hasher,
	}
}

// PartialTree represents a part of the original tree that is enough to calculate the root.
// Used in to extract the root in a merkle proof, to apply diff to a tree or to merge
// multiple trees into one.
// It is a rare case when you need to use this struct on it's own. It's mostly used inside
type PartialTree struct {
	layers [][]Leaf
	hasher Hasher
}

// NewPartialTree Takes hasher as an argument and build a Merkle Tree from them.
// Since it's a partial tree, hashes must be accompanied by their index in the original tree.
func NewPartialTree(hasher Hasher) PartialTree {
	return PartialTree{
		layers: [][]Leaf{},
		hasher: hasher,
	}
}

// Proof is used to parse, verify, calculate a root for Merkle proofs.
// Proof requires specifying hashing algorithm and hash size in order to work.
// The hashing algorithm is set through the Hasher interface, which is supplied as a generic
// parameter to the Proof.
type Proof struct {
	proofHashes      []Hash
	leaves           []Leaf
	totalLeavesCount uint32
	hasher           Hasher
}

// NewProof create new instance of merkle proof
func NewProof(leaves []Leaf, proofHashes []Hash, totalLeavesCount uint32, hasher Hasher) Proof {
	return Proof{
		leaves:           leaves,
		proofHashes:      proofHashes,
		totalLeavesCount: totalLeavesCount,
		hasher:           hasher,
	}
}

// Hasher is an interface used to provide a hashing algorithm for the library.
type Hasher interface {
	Hash(data []byte) (Hash, error)
}

// Leaf is a represention of leaf index and its hash
type Leaf struct {
	Index uint32
	Hash  Hash
}
