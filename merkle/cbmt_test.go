package merkle_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/ComposableFi/merkle-go/merkle"
)

type MergeInt32 struct{}

func (m MergeInt32) Merge(left, right []byte) []byte {
	var r, l int32
	r = b2i(right)
	l = b2i(left)
	merged := r - l
	return i2b(merged)
}

type MergeKeccak256 struct{}

func (m MergeKeccak256) Merge(left, right []byte) []byte {
	h := crypto.Keccak256Hash(
		// []byte(s),
		append(left[:], right[:]...),
	)

	// xx := h.Hex()
	// fmt.Println(xx)

	return h.Bytes()
}

func b2i(b []byte) int32 {
	i := int32(binary.LittleEndian.Uint64(b))
	return i
}

func i2b(i int32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func TestVerifyProof(t *testing.T) {
	var leaves [][]byte
	for _, v := range testAddresses {
		b := common.FromHex(v)
		h := crypto.Keccak256Hash(b).Hex()
		// vv := common.Bytes2Hex(h)
		// require.True(t, common.IsHexAddress(vv))
		leaves = append(leaves, []byte(h))
	}
	cbmt := merkle.CBMT{
		Merge: MergeKeccak256{},
	}

	tree := cbmt.BuildMerkleTree(leaves)
	root := tree.GetRoot()
	xx := crypto.Keccak256Hash(root).Hex()

	// fmt.Println(xx)
	// hexRoot := common.Bytes2Hex(crypto.Keccak256(root))
	require.Equal(t, expectedRoot, xx)

	for i := 0; i < len(testAddresses); i++ {
		leaf := testAddresses[i]
		b := common.FromHex(leaf)
		hashedLeaf := crypto.Keccak256Hash(b)

		hashedLeafHex := hashedLeaf.Hex()
		fmt.Println(hashedLeafHex)
		// proof := tree.getHexProof(hashedLeafHex)

		// const solidityRoot = await this.merkleProof.computeRootFromProofAtPosition(
		// 	hashedLeafHex,
		// 	i,
		// 	leaves.length,
		// 	proof
		// );

		// expect(solidityRoot).to.be.equal(expectedRoot);
	}
}

func TestKeccack(t *testing.T) {
	leaves := []string{"a", "b", "c"}
	var keccackLeaves [][]byte

	for _, x := range leaves {
		keccackLeaves = append(keccackLeaves, crypto.Keccak256([]byte(x)))
	}

	aHash := "3ac225168df54212a25c1c01fd35bebfea408fdac2e31ddd6f80a4bbf9a5f1cb"
	bHash := "b5553de315e0edf504d9150af82dafa5c4667fa618ed0a6f19c69b41166c5510"
	cHash := "0b42b6393c1f53060fe3ddbfcd7aadcca894465a5a438f69c87d790b2299b9b2"
	var hashes []string
	for _, l := range keccackLeaves {
		hashes = append(hashes, common.Bytes2Hex(l))
	}
	require.Equal(t, []string{aHash, bHash, cHash}, hashes)
	cbmt := merkle.CBMT{
		Merge: MergeKeccak256{},
	}
	tree := cbmt.BuildMerkleTree(keccackLeaves)

	root := tree.GetRoot()
	require.Equal(t, "aff1208e69c9e8be9b584b07ebac4e48a1ee9d15ce3afe20b77a4d29e4175aa3", common.Bytes2Hex(root))

}

func TestBuildEmpty(t *testing.T) {
	var leaves [][]byte

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, int(0), int(len(tree.Nodes)))
	require.Equal(t, []byte{0}, tree.GetRoot())
}

func TestBuildOne(t *testing.T) {
	var leaves = [][]byte{i2b(1)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1)}, tree.Nodes)
}

func TestBuildTwo(t *testing.T) {
	var leaves = [][]byte{i2b(1), i2b(2)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1), i2b(1), i2b(2)}, tree.Nodes)
}

func TestBuildFive(t *testing.T) {
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	require.Equal(t, [][]byte{i2b(1), i2b(1), i2b(2), i2b(2), i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}, tree.Nodes)
}

func TestBuildRootDirectly(t *testing.T) {
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, i2b(1), root)
}

func TestBuiltRootIsSameAsTreeRoot(t *testing.T) {
	var leaves [][]byte
	var start int
	var end = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i2b(int32(i)))
	}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, tree.GetRoot(), root)
}

