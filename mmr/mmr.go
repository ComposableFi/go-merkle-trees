package mmr

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ComposableFi/merkle-go/merkle"
)

type MMR struct {
	size  uint64
	batch *Batch
	merge merkle.Merge
}

func NewMMR(mmrSize uint64, s Store, m merkle.Merge) *MMR {
	return &MMR{
		size:  mmrSize,
		batch: NewBatch(s),
		merge: m,
	}
}

func (m *MMR) findElem(pos uint64, hashes []interface{}) (interface{}, error) {
	checkSub := func(left, right uint64) (bool, uint64) {
		if left >= right {
			return true, left - right
		}
		return false, 0
	}

	notOverflow, posOffset := checkSub(pos, m.size)
	if notOverflow && uint64(len(hashes)) > posOffset {
		return hashes[posOffset], nil
	}

	elem := m.batch.getElem(pos)
	if elem == nil {
		return nil, ErrInconsistentStore
	}

	return elem, nil
}

func (m *MMR) MMRSize() uint64 {
	return m.size
}

func (m *MMR) IsEmpty() bool {
	return m.size == 0
}

// push a element and return position
func (m *MMR) Push(elem interface{}) (interface{}, error) {
	var elems []interface{}
	// position of new elems
	elemPos := m.size
	elems = append(elems, elem)

	var height uint32 = 0
	var pos = elemPos
	// continue to merge tree node if next pos higher than current
	for posHeightInTree(pos+1) > height {
		pos += 1
		leftPos := pos - parentOffset(height)
		rightPos := leftPos + siblingOffset(height)
		leftElem, err := m.findElem(leftPos, elems)
		if err != nil {
			return nil, err
		}

		rightElem, err := m.findElem(rightPos, elems)
		if err != nil {
			return nil, err
		}

		parentElem := m.merge.Merge(leftElem, rightElem)
		elems = append(elems, parentElem)
		height += 1
	}
	// store hashes
	m.batch.append(elemPos, elems)
	// update mmrSize
	m.size = pos + 1
	return elemPos, nil
}

func (m *MMR) GetRoot() (interface{}, error) {
	if m.size == 0 {
		return nil, ErrGetRootOnEmpty
	} else if m.size == 1 {
		e := m.batch.getElem(0)
		if e == nil {
			return nil, ErrInconsistentStore
		}
		return e, nil
	}

	var peaks []interface{}
	for _, peakPos := range getPeaks(m.size) {
		elem := m.batch.getElem(peakPos)
		if elem == nil {
			return nil, ErrInconsistentStore
		}
		peaks = append(peaks, elem)
	}

	var p interface{}
	if p, peaks = m.bagRHSPeaks(peaks); p == nil {
		return nil, ErrInconsistentStore
	}

	return p, nil
}

func (m *MMR) bagRHSPeaks(rhsPeaks []interface{}) (interface{}, []interface{}) {
	for len(rhsPeaks) > 1 {
		var rp, lp interface{}
		if rp, rhsPeaks = pop(rhsPeaks); rp == nil {
			panic("pop")
		}

		if lp, rhsPeaks = pop(rhsPeaks); lp == nil {
			panic("pop")
		}
		rhsPeaks = append(rhsPeaks, m.merge.Merge(rp, lp))
	}

	if len(rhsPeaks) > 0 {
		return rhsPeaks[len(rhsPeaks)-1], rhsPeaks
	}
	return nil, rhsPeaks[:]
}

/// generate merkle proof for a peak
/// the pos_list must be sorted, otherwise the behaviour is undefined
///
/// 1. find a lower tree in peak that can generate a complete merkle proof for position
/// 2. find that tree by compare positions
/// 3. generate proof for each positions
func (m *MMR) genProofForPeak(proof *Iterator, posList []uint64, peakPos uint64) error {
	if len(posList) == 1 && reflect.DeepEqual(posList, []uint64{peakPos}) {
		return nil
	}
	// take peak root from store if no positions need to be proof
	if len(posList) == 0 {
		elem := m.batch.getElem(peakPos)
		if elem == nil {
			return ErrInconsistentStore
		}
		proof.push(elem)
		return nil
	}

	var queue []peak
	for _, p := range posList {
		queue = append(queue, peak{pos: p, height: 0})
	}

	for len(queue) > 0 {
		pos, height := queue[0].pos, queue[0].height
		// pop front
		queue = queue[1:]
		if pos <= peakPos {
			panic("pos is less or equal to peak position")
		}

		if pos == peakPos {
			break
		}

		// calculate sibling
		sibPos, parentPos := func() (uint64, uint64) {
			var nextHeight = posHeightInTree(pos + 1)
			var siblingOffset = siblingOffset(height)
			if nextHeight > height {
				return pos - siblingOffset, pos + 1
			} else {
				return pos + siblingOffset, pos + parentOffset(height)
			}
		}()

		if len(queue) > 0 && sibPos == queue[0].pos {
			// drop sibling
			queue = queue[1:]
		} else {
			p := m.batch.getElem(sibPos)
			if p == nil {
				return ErrCorruptedProof
			}
			proof.push(p)
		}
		if parentPos < peakPos {
			queue = append(queue, peak{height + 1, parentPos})
		}
	}
	return nil
}

