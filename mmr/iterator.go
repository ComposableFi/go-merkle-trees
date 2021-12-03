package mmr

type Iterator struct {
	item  []interface{}
	index uint64
}

func NewIterator() *Iterator {
	return &Iterator{
		item: make([]interface{}, 0),
	}
}

func (i *Iterator) next() interface{} {
	if len(i.item) > int(i.index) {
		in := i.index
		i.index++
		return i.item[in]
	}
	return nil
}

func (i *Iterator) isEmpty() bool {
	if len(i.item) == 0 {
		return true
	}
	return false
}