func TestVerifyRetrieveLeaves(t *testing.T) {
	var leaves = [][]byte{i2b(2), i2b(3), i2b(5), i2b(7), i2b(11), i2b(13)}

	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	proof, err := cbmt.BuildMerkleProof(leaves, []uint32{0, 3})
	require.NoError(t, err)
	retrievedLeaves, err := cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)

	require.Equal(t, [][]byte{i2b(2), i2b(7)}, retrievedLeaves)

	retreivedRoot, err := proof.CalculateRootHash()
	require.NoError(t, err)
	mroot := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, mroot, retreivedRoot)

	proof.Leaves = []merkle.LeafData{}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.Error(t, err)
	require.Equal(t, [][]byte{}, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{{Index: 0, Leaf: i2b(4)}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)

	proof.Leaves = []merkle.LeafData{{Index: 0, Leaf: i2b(11)}}
	retrievedLeaves, err = cbmt.RetriveLeaves(proof, leaves)
	require.NoError(t, err)
	require.Nil(t, retrievedLeaves)
}

func TestRebuildProof(t *testing.T) {
	var leaves = [][]byte{i2b(2), i2b(3), i2b(5), i2b(7), i2b(11)}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	tree := cbmt.BuildMerkleTree(leaves)
	root := tree.GetRoot()

	//build proof
	proof, err := tree.BuildProof([]uint32{0, 3})
	require.NoError(t, err)
	lemmas := proof.Proofs
	leafDataList := proof.Leaves

	rebuildProof := merkle.Proof{
		Leaves: leafDataList,
		Proofs: lemmas,
		Merge:  MergeInt32{},
	}

	isValid, err := rebuildProof.VerifyRootHash(root)
	require.NoError(t, err)
	require.Equal(t, true, isValid)

	rebuiltRoot, err := rebuildProof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, root, rebuiltRoot)
}

func TestBuildProof(t *testing.T) {
	var leaves = [][]byte{i2b(3), i2b(5), i2b(7), i2b(11), i2b(13), i2b(17)}
	leafIndecies := []uint32{0, 5}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}

	//build proof
	proof, err := cbmt.BuildMerkleProof(leaves, leafIndecies)
	require.NoError(t, err)
	require.Equal(t, [][]byte{i2b(13), i2b(5), i2b(4)}, proof.Proofs)
	root, err := proof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, i2b(2), root)

	leaves = [][]byte{i2b(2)}
	proof, err = cbmt.BuildMerkleProof(leaves, []uint32{0})
	require.NoError(t, err)
	require.Equal(t, 0, len(proof.Proofs))
	root, err = proof.CalculateRootHash()
	require.NoError(t, err)
	require.Equal(t, i2b(2), root)
}

func TestTreeRootIsTheSameAsProofRoot(t *testing.T) {
	var leaves [][]byte
	var leafIndices []uint32
	var start uint32 = 2
	var end uint32 = 1000
	for i := start; i < end; i++ {
		leaves = append(leaves, i2b(int32(i)))
		leafIndices = append(leafIndices, i-start)
	}
	cbmt := merkle.CBMT{
		Merge: MergeInt32{},
	}
	proof, err := cbmt.BuildMerkleProof(leaves, leafIndices)
	require.NoError(t, err)

	proofRoot, err := proof.CalculateRootHash()
	require.NoError(t, err)

	treeRoot := cbmt.BuildMerkleRoot(leaves)
	require.Equal(t, proofRoot, treeRoot)
}

