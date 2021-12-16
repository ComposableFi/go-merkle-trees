package mmr

import (
	"golang.org/x/crypto/blake2b"
)

type MemStore map[uint64]interface{}

func NewMemStore() MemStore {
	return make(MemStore, 0)
}

func (m MemStore) append(pos uint64, elem []interface{}) {
	for index, value := range elem {
		m[pos+uint64(index)] = value
	}
}

func (m MemStore) getElem(pos uint64) interface{} {
	return m[pos]
}

type MemMMR struct {
	store   MemStore
	mmrSize uint64
}

func NewMemMMR(mmrSize uint64, store MemStore) *MemMMR {
	return &MemMMR{
		mmrSize: mmrSize,
		store:   store,
	}
}

func (m *MemMMR) Store() MemStore {
	return m.store
}

func (m *MemMMR) GetRoot() (interface{}, error) {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, merge)
	return mmr.GetRoot()
}

func (m *MemMMR) Push(elem interface{}) interface{} {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, merge)
	pos, err := mmr.Push(elem)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if mmr.Commit() == nil {
		return nil
	}
	return pos
}

func (m *MemMMR) GenProof(posList []uint64) (*MerkleProof, error) {
	merge := &Merge{}
	mmr := NewMMR(m.mmrSize, m.store, merge)
	return mmr.GenProof(posList)
}

type Merge struct{}

func (m *Merge) Merge(left, right interface{}) interface{} {
	l := left.([]byte)
	r := right.([]byte)
	hash := blake2b.Sum256(append(l, r...))
	return hash[:]
}
