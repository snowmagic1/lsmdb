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
