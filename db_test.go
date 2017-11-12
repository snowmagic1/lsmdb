package lsmdb_test

import (
	"reflect"
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

	key := []byte("key1")
	val := []byte("val1")

	err = db.Put(key, val, nil)
	if err != nil {
		t.Errorf("failed to put, %v", err)
	}

	err = db.Put(key, val, nil)
	if err != nil {
		t.Errorf("failed to put, %v", err)
	}

	rval, err := db.Get(key, nil)
	if err != nil || !reflect.DeepEqual(val, rval) {
		t.Errorf("failed to get, %v", err)
	}
}
