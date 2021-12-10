package mmr

type Store interface {
	getElem(pos uint64) interface{}
	append(pos uint64, elems []interface{})
}

type BatchElem struct {
	pos  uint64
	elem []interface{}
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
	memoryBatch := b.memoryBatch
	reverse(memoryBatch)
	for _, mb := range memoryBatch {
		startPos, elems := mb.pos, mb.elem
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

// TODO: implement batch commit
func (b *Batch) commit() struct{} {
	return struct{}{}
}
