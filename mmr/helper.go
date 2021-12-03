package mmr

import (
	"fmt"
	"strconv"
)

func getPeakPosByHeight(height uint32) uint64 {
	return (1 << (height + 1)) - 2
}

func leftPeakHeightPos(mmrSize uint64) (uint32, uint64) {
	var height uint32 = 1
	var prevPos uint64 = 0
	var pos = getPeakPosByHeight(height)
	for pos < mmrSize {
		height += 1
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
	// move to right sibling pos
	pos += siblingOffset(height)
	// loop until we find a pos in mmr
	for pos > mmrSize-1 {
		if height == 0 {
			return nil
		}
		// move to left child
		pos -= parentOffset(height - 1)
		height -= 1
	}
	return &peak{height, pos}
}

func getPeaks(mmrSize uint64) (pos_s []uint64) {
	var height, pos = leftPeakHeightPos(mmrSize)
	pos_s = append(pos_s, pos)

	for height > 0 {
		p := getRightPeak(height, pos, mmrSize)
		if p == nil {
			break
		}
		height = p.height
		pos = p.pos
		pos_s = append(pos_s, pos)
	}

	return pos_s
}

func posHeightInTree(pos uint64) uint32 {
	pos += 1
	allOnes := func(num uint64) bool { return num != 0 && countZeros(num) == countLeadingZeros(num) }
	jumpLeft := func(pos uint64) uint64 {
		var bitLength = uint32(64 - countLeadingZeros(pos))
		var mostSignificantBits uint64 = 1 << (bitLength - 1)
		return pos - (mostSignificantBits - 1)
	}

	for !allOnes(pos) {
		fmt.Printf("pos %v \n", pos)
		pos = jumpLeft(pos)
	}

	return uint32(64 - countLeadingZeros(pos) - 1)
}

func countLeadingZeros(num uint64) int {
	str := strconv.FormatInt(int64(num), 2)
	counter := 0
	for _, value := range str {
		if value == '0' {
			counter++
		} else {
			break
		}
	}
	return counter
}

func countZeros(num uint64) int {
	str := strconv.FormatInt(int64(num), 2)
	counter := 0
	for _, value := range str {
		if value == '0' {
			counter++
		}
	}
	return counter
}

func pop(ph []interface{}) (interface{}, []interface{}) {
	if len(ph) == 0 {
		return nil, ph[:]
	}
	// return the last item in the slice and the rest of the slice excluding the last item
	return ph[len(ph)-1], ph[:len(ph)-1]
}
