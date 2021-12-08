package mmr

type Iterator struct {
	Items []interface{}
	index uint64
}

func NewIterator() *Iterator {
	return &Iterator{
		Items: make([]interface{}, 0),
	}
}

func (i *Iterator) next() interface{} {
	if len(i.Items) > int(i.index) {
		in := i.index
		i.index++
		return i.Items[in]
	}
	return nil
}

func (i *Iterator) isEmpty() bool {
	return len(i.Items) == 0
}
