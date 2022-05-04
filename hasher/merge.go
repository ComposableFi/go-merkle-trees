// Package hasher is responsible for hashing and merging the nodes
package hasher

import "github.com/ComposableFi/go-merkle-trees/types"

// MergeAndHash appends two bytes and the uses hasher to hash the appended bytes
func MergeAndHash(hasher types.Hasher, left []byte, right []byte) ([]byte, error) {
	if right == nil {
		return left, nil
	}
	return hasher.Hash(append(left, right...))
}
