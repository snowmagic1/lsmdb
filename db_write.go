package lsmdb

import (
	"github.com/snowmagic1/lsmdb/db"
)

func (db *DB) Put(key, val []byte, wo *db.WriteOptions) error {
	return db.putRecord(KeyTypeSet, key, val, wo)
}

func (db *DB) putRecord(kt internalKeyType, key, val []byte, wo *db.WriteOptions) error {
	if err := db.ok(); err != nil {
		return err
	}

	merge := wo.GetNoWriteMerge() && db.s.o.GetNoWriteMerge()
	sync := wo.GetSync() && !db.s.o.GetNoSync()

	if merge {

	} else {
		select {
		case db.writeLockC <- struct{}{}:
		}
	}

	batch := db.batchPool.Get().(*Batch)
	batch.Reset()
	batch.appendRecord(kt, key, val)

	return db.writeLocked(batch, batch, merge, sync)
}

func (db *DB) unlockWrite() {
	<-db.writeLockC
}

func (db *DB) writeLocked(batch, ourBatch *Batch, merge, sync bool) error {
	// flush memdb

	var (
		batches = []*Batch{batch}
	)

	seq := db.seq + 1

	// write journal

	// put batches
	for _, batch := range batches {
		if err := batch.putMem(seq, db.memdb.DB); err != nil {
			panic(err)
		}
		seq += uint64(batch.Len())
	}

	db.addSeq(uint64(batchesLen(batches)))

	db.unlockWrite()

	return nil
}

func (db *DB) writeJournal(batches []*Batch, seq uint64, sync bool) error {
	return nil
}
