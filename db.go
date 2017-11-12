package lsmdb

import (
	"container/list"
	"io"
	"log"
	"sync"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/memdb"
	"github.com/snowmagic1/lsmdb/storage"

	"github.com/snowmagic1/lsmdb/journal"
)

type DB struct {
	seq uint64

	// session
	s *session

	// memdb
	memMu       sync.RWMutex
	memPool     chan *memdb.DB
	memdb       *memDB
	frozenMemdb *memDB
	journal     *journal.Writer
	frozenSeq   uint64

	// snapshot
	snapsMu   sync.Mutex
	snapsList *list.List

	// write
	batchPool  sync.Pool
	writeLockC chan struct{}
	writeAckC  chan error

	// close
	closer io.Closer
}

func openDB(s *session) (*DB, error) {
	db := &DB{
		s:   s,
		seq: s.stSeqNum,
		// memdb
		memPool: make(chan *memdb.DB, 1),
		//snapshot
		snapsList: list.New(),
		//write
		batchPool:  sync.Pool{New: newBatch},
		writeLockC: make(chan struct{}, 1),
		writeAckC:  make(chan error),
	}

	// recover journals
	if err := db.recoverJournal(); err != nil {
		log.Println("failed to recover journal, ", err)
		return nil, err
	}

	// remove obsoletes file

	return db, nil
}

func OpenFile(path string, o *db.Options) (db *DB, err error) {
	stor, err := storage.OpenDiskStorage(path, o.GetReadOnly())
	if err != nil {
		return
	}

	db, err = Open(stor, o)
	if err != nil {
		stor.Close()
	} else {
		db.closer = stor
	}

	return
}

func Open(stor storage.Storage, o *db.Options) (db *DB, err error) {
	s, err := newSession(stor, o)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			// s.Close()
		}
	}()

	// recover

	return openDB(s)
}

func (db *DB) recoverJournal() error {
	if _, err := db.newMem(0); err != nil {
		return err
	}

	return nil
}

func (db *DB) Get(key []byte, ro *db.ReadOptions) (val []byte, err error) {
	if err = db.ok(); err != nil {
		return
	}

	se := db.acquireSnapshot()
	defer db.releaseSnapshot(se)

	return db.get(key, se.seq, ro)
}

func (db *DB) get(key []byte, seq uint64, ro *db.ReadOptions) (val []byte, err error) {
	ikey := makeInternalKey(key, KeyTypeMax, seq)

	memdb, fmemdb := db.getMems()
	for _, m := range []*memDB{memdb, fmemdb} {
		if m == nil {
			continue
		}
		defer m.decref()

		if ok, mv, me := memGet(m.DB, ikey, db.s.keycmp); ok {
			return append([]byte{}, mv...), me
		}
	}

	return
}

func memGet(mdb *memdb.DB, ikey internalKey, keycmp db.Comparer) (ok bool, mv []byte, err error) {
	mk, mv, err := mdb.Find(ikey)
	if err == nil {
		ik := internalKey(mk)
		if keycmp.Compare(ik.userKey(), ikey.userKey()) == 0 {
			if ik.keyType() == KeyTypeDelete {
				return true, nil, ErrNotFound
			}
			return true, mv, nil
		}
	} else if err != ErrNotFound {
		return true, nil, err
	}

	return
}
