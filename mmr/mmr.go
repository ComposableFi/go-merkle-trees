package mmr

import (
	"reflect"
	"sort"

	"github.com/ComposableFi/merkle-go/merkle"
)

// MMR contains fields for computing the MMR tree.
type MMR struct {
	// size is the MMR size of the tree
	size  uint64
	batch *Batch
	// merge accepts any type that satisfies the Merge interface
	merge  merkle.Merge
	leaves []Leaf
}

// NewMMR returns a new MMR type. It takes three arguments. It takes the mmrSize, Store and Merge interfaces. It accepts
// any type that satisfies both the Store and Merge interfaces.
func NewMMR(mmrSize uint64, s Store, leaves []Leaf, m merkle.Merge) *MMR {
	return &MMR{
		size:   mmrSize,
		batch:  NewBatch(s),
		merge:  m,
		leaves: leaves,
	}
}

func (m *MMR) findElem(pos uint64, hashes [][]byte) ([]byte, error) {
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

// MMRSize returns the size of the mmr tree
func (m *MMR) MMRSize() uint64 {
	return m.size
}

// IsEmpty returns true if the MMR is empty and false if it is not.
func (m *MMR) IsEmpty() bool {
	return m.size == 0
}

// Push adds an element to the store and returns its position
func (m *MMR) Push(elem []byte) (uint64, error) {
	var elems [][]byte
	// position of new elems
	elemPos := m.size
	elems = append(elems, elem)

	var height uint32
	var pos = elemPos
	// continue to merge tree node if next Pos higher than current
	for posHeightInTree(pos+1) > height {
		pos++
		leftPos := pos - parentOffset(height)
		rightPos := leftPos + siblingOffset(height)
		leftElem, err := m.findElem(leftPos, elems)
		if err != nil {
			return 0, err
		}

		rightElem, err := m.findElem(rightPos, elems)
		if err != nil {
			return 0, err
		}

		parentElem := m.merge.Merge(leftElem, rightElem)
		elems = append(elems, parentElem.([]byte))
		height++
	}
	// store hashes
	m.batch.append(elemPos, elems)
	// update mmrSize
	m.size = pos + 1
	return elemPos, nil
}

// GetRoot returns the root of the MMR tree
func (m *MMR) GetRoot() ([]byte, error) {
	if m.size == 0 {
		return nil, ErrGetRootOnEmpty
	} else if m.size == 1 {
		e := m.batch.getElem(0)
		if e == nil {
			return nil, ErrInconsistentStore
		}
		return e, nil
	}

	var peaks [][]byte
	for _, peakPos := range getPeaks(m.size) {
		elem := m.batch.getElem(peakPos)
		if elem == nil {
			return nil, ErrInconsistentStore
		}
		peaks = append(peaks, elem)
	}

	var p []byte
	if p = m.bagRHSPeaks(peaks); p == nil {
		return nil, ErrInconsistentStore
	}

	return p, nil
}

func (m *MMR) bagRHSPeaks(rhsPeaks [][]byte) []byte {
	for len(rhsPeaks) > 1 {
		var rp, lp []byte
		if rp, rhsPeaks = pop(rhsPeaks); rp == nil {
			panic("pop")
		}

		if lp, rhsPeaks = pop(rhsPeaks); lp == nil {
			panic("pop")
		}
		rhsPeaks = append(rhsPeaks, m.merge.Merge(rp, lp).([]byte))
	}

	if len(rhsPeaks) > 0 {
		return rhsPeaks[len(rhsPeaks)-1]
	}
	return nil
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
		if !(pos <= peakPos) {
			panic("Pos is not less than or equal to peak position")
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
			}
			return pos + siblingOffset, pos + parentOffset(height)
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

// GenProof generates merkle proof for positions. It sorts positions, pushes merkle proof to proof by peak from left to
// right. It then pushes bagged right hand side root
func (m *MMR) GenProof(posList []uint64) (*Proof, error) {
	if len(posList) == 0 {
		return nil, ErrGenProofForInvalidLeaves
	}
	if m.size == 1 && reflect.DeepEqual(posList, []uint64{0}) {
		return NewProof(m.size, [][]byte{}, m.leaves, m.merge), nil
	}

	sort.Slice(posList, func(i, j int) bool {
		return posList[i] < posList[j]
	})
	var peaks = getPeaks(m.size)
	var proof = NewIterator()
	// generate merkle proof for each peaks
	var baggingTrack uint
	for _, peakPos := range peaks {
		pl := takeWhileVecUint64(&posList, func(u uint64) bool {
			return u <= peakPos
		})
		if len(pl) == 0 {
			baggingTrack++
		} else {
			baggingTrack = 0
		}

		err := m.genProofForPeak(proof, pl, peakPos)
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
		var p []byte

		if p = m.bagRHSPeaks(rhsPeaks); p == nil {
			panic("bagging rhs peaks")
		}
		proof.push(p)
	}

	return NewProof(m.size, proof.Items, m.leaves, m.merge), nil
}

// Commit calls the commit method on the batch property. It adds a batch element to the store
func (m *MMR) Commit() {
	m.batch.commit()
}

// Proof is the mmr proof. It is constructed to verify an MMR leaf.
type Proof struct {
	mmrSize uint64
	proof   *Iterator
	Merge   merkle.Merge
	Leaves  []Leaf
}

// NewProof creates and returns new Proof. It takes the mmrSize, proof which is of type *Iterator and any type
// that satisfies the Merge interface.
func NewProof(mmrSize uint64, proofItems [][]byte, mmrLeaves []Leaf, m merkle.Merge) *Proof {
	return &Proof{
		mmrSize: mmrSize,
		proof:   &Iterator{Items: proofItems},
		Merge:   m,
		Leaves:  mmrLeaves,
	}
}

// MMRSize returns the mmr size
func (m *Proof) MMRSize() uint64 {
	return m.mmrSize
}

// ProofItems returns all the proof items from the Iterator.
func (m *Proof) ProofItems() [][]byte {
	return m.proof.Items
}

func (m *Proof) calculatePeakRoot(leaves []Leaf, peakPos uint64, proofs *Iterator) ([]byte, error) {
	if len(leaves) == 0 {
		panic("can't be empty")
	}

	// (position, Hash, height)
	var queue []leafWithHash
	for _, l := range leaves {
		queue = append(queue, leafWithHash{LeafIndexToPos(l.Index), l.Hash, 0})
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
				// implies Pos is right sibling
				return pos - siblingOffset, pos + 1
			}
			// Pos is left sibling
			return pos + siblingOffset, pos + parentOffset(height)
		}()

		var siblingItem []byte
		if len(queue) > 0 && queue[0].pos == sibPos {
			siblingItem, queue = queue[0].hash, queue[1:]
		} else {
			if siblingItem = proofs.Next(); siblingItem == nil {
				return nil, ErrCorruptedProof
			}
		}

		var parentItem []byte
		if nextHeight > height {
			parentItem = m.Merge.Merge(siblingItem, item).([]byte)
		} else {
			parentItem = m.Merge.Merge(item, siblingItem).([]byte)
		}

		if parentPos < peakPos {
			queue = append(queue, leafWithHash{parentPos, parentItem, height + 1})
		} else {
			return parentItem, nil
		}
	}

	return nil, ErrCorruptedProof
}

