package mmr

import (
	"golang.org/x/crypto/blake2b"
)

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

func (m MemStore) getElem(pos uint64) []byte {
	return m[pos]
}

// MemMMR is MMR implementation that uses a memory store. It has the fields store which takes the MemStore as an argument
// and the mmrSize of the tree.
type MemMMR struct {
	store   MemStore
	mmrSize uint64
	leaves  []Leaf
}

// NewMemMMR returns an object of the MemMMR type. It takes in the mmrSize and Memory Store (MemStore) as arguments.
func NewMemMMR(mmrSize uint64, store MemStore) *MemMMR {
	return &MemMMR{
		mmrSize: mmrSize,
		store:   store,
	}
}

// Store returns the MemStore
func (m *MemMMR) Store() MemStore {
	return m.store
}

// GetRoot returns the root of the MMR tree
func (m *MemMMR) GetRoot() (interface{}, error) {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, m.leaves, merge)
	return mmr.GetRoot()
}

// Push adds an element to the store and returns the position of the element
func (m *MemMMR) Push(elem []byte) uint64 {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, m.leaves, merge)
	pos, err := mmr.Push(elem)
	if err != nil {
		log.Error(err.Error())
		return 0
	}
	mmr.Commit()
	return pos
}

// GenProof generates proofs for validating an MMR leaf
func (m *MemMMR) GenProof(posList []uint64) (*Proof, error) {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, m.leaves, merge)
	return mmr.GenProof(posList)
}

// Merge is an empty struct that satisfies the Merge interface by implementing the Merge method.
type Merge struct{}

// Merge takes two arguments of type []byte, appends both into a single byte slice, hashes the result using blake2b
// and returns the resultant hash.
func (m *Merge) Merge(left, right interface{}) interface{} {
	l := left.([]byte)
	r := right.([]byte)
	hash := blake2b.Sum256(append(l, r...))
	return hash[:]
}
