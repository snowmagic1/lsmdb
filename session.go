package lsmdb

import (
	"os"
	"sync"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/storage"
)

type session struct {
	stNextFileNum int64
	stSeqNum      uint64

	stor     storage.Storage
	storLock storage.Locker
	keycmp   db.Comparer
	o        *cachedOptions
	tops     *tOps

	stVersion *version
	vMu       sync.Mutex
}

type cachedOptions struct {
	*db.Options
}

func newSession(stor storage.Storage, o *db.Options) (s *session, err error) {
	if stor == nil {
		return nil, os.ErrInvalid
	}
	storLock, err := stor.Lock()
	if err != nil {
		return
	}

	s = &session{
		stor:     stor,
		storLock: storLock,
		keycmp:   o.GetComparer(),
	}

	if o != nil {
		*s.o.Options = *o
	} else {
		no := &db.Options{}
		s.o = &cachedOptions{Options: no}
	}

	return
}
