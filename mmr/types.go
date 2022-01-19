package mmr

// Leaf is an mmr leaf. It also holds the field Leaf which is a byte representation of the leaf and Pos, the leaf
// position.
type Leaf struct {
	Index uint64
	Hash  []byte
}

type leafWithashOfH struct {
	pos    uint64
	hash   []byte
	height uint32
}

type peak struct {
	height uint32
	pos    uint64
}
