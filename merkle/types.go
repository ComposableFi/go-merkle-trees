package merkle

type Hash []byte

type Tree struct {
	currentWorkingTree PartialTree
	history            []PartialTree
	UncommittedLeaves  []Hash
	hasher             hasher
}

func NewTree(hasher hasher) Tree {
	return Tree{
		currentWorkingTree: NewPartialTree(hasher),
		history:            []PartialTree{},
		UncommittedLeaves:  []Hash{},
		hasher:             hasher,
	}
}

type PartialTree struct {
	layers [][]leafIndexAndHash
	hasher hasher
}

func NewPartialTree(hasher hasher) PartialTree {
	return PartialTree{
		layers: [][]leafIndexAndHash{},
		hasher: hasher,
	}
}

type hasher interface {
	Hash(data []byte) (Hash, error)
	ConcatAndHash(left, right []byte) (Hash, error)
}

type leafIndexAndHash struct {
	index uint32
	hash  Hash
}
