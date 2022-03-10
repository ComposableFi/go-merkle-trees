package mmr

import (
	"math/bits"
	"reflect"

	"github.com/ComposableFi/go-merkle-trees/types"
)

func getPeakPosByHeight(height uint32) uint64 {
	return (1 << (height + 1)) - 2
}

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

func PosHeightInTree(pos uint64) uint32 {
	pos++
	allOnes := func(num uint64) bool { return num != 0 && zerosCount64(num) == bits.LeadingZeros64(num) }
	jumpLeft := func(pos uint64) uint64 {
		var bitLength = uint32(64 - bits.LeadingZeros64(pos))
		var mostSignificantBits uint64 = 1 << (bitLength - 1)
		return pos - (mostSignificantBits - 1)
	}

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

func pop(ph [][]byte) ([]byte, [][]byte) {
	if len(ph) == 0 {
		return nil, ph[:]
	}
	// return the last Items in the slice and the rest of the slice excluding the last Items
	return ph[len(ph)-1], ph[:len(ph)-1]
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func pushLeaf(leaves *[]types.Leaf, l types.Leaf) {
	*leaves = append(*leaves, l)
}
