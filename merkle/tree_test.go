package merkle_test

import (
	"crypto/sha256"
	"testing"

	"github.com/ComposableFi/merkle-go/merkle"
	"github.com/stretchr/testify/require"
)

type Sha256Hasher struct{}

func (hr Sha256Hasher) Hash(b []byte) (merkle.Hash, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (hr Sha256Hasher) ConcatAndHash(left, right []byte) (merkle.Hash, error) {
	return hr.Hash(append(left[:], right[:]...))
}

func TestNewMerkleTree(t *testing.T) {
	merkle.NewTree(Sha256Hasher{})
}

func TestFromLeaves(t *testing.T) {
	aHash, err := Sha256Hasher{}.Hash([]byte("a"))
	require.NoError(t, err)
	bHash, err := Sha256Hasher{}.Hash([]byte("b"))
	require.NoError(t, err)
	cHash, err := Sha256Hasher{}.Hash([]byte("c"))
	require.NoError(t, err)

	leaves := []merkle.Hash{aHash, bHash, cHash}
	mtree := merkle.NewTree(Sha256Hasher{})
	mtree = mtree.FromLeaves(leaves)
	require.Equal(t, []merkle.Hash{}, mtree.UncommittedLeaves)
}
