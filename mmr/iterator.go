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

func (i *Iterator) push(item interface{}) {
	i.Items = append(i.Items, item)
}

func (i *Iterator) get(index int) interface{} {
	return i.Items[index]
}

func (i *Iterator) length() int {
	return len(i.Items)
}

func (i *Iterator) splitOff(index int) []interface{} {
	split := i.Items[i.length()-index:]
	i.Items = i.Items[:index]
	return split
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
