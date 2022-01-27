package mmr_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/ComposableFi/go-merkle-trees/hasher"
	merkleMmr "github.com/ComposableFi/go-merkle-trees/mmr"
	"github.com/ComposableFi/go-merkle-trees/types"
)

func uint32ToHash(num uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, num)
	h, err := hasher.Sha256Hasher{}.Hash(b)
	if err != nil {
		panic(err)
	}
	return h[:]
}

func testMMR(count uint32, proofElem []uint32) error {
	leaves := func() []types.Leaf {
		var leaves []types.Leaf
		for _, e := range proofElem {
			leaves = append(leaves, types.Leaf{Index: uint64(e), Hash: uint32ToHash(e)})
		}
		return leaves
	}()

	mmrTree := merkleMmr.NewMMR(0, merkleMmr.NewMemStore(), leaves, hasher.Keccak256Hasher{})
	var positions []uint64
	for i := uint32(0); i < count; i++ {
		position, err := mmrTree.Push(uint32ToHash(i))
		if err != nil {
			return err
		}
		positions = append(positions, position)
	}

	root, err := mmrTree.Root()
	if err != nil {
		return err
	}

	proof, err := mmrTree.GenProof(func() []uint64 {
		var elem []uint64
		for _, p := range proofElem {
			elem = append(elem, positions[uint64(p)])
		}
		return elem
	}())
	if err != nil {
		return err
	}

	mmrTree.Commit()

	result := proof.Verify(root)

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
		"node 0":                         {7, []uint32{0}},
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
		if name == "node 0" {
			err := testMMR(test.count, test.proofElem)
			if err != nil {
				t.Errorf("%s: %s", name, err.Error())
			}
		}
	}
}

func TestGenRootFromProof(t *testing.T) {
	store := merkleMmr.NewMemStore()

	count := 11
	elem := uint32(count - 1)
	leaves := []types.Leaf{{Index: uint64(elem), Hash: uint32ToHash(elem)}}
	mmr := merkleMmr.NewMMR(0, store, leaves, hasher.Keccak256Hasher{})

	var positions []uint64
	for i := 0; i < 11; i++ {
		position, err := mmr.Push(uint32ToHash(uint32(i)))
		if err != nil {
			t.Errorf("%s: %s", "merkleMmr root", err.Error())
			return
		}
		positions = append(positions, position)
	}

	var pos = positions[uint64(elem)]
	proof, err := mmr.GenProof([]uint64{pos})
	if err != nil {
		t.Errorf("%s: %s", "merkleMmr gen proof", err.Error())
		return
	}

	newElem := count
	_, err = mmr.Push(uint32ToHash(uint32(newElem)))
	if err != nil {
		t.Errorf("%s: %s", "merkleMmr gen proof", err.Error())
		return
	}

	root, err := mmr.Root()
	if err != nil {
		t.Errorf("%s: %s", "merkleMmr root", err.Error())
		return
	}

	mmr.Commit()
	calculatedRoot, err := proof.CalculateRootWithNewLeaf(
		leaves,
		uint64(newElem),
		uint32ToHash(uint32(newElem)),
		merkleMmr.LeafIndexToMMRSize(uint64(newElem)),
		//merkleMmr.LeafIndexToMMRSize(uint64(newElem)),
	)
	if err != nil {
		t.Errorf("%s: %s", "merkleMmr root calculateRootWithNewLeaf", err.Error())
		return
	}

	if !reflect.DeepEqual(calculatedRoot, root) {
		t.Errorf("%s: want :%v  got %v", "empty merkleMmr root", root, calculatedRoot)
	}
}

func TestEmptyMMRRoot(t *testing.T) {
	store := merkleMmr.NewMemStore()
	mmr := merkleMmr.NewMMR(0, store, []types.Leaf{}, hasher.Keccak256Hasher{})
	_, err := mmr.Root()
	if err != merkleMmr.ErrGetRootOnEmpty {
		t.Errorf("%s: want :%v  got %v", "empty merkleMmr root", merkleMmr.ErrGetRootOnEmpty, err)
	}
}

func TestMMRRoot(t *testing.T) {
	store := merkleMmr.NewMemStore()
	mmr := merkleMmr.NewMMR(0, store, []types.Leaf{}, hasher.Keccak256Hasher{})
	for i := 0; i < 11; i++ {
		_, err := mmr.Push(uint32ToHash(uint32(i)))
		if err != nil {
			t.Errorf("%s: %s", "merkleMmr root", err.Error())
			return
		}
	}

	root, err := mmr.Root()
	rootHex := hex.EncodeToString(root)
	if err != nil {
		t.Errorf("%s: %s", "merkleMmr root", err.Error())
	}

	want := "285f5038cc67c811a4b2a470da53407afdf8ff673b18860f1154b55b974d55e2"
	if !reflect.DeepEqual(rootHex, want) {
		t.Errorf("%s: want :%v  got %v", "empty merkleMmr root", want, rootHex)
	}
}

