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
	var leavesI [][]byte
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
	fmt.Printf("merkle proofs are:\n")
	for _, v := range proof.Proofs {
		fmt.Printf(" - %v\n", HashToStr(v))
	}
	fmt.Printf("merkle proof indices are %v\n", proof.Leaves)

	// verify merkle proof
	verifyResult, err := proof.VerifyRootHash(root)
	if err != nil {
		panic(err)
	} else if !verifyResult {
		panic("merkle proof verify result is false")
	}
	fmt.Printf("merkle proof verify result is %v\n", verifyResult)

	// build merkle proof for 42 and 20191116 (indices are 1 and 4);
	proof, err = cbmt.BuildMerkleProof(leavesI, []uint32{1, 4})
	if err != nil {
		panic(err)
	}

	fmt.Printf("merkle proofs are:\n")
	for _, v := range proof.Proofs {
		fmt.Printf(" - %v\n", HashToStr(v))
	}
	fmt.Printf("merkle proof indices are %v\n", proof.Leaves)

	// retrieve leaves
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leavesI)
	if err != nil {
		panic(err)
	}

	fmt.Printf("retrieved leaves are:\n")
	for _, v := range retrievedLeaves {
		fmt.Printf(" - %v\n", HashToStr(v))
	}

}