func (m *Proof) baggingPeaksHashes(peaksHashes [][]byte) ([]byte, error) {
	var rightPeak, leftPeak []byte
	for len(peaksHashes) > 1 {
		if rightPeak, peaksHashes = pop(peaksHashes); rightPeak == nil {
			panic("pop")
		}

		if leftPeak, peaksHashes = pop(peaksHashes); leftPeak == nil {
			panic("pop")
		}
		peaksHashes = append(peaksHashes, m.Merge.Merge(rightPeak, leftPeak).([]byte))
	}

	if len(peaksHashes) == 0 {
		return nil, ErrCorruptedProof
	}
	return peaksHashes[len(peaksHashes)-1], nil
}

// CalculateRoot calculates and returns the root of the MMR tree using the leaves, mmrSize and proofs. It sorts the leaves
// by position, calculates the root of each peak and bags the peaks
func (m *Proof) CalculateRoot(leaves []Leaf, mmrSize uint64, proofs *Iterator) ([]byte, error) {
	var peaksHashes, err = m.calculatePeaksHashes(leaves, mmrSize, proofs)
	if err != nil {
		return nil, err
	}

	return m.baggingPeaksHashes(peaksHashes)
}

// CalculateRootWithNewLeaf calculates and returns a new root provided a new leaf element, new position and new MMRsize.
// from merkle proof of leaf n to calculate merkle root of n + 1 leaves. By observing the MMR construction graph we know
// it is possible. https://github.com/jjyr/merkle-mountain-range#construct this is kinda tricky, but it works, and useful
func (m *Proof) CalculateRootWithNewLeaf(leaves []Leaf, newIndex uint64, newElem []byte, newMMRSize uint64) ([]byte, error) {
	newPos := LeafIndexToPos(newIndex)
	posHeight := posHeightInTree(newPos)
	nextHeight := posHeightInTree(newPos + 1)
	if nextHeight > posHeight {
		peaksHashes, err := m.calculatePeaksHashes(leaves, m.mmrSize, m.proof)
		if err != nil {
			return nil, err
		}
		peaksPos := getPeaks(newMMRSize)
		// reverse touched peaks
		var i uint
		for peaksPos[i] < newPos {
			i++
		}

		reversePeakHashes := peaksHashes[i:]
		reverse(reversePeakHashes)
		peaksHashes = append(peaksHashes[:i], reversePeakHashes...)
		iter := NewIterator()
		iter.Items = peaksHashes
		return m.CalculateRoot([]Leaf{{Index: newIndex, Hash: newElem}}, newMMRSize, iter)
	}
	pushLeaf(&leaves, Leaf{Index: newIndex, Hash: newElem})
	return m.CalculateRoot(leaves, newMMRSize, m.proof)
}

