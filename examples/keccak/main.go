package main

import (
	"fmt"

	"github.com/ComposableFi/go-merkle-trees/hasher"
	"github.com/ComposableFi/go-merkle-trees/merkle"
	"github.com/ComposableFi/go-merkle-trees/mmr"
	"github.com/ComposableFi/go-merkle-trees/types"
)

func main() {
	leaves := []string{
		"Hello",
		"Dorood",
		"Hi",
		"Hey",
		"Hola",
	}
	var leavesI [][]byte
	for _, l := range leaves {
		h, _ := hasher.Keccak256Hasher{}.Hash([]byte(l))
		leavesI = append(leavesI, h)
	}

	merkleTree := merkle.NewTree(hasher.Keccak256Hasher{})
	merkleTree, err := merkleTree.FromLeaves(leavesI)
	if err != nil {
		panic(err)
	}

	root := merkleTree.Root()
	fmt.Printf("Merkle root is %v \n", merkleTree.RootHex())

	// build merkle proof for 42 (its index is 1);
	proof := merkleTree.Proof([]uint64{1})

	fmt.Printf("merkle proof hashes are:\n")
	for _, v := range proof.ProofHashesHex() {
		fmt.Printf(" - %v\n", v)
	}

	// verify merkle proof
	verifyResult, err := proof.Verify(root)
	if err != nil {
		panic(err)
	} else if !verifyResult {
		panic("merkle proof verify result is false")
	}
	fmt.Printf("merkle proof verify result is %v\n", verifyResult)

	var leavesMmr []types.Leaf
	for i := 0; i < len(leavesI); i++ {
		h, _ := hasher.Keccak256Hasher{}.Hash(leavesI[i])
		leavesMmr = append(leavesMmr, types.Leaf{Index: uint64(i), Hash: h})
	}
	mmrTree := mmr.NewMMR(0, mmr.NewMemStore(), leavesMmr, hasher.Keccak256Hasher{})
	var positions []uint64
	for i := 0; i < len(leavesI); i++ {
		pos, err := mmrTree.Push(leavesI[i])
		if err != nil {
			panic(err)
		}
		positions = append(positions, pos)
	}

	mmrRoot, err := mmrTree.Root()
	if err != nil {
		panic(err)
	}
	mmrRootHex, err := mmrTree.RootHex()
	if err != nil {
		panic(err)
	}
	fmt.Printf("MMR root is %v \n", mmrRootHex)

	mmrProof, err := mmrTree.GenProof(func() []uint64 {
		var elem []uint64
		elem = append(elem, positions[uint64(0)])
		return elem
	}())
	if err != nil {
		panic(err)
	}

	mmrTree.Commit()

	verifyResult = mmrProof.Verify(mmrRoot)

	fmt.Printf("MMR verify result is %v\n", verifyResult)
}
