package merkle

type Hash []byte

type Tree struct {
	currentWorkingTree PartialTree
	history            []PartialTree
	UncommittedLeaves  []Hash
	hasher             Hasher
}

func NewTree(hasher Hasher) Tree {
	return Tree{
		currentWorkingTree: NewPartialTree(hasher),
		history:            []PartialTree{},
		UncommittedLeaves:  []Hash{},
		hasher:             hasher,
	}
}

type PartialTree struct {
	layers [][]Leaf
	hasher Hasher
}

func NewPartialTree(hasher Hasher) PartialTree {
	return PartialTree{
		layers: [][]Leaf{},
		hasher: hasher,
	}
}

type Proof struct {
	proofHashes []Hash
	// leaves           []Leaf
	// totalLeavesCount uint32
	hasher Hasher
}

// merkle.NewProof(authorityLeaves, proofHashes, totalLeavesCount, Keccak256{})

func NewProof(proofHashes []Hash, hasher Hasher) Proof {
	return Proof{
		// leaves:           leaves,
		proofHashes: proofHashes,
		// totalLeavesCount: totalLeavesCount,
		hasher: hasher,
	}
}

type Hasher interface {
	Hash(data []byte) (Hash, error)
}

func ConcatAndHash(hasher Hasher, left []byte, right []byte) (Hash, error) {
	if right == nil {
		return left, nil
	}
	return hasher.Hash(append(left[:], right[:]...))
}

type Leaf struct {
	Index uint32
	Hash  Hash
}
