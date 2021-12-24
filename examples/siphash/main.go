package main

import (
	"fmt"

	"github.com/ComposableFi/merkle-go/merkle"
)

func main() {
	leaves := []uint64{
		3584654056691428718,
		42,
		11643453954163878810,
		11177097603989645559,
		20191116,
		10289152030157698709,
	}
	var leavesI [][]byte
	for _, l := range leaves {
		leavesI = append(leavesI, i2b(l))
	}
	cbmt := merkle.CBMT{
		Merge: MergeUint64{},
	}
	root := cbmt.BuildMerkleRoot(leavesI)
	fmt.Printf("Merkle root is %v \n", root)

	// build merkle proof for 42 (its index is 1);
	proof, err := cbmt.BuildMerkleProof(leavesI, []uint32{1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("merkle proof lemmas are %v, indices are %v\n", proof.Lemmas, proof.Leaves)

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
	fmt.Printf("merkle proof lemmas are %v, indices are %v\n", proof.Lemmas, proof.Leaves)

	// retrieve leaves
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leavesI)
	if err != nil {
		panic(err)
	}
	fmt.Printf("retrieved leaves are %v\n", retrievedLeaves)

}
