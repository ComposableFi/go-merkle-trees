package mmr

// MemStore is a map of bytes with uint64 values as its key
type MemStore map[uint64][]byte

// NewMemStore creates an returns a map of the MemStore type
func NewMemStore() MemStore {
	return make(MemStore)
}

func (m MemStore) append(pos uint64, elem [][]byte) {
	for index, value := range elem {
		m[pos+uint64(index)] = value
	}
}

func (m MemStore) GetElem(pos uint64) []byte {
	return m[pos]
}
