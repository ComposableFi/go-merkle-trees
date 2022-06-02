package mmr

// MemStore is a map of bytes with uint64 values as its key
type MemStore map[uint64][]byte

// NewMemStore creates an returns a map of the MemStore type
func NewMemStore() MemStore {
	return make(MemStore)
}

func (m MemStore) append(pos uint64, elem [][]byte) {
	for i := 0; i < len(elem); i++ {
		m[pos+uint64(i)] = elem[i]
	}
}

func (m MemStore) GetElem(pos uint64) []byte {
	return m[pos]
}
