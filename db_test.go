package lsmdb_test

import (
	"testing"

	"github.com/snowmagic1/lsmdb"
)

func TestMakeKey(t *testing.T) {
	dbname := "dbtest"

	_, err := lsmdb.OpenFile(dbname, nil)
	if err != nil {
		t.Errorf("failed to open file, %v", err)
	}
}