/// Generate merkle proof for positions
/// 1. sort positions
/// 2. push merkle proof to proof by peak from left to right
/// 3. push bagged right hand side root
func (m *MMR) GenProof(posList []uint64) (*MerkleProof, error) {
	if len(posList) == 0 {
		return nil, ErrGenProofForInvalidLeaves
	}
	if m.size == 1 && reflect.DeepEqual(posList, []uint64{0}) {
		fmt.Printf("returning empty proof \n")
		return NewMerkleProof(m.size, NewIterator(), m.merge), nil
	}

	sort.Slice(posList, func(i, j int) bool {
		return posList[i] < posList[j]
	})
	var peaks = getPeaks(m.size)
	var proof = NewIterator()
	// generate merkle proof for each peaks
	var baggingTrack uint = 0
	for _, peakPos := range peaks {
		var pl []uint64
		pl = takeWhileVecUint64(&posList, func(u uint64) bool {
			return u <= peakPos
		})
		if len(pl) == 0 {
			baggingTrack += 1
		} else {
			baggingTrack = 0
		}
		var err error
		err = m.genProofForPeak(proof, pl, peakPos)
		if err != nil {
			return nil, err
		}
	}

	// ensure there are no remaining positions
	if len(posList) != 0 {
		return nil, ErrGenProofForInvalidLeaves
	}

	if baggingTrack > 1 {
		var rhsPeaks = proof.splitOff(proof.length() - int(baggingTrack))

		var p interface{}
		p, rhsPeaks = m.bagRHSPeaks(rhsPeaks)
		if p != nil {
			// TODO: handle error properly
			panic("bagging rhs peaks")
		}
		proof.push(p)
	}

	return NewMerkleProof(m.size, proof, m.merge), nil
}

func (m *MMR) Commit() interface{} {
	return m.batch.commit()
}

type MerkleProof struct {
	mmrSize uint64
	proof   *Iterator
	Merge   merkle.Merge
}

func NewMerkleProof(mmrSize uint64, proof *Iterator, m merkle.Merge) *MerkleProof {
	return &MerkleProof{
		mmrSize: mmrSize,
		proof:   proof,
		Merge:   m,
	}
}

func (m *MerkleProof) MMRSize() uint64 {
	return m.mmrSize
}

func (m *MerkleProof) ProofItems() []interface{} {
	return m.proof.Items
}

func (m *MerkleProof) calculatePeakRoot(leaves []Leaf, peakPos uint64, proofs *Iterator) (interface{}, error) {
	if len(leaves) == 0 {
		panic("can't be empty")
	}

	// (position, hash, height)
	var queue []leafWithHash
	for _, l := range leaves {
		queue = append(queue, leafWithHash{l.pos, l.hash, 0})
	}

	// calculate tree root from each items
	for len(queue) > 0 {
		pop := queue[0]
		// pop from front
		queue = queue[1:]

		pos, item, height := pop.pos, pop.hash, pop.height
		if pos == peakPos {
			return item, nil
		}
		// calculate sibling
		var nextHeight = posHeightInTree(pos + 1)
		var sibPos, parentPos = func() (uint64, uint64) {
			var siblingOffset = siblingOffset(height)
			if nextHeight > height {
				// implies pos is right sibling
				return pos - siblingOffset, pos + 1
			} else {
				// pos is left sibling
				return pos + siblingOffset, pos + parentOffset(height)
			}
		}()

		var siblingItem interface{}
		if len(queue) > 0 && queue[0].pos == sibPos {
			siblingItem, queue = queue[0].hash, queue[1:]
		} else {
			if siblingItem = proofs.next(); siblingItem == nil {
				return nil, ErrCorruptedProof
			}
		}

		var parentItem interface{}
		if nextHeight > height {
			parentItem = m.Merge.Merge(siblingItem, item)
		} else {
			parentItem = m.Merge.Merge(item, siblingItem)
		}

		if parentPos < peakPos {
			queue = append(queue, leafWithHash{parentPos, parentItem, height + 1})
		} else {
			return parentItem, nil
		}
	}

	return nil, ErrCorruptedProof
}

func (m *MerkleProof) baggingPeaksHashes(peaksHashes []interface{}) (interface{}, error) {
	var rightPeak, leftPeak interface{}
	for len(peaksHashes) > 1 {
		if rightPeak, peaksHashes = pop(peaksHashes); rightPeak == nil {
			panic("pop")
		}

		if leftPeak, peaksHashes = pop(peaksHashes); leftPeak == nil {
			panic("pop")
		}
		peaksHashes = append(peaksHashes, m.Merge.Merge(rightPeak, leftPeak))
	}

	if len(peaksHashes) == 0 {
		fmt.Printf("length of peaksHashes is 0 \n")
		return nil, ErrCorruptedProof
	}
	return peaksHashes[len(peaksHashes)-1], nil
}

