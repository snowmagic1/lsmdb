package table

import (
	"os"
	"sync"

	"github.com/snowmagic1/lsmdb/db"
)

type Reader struct {
	mu     sync.RWMutex
	file   *os.File
	err    error
	keyCmp db.Comparer

	metaBH   blockHandle
	indexBH  blockHandle
	filterBH blockHandle
}

func (r *Reader) Get(key []byte) (val []byte, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.err != nil {
		return nil, r.err
	}

}

func (r *Reader) find(key []byte, useFilter bool) (rKey, value []byte, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.err != nil {
		err = r.err
		return
	}

}

func NewReader(f *os.File) *Reader {

}
