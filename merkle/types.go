package merkle

// Merge is the interface for merge function of tree
type Merge interface {
	Merge(left, right interface{}) interface{}
}

// Tree is representation type for the merkle tree
type Tree struct {
	Nodes []interface{}
	Merge Merge
}

// Proof is the representation of a merkle proof
type Proof struct {
	Indices []uint32
	Lemmas  []interface{}
	Merge   Merge
}

// LeafIndex is the representation of a leaf index
type LeafIndex struct {
	Index uint32
	Leaf  interface{}
}
