package mmr

// Iterator is a wrapper for a slice of bytes. It exposes helper methods for accessing about the slice and storing data
// to it
type Iterator struct {
	Items [][]byte
	index uint64
}

// NewIterator creates a new object of the Iterator
func NewIterator() *Iterator {
	return &Iterator{
		Items: make([][]byte, 0),
	}
}

func (i *Iterator) push(item []byte) {
	i.Items = append(i.Items, item)
}

func (i *Iterator) length() int {
	return len(i.Items)
}

func (i *Iterator) splitOff(index int) [][]byte {
	split := i.Items[index:]
	i.Items = i.Items[:index]
	return split
}

// Next returns the next item from the slice of items and increases the index. It returns nil when the last item in
// the slice has already been returned.
func (i *Iterator) Next() []byte {
	if len(i.Items) > int(i.index) {
		in := i.index
		i.index++
		return i.Items[in]
	}
	return nil
}
