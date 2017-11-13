package lsmdb

import (
	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/iterator"
	"github.com/snowmagic1/lsmdb/util"
)

func (db *DB) newRawIterator(slice *util.Range, ro *db.ReadOptions) iterator.Iterator {
	em, _ := db.getMems()
	return em.NewIterator(slice)
}

func (db *DB) newIterator(seq uint64, ro *db.ReadOptions) iterator.Iterator {
	slice := &util.Range{Start: nil, Limit: nil}
	return db.newRawIterator(slice, ro)
}
