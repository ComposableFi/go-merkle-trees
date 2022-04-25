package mmr

// Store defines the required method on any store passed to the Batch struct
type Store interface {
	GetElem(pos uint64) []byte
	append(pos uint64, elems [][]byte)
}

// BatchElem holds the fields of data for a Batch Element
type BatchElem struct {
	pos   uint64
	elems [][]byte
}

// Batch contains the a slice of Batch elements and a Store
type Batch struct {
	memoryBatch []BatchElem
	store       Store
}

// NewBatch returns an object of the Batch type
func NewBatch(store Store) *Batch {
	return &Batch{
		memoryBatch: []BatchElem{},
		store:       store,
	}
}

func (b *Batch) append(pos uint64, elems [][]byte) {
	b.memoryBatch = append(b.memoryBatch, BatchElem{pos, elems})
}

// GetElem returns an element in a store implementation using its position.
func (b *Batch) GetElem(pos uint64) []byte {
	i := len(b.memoryBatch)
	for i > 0 {
		mb := b.memoryBatch[i-1]
		startPos, elems := mb.pos, mb.elems
		if pos < startPos {
			i -= 1
			continue
		} else if pos < startPos+uint64(len(elems)) {
			return elems[pos-startPos]
		} else {
			break
		}
	}
	return b.store.GetElem(pos)
}

func (b *Batch) commit() {
	for _, mb := range b.memoryBatch {
		b.store.append(mb.pos, mb.elems)
	}
}
