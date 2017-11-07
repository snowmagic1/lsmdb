package table

type blockIter struct {
	tr    *Reader
	block *block
	err   error

	restartIndex int

	restartStart int
	restartLen   int
}

func (bi *blockIter) Seek(key []byte) bool {
	if bi.err != nil {
		return false
	}

	ri, offset, err := bi.block.seek(bi.tr.keyCmp, bi.restartStart, bi.restartLen, key)
	if err != nil {
		bi.err = err
		return false
	}

	bi.restartIndex = ri

	if offset <= 0 {
		return false
	}

	return true
}
