package mmr

import (
	"golang.org/x/crypto/blake2b"
	"reflect"
	"testing"
)

func merge(left, right interface{}) interface{} {
	l := left.([]byte)
	r := right.([]byte)
	hash := blake2b.Sum256(append(l, r...))
	return hash[:]
}

func TestCalculateRoot3Peaks(t *testing.T) {
	leaves := []leaf{{pos: 8, hash: []byte("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3")}}
	root := []byte("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3")
	mmrSize := uint64(19)
	proofIter := &Iterator{item: []interface{}{
		[]byte("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249"),
		[]byte("64fa1a16b569918daf33bf20fc82cab12b506357dbf176a6f6dac3f14108d45c"),
		[]byte("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
		[]byte("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
	}}
	m := MMR{merge: merge}
	got, err := m.CalculateRoot(leaves, mmrSize, proofIter)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if !reflect.DeepEqual(got, root) {
		t.Errorf("want %v  got %v", root, got)
	}
}


