package mmr

// Leaf is an mmr leaf. It also holds the field Leaf which is a byte representation of the leaf and Pos, the leaf
// position.
type Leaf struct {
	Pos  uint64
	Leaf []byte
}

type leafWithHash struct {
	pos    uint64
	hash   []byte
	height uint32
}

type peak struct {
	height uint32
	pos    uint64
}
