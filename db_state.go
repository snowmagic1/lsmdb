package lsmdb

import (
	"errors"
	"sync/atomic"

	"github.com/snowmagic1/lsmdb/memdb"
	"github.com/snowmagic1/lsmdb/storage"
	"github.com/snowmagic1/lsmdb/util"
)

var (
	errHasFrozenMem = errors.New("has frozen mem")
)

type memDB struct {
	db *DB
	*memdb.DB
	ref int32
}

func (m *memDB) getref() int32 {
	return atomic.LoadInt32(&m.ref)
}

func (m *memDB) incref() {
	atomic.AddInt32(&m.ref, 1)
}

func (m *memDB) decref() {
	if ref := atomic.AddInt32(&m.ref, -1); ref == 0 {
		// put back to mempool
	} else if ref < 0 {
		panic("negative memdb ref")
	}
}

func (db *DB) ok() error {
	return nil
}

func (db *DB) addSeq(delta uint64) {
	atomic.AddUint64(&db.seq, delta)
}

func (db *DB) newMem(n int) (mem *memDB, err error) {
	fd := storage.FileDesc{Type: storage.TypeJournal, Num: db.s.allocFileNum()}
	_, err = db.s.stor.Create(fd)
	if err != nil {
		return
	}

	db.memMu.Lock()
	defer db.memMu.Unlock()

	if db.frozenMemdb != nil {
		return nil, errHasFrozenMem
	}

	if db.journal == nil {
		// db.journal == journal.NewWriter(w)
	}

	db.frozenMemdb = db.memdb
	mem = db.mpoolGet(n)
	mem.incref() // self
	mem.incref() // caller
	db.memdb = mem

	db.frozenSeq = db.seq

	return
}

func (db *DB) mpoolGet(n int) *memDB {
	var mdb *memdb.DB
	select {
	case mdb = <-db.memPool:
	default:
	}

	if mdb == nil || mdb.Capacity() < n {
		icap := util.MaxInt(db.s.o.GetWriteBuffer(), n)
		mdb = memdb.New(db.s.keycmp, icap)
	}

	return &memDB{
		db: db,
		DB: mdb,
	}
}
