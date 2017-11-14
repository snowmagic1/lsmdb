package lsmdb

import "github.com/snowmagic1/lsmdb/memdb"

func (s *session) pickMemdbLevel(umin, umax []byte, maxLevel int) int {
	v := s.version()
	defer v.release()

	return v.pickMemdbLevel(umin, umax, maxLevel)
}

func (s *session) flushMemdb(rec *sessionRecord, mdb *memdb.DB) (int, error) {
	// create sorted table
	iter := mdb.NewIterator(nil)
	defer iter.Release()

	t, _, err := s.tops.createFrom(iter)
	if err != nil {
		return 0, err
	}

	flushLevel := s.pickMemdbLevel(t.imin.userKey(), t.imax.userKey(), 0)
	rec.addTableFile(flushLevel, t)

	return 0, nil
}
