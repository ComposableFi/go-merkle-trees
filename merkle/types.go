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
	Leaves []LeafData
	Lemmas []interface{}
	Merge  Merge
}

// LeafData is the representation of a leaf index
type LeafData struct {
	Index uint32
	Leaf  interface{}
}

// CBMT is representation type for the complete binary merkle tree
type CBMT struct {
	Merge Merge
}
