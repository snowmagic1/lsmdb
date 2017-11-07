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

func (b *block) seek(cmp db.Comparer, rstart, rlimit int, key []byte) (index, offset int, err error) {
	index = sort.Search(b.restartsLen-rstart-(b.restartsLen-rlimit), func(i int) bool {
		restartKeyOffset := int(binary.LittleEndian.Uint32(b.data[b.restartsOffset+4*(rstart+i):]))

		// shared size zero, since this is a restart point
		restartKeyOffset++
		v1, n1 := binary.Uvarint(b.data[offset:])   // key length
		_, n2 := binary.Uvarint(b.data[offset+n1:]) // value length

		// offset for key
		m := offset + n1 + n2
		return cmp.Compare(b.data[m:m+int(v1)], key) > 0
	}) + rstart - 1

	if index < rstart {
		index = rstart
	}

	offset = int(binary.LittleEndian.Uint32(b.data[b.restartsOffset+4*index:]))

	return
}
