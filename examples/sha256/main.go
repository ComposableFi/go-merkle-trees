package main

import (
	"fmt"

	"github.com/ComposableFi/merkle-go/merkle"
)

func main() {
	leaves := []string{
		"Hello",
		"Dorood",
		"Hi",
		"Hey",
		"Hola",
	}
	var leavesI []interface{}
	for _, l := range leaves {
		leavesI = append(leavesI, []byte(l))
	}
	cbmt := merkle.CBMT{
		Merge: MergeByteArray{},
	}
	root := cbmt.BuildMerkleRoot(leavesI)
	fmt.Printf("Merkle root is %v \n", HashToStr(root))

	// build merkle proof for 42 (its index is 1);
	proof, err := cbmt.BuildMerkleProof(leavesI, []uint32{1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("merkle proof lemmas are:\n")
	for _, v := range proof.Lemmas {
		fmt.Printf(" - %v\n", HashToStr(v))
	}
	fmt.Printf("merkle proof indices are %v\n", proof.Indices)

	// TODO: make []byte conversion and compare possible
	// verify merkle proof
	// verifyResult, err := proof.Verify(root, []interface{}{"Hi"})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("merkle proof verify result is %v\n", verifyResult)

	// build merkle proof for 42 and 20191116 (indices are 1 and 4);
	proof, err = cbmt.BuildMerkleProof(leavesI, []uint32{1, 4})
	if err != nil {
		panic(err)
	}

	fmt.Printf("merkle proof lemmas are:\n")
	for _, v := range proof.Lemmas {
		fmt.Printf(" - %v\n", HashToStr(v))
	}
	fmt.Printf("merkle proof indices are %v\n", proof.Indices)

	// retrieve leaves
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leavesI)
	if err != nil {
		panic(err)
	}

	fmt.Printf("retrieved leaves are:\n")
	for _, v := range retrievedLeaves {
		fmt.Printf(" - %v\n", HashToStr(v))
	}
	root, err = proof.GetRoot(retrievedLeaves)
	if err != nil {
		panic(err)
	}
	fmt.Printf("calculated root of proof is %v\n", HashToStr(root))

}
