package lsmdb

import (
	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/memdb"
	"github.com/snowmagic1/lsmdb/util"
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

func (db *DB) CompactRange(r util.Range) error {
	if err := db.ok(); err != nil {
		return err
	}

	select {
	case db.writeLockC <- struct{}{}:
	case <-db.closeC:
		return ErrClosed
	}

	mdb := db.getEffectiveMem()
	if mdb == nil {
		return ErrClosed
	}
	defer mdb.decref()

	if isMemOverlaps(db.s.keycmp, mdb.DB, r.Start, r.Limit) {
		if _, err := db.rotateMem(0, false); err != nil {
			<-db.writeLockC
			return err
		}
		<-db.writeLockC
		if err := db.compTriggerWait(db.mcompCmdC); err != nil {
			return err
		}
	} else {
		<-db.writeLockC
	}

	// table compaction

	return nil
}

func isMemOverlaps(keycmp db.Comparer, mem *memdb.DB, min, max []byte) bool {
	iter := mem.NewIterator(nil)
	defer iter.Release()

	lessThanMax := (max == nil || (iter.First() && keycmp.Compare(max, internalKey(iter.Key()).userKey()) >= 0))
	moreThanMin := (min == nil || (iter.Last() && keycmp.Compare(min, internalKey(iter.Key()).userKey()) <= 0))

	return lessThanMax && moreThanMin
}

func (db *DB) rotateMem(n int, wait bool) (mem *memDB, err error) {
	retryLimit := 3
retry:
	// wait for pending memdb compaction
	err = db.compTriggerWait(db.mcompCmdC)
	if err != nil {
		return
	}
	retryLimit--

	// create new memdb and journal
	mem, err = db.newMem(n)
	if err != nil {
		if retryLimit <= 0 {
			panic("still has frozen memdb")
		}
		goto retry
	}

	if wait {
		err = db.compTriggerWait(db.mcompCmdC)
	} else {
		db.compTrigger(db.mcompCmdC)
	}

	return
}
