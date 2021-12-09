package mmr

type Store interface {
	getElement(pos uint64) interface{}
	append(pos uint64, elems []interface{}) interface{}
}

type Batch struct {
	memoryBatch []Leaf
	store       Store
}

func NewBatch(store Store) *Batch {
	return &Batch{
		memoryBatch: []Leaf{},
		store:       store,
	}
}

func (b *Batch) append(pos uint64, elems []interface{}) {
	b.memoryBatch = append(b.memoryBatch, Leaf{pos, elems})
}

func (b *Batch) getElem(pos uint64) (interface{}, error) {
	// TODO:  only search memoryBatch for elem.
	return nil, nil
}

func (b *Batch) commit() struct{} {
	return struct{}{}
}
