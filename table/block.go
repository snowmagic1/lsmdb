package table

import (
	"encoding/binary"
	"sort"

	"github.com/snowmagic1/lsmdb/db"
)

type block struct {
	bh             blockHandle
	data           []byte
	restartsLen    int
	restartsOffset int
}

func (b *block) seek(cmp db.Comparer, rsSearchStart, rsSearchEnd int, key []byte) (restartIndex, recordStartOffset int, err error) {
	rsSearchLen := rsSearchEnd - rsSearchStart
	rsFirstGreaterThan := sort.Search(rsSearchLen, func(i int) bool {
		restartOffset := b.restartsOffset + 4*(rsSearchStart+i)
		restartKeyOffset := int(binary.LittleEndian.Uint32(b.data[restartOffset:]))

		// shared size zero, since this is a restart point
		restartKeyOffset++
		v1, n1 := binary.Uvarint(b.data[restartKeyOffset:])   // key length
		_, n2 := binary.Uvarint(b.data[restartKeyOffset+n1:]) // value length

		// offset for key
		m := restartKeyOffset + n1 + n2
		return cmp.Compare(b.data[m:m+int(v1)], key) > 0
	}) + rsSearchStart - 1

	restartIndex = rsFirstGreaterThan
	if restartIndex < rsSearchStart {
		restartIndex = rsSearchStart
	}

	recordStartOffset = int(binary.LittleEndian.Uint32(b.data[b.restartsOffset+4*restartIndex:]))

	return
}

func (b *block) entry(offset int) (key, val []byte, sharedLen, entryLen int, err error) {
	if offset >= b.restartsOffset {
		if offset != b.restartsOffset {
			err = &ErrCorrupted{Reason: "offset not aligned"}
		}
		return
	}
	v0, n0 := binary.Uvarint(b.data[offset:])       // shared len
	v1, n1 := binary.Uvarint(b.data[offset+n0:])    // non-shared len
	v2, n2 := binary.Uvarint(b.data[offset+n0+n1:]) // val len

	headerLen := n0 + n1 + n2
	entryLen = headerLen + int(v1) + int(v2)

	if n0 <= 0 || n1 <= 0 || n2 <= 0 || offset+entryLen > b.restartsOffset {
		err = &ErrCorrupted{Reason: "entry corrupted"}
		return
	}

	key = b.data[offset+headerLen : offset+headerLen+int(v1)]
	val = b.data[offset+headerLen+int(v1) : offset+entryLen]

	sharedLen = int(v0)

	return
}