// Verify takes a root and leaves as arguments. It calculates a root from the leaves using the CalculateRoot method and
// compares it with the supplied root. It returns tree if the roots are equal and false if they are not.
func (m *Proof) Verify(root []byte) bool {
	calculatedRoot, err := m.CalculateRoot(m.Leaves, m.mmrSize, m.proof)
	if err != nil {
		log.Errorf("root verification: %s \n", err.Error())
		return false
	}

	return reflect.DeepEqual(calculatedRoot, root)
}

func (m *Proof) calculatePeaksHashes(leaves []Leaf, mmrSize uint64, proofs *Iterator) ([][]byte, error) {
	// special handle the only 1 Hash Proof
	if mmrSize == 1 && len(leaves) == 1 && LeafIndexToPos(leaves[0].Index) == 0 {
		var items [][]byte
		for _, l := range leaves {
			items = append(items, l.Hash)
		}
		return items, nil
	}

	// sort items by position
	sort.SliceStable(leaves, func(i, j int) bool {
		return LeafIndexToPos(leaves[i].Index) < LeafIndexToPos(leaves[j].Index)
	})

	peaks := getPeaks(mmrSize)
	var peaksHashes [][]byte
	for _, peaksPos := range peaks {
		lvs := takeWhileVec(&leaves, func(l Leaf) bool {
			return LeafIndexToPos(l.Index) <= peaksPos
		})

		var peakRoot []byte
		if len(lvs) == 1 && LeafIndexToPos(lvs[0].Index) == peaksPos {
			// Hash is the peak
			peakRoot = lvs[0].Hash
		} else if len(lvs) == 0 {
			// if empty, means the next proof is a peak root or rhs bagged root
			if proof := proofs.Next(); proof != nil {
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
	if rhsPeaksHashes := proofs.Next(); rhsPeaksHashes != nil {
		peaksHashes = append(peaksHashes, rhsPeaksHashes)
	}
	// ensure nothing left in proof_iter
	if proofs.Next() != nil {
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
