package lsmdb

import (
	"container/list"
	"io"
	"sync"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/storage"

	"github.com/snowmagic1/lsmdb/journal"
)

type DB struct {
	seq uint64

	// session
	s *session

	// memdb
	memMu       sync.RWMutex
	memdb       *memDB
	frozenMemdb *memDB
	journal     *journal.Writer

	// snapshot
	snapsMu   sync.Mutex
	snapsList *list.List

	// close
	closer io.Closer
}

func OpenFile(path string, o *db.Options) (db *DB, err error) {
	stor, err := storage.OpenFile(path, o.GetReadOnly())
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

func openDB(s *session) (*DB, error) {
	db := &DB{
		s:   s,
		seq: s.stSeqNum,
	}

	// recover journals

	// remove obsoletes file

	return db, nil
}
