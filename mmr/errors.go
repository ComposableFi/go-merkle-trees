package mmr

import (
	"errors"
)

// proof items is not enough to build a tree
var ErrCorruptedProof = errors.New("corrupted proof: proof items is not enough to build a tree")
