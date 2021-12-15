package mmr

type Leaf struct {
	Pos  uint64
	Hash interface{}
}

type leafWithHash struct {
	pos    uint64
	hash   interface{}
	height uint32
}

type peak struct {
	height uint32
	pos    uint64
}
