package mmr

import (
	"github.com/ComposableFi/go-merkle-trees/types"
	"math/bits"
)

func getPeakPosByHeight(height uint32) uint64 {
	return (1 << (height + 1)) - 2
}

// leftPeakHeightPos derives and returns the height and position of the leftmost peak of an MMR from
// the MMR size.
func leftPeakHeightPos(mmrSize uint64) (uint32, uint64) {
	var prevPos uint64
	var height uint32 = 1
	var pos = getPeakPosByHeight(height)
	for pos < mmrSize {
		height++
		prevPos = pos
		pos = getPeakPosByHeight(height)
	}
	return height - 1, prevPos
}

func siblingOffset(height uint32) uint64 {
	return (2 << height) - 1
}

func parentOffset(height uint32) uint64 {
	return 2 << height
}

func getRightPeak(height uint32, pos, mmrSize uint64) *peak {
	// move to right sibling Pos
	pos += siblingOffset(height)
	// loop until we find a Pos in mmr
	for pos > mmrSize-1 {
		if height == 0 {
			return nil
		}
		// move to left child
		pos -= parentOffset(height - 1)
		height--
	}
	return &peak{height, pos}
}

// GetPeaks returns the positions of the peaks of the MMR using the MMR size.
// 1. It starts by finding the leftmost peak.
// 2. It then finds the next peak (right peak) by moving to the right sibling. If that node isn't in the MMR (which it won't),
//    it take its left child. If that child is not in the MMR either, it keeps taking its left child until it finds a node
//    that exists in the MMR.
// 3. The process is repeated until it is at the last node.
func GetPeaks(mmrSize uint64) (positions []uint64) {
	var height, pos = leftPeakHeightPos(mmrSize)
	positions = append(positions, pos)

	for height > 0 {
		p := getRightPeak(height, pos, mmrSize)
		if p == nil {
			break
		}
		height = p.height
		pos = p.pos
		positions = append(positions, pos)
	}

	return positions
}

// PosHeightInTree calculates and returns the height of a node in the tree using its position.
func PosHeightInTree(pos uint64) uint32 {
	// increase position by 1 since this algorithm starts the node position with 0
	pos++
	allOnes := func(num uint64) bool { return num != 0 && zerosCount64(num) == bits.LeadingZeros64(num) }
	jumpLeft := func(pos uint64) uint64 {
		var bitLength = uint32(64 - bits.LeadingZeros64(pos))
		var mostSignificantBits uint64 = 1 << (bitLength - 1)
		return pos - (mostSignificantBits - 1)
	}

	// in merkle mountain ranges, the leftmost nodes usually have a position (in binary) of all ones.
	// keep jumping to the left until the leftmost node is obtained.
	for !allOnes(pos) {
		pos = jumpLeft(pos)
	}

	return uint32(64 - bits.LeadingZeros64(pos) - 1)
}

func zerosCount64(num uint64) int {
	return 64 - bits.OnesCount64(num)
}

// LeafIndexToPos returns the position of a leaf from its index.
func LeafIndexToPos(index uint64) uint64 {
	// mmr_size - H - 1, H is the height(intervals) of last peak
	return LeafIndexToMMRSize(index) - uint64(bits.TrailingZeros64(index+1)) - 1
}

// LeafIndexToMMRSize returns the mmr size of an mmr tree provided the leaf index passed as an argument is the last leaf in
// the tree.
func LeafIndexToMMRSize(index uint64) uint64 {
	// leaf index start with 0
	var leavesCount = index + 1

	// the peak count(k) is actually the count of 1 in leaves count's binary representation
	var peakCount = bits.OnesCount64(leavesCount)
	return 2*leavesCount - uint64(peakCount)
}

// pop removes the last item from a slice and returns it
func pop(slice *[][]byte) []byte {
	var sliceCopy = *slice
	*slice = sliceCopy[:len(sliceCopy)-1]
	return sliceCopy[len(sliceCopy)-1]
}

func pushLeaf(leaves *[]types.Leaf, l types.Leaf) {
	*leaves = append(*leaves, l)
}
