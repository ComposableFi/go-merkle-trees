package merkle

// Merge is the interface for merge function of tree
type Merge interface {
	Merge(left []byte, right []byte) []byte
}

// Tree is representation type for the merkle tree
type Tree struct {
	Nodes [][]byte
	Merge Merge
}

// Proof is the representation of a merkle proof
type Proof struct {
	Leaves []LeafData
	Proofs [][]byte
	Merge  Merge
}

// LeafData is the representation of a leaf index
type LeafData struct {
	Index uint32
	Leaf  []byte
}

// CBMT is representation type for the complete binary merkle tree
type CBMT struct {
	Merge Merge
}
