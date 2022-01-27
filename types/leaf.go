package types

// Leaf is the type of leaf used in the merkle tree and mmr
type Leaf struct {
	Index uint64
	Hash  []byte
}