/// merkle proof
/// 1. sort items by position
/// 2. calculate root of each peak
/// 3. bagging peaks
func (m *MerkleProof) CalculateRoot(leaves []Leaf, mmrSize uint64, proofs *Iterator) (interface{}, error) {
	var peaksHashes, err = m.calculatePeaksHashes(leaves, mmrSize, proofs)
	if err != nil {
		return nil, err
	}

	return m.baggingPeaksHashes(peaksHashes)
}

/// from merkle proof of leaf n to calculate merkle root of n + 1 leaves.
/// by observe the MMR construction graph we know it is possible.
/// https://github.com/jjyr/merkle-mountain-range#construct
/// this is kinda tricky, but it works, and useful
func (m *MerkleProof) CalculateRootWithNewLeaf(leaves []Leaf, newPos uint64, newElem interface{}, newMMRSize uint64) (interface{}, error) {
	posHeight := posHeightInTree(newPos)
	nextHeight := posHeightInTree(newPos + 1)
	if nextHeight > posHeight {
		peaksHashes, err := m.calculatePeaksHashes(leaves, m.mmrSize, m.proof)
		if err != nil {
			return nil, err
		}
		peaksPos := getPeaks(newMMRSize)
		// reverse touched peaks
		var i uint = 0
		for peaksPos[i] < newPos {
			i += 1
		}

		reversePeakHashes := peaksHashes[i:]
		reverse(reversePeakHashes)
		peaksHashes = append(peaksHashes[:i], reversePeakHashes...)
		iter := NewIterator()
		iter.Items = peaksHashes
		return m.CalculateRoot([]Leaf{{newPos, newElem}}, newMMRSize, iter)
	} else {
		pushLeaf(&leaves, Leaf{newPos, newElem})
		return m.CalculateRoot(leaves, newMMRSize, m.proof)
	}
}

func (m *MerkleProof) calculatePeaksHashes(leaves []Leaf, mmrSize uint64, proofs *Iterator) ([]interface{}, error) {
	// special handle the only 1 Leaf MerkleProof
	if mmrSize == 1 && len(leaves) == 1 && leaves[0].pos == 0 {
		var items []interface{}
		for _, l := range leaves {
			items = append(items, l.hash)
		}
		return items, nil
	}

	// sort items by position
	sort.SliceStable(leaves, func(i, j int) bool {
		return leaves[i].pos < leaves[j].pos
	})

	peaks := getPeaks(mmrSize)
	var peaksHashes []interface{}
	for _, peaksPos := range peaks {
		var lvs []Leaf
		lvs = takeWhileVec(&leaves, func(l Leaf) bool {
			return l.pos <= peaksPos
		})

		var peakRoot interface{}
		if len(lvs) == 1 && lvs[0].pos == peaksPos {
			// Leaf is the peak
			peakRoot = lvs[0].hash
			// remove Leaf
			lvs = append(lvs[:0], lvs[0+1:]...)
		} else if len(lvs) == 0 {
			// if empty, means the next proof is a peak root or rhs bagged root
			if proof := proofs.next(); proof != nil {
				peakRoot = proof
			} else {
				// means that either all right peaks are bagged, or proof is corrupted
				// so we break loop and check no items left
				break
			}
		} else {
			var err error
			peakRoot, err = m.calculatePeakRoot(lvs, peaksPos, proofs)
			if err != nil {
				return nil, err
			}
		}
		peaksHashes = append(peaksHashes, peakRoot)
	}
	// ensure nothing left in leaves
	if len(leaves) != 0 {
		return nil, ErrCorruptedProof
	}

	// check rhs peaks
	if rhsPeaksHashes := proofs.next(); rhsPeaksHashes != nil {
		peaksHashes = append(peaksHashes, rhsPeaksHashes)
	}
	// ensure nothing left in proof_iter
	if proofs.next() != nil {
		fmt.Printf("something is left in proof_iter")
		return nil, ErrCorruptedProof
	}

	return peaksHashes, nil
}

func takeWhileVec(v *[]Leaf, p func(Leaf) bool) []Leaf {
	vCopy := *v
	for i := 0; i < len(vCopy); i++ {
		if !p(vCopy[i]) {
			*v = vCopy[i:]
			return vCopy[:i]
		}
	}
	*v = vCopy[:0]
	return vCopy[:]
}

func takeWhileVecUint64(v *[]uint64, p func(uint64) bool) []uint64 {
	vCopy := *v
	for i := 0; i < len(vCopy); i++ {
		if !p(vCopy[i]) {
			*v = vCopy[i:]
			return vCopy[:i]
		}
	}
	*v = vCopy[:0]
	return vCopy[:]
}
