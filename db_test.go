package lsmdb_test

import (
	"testing"

	"github.com/snowmagic1/lsmdb"
)

func TestBasics(t *testing.T) {
	dbname := "dbtest"

	db, err := lsmdb.OpenFile(dbname, nil)
	if err != nil {
		t.Errorf("failed to open file, %v", err)
		return
	}

	key := "key1"
	val := "val1"

	err = db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		t.Errorf("failed to put, %v", err)
	}
}
