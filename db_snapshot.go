package lsmdb

import "container/list"

type snapshotElement struct {
	seq uint64
	ref int
	e   *list.Element
}

func (db *DB) acquireSnapshot() *snapshotElement {
	db.snapsMu.Lock()
	defer db.snapsMu.Unlock()

	seq := db.getSeq()

	se := &snapshotElement{seq: seq, ref: 1}
	se.e = db.snapsList.PushBack(se)

	return se
}

func (db *DB) releaseSnapshot(se *snapshotElement) {
	db.snapsMu.Lock()
	defer db.snapsMu.Unlock()

	se.ref--
	if se.ref == 0 {

	} else if se.ref < 0 {
		panic("leveldb: negtive ref count")
	}
}
