package lsmdb

import (
	"encoding/binary"

	"github.com/snowmagic1/lsmdb/memdb"
)

type batchIndex struct {
	keyType   internalKeyType
	keyOffset int
	keyLen    int
	valOffset int
	valLen    int
}

func (index *batchIndex) k(data []byte) []byte {
	return data[index.keyOffset : index.keyOffset+index.keyLen]
}

func (index *batchIndex) v(data []byte) []byte {
	if index.valLen != 0 {
		return data[index.valOffset : index.valOffset+index.valLen]
	}
	return nil
}

type Batch struct {
	data  []byte
	index []batchIndex

	internalLen int
}

func (b *Batch) Reset() {
	b.data = b.data[:0]
	b.index = b.index[:0]
	b.internalLen = 0
}

func newBatch() interface{} {
	return &Batch{}
}

func (b *Batch) appendRecord(kt internalKeyType, key, val []byte) {
	n := 1 + binary.MaxVarintLen32 + len(key)
	if kt == KeyTypeSet {
		n += binary.MaxVarintLen32 + len(val)
	}

	b.grow(n)

	index := batchIndex{keyType: kt}
	curr := len(b.data)
	data := b.data[:curr+n]
	data[curr] = byte(kt)
	curr++

	curr += binary.PutUvarint(data[curr:], uint64(len(key)))
	index.keyOffset = curr
	index.keyLen = len(key)
	curr += copy(data[curr:], key)

	if kt == KeyTypeSet {
		curr += binary.PutUvarint(data[curr:], uint64(len(val)))
		index.valOffset = curr
		index.valLen = len(val)
		curr += copy(data[curr:], val)
	}

	b.data = data[:curr]
	b.index = append(b.index, index)
	b.internalLen += index.keyLen + index.valLen + 8
}

func (b *Batch) grow(n int) {
	o := len(b.data)
	if cap(b.data)-o < n {
		ndata := make([]byte, o, o+n)
		copy(ndata, b.data)
		b.data = ndata
	}
}

func (b *Batch) putMem(seq uint64, mdb *memdb.DB) error {
	for i, index := range b.index {
		ik := makeInternalKey(index.k(b.data), KeyTypeSet, seq+uint64(i))
		if err := mdb.Put(ik, index.v(b.data)); err != nil {
			return err
		}
	}

	return nil
}

func (b *Batch) Len() int {
	return len(b.index)
}

func batchesLen(batches []*Batch) int {
	batchLen := 0
	for _, batch := range batches {
		batchLen += batch.Len()
	}
	return batchLen
}
