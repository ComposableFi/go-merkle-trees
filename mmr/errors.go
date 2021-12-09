package mmr

import (
	"errors"
)

// proof items is not enough to build a tree
var ErrCorruptedProof = errors.New("corrupted proof: proof items is not enough to build a tree")
var ErrGenProofForInvalidLeaves = errors.New("leaves is an empty list, or beyond the mmr range")
var ErrInconsistentStore = errors.New("inconsistent store")
