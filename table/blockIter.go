package table

type dir int

const (
	dirReleased dir = iota - 1
	dirSOI
	dirEOI
	dirBackward
	dirForward
)

type blockIter struct {
	tr    *Reader
	block *block
	err   error
	dir   dir

	// current state
	currKey     []byte
	currVal     []byte
	currOffset  int
	currRsIndex int

	// block meta data
	rsStartIdx int
	rsEndIdx   int

	offsetStart int
	offsetEnd   int
}

func (bi *blockIter) Key() []byte {
	if bi.err != nil || bi.dir <= dirEOI {
		return nil
	}

	return bi.currKey
}

func (bi *blockIter) Val() []byte {
	if bi.err != nil || bi.dir <= dirEOI {
		return nil
	}

	return bi.currVal
}

func (bi *blockIter) Seek(key []byte) bool {
	if bi.err != nil {
		return false
	}

	rsIndex, recordStartOffset, err := bi.block.seek(bi.tr.keyCmp, bi.rsStartIdx, bi.rsEndIdx, key)
	if err != nil {
		bi.err = err
		return false
	}

	bi.currRsIndex = rsIndex
	bi.currOffset = recordStartOffset

	if bi.dir == dirSOI || bi.dir == dirEOI {
		bi.dir = dirForward
	}

	for bi.Next() {
		if bi.tr.keyCmp.Compare(bi.currKey, key) >= 0 {
			return true
		}
	}

	return false
}

func (bi *blockIter) Next() bool {
	if bi.currOffset >= bi.offsetEnd {
		bi.dir = dirEOI
		if bi.currOffset != bi.offsetEnd {
			bi.err = bi.tr.newErrCorruptedBH(bi.block.bh, "entry offset not aligned")
		}

		return false
	}

	key, val, sharedLen, entryLen, err := bi.block.entry(bi.currOffset)
	if err != nil {
		bi.err = err
		return false
	}

	if entryLen == 0 {
		bi.dir = dirEOI
		return false
	}

	// for restart key, sharedLen is 0
	bi.currKey = append(bi.currKey[:sharedLen], key...)
	bi.currVal = val
	bi.currOffset += entryLen
	bi.dir = dirForward

	return true
}
