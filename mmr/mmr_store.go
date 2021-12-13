package mmr

type Store interface {
	getElem(pos uint64) interface{}
	append(pos uint64, elems []interface{})
}

type BatchElem struct {
	pos   uint64
	elems []interface{}
}

type Batch struct {
	memoryBatch []BatchElem
	store       Store
}

func NewBatch(store Store) *Batch {
	return &Batch{
		memoryBatch: []BatchElem{},
		store:       store,
	}
}

func (b *Batch) append(pos uint64, elems []interface{}) {
	b.memoryBatch = append(b.memoryBatch, BatchElem{pos, elems})
}

func (b *Batch) getElem(pos uint64) interface{} {
	memoryBatch := make([]BatchElem, len(b.memoryBatch))
	copy(memoryBatch, b.memoryBatch)
	reverse(memoryBatch)

	for _, mb := range memoryBatch {
		startPos, elems := mb.pos, mb.elems
		if pos < startPos {
			continue
		} else if pos < startPos+uint64(len(elems)) {
			return elems[pos-startPos]
		} else {
			break
		}
	}
	return b.store.getElem(pos)
}

func (b *Batch) commit() struct{} {
	for _, mb := range b.memoryBatch {
		b.store.append(mb.pos, mb.elems)
	}
	return struct{}{}
}
