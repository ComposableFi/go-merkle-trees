package merkle

import (
	"github.com/ComposableFi/go-merkle-trees/types"
)

// Leaves is a representation of slice of leaf
type Leaves []types.Leaf

// Layers is a representation of slice of Leaves slice
type Layers []Leaves

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
	UncommittedLeaves  [][]byte
	hasher             types.Hasher
}

// NewTree creates a new instance of merkle tree. requires a hash algorithm to be specified.
func NewTree(hasher types.Hasher) Tree {
	return Tree{
		currentWorkingTree: NewPartialTree(hasher),
		UncommittedLeaves:  [][]byte{},
		hasher:             hasher,
	}
}

// PartialTree represents a part of the original tree that is enough to calculate the root.
// Used in to extract the root in a merkle proof, to apply diff to a tree or to merge
// multiple trees into one.
// It is a rare case when you need to use this struct on it's own. It's mostly used inside
type PartialTree struct {
	layers Layers
	hasher types.Hasher
}

// NewPartialTree Takes hasher as an argument and build a Merkle Tree from them.
// Since it's a partial tree, hashes must be accompanied by their index in the original tree.
func NewPartialTree(hasher types.Hasher) PartialTree {
	return PartialTree{
		layers: Layers{},
		hasher: hasher,
	}
}

// Proof is used to parse, verify, calculate a root for Merkle proofs.
// Proof requires specifying hashing algorithm and hash size in order to work.
// The hashing algorithm is set through the Hasher interface, which is supplied as a generic
// parameter to the Proof.
type Proof struct {
	proofHashes      [][]byte
	leaves           Leaves
	totalLeavesCount uint64
	hasher           types.Hasher
}

// NewProof create new instance of merkle proof
func NewProof(leaves Leaves, proofHashes [][]byte, totalLeavesCount uint64, hasher types.Hasher) Proof {
	return Proof{
		leaves:           leaves,
		proofHashes:      proofHashes,
		totalLeavesCount: totalLeavesCount,
		hasher:           hasher,
	}
}
