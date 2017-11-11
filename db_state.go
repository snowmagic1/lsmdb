package lsmdb

import "github.com/snowmagic1/lsmdb/memdb"

type memDB struct {
	db *DB
	*memdb.DB
	ref int32
}
