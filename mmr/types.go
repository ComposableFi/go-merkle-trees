package mmr

type leafWithashOfH struct {
	pos    uint64
	hash   []byte
	height uint32
}

type peak struct {
	height uint32
	pos    uint64
}
