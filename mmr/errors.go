package mmr

import (
	"errors"
)

// ErrCorruptedProof is of the type error. It is returned when proof is considered corrupt
var ErrCorruptedProof = errors.New("corrupted proof: proof items is not enough to build a tree")

// ErrGenProofForInvalidLeaves is of the type error. It is returned when the list of leaves is empty or beyond mmr range
var ErrGenProofForInvalidLeaves = errors.New("leaves is an empty list, or beyond the mmr range")

// ErrInconsistentStore is of the type error. It is returned when the store is considered inconsistent
var ErrInconsistentStore = errors.New("inconsistent store")

// ErrGetRootOnEmpty is of the type error. It is returned when the MMR is empty
var ErrGetRootOnEmpty = errors.New("get root on an empty MMR")
