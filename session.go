package lsmdb

import (
	"os"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/storage"
)

type session struct {
	stSeqNum uint64

	stor     storage.Storage
	storLock storage.Locker

	o *cachedOptions
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
	}

	if o != nil {
		*s.o.Options = *o
	}

	return
}