func mergeKeccak256(left, right []byte) []byte {
	hash, _ := hasher.MergeAndHash(hasher.Keccak256Hasher{}, left, right)
	return hash
}

func hexToByte(h string) []byte {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return b
}

func Test7LeafVerify(t *testing.T) {
	// print tree structure
	fmt.Println("                 7-leaf MMR:           ")
	fmt.Println()
	fmt.Println("    Height 3 |      7")
	fmt.Println("    Height 2 |   3      6     10")
	fmt.Println("    Height 1 | 1  2   4  5   8  9    11")
	fmt.Println("             | |--|---|--|---|--|-----|-")
	fmt.Println("Hash indices | 0  1   2  3   4  5     6")

	// ---------------------------- Tree contents ----------------------------
	//  - For leaf nodes, node hash is the SCALE-encoding of the leaf data.
	//  - For parent nodes, node hash is the hash of it"s children (left, right).
	//
	// 0xda5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed // 1  LEAF NODE
	// 0xff5d891b28463a3440e1b650984685efdf260e482cb3807d53c49090841e755f // 2  LEAF NODE
	// 0xbc54778fab79f586f007bd408dca2c4aa07959b27d1f2c8f4f2549d1fcfac8f8 // 3  PARENT[1, 2] NODE
	// 0x7a84d84807ce4bbff8fb84667edf82aff5f2c5eb62e835f32093ee19a43c2de7 // 4  LEAF NODE
	// 0x27d8f4221cd6f7fc141ea20844c92aa8f647ac520853fbded619a46b1146ab8a // 5  LEAF NODE
	// 0x00b0046bd2d63fcb760cf50a262448bb2bbf9a264b0b0950d8744044edf00dc3 // 6  PARENT[4, 5] NODE
	// 0xe53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b // 7  PARENT[3, 6] NODE
	// 0x99af07747700389aba6e6cb0ee5d553fa1241688d9f96e48987bca1d7f275cbe // 8  LEAF NODE
	// 0xc09d4a008a0f1ef37860bef33ec3088ccd94268c0bfba7ff1b3c2a1075b0eb92 // 9  LEAF NODE
	// 0xdad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e // 10 PARENT[8, 9] NODE
	// 0xaf3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c // 11 LEAF NODE

	tests := map[string]struct {
		leafIndex uint64
		leafCount uint64
		leaf      []byte
		proofs    [][]byte
	}{
		"leaf index 0 (node 1)": {0, 7,
			hexToByte("da5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed"),
			[][]byte{
				hexToByte("ff5d891b28463a3440e1b650984685efdf260e482cb3807d53c49090841e755f"),
				hexToByte("00b0046bd2d63fcb760cf50a262448bb2bbf9a264b0b0950d8744044edf00dc3"),
				mergeKeccak256(hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"), // bag right hand side peaks keccak(right, left)
					hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"))}},
		"leaf index 1 (node 2)": {1, 7,
			hexToByte("ff5d891b28463a3440e1b650984685efdf260e482cb3807d53c49090841e755f"),
			[][]byte{
				hexToByte("da5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed"),
				hexToByte("00b0046bd2d63fcb760cf50a262448bb2bbf9a264b0b0950d8744044edf00dc3"),
				mergeKeccak256(hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"), // bag right hand side peaks keccak(right, left)
					hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"))}},
		"leaf index 2 (node 4)": {2, 7,
			hexToByte("7a84d84807ce4bbff8fb84667edf82aff5f2c5eb62e835f32093ee19a43c2de7"),
			[][]byte{
				hexToByte("27d8f4221cd6f7fc141ea20844c92aa8f647ac520853fbded619a46b1146ab8a"),
				hexToByte("bc54778fab79f586f007bd408dca2c4aa07959b27d1f2c8f4f2549d1fcfac8f8"),
				mergeKeccak256(hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"), // bag right hand side peaks keccak(right, left)
					hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"))}},
		"leaf index 3 (node 5)": {3, 7,
			hexToByte("27d8f4221cd6f7fc141ea20844c92aa8f647ac520853fbded619a46b1146ab8a"),
			[][]byte{
				hexToByte("7a84d84807ce4bbff8fb84667edf82aff5f2c5eb62e835f32093ee19a43c2de7"),
				hexToByte("bc54778fab79f586f007bd408dca2c4aa07959b27d1f2c8f4f2549d1fcfac8f8"),
				mergeKeccak256(hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"), // bag right hand side peaks keccak(right, left)
					hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"))}},
		"leaf index 4 (node 8)": {4, 7,
			hexToByte("99af07747700389aba6e6cb0ee5d553fa1241688d9f96e48987bca1d7f275cbe"),
			[][]byte{
				hexToByte("e53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b"),
				hexToByte("c09d4a008a0f1ef37860bef33ec3088ccd94268c0bfba7ff1b3c2a1075b0eb92"),
				hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"),
			}},
		"leaf index 5 (node 9)": {5, 7,
			hexToByte("c09d4a008a0f1ef37860bef33ec3088ccd94268c0bfba7ff1b3c2a1075b0eb92"),
			[][]byte{
				hexToByte("e53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b"),
				hexToByte("99af07747700389aba6e6cb0ee5d553fa1241688d9f96e48987bca1d7f275cbe"),
				hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"),
			}},
		"leaf index 6 (node 11)": {6, 7,
			hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"),
			[][]byte{
				hexToByte("e53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b"),
				hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"),
			}},
	}

	root := hexToByte("fc4f9042bd2f73feb26f3fc42db834c5f1943fa20070ddf106c486a478a0d561")
	for desc, test := range tests {
		leaves := []types.Leaf{{Index: test.leafIndex, Hash: test.leaf}}
		proof := merkleMmr.NewProof(merkleMmr.LeafIndexToMMRSize(test.leafCount-1), test.proofs, leaves, hasher.Keccak256Hasher{})
		if !proof.Verify(root) {
			t.Errorf("%s: failed to verify leaf inclusion", desc)
		}
	}

	tests = map[string]struct {
		leafIndex uint64
		leafCount uint64
		leaf      []byte
		proofs    [][]byte
	}{
		"invalid proofs": {5, 7,
			hexToByte("0000000000000000000000000000000000000000000000000000000000123456"),
			[][]byte{
				hexToByte("e53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b"),
				hexToByte("99af07747700389aba6e6cb0ee5d553fa1241688d9f96e48987bca1d7f275cbe"),
				hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"),
			}},
	}

	for desc, test := range tests {
		leaves := []types.Leaf{{Index: test.leafIndex, Hash: test.leaf}}
		proof := merkleMmr.NewProof(merkleMmr.LeafIndexToMMRSize(test.leafCount-1), test.proofs, leaves, hasher.Keccak256Hasher{})
		if proof.Verify(root) {
			t.Errorf("%s: verified a leaf inclusion with invalid proofs", desc)
		}
	}
}

func Test15LeafVerify(t *testing.T) {
	fmt.Println("                                    15-leaf MMR:                            ")
	fmt.Println("                                                                            ")
	fmt.Println("    Height 4 |             15                                               ")
	fmt.Println("    Height 3 |      7             14                22                      ")
	fmt.Println("    Height 2 |   3      6     10      13       18        21       25        ")
	fmt.Println("    Height 1 | 1  2   4  5   8  9   11  12   16  17   19   20   23  24  26  ")
	fmt.Println("             | |--|---|--|---|--|-----|---|---|---|----|---|----|---|---|---")
	fmt.Println("Hash indices | 0  1   2  3   4  5     6   7   8   9   10   11   12  13  14  ")

	// ---------------------------- Tree contents ----------------------------
	//  - For leaf nodes, node hash is the SCALE-encoding of the leaf data.
	//  - For parent nodes, node hash is the hash of it's children (left, right).
	//
	// 0xda5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed // 1  LEAF NODE
	// 0xff5d891b28463a3440e1b650984685efdf260e482cb3807d53c49090841e755f // 2  LEAF NODE
	// 0xbc54778fab79f586f007bd408dca2c4aa07959b27d1f2c8f4f2549d1fcfac8f8 // 3  PARENT[1, 2] NODE
	// 0x7a84d84807ce4bbff8fb84667edf82aff5f2c5eb62e835f32093ee19a43c2de7 // 4  LEAF NODE
	// 0x27d8f4221cd6f7fc141ea20844c92aa8f647ac520853fbded619a46b1146ab8a // 5  LEAF NODE
	// 0x00b0046bd2d63fcb760cf50a262448bb2bbf9a264b0b0950d8744044edf00dc3 // 6  PARENT[4, 5] NODE
	// 0xe53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b // 7  PARENT[3, 6] NODE
	// 0x99af07747700389aba6e6cb0ee5d553fa1241688d9f96e48987bca1d7f275cbe // 8  LEAF NODE
	// 0xc09d4a008a0f1ef37860bef33ec3088ccd94268c0bfba7ff1b3c2a1075b0eb92 // 9  LEAF NODE
	// 0xdad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e // 10 PARENT[8, 9] NODE
	// 0xaf3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c // 11 LEAF NODE
	// 0x643609ae1433f1d6caf366bb917873c3a3d82d7dc30e1c5e9a224d537f630dab // 12 LEAF NODE
	// 0x7fde31376facc58f621bacd80dfd77166544c84155bf1b82bf32281b93feaf78 // 13 PARENT[11, 12] NODE
	// 0xa63c4ec7ed257b6b4ab4fab3676f70b3b7c717357b537c0321d766de0e9e5312 // 14 PARENT[10, 13] NODE
	// 0xea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954 // 15 PARENT[7, 14] NODE
	// 0xbf5f579a06beced3256538b161b5096839db4b94ea1d3862bbe1fa5a2182e074 // 16 LEAF NODE
	// 0x7d8a0fe1021702eada6c608f3e09f833b63f21fdfe60f3bbb3401d5add4479af // 17 LEAF NODE
	// 0xa9ef6dd0b19d56f48a05c2475629c59713d0a992d335917135029432d611533d // 18 PARENT[16, 17] NODE
	// 0x2fd49d6e84591c6cc1fc38189b806dec1a1cb00c62727b63ac1cb9a37022c0fe // 19 LEAF NODE
	// 0x365f9e095800bd03add9be88b7f7bb06ff644ac2b77ce5da6a7c77e2fb19f1fb // 20 LEAF NODE
	// 0x3f7b0534bf60f62057a1ab9a0bf4751014d4d464245b5a7ad86801c9bac21b15 // 21 PARENT[19, 20] NODE
	// 0x16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f // 22 PARENT[18, 21] NODE
	// 0x94014b81bc56d64cac8dcde8eee47da0ed9b1319dccd9e86ad8d2266d8ef060a // 23 LEAF NODE
	// 0x883f1aca23002690575957cc85663774bbd3b9549ba5f0ee0fcc8aed9c88cf99 // 24 LEAF NODE
	// 0x1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4 // 25 PARENT[23, 24] NODE
	// 0x0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68 // 26 LEAF NODE

	tests := map[string]struct {
		leafIndex uint64
		leafCount uint64
		leaf      []byte
		proofs    [][]byte
	}{
		"leaf index 7 (node 12)": {7, 15,
			hexToByte("643609ae1433f1d6caf366bb917873c3a3d82d7dc30e1c5e9a224d537f630dab"),
			[][]byte{
				hexToByte("af3327deed0515c8d1902c9b5cd375942d42f388f3bfe3d1cd6e1b86f9cc456c"), // 11
				hexToByte("dad09f50b41822fc5ecadc25b08c3a61531d4d60e962a5aa0b6998fad5c37c5e"), // 10
				hexToByte("e53ee36ba6c068b1a6cfef7862fed5005df55615e1c9fa6eeefe08329ac4b94b"), // 7
				mergeKeccak256( // bag right hand side peaks keccak(right, left)
					mergeKeccak256(
						hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
						hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
					),
					hexToByte("16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f"), // 22
				),
			}},
		"leaf index 8 (node 16)": {8, 15,
			hexToByte("bf5f579a06beced3256538b161b5096839db4b94ea1d3862bbe1fa5a2182e074"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("7d8a0fe1021702eada6c608f3e09f833b63f21fdfe60f3bbb3401d5add4479af"), // 17
				hexToByte("3f7b0534bf60f62057a1ab9a0bf4751014d4d464245b5a7ad86801c9bac21b15"), // 21
				mergeKeccak256( // bag right hand side peaks keccak(right, left)
					hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
					hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
				),
			}},
		"leaf index 9 (node 17)": {9, 15,
			hexToByte("7d8a0fe1021702eada6c608f3e09f833b63f21fdfe60f3bbb3401d5add4479af"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("bf5f579a06beced3256538b161b5096839db4b94ea1d3862bbe1fa5a2182e074"), // 16
				hexToByte("3f7b0534bf60f62057a1ab9a0bf4751014d4d464245b5a7ad86801c9bac21b15"), // 21
				mergeKeccak256( // bag right hand side peaks keccak(right, left)
					hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
					hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
				),
			}},
		"leaf index 10 (node 19)": {10, 15,
			hexToByte("2fd49d6e84591c6cc1fc38189b806dec1a1cb00c62727b63ac1cb9a37022c0fe"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("365f9e095800bd03add9be88b7f7bb06ff644ac2b77ce5da6a7c77e2fb19f1fb"), // 20
				hexToByte("a9ef6dd0b19d56f48a05c2475629c59713d0a992d335917135029432d611533d"), // 18
				mergeKeccak256( // bag right hand side peaks keccak(right, left)
					hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
					hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
				),
			}},
		"leaf index 11 (node 20)": {11, 15,
			hexToByte("365f9e095800bd03add9be88b7f7bb06ff644ac2b77ce5da6a7c77e2fb19f1fb"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("2fd49d6e84591c6cc1fc38189b806dec1a1cb00c62727b63ac1cb9a37022c0fe"), // 19
				hexToByte("a9ef6dd0b19d56f48a05c2475629c59713d0a992d335917135029432d611533d"), // 18
				mergeKeccak256( // bag right hand side peaks keccak(right, left)
					hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
					hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
				),
			}},
		"leaf index 12 (node 23)": {12, 15,
			hexToByte("94014b81bc56d64cac8dcde8eee47da0ed9b1319dccd9e86ad8d2266d8ef060a"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f"), // 22
				hexToByte("883f1aca23002690575957cc85663774bbd3b9549ba5f0ee0fcc8aed9c88cf99"), // 24
				hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
			}},
		"leaf index 13 (node 24)": {13, 15,
			hexToByte("883f1aca23002690575957cc85663774bbd3b9549ba5f0ee0fcc8aed9c88cf99"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f"), // 22
				hexToByte("94014b81bc56d64cac8dcde8eee47da0ed9b1319dccd9e86ad8d2266d8ef060a"), // 23
				hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"), // 26
			}},
		"leaf index 14 (node 26)": {14, 15,
			hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"),
			[][]byte{
				hexToByte("ea97f06e80ac768687e72d4224999a51d272e1b4cafcbc64bd3ce63357119954"), // 15
				hexToByte("16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f"), // 22
				hexToByte("1ce766309c74f07f3dc0839080f518ddcb6500d31fc4e0cf21534bad0785dfc4"), // 25
			}},
	}

	root := hexToByte("197fbc87461398680c858f1daf61e719a1865edd96db34cca3b48c4b43d82e74")
	for desc, test := range tests {
		leaves := []types.Leaf{{Index: test.leafIndex, Hash: test.leaf}}
		proof := merkleMmr.NewProof(merkleMmr.LeafIndexToMMRSize(test.leafCount-1), test.proofs, leaves, hasher.Keccak256Hasher{})
		if !proof.Verify(root) {
			t.Errorf("%s: failed to verify leaf inclusion", desc)
		}
	}

	tests = map[string]struct {
		leafIndex uint64
		leafCount uint64
		leaf      []byte
		proofs    [][]byte
	}{
		"invalid proof missing a left root proof item for leaf index 13 (node 24)": {13, 15,
			hexToByte("883f1aca23002690575957cc85663774bbd3b9549ba5f0ee0fcc8aed9c88cf99"),
			[][]byte{
				hexToByte("16c5d5eb80eec816ca1804cd15705ac2418325b51b57a272e5e7f119e197c31f"),
				hexToByte("94014b81bc56d64cac8dcde8eee47da0ed9b1319dccd9e86ad8d2266d8ef060a"),
				hexToByte("0a73e5a8443de3fcb6f918d786ad6dece6733ec936aa6b1b79beaab19e269d68"),
			}},
	}

	for desc, test := range tests {
		leaves := []types.Leaf{{Index: test.leafIndex, Hash: test.leaf}}
		proof := merkleMmr.NewProof(merkleMmr.LeafIndexToMMRSize(test.leafCount-1), test.proofs, leaves, hasher.Keccak256Hasher{})
		if proof.Verify(root) {
			t.Errorf("%s: verified a leaf inclusion with invalid proofs", desc)
		}
	}
}

func Test1LeafVerify(t *testing.T) {
	fmt.Println("                 1-leaf MMR:           ")
	fmt.Println("                                       ")
	fmt.Println("    Height 1 | 1                       ")
	fmt.Println("             | |                       ")
	fmt.Println("Hash indices | 0                       ")

	// ---------------------------- Tree contents ----------------------------
	//  - For leaf nodes, node hash is the SCALE-encoding of the leaf data.
	//  - For parent nodes, node hash is the hash of it's children (left, right).
	//
	// 0xda5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed // 1  LEAF NODE

	tests := map[string]struct {
		leafIndex uint64
		leafCount uint64
		leaf      []byte
		proofs    [][]byte
	}{
		"leaf index 0 (node 1)": {0, 1,
			hexToByte("da5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed"),
			[][]byte{},
		},
	}

	root := hexToByte("da5e6d0616e05c6a6348605a37ca33493fc1a15ad1e6a405ee05c17843fdafed")
	for desc, test := range tests {
		leaves := []types.Leaf{{Index: test.leafIndex, Hash: test.leaf}}
		proof := merkleMmr.NewProof(merkleMmr.LeafIndexToMMRSize(test.leafCount-1), test.proofs, leaves, hasher.Keccak256Hasher{})
		if !proof.Verify(root) {
			t.Errorf("%s: failed to verify leaf inclusion", desc)
		}
	}
}

func TestFixture7Leaves(t *testing.T) {
	type proof struct {
		leafIndex uint64
		leafCount uint64
		items     [][]byte
	}

	var fixture7Leaves = struct {
		leaves   [][]byte
		rootHash []byte
		proofs   []proof
	}{
		leaves: [][]byte{
			hexToByte("4320435e8c3318562dba60116bdbcc0b82ffcecb9bb39aae3300cfda3ad0b8b0"),
			hexToByte("ad4cbc033833612ccd4626d5f023b9dfc50a35e838514dd1f3c86f8506728705"),
			hexToByte("9ba3bd51dcd2547a0155cf13411beeed4e2b640163bbea02806984f3fcbf822e"),
			hexToByte("1b14c1dc7d3e4def11acdf31be0584f4b85c3673f1ff72a3af467b69a3b0d9d0"),
			hexToByte("3b031d22e24f1126c8f7d2f394b663f9b960ed7abbedb7152e17ce16112656d0"),
			hexToByte("8ed25570209d8f753d02df07c1884ddb36a3d9d4770e4608b188322151c657fe"),
			hexToByte("611c2174c6164952a66d985cfe1ec1a623794393e3acff96b136d198f37a648c"),
		},
		rootHash: hexToByte("e45e25259f7930626431347fa4dd9aae7ac83b4966126d425ca70ab343709d2c"),
		proofs: []proof{
			{
				leafIndex: 0, leafCount: 7,
				items: [][]byte{
					hexToByte("ad4cbc033833612ccd4626d5f023b9dfc50a35e838514dd1f3c86f8506728705"),
					hexToByte("cb24f4614ad5b2a5430344c99545b421d9af83c46fd632d70a332200884b4d46"),
					hexToByte("dca421199bdcc55bb773c6b6967e8d16675de69062b52285ca63685241fdf626"),
				},
			},
			{
				leafIndex: 1, leafCount: 7,
				items: [][]byte{
					hexToByte("4320435e8c3318562dba60116bdbcc0b82ffcecb9bb39aae3300cfda3ad0b8b0"),
					hexToByte("cb24f4614ad5b2a5430344c99545b421d9af83c46fd632d70a332200884b4d46"),
					hexToByte("dca421199bdcc55bb773c6b6967e8d16675de69062b52285ca63685241fdf626"),
				},
			},
			{
				leafIndex: 2, leafCount: 7,
				items: [][]byte{
					hexToByte("1b14c1dc7d3e4def11acdf31be0584f4b85c3673f1ff72a3af467b69a3b0d9d0"),
					hexToByte("672c04a9cd05a644789d769daa552d35d8de7c33129f8a7cbf49e595234c4854"),
					hexToByte("dca421199bdcc55bb773c6b6967e8d16675de69062b52285ca63685241fdf626"),
				},
			},
			{
				leafIndex: 3, leafCount: 7,
				items: [][]byte{
					hexToByte("9ba3bd51dcd2547a0155cf13411beeed4e2b640163bbea02806984f3fcbf822e"),
					hexToByte("672c04a9cd05a644789d769daa552d35d8de7c33129f8a7cbf49e595234c4854"),
					hexToByte("dca421199bdcc55bb773c6b6967e8d16675de69062b52285ca63685241fdf626"),
				},
			},
			{
				leafIndex: 4, leafCount: 7,
				items: [][]byte{
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("8ed25570209d8f753d02df07c1884ddb36a3d9d4770e4608b188322151c657fe"),
					hexToByte("611c2174c6164952a66d985cfe1ec1a623794393e3acff96b136d198f37a648c"),
				},
			},
			{
				leafIndex: 5, leafCount: 7,
				items: [][]byte{
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("3b031d22e24f1126c8f7d2f394b663f9b960ed7abbedb7152e17ce16112656d0"),
					hexToByte("611c2174c6164952a66d985cfe1ec1a623794393e3acff96b136d198f37a648c"),
				},
			},
			{
				leafIndex: 6, leafCount: 7,
				items: [][]byte{
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("7e4316ae2ebf7c3b6821cb3a46ca8b7a4f9351a9b40fcf014bb0a4fd8e8f29da"),
				},
			},
		},
	}

	for i, p := range fixture7Leaves.proofs {
		leaves := []types.Leaf{{Index: p.leafIndex, Hash: fixture7Leaves.leaves[i]}}
		merkleProof := merkleMmr.NewProof(
			merkleMmr.LeafIndexToMMRSize(p.leafCount-1),
			p.items,
			leaves,
			hasher.Keccak256Hasher{},
		)
		if !merkleProof.Verify(fixture7Leaves.rootHash) {
			t.Errorf("failed to verify leaf inclusion for leaf index %v", p.leafIndex)
		}
	}
}

func TestFixture15Leaves(t *testing.T) {
	type proof struct {
		leafIndex uint64
		leafCount uint64
		items     [][]byte
	}

	var fixture15Leaves = struct {
		leaves   [][]byte
		rootHash []byte
		proofs   []proof
	}{
		leaves: [][]byte{
			hexToByte("4320435e8c3318562dba60116bdbcc0b82ffcecb9bb39aae3300cfda3ad0b8b0"),
			hexToByte("ad4cbc033833612ccd4626d5f023b9dfc50a35e838514dd1f3c86f8506728705"),
			hexToByte("9ba3bd51dcd2547a0155cf13411beeed4e2b640163bbea02806984f3fcbf822e"),
			hexToByte("1b14c1dc7d3e4def11acdf31be0584f4b85c3673f1ff72a3af467b69a3b0d9d0"),
			hexToByte("3b031d22e24f1126c8f7d2f394b663f9b960ed7abbedb7152e17ce16112656d0"),
			hexToByte("8ed25570209d8f753d02df07c1884ddb36a3d9d4770e4608b188322151c657fe"),
			hexToByte("611c2174c6164952a66d985cfe1ec1a623794393e3acff96b136d198f37a648c"),
			hexToByte("1e959bd2b05d662f179a714fbf58928730380ad8579a966a9314c8e13b735b13"),
			hexToByte("1c69edb31a1f805991e8e0c27d9c4f5f7fbb047c3313385fd9f4088d60d3d12b"),
			hexToByte("0a4098f56c2e74557cf95f4e9bdc32e7445dd3c7458766c807cd6b54b89e8b38"),
			hexToByte("79501646d325333e636b557abefdfb6fa688012eef0b57bd0b93ef368ff86833"),
			hexToByte("251054c04fcdeca1058dd511274b5eeb22c04b76a3c80f92a989cec535abbd5e"),
			hexToByte("9b2645185bbf36ecfd425c4f99596107d78d160cea01b428be0b079ec8bf2a85"),
			hexToByte("9a9ca4381b27601fe46fe517eb2eedffd8b14d7140cb10fec111337968c0dd28"),
			hexToByte("c43faffd065ac4fc5bc432ad45c13de341b233dcc55afe99ac05eef2fbb8a583"),
		},
		rootHash: hexToByte("3e81e73a77ddf45c0252bba8d1195d1076003d8387df373a46a3a559bc06acca"),
		proofs: []proof{
			{
				leafIndex: 0, leafCount: 15,
				items: [][]byte{
					hexToByte("ad4cbc033833612ccd4626d5f023b9dfc50a35e838514dd1f3c86f8506728705"),
					hexToByte("cb24f4614ad5b2a5430344c99545b421d9af83c46fd632d70a332200884b4d46"),
					hexToByte("441bf63abc7cf9b9e82eb57b8111c883d50ae468d9fd7f301e12269fc0fa1e75"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 1, leafCount: 15,
				items: [][]byte{
					hexToByte("4320435e8c3318562dba60116bdbcc0b82ffcecb9bb39aae3300cfda3ad0b8b0"),
					hexToByte("cb24f4614ad5b2a5430344c99545b421d9af83c46fd632d70a332200884b4d46"),
					hexToByte("441bf63abc7cf9b9e82eb57b8111c883d50ae468d9fd7f301e12269fc0fa1e75"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 2, leafCount: 15,
				items: [][]byte{
					hexToByte("1b14c1dc7d3e4def11acdf31be0584f4b85c3673f1ff72a3af467b69a3b0d9d0"),
					hexToByte("672c04a9cd05a644789d769daa552d35d8de7c33129f8a7cbf49e595234c4854"),
					hexToByte("441bf63abc7cf9b9e82eb57b8111c883d50ae468d9fd7f301e12269fc0fa1e75"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 3, leafCount: 15,
				items: [][]byte{
					hexToByte("9ba3bd51dcd2547a0155cf13411beeed4e2b640163bbea02806984f3fcbf822e"),
					hexToByte("672c04a9cd05a644789d769daa552d35d8de7c33129f8a7cbf49e595234c4854"),
					hexToByte("441bf63abc7cf9b9e82eb57b8111c883d50ae468d9fd7f301e12269fc0fa1e75"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 4, leafCount: 15,
				items: [][]byte{
					hexToByte("8ed25570209d8f753d02df07c1884ddb36a3d9d4770e4608b188322151c657fe"),
					hexToByte("421865424d009fee681cc1e439d9bd4cce0a6f3e79cce0165830515c644d95d4"),
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 5, leafCount: 15,
				items: [][]byte{
					hexToByte("3b031d22e24f1126c8f7d2f394b663f9b960ed7abbedb7152e17ce16112656d0"),
					hexToByte("421865424d009fee681cc1e439d9bd4cce0a6f3e79cce0165830515c644d95d4"),
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 6, leafCount: 15,
				items: [][]byte{
					hexToByte("1e959bd2b05d662f179a714fbf58928730380ad8579a966a9314c8e13b735b13"),
					hexToByte("7e4316ae2ebf7c3b6821cb3a46ca8b7a4f9351a9b40fcf014bb0a4fd8e8f29da"),
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 7, leafCount: 15,
				items: [][]byte{
					hexToByte("611c2174c6164952a66d985cfe1ec1a623794393e3acff96b136d198f37a648c"),
					hexToByte("7e4316ae2ebf7c3b6821cb3a46ca8b7a4f9351a9b40fcf014bb0a4fd8e8f29da"),
					hexToByte("ae88a0825da50e953e7a359c55fe13c8015e48d03d301b8bdfc9193874da9252"),
					hexToByte("de783edd9fe65db4ce28c56687da424218086b4948185bdd9f685a42506e3ba2"),
				},
			},
			{
				leafIndex: 8, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("0a4098f56c2e74557cf95f4e9bdc32e7445dd3c7458766c807cd6b54b89e8b38"),
					hexToByte("7d1f24a6c60769cc6bdc9fc123848d36ef2c6c48e84d9dd464d153cbb0e7ae76"),
					hexToByte("24a44d3d08fbb13a1902e9fa3995456e9a141e0960a2f59725e65a37d474f2c0"),
				},
			},
			{
				leafIndex: 9, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("1c69edb31a1f805991e8e0c27d9c4f5f7fbb047c3313385fd9f4088d60d3d12b"),
					hexToByte("7d1f24a6c60769cc6bdc9fc123848d36ef2c6c48e84d9dd464d153cbb0e7ae76"),
					hexToByte("24a44d3d08fbb13a1902e9fa3995456e9a141e0960a2f59725e65a37d474f2c0"),
				},
			},
			{
				leafIndex: 10, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("251054c04fcdeca1058dd511274b5eeb22c04b76a3c80f92a989cec535abbd5e"),
					hexToByte("2c6280fdcaf131531fe103e0e7353a77440333733c68effa4d3c49413c00b55f"),
					hexToByte("24a44d3d08fbb13a1902e9fa3995456e9a141e0960a2f59725e65a37d474f2c0"),
				},
			},
			{
				leafIndex: 11, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("79501646d325333e636b557abefdfb6fa688012eef0b57bd0b93ef368ff86833"),
					hexToByte("2c6280fdcaf131531fe103e0e7353a77440333733c68effa4d3c49413c00b55f"),
					hexToByte("24a44d3d08fbb13a1902e9fa3995456e9a141e0960a2f59725e65a37d474f2c0"),
				},
			},
			{
				leafIndex: 12, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("f323ac1a7f56de5f40ed8df3e97af74eec0ee9d72883679e49122ffad2ffd03b"),
					hexToByte("9a9ca4381b27601fe46fe517eb2eedffd8b14d7140cb10fec111337968c0dd28"),
					hexToByte("c43faffd065ac4fc5bc432ad45c13de341b233dcc55afe99ac05eef2fbb8a583"),
				},
			},
			{
				leafIndex: 13, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("f323ac1a7f56de5f40ed8df3e97af74eec0ee9d72883679e49122ffad2ffd03b"),
					hexToByte("9b2645185bbf36ecfd425c4f99596107d78d160cea01b428be0b079ec8bf2a85"),
					hexToByte("c43faffd065ac4fc5bc432ad45c13de341b233dcc55afe99ac05eef2fbb8a583"),
				},
			},
			{
				leafIndex: 14, leafCount: 15,
				items: [][]byte{
					hexToByte("73d1bf5a0b1329cd526fba68bb89504258fec5a2282001167fd51c89f7ef73d3"),
					hexToByte("f323ac1a7f56de5f40ed8df3e97af74eec0ee9d72883679e49122ffad2ffd03b"),
					hexToByte("a0d0a78fe68bd0af051c24c6f0ddd219594b582fa3147570b8fd60cf1914efb4"),
				},
			},
		},
	}

	for i, p := range fixture15Leaves.proofs {
		leaves := []types.Leaf{{Index: p.leafIndex, Hash: fixture15Leaves.leaves[i]}}
		merkleProof := merkleMmr.NewProof(
			merkleMmr.LeafIndexToMMRSize(p.leafCount-1),
			p.items,
			leaves,
			hasher.Keccak256Hasher{},
		)
		if !merkleProof.Verify(fixture15Leaves.rootHash) {
			t.Errorf("failed to verify leaf inclusion for leaf index %v", p.leafIndex)
		}
	}
}
