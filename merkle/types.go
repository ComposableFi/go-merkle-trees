package merkle

// Merge is the interface for merge function of tree
type Merge interface {
	Merge(left, right uint64) uint64
}

// Tree is representation type for the merkle tree
type Tree struct {
	Nodes []uint64
	Merge Merge
}

// Proof is the representation of a merkle proof
type Proof struct {
	Indices []uint64
	Lemmas  []uint64
	Merge   Merge
}

// LeafIndex is the representation of a leaf index
type LeafIndex struct {
	Index uint64
	Leaf  uint64
}