var testAddresses = []string{
	"0x9aF1Ca5941148eB6A3e9b9C741b69738292C533f",
	"0xDD6ca953fddA25c496165D9040F7F77f75B75002",
	"0x60e9C47B64Bc1C7C906E891255EaEC19123E7F42",
	"0xfa4859480Aa6D899858DE54334d2911E01C070df",
	"0x19B9b128470584F7209eEf65B69F3624549Abe6d",
	"0xC436aC1f261802C4494504A11fc2926C726cB83b",
	"0xc304C8C2c12522F78aD1E28dD86b9947D7744bd0",
	"0xDa0C2Cba6e832E55dE89cF4033affc90CC147352",
	"0xf850Fd22c96e3501Aad4CDCBf38E4AEC95622411",
	"0x684918D4387CEb5E7eda969042f036E226E50642",
	"0x963F0A1bFbb6813C0AC88FcDe6ceB96EA634A595",
	"0x39B38ad74b8bCc5CE564f7a27Ac19037A95B6099",
	"0xC2Dec7Fdd1fef3ee95aD88EC8F3Cd5bd4065f3C7",
	"0x9E311f05c2b6A43C2CCF16fB2209491BaBc2ec01",
	"0x927607C30eCE4Ef274e250d0bf414d4a210b16f0",
	"0x98882bcf85E1E2DFF780D0eB360678C1cf443266",
	"0xFBb50191cd0662049E7C4EE32830a4Cc9B353047",
	"0x963854fc2C358c48C3F9F0A598B9572c581B8DEF",
	"0xF9D7Bc222cF6e3e07bF66711e6f409E51aB75292",
	"0xF2E3fd32D063F8bBAcB9e6Ea8101C2edd899AFe6",
	"0x407a5b9047B76E8668570120A96d580589fd1325",
	"0xEAD9726FAFB900A07dAd24a43AE941d2eFDD6E97",
	"0x42f5C8D9384034A9030313B51125C32a526b6ee8",
	"0x158fD2529Bc4116570Eb7C80CC76FEf33ad5eD95",
	"0x0A436EE2E4dEF3383Cf4546d4278326Ccc82514E",
	"0x34229A215db8FeaC93Caf8B5B255e3c6eA51d855",
	"0xEb3B7CF8B1840242CB98A732BA464a17D00b5dDF",
	"0x2079692bf9ab2d6dc7D79BBDdEE71611E9aA3B72",
	"0x46e2A67e5d450e2Cf7317779f8274a2a630f3C9B",
	"0xA7Ece4A5390DAB18D08201aE18800375caD78aab",
	"0x15E1c0D24D62057Bf082Cb2253dA11Ef0d469570",
	"0xADDEF4C9b5687Eb1F7E55F2251916200A3598878",
	"0xe0B16Fb96F936035db2b5A68EB37D470fED2f013",
	"0x0c9A84993feaa779ae21E39F9793d09e6b69B62D",
	"0x3bc4D5148906F70F0A7D1e2756572655fd8b7B34",
	"0xFf4675C26903D5319795cbd3a44b109E7DDD9fDe",
	"0xCec4450569A8945C6D2Aba0045e4339030128a92",
	"0x85f0584B10950E421A32F471635b424063FD8405",
	"0xb38bEe7Bdc0bC43c096e206EFdFEad63869929E3",
	"0xc9609466274Fef19D0e58E1Ee3b321D5C141067E",
	"0xa08EA868cF75268E7401021E9f945BAe73872ecc",
	"0x67C9Cb1A29E964Fe87Ff669735cf7eb87f6868fE",
	"0x1B6BEF636aFcdd6085cD4455BbcC93796A12F6E2",
	"0x46B37b243E09540b55cF91C333188e7D5FD786dD",
	"0x8E719E272f62Fa97da93CF9C941F5e53AA09e44a",
	"0xa511B7E7DB9cb24AD5c89fBb6032C7a9c2EfA0a5",
	"0x4D11FDcAeD335d839132AD450B02af974A3A66f8",
	"0xB8cf790a5090E709B4619E1F335317114294E17E",
	"0x7f0f57eA064A83210Cafd3a536866ffD2C5eDCB3",
	"0xC03C848A4521356EF800e399D889e9c2A25D1f9E",
	"0xC6b03DF05cb686D933DD31fCa5A993bF823dc4FE",
	"0x58611696b6a8102cf95A32c25612E4cEF32b910F",
	"0x2ed4bC7197AEF13560F6771D930Bf907772DE3CE",
	"0x3C5E58f334306be029B0e47e119b8977B2639eb4",
	"0x288646a1a4FeeC560B349d210263c609aDF649a6",
	"0xb4F4981E0d027Dc2B3c86afA0D0fC03d317e83C0",
	"0xaAE4A87F8058feDA3971f9DEd639Ec9189aA2500",
	"0x355069DA35E598913d8736E5B8340527099960b8",
	"0x3cf5A0F274cd243C0A186d9fCBdADad089821B93",
	"0xca55155dCc4591538A8A0ca322a56EB0E4aD03C4",
	"0xE824D0268366ec5C4F23652b8eD70D552B1F2b8B",
	"0x84C3e9B25AE8a9b39FF5E331F9A597F2DCf27Ca9",
	"0xcA0018e278751De10d26539915d9c7E7503432FE",
	"0xf13077dE6191D6c1509ac7E088b8BE7Fe656c28b",
	"0x7a6bcA1ec9Db506e47ac6FD86D001c2aBc59C531",
	"0xeA7f9A2A9dd6Ba9bc93ca615C3Ddf26973146911",
	"0x8D0d8577e16F8731d4F8712BAbFa97aF4c453458",
	"0xB7a7855629dF104246997e9ACa0E6510df75d0ea",
	"0x5C1009BDC70b0C8Ab2e5a53931672ab448C17c89",
	"0x40B47D1AfefEF5eF41e0789F0285DE7b1C31631C",
	"0x5086933d549cEcEB20652CE00973703CF10Da373",
	"0xeb364f6FE356882F92ae9314fa96116Cf65F47d8",
	"0xdC4D31516A416cEf533C01a92D9a04bbdb85EE67",
	"0x9b36E086E5A274332AFd3D8509e12ca5F6af918d",
	"0xBC26394fF36e1673aE0608ce91A53B9768aD0D76",
	"0x81B5AB400be9e563fA476c100BE898C09966426c",
	"0x9d93C8ae5793054D28278A5DE6d4653EC79e90FE",
	"0x3B8E75804F71e121008991E3177fc942b6c28F50",
	"0xC6Eb5886eB43dD473f5BB4e21e56E08dA464D9B4",
	"0xfdf1277b71A73c813cD0e1a94B800f4B1Db66DBE",
	"0xc2ff2cCc98971556670e287Ff0CC39DA795231ad",
	"0x76b7E1473f0D0A87E9B4a14E2B179266802740f5",
	"0xA7Bc965660a6EF4687CCa4F69A97563163A3C2Ef",
	"0xB9C2b47888B9F8f7D03dC1de83F3F55E738CebD3",
	"0xEd400162E6Dd6bD2271728FFb04176bF770De94a",
	"0xE3E8331156700339142189B6E555DCb2c0962750",
	"0xbf62e342Bc7706a448EdD52AE871d9C4497A53b1",
	"0xb9d7A1A111eed75714a0AcD2dd467E872eE6B03D",
	"0x03942919DFD0383b8c574AB8A701d89fd4bfA69D",
	"0x0Ef4C92355D3c8c7050DFeb319790EFCcBE6fe9e",
	"0xA6895a3cf0C60212a73B3891948ACEcF1753f25E",
	"0x0Ed509239DB59ef3503ded3d31013C983d52803A",
	"0xc4CE8abD123BfAFc4deFf37c7D11DeCd5c350EE4",
	"0x4A4Bf59f7038eDcd8597004f35d7Ee24a7Bdd2d3",
	"0x5769E8e8A2656b5ed6b6e6fa2a2bFAeaf970BB87",
	"0xf9E15cCE181332F4F57386687c1776b66C377060",
	"0xc98f8d4843D56a46C21171900d3eE538Cc74dbb5",
	"0x3605965B47544Ce4302b988788B8195601AE4dEd",
	"0xe993BDfdcAac2e65018efeE0F69A12678031c71d",
	"0x274fDf8801385D3FAc954BCc1446Af45f5a8304c",
	"0xBFb3f476fcD6429F4a475bA23cEFdDdd85c6b964",
	"0x806cD16588Fe812ae740e931f95A289aFb4a4B50",
	"0xa89488CE3bD9C25C3aF797D1bbE6CA689De79d81",
	"0xd412f1AfAcf0Ebf3Cd324593A231Fc74CC488B12",
	"0xd1f715b2D7951d54bc31210BbD41852D9BF98Ed1",
	"0xf65aD707c344171F467b2ADba3d14f312219cE23",
	"0x2971a4b242e9566dEF7bcdB7347f5E484E11919B",
	"0x12b113D6827E07E7D426649fBd605f427da52314",
	"0x1c6CA45171CDb9856A6C9Dba9c5F1216913C1e97",
	"0x11cC6ee1d74963Db23294FCE1E3e0A0555779CeA",
	"0x8Aa1C721255CDC8F895E4E4c782D86726b068667",
	"0xA2cDC1f37510814485129aC6310b22dF04e9Bbf0",
	"0xCf531b71d388EB3f5889F1f78E0d77f6fb109767",
	"0xBe703e3545B2510979A0cb0C440C0Fba55c6dCB5",
	"0x30a35886F989db39c797D8C93880180Fdd71b0c8",
	"0x1071370D981F60c47A9Cd27ac0A61873a372cBB2",
	"0x3515d74A11e0Cb65F0F46cB70ecf91dD1712daaa",
	"0x50500a3c2b7b1229c6884505D00ac6Be29Aecd0C",
	"0x9A223c2a11D4FD3585103B21B161a2B771aDA3d1",
	"0xd7218df03AD0907e6c08E707B15d9BD14285e657",
	"0x76CfD72eF5f93D1a44aD1F80856797fBE060c70a",
	"0x44d093cB745944991EFF5cBa151AA6602d6f5420",
	"0x626516DfF43bf09A71eb6fd1510E124F96ED0Cde",
	"0x6530824632dfe099304E2DC5701cA99E6d031E08",
	"0x57e6c423d6a7607160d6379A0c335025A14DaFC0",
	"0x3966D4AD461Ef150E0B10163C81E79b9029E69c3",
	"0xF608aCfd0C286E23721a3c347b2b65039f6690F1",
	"0xbfB8FAac31A25646681936977837f7740fCd0072",
	"0xd80aa634a623a7ED1F069a1a3A28a173061705c7",
	"0x9122a77B36363e24e12E1E2D73F87b32926D3dF5",
	"0x62562f0d1cD31315bCCf176049B6279B2bfc39C2",
	"0x48aBF7A2a7119e5675059E27a7082ba7F38498b2",
	"0xb4596983AB9A9166b29517acD634415807569e5F",
	"0x52519D16E20BC8f5E96Da6d736963e85b2adA118",
	"0x7663893C3dC0850EfC5391f5E5887eD723e51B83",
	"0x5FF323a29bCC3B5b4B107e177EccEF4272959e61",
	"0xee6e499AdDf4364D75c05D50d9344e9daA5A9AdF",
	"0x1631b0BD31fF904aD67dD58994C6C2051CDe4E75",
	"0xbc208e9723D44B9811C428f6A55722a26204eEF2",
	"0xe76103a222Ee2C7Cf05B580858CEe625C4dc00E1",
	"0xC71Bb2DBC51760f4fc2D46D84464410760971B8a",
	"0xB4C18811e6BFe564D69E12c224FFc57351f7a7ff",
	"0xD11DB0F5b41061A887cB7eE9c8711438844C298A",
	"0xB931269934A3D4432c084bAAc3d0de8143199F4f",
	"0x070037cc85C761946ec43ea2b8A2d5729908A2a1",
	"0x2E34aa8C95Ffdbb37f14dCfBcA69291c55Ba48DE",
	"0x052D93e8d9220787c31d6D83f87eC7dB088E998f",
	"0x498dAC6C69b8b9ad645217050054840f1D91D029",
	"0xE4F7D60f9d84301e1fFFd01385a585F3A11F8E89",
	"0xEa637992f30eA06460732EDCBaCDa89355c2a107",
	"0x4960d8Da07c27CB6Be48a79B96dD70657c57a6bF",
	"0x7e471A003C8C9fdc8789Ded9C3dbe371d8aa0329",
	"0xd24265Cc10eecb9e8d355CCc0dE4b11C556E74D7",
	"0xDE59C8f7557Af779674f41CA2cA855d571018690",
	"0x2fA8A6b3b6226d8efC9d8f6EBDc73Ca33DDcA4d8",
	"0xe44102664c6c2024673Ff07DFe66E187Db77c65f",
	"0x94E3f4f90a5f7CBF2cc2623e66B8583248F01022",
	"0x0383EdBbc21D73DEd039E9C1Ff6bf56017b4CC40",
	"0x64C3E49898B88d1E0f0d02DA23E0c00A2Cd0cA99",
	"0xF4ccfB67b938d82B70bAb20975acFAe402E812E1",
	"0x4f9ee5829e9852E32E7BC154D02c91D8E203e074",
	"0xb006312eF9713463bB33D22De60444Ba95609f6B",
	"0x7Cbe76ef69B52110DDb2e3b441C04dDb11D63248",
	"0x70ADEEa65488F439392B869b1Df7241EF317e221",
	"0x64C0bf8AA36Ba590477585Bc0D2BDa7970769463",
	"0xA4cDc98593CE52d01Fe5Ca47CB3dA5320e0D7592",
	"0xc26B34D375533fFc4c5276282Fa5D660F3d8cbcB",
}

const expectedRoot = "0x72b0acd7c302a84f1f6b6cefe0ba7194b7398afb440e1b44a9dbbe270394ca53"
