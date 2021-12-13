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
	var leavesI []interface{}
	for _, l := range leaves {
		leavesI = append(leavesI, l)
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
	fmt.Printf("merkle proof lemmas are %v, indices are %v\n", proof.Lemmas, proof.Indices)

	// verify merkle proof
	verifyResult, err := proof.Verify(root, []interface{}{uint64(42)})
	if err != nil {
		panic(err)
	}
	fmt.Printf("merkle proof verify result is %v\n", verifyResult)

	// build merkle proof for 42 and 20191116 (indices are 1 and 4);
	proof, err = cbmt.BuildMerkleProof(leavesI, []uint32{1, 4})
	if err != nil {
		panic(err)
	}
	fmt.Printf("merkle proof lemmas are %v, indices are %v\n", proof.Lemmas, proof.Indices)

	// retrieve leaves
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leavesI)
	if err != nil {
		panic(err)
	}
	fmt.Printf("retrieved leaves are %v\n", retrievedLeaves)
	root, err = proof.GetRoot(retrievedLeaves)
	if err != nil {
		panic(err)
	}
	fmt.Printf("calculated root of proof is %v\n", root)

}
