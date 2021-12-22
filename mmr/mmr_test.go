package mmr

import (
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"testing"

	"golang.org/x/crypto/blake2b"
)

type NumberHash []byte

func (n NumberHash) From(num uint32) interface{} {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	hash := blake2b.Sum256(b)
	return hash[:]
}

func testMMR(count uint32, proofElem []uint32) error {
	mmr := NewMMR(0, NewMemStore(), &Merge{})
	var positions []uint64
	for i := uint32(0); i < count; i++ {
		position, err := mmr.Push(NumberHash{}.From(i))
		if err != nil {
			return err
		}
		positions = append(positions, position)
	}

	root, err := mmr.GetRoot()
	if err != nil {
		return err
	}

	proof, err := mmr.GenProof(func() []uint64 {
		var elem []uint64
		for _, p := range proofElem {
			elem = append(elem, positions[uint(p)])
		}
		return elem
	}())
	if err != nil {
		return err
	}

	mmr.Commit()

	result := proof.Verify(root, func() []Leaf {
		var leaves []Leaf
		for _, e := range proofElem {
			leaves = append(leaves, Leaf{positions[e], NumberHash{}.From(e)})
		}
		return leaves
	}())

	if !result {
		return err
	}

	return nil
}

func TestMMR(t *testing.T) {
	tests := map[string]struct {
		count     uint32
		proofElem []uint32
	}{
		"3 peaks":                        {11, []uint32{5}},
		"2 peaks":                        {10, []uint32{5}},
		"1 peak":                         {8, []uint32{5}},
		"first elem proof":               {11, []uint32{0}},
		"last elem proof":                {11, []uint32{10}},
		"1 elem":                         {1, []uint32{0}},
		"2 elems":                        {2, []uint32{0}},
		"2 elems*":                       {2, []uint32{1}},
		"2 leaves merkle proof":          {11, []uint32{3, 7}},
		"2 leaves merkle proof*":         {11, []uint32{3, 4}},
		"2 sibling leaves merkle proof":  {11, []uint32{4, 5}},
		"2 sibling leaves merkle proof*": {11, []uint32{5, 6}},
		"3 leaves merkle proof":          {11, []uint32{4, 5, 6}},
		"3 leaves merkle proof*":         {11, []uint32{3, 5, 7}},
		"3 leaves merkle proof**":        {11, []uint32{3, 4, 5}},
		"3 leaves merkle proof***":       {100, []uint32{3, 5, 13}},
	}

	for name, test := range tests {
		err := testMMR(test.count, test.proofElem)
		if err != nil {
			t.Errorf("%s: %s", name, err.Error())
		}
	}
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
		positions = append(positions, position)
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

	mmr.Commit()
	calculatedRoot, err := proof.CalculateRootWithNewLeaf(
		[]Leaf{{pos, NumberHash{}.From(elem)}},
		newPos,
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

func TestEmptyMMRRoot(t *testing.T) {
	merge := &Merge{}
	store := NewMemStore()
	mmr := NewMMR(0, store, merge)
	_, err := mmr.GetRoot()
	if err != ErrGetRootOnEmpty {
		t.Errorf("%s: want :%v  got %v", "empty mmr root", ErrGetRootOnEmpty, err)
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
