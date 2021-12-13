package mmr

import (
	"encoding/binary"
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"reflect"
	"testing"
)

func hexDecode(h string) []byte {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return b
}

func TestCalculateRoot(t *testing.T) {
	type input struct {
		leaves    []Leaf
		mmrSize   uint64
		proofIter *Iterator
	}

	tests := map[string]struct {
		input input
		want  []byte
	}{
		"3 peaks": {
			input: input{
				leaves:  []Leaf{{pos: 8, hash: hexDecode("8c35d22f459d77ca4c0b0b5035869766d60d182b9716ab3e8879e066478899a8")}},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249"),
					hexDecode("64fa1a16b569918daf33bf20fc82cab12b506357dbf176a6f6dac3f14108d45c"),
					hexDecode("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
					hexDecode("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
		"2 peaks": {
			input: input{
				leaves:  []Leaf{{pos: 8, hash: hexDecode("8c35d22f459d77ca4c0b0b5035869766d60d182b9716ab3e8879e066478899a8")}},
				mmrSize: 18,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249"),
					hexDecode("64fa1a16b569918daf33bf20fc82cab12b506357dbf176a6f6dac3f14108d45c"),
					hexDecode("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
					hexDecode("a8682c7cd2e2d29666a2597c25cebc0aa23c8e3e6b3ee3e04f2bc1713da64785"),
				}},
			},
			want: hexDecode("008bc96cb3da2098e7ccacd0548b36abadb91b153b28cbc1cb513315ea60da46"),
		},
		"1 peak": {
			input: input{
				leaves:  []Leaf{{pos: 8, hash: hexDecode("8c35d22f459d77ca4c0b0b5035869766d60d182b9716ab3e8879e066478899a8")}},
				mmrSize: 18,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249"),
					hexDecode("64fa1a16b569918daf33bf20fc82cab12b506357dbf176a6f6dac3f14108d45c"),
					hexDecode("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
				}},
			},
			want: hexDecode("9a42e0efd142d8dadb23fe5eec4b5078c9ec740d9927109646bcad1c9939be80"),
		},
		"first element proof": {
			input: input{
				leaves:  []Leaf{{pos: 0, hash: hexDecode("11da6d1f761ddf9bdb4c9d6e5303ebd41f61858d0a5647a1a7bfe089bf921be9")}},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("e12c22d4f162d9a012c9319233da5d3e923cc5e1029b8f90e47249c9ab256b35"),
					hexDecode("ea750bdb0a08f96991f00ceaf9c3517805b1844866091df48b3612a24225429a"),
					hexDecode("405a3e50f8d864f3a15f28b7f290363cf26727760a9144ff364b3aa8ccd2f839"),
					hexDecode("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
		"last element proof": {
			input: input{
				leaves:  []Leaf{{pos: 18, hash: hexDecode("2c088bf3b4e7853c99e49636d9e7c9a351918d70bd6cdf6148b81e68f5706f68")}},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("9a42e0efd142d8dadb23fe5eec4b5078c9ec740d9927109646bcad1c9939be80"),
					hexDecode("a8682c7cd2e2d29666a2597c25cebc0aa23c8e3e6b3ee3e04f2bc1713da64785"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
		"1 element": {
			input: input{
				leaves:    []Leaf{{pos: 0, hash: hexDecode("11da6d1f761ddf9bdb4c9d6e5303ebd41f61858d0a5647a1a7bfe089bf921be9")}},
				mmrSize:   1,
				proofIter: &Iterator{Items: []interface{}{}},
			},
			want: hexDecode("11da6d1f761ddf9bdb4c9d6e5303ebd41f61858d0a5647a1a7bfe089bf921be9"),
		},
		"2 elements pos 0": {
			input: input{
				leaves:  []Leaf{{pos: 0, hash: hexDecode("11da6d1f761ddf9bdb4c9d6e5303ebd41f61858d0a5647a1a7bfe089bf921be9")}},
				mmrSize: 3,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("e12c22d4f162d9a012c9319233da5d3e923cc5e1029b8f90e47249c9ab256b35"),
				}},
			},
			want: hexDecode("dd1445ec419376975790d7d4e487dfa5fca42a75f41e8893bc2d8b02c527f8f4"),
		},
		"2 elements pos 1": {
			input: input{
				leaves:  []Leaf{{pos: 1, hash: hexDecode("e12c22d4f162d9a012c9319233da5d3e923cc5e1029b8f90e47249c9ab256b35")}},
				mmrSize: 3,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("11da6d1f761ddf9bdb4c9d6e5303ebd41f61858d0a5647a1a7bfe089bf921be9"),
				}},
			},
			want: hexDecode("dd1445ec419376975790d7d4e487dfa5fca42a75f41e8893bc2d8b02c527f8f4"),
		},
		"2 leaves merkle proof": {
			input: input{
				leaves: []Leaf{
					{pos: 4, hash: hexDecode("8c039ff7caa17ccebfcadc44bd9fce6a4b6699c4d03de2e3349aa1dc11193cd7")},
					{pos: 11, hash: hexDecode("5b8f29db76cf4e676e4fc9b17040312debedafcd5637fb3c7badd2cddce6a445")},
				},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("7b0aa1735e5ba58d3236316c671fe4f00ed366ee72417c9ed02a53a8019e85b8"),
					hexDecode("f4aac2fbe33f03554bfeb559ea2690ed8521caa4be961e61c91ac9a1530dce7a"),
					hexDecode("dd1445ec419376975790d7d4e487dfa5fca42a75f41e8893bc2d8b02c527f8f4"),
					hexDecode("86e27cc779fc3f19c1cf4ece5f9ae8a2b7cc24301e2cd17fad15342e495c187d"),
					hexDecode("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
		"2 sibling leaves merkle proof": {
			input: input{
				leaves: []Leaf{
					{pos: 7, hash: hexDecode("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249")},
					{pos: 8, hash: hexDecode("8c35d22f459d77ca4c0b0b5035869766d60d182b9716ab3e8879e066478899a8")},
				},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("64fa1a16b569918daf33bf20fc82cab12b506357dbf176a6f6dac3f14108d45c"),
					hexDecode("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
					hexDecode("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
		"3 leaves merkle proof": {
			input: input{
				leaves: []Leaf{
					{pos: 7, hash: hexDecode("26a08e4d0c5190f01871e0569b6290b86760085d99f17eb4e7e6b58feb8d6249")},
					{pos: 8, hash: hexDecode("8c35d22f459d77ca4c0b0b5035869766d60d182b9716ab3e8879e066478899a8")},
					{pos: 10, hash: hexDecode("f4aac2fbe33f03554bfeb559ea2690ed8521caa4be961e61c91ac9a1530dce7a")},
				},
				mmrSize: 19,
				proofIter: &Iterator{Items: []interface{}{
					hexDecode("5b8f29db76cf4e676e4fc9b17040312debedafcd5637fb3c7badd2cddce6a445"),
					hexDecode("f0c1d8dd595c1e705ff3e42cb104b5b558b2fe09577b1ec44c0e0ea67982a884"),
					hexDecode("5e7bc66323e34ccbbe88ba9172c6dadbb050f0fb60a04b761e95ab46580ddc5e"),
				}},
			},
			want: hexDecode("f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"),
		},
	}

	for name, test := range tests {
		m := MerkleProof{Merge: &Merge{}}
		got, err := m.CalculateRoot(test.input.leaves, test.input.mmrSize, test.input.proofIter)
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want %x  got %x", name, test.want, got)
		}
	}
}

func TestEmptyMMRRoot(t *testing.T) {
	merge := &Merge{}
	store := NewMemStore()
	mmr := NewMMR(0, store, merge)
	_, err := mmr.GetRoot()
	if err != ErrGetRootOnEmpty {
		t.Errorf("%s: want :%v  got %v", "empty mmr root", ErrGetRootOnEmpty, err)
	}
}

type NumberHash []byte

func (n NumberHash) From(num uint32) interface{} {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	hash := blake2b.Sum256(b)
	return hash[:]
}

func TestGenRootFromProof(t *testing.T) {
	merge := &Merge{}
	store := NewMemStore()
	mmr := NewMMR(0, store, merge)

	count := 11
	var positions []uint64
	for i := 0; i < 11; i++ {
		position, err := mmr.Push(NumberHash{}.From(uint32(i)))
		if err != nil {
			t.Errorf("%s: %s", "mmr root", err.Error())
			return
		}
		positions = append(positions, position.(uint64))
	}

	var elem = uint32(count - 1)
	var pos = positions[uint(elem)]
	proof, err := mmr.GenProof([]uint64{pos})
	if err != nil {
		t.Errorf("%s: %s", "mmr gen proof", err.Error())
		return
	}

	newElem := count
	newPos, err := mmr.Push(NumberHash{}.From(uint32(newElem)))
	if err != nil {
		t.Errorf("%s: %s", "mmr gen proof", err.Error())
		return
	}

	root, err := mmr.GetRoot()
	if err != nil {
		t.Errorf("%s: %s", "mmr root", err.Error())
		return
	}

	commit := mmr.Commit()
	if commit == nil {
		t.Errorf("%s: %s", "mmr root", "commit changes")
		return
	}

	calculatedRoot, err := proof.CalculateRootWithNewLeaf(
		[]Leaf{{pos, NumberHash{}.From(elem)}},
		newPos.(uint64),
		NumberHash{}.From(uint32(newElem)),
		LeafIndexToMMRSize(uint64(newElem)),
	)
	if err != nil {
		t.Errorf("%s: %s", "mmr root calculateRootWithNewLeaf", err.Error())
		return
	}

	if !reflect.DeepEqual(calculatedRoot, root) {
		t.Errorf("%s: want :%v  got %v", "empty mmr root", root, calculatedRoot)
	}
}

func TestMMRRoot(t *testing.T) {
	merge := &Merge{}
	store := NewMemStore()
	mmr := NewMMR(0, store, merge)
	for i := 0; i < 11; i++ {
		_, err := mmr.Push(NumberHash{}.From(uint32(i)))
		if err != nil {
			t.Errorf("%s: %s", "mmr root", err.Error())
			return
		}
	}

	root, err := mmr.GetRoot()
	if err != nil {
		t.Errorf("%s: %s", "mmr root", err.Error())
	}

	want := "f6794677f37a57df6a5ec36ce61036e43a36c1a009d05c81c9aa685dde1fd6e3"
	if !reflect.DeepEqual(hex.EncodeToString(root.([]byte)), want) {
		t.Errorf("%s: want :%v  got %v", "empty mmr root", want, hex.EncodeToString(root.([]byte)))
	}
}
