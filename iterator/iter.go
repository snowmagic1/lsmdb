package iterator

import "github.com/snowmagic1/lsmdb/util"

type IteratorSeeker interface {
	First() bool
	Last() bool
	Seek(key []byte) bool
	Next() bool
	Prev() bool
}

type CommonIterator interface {
	IteratorSeeker
	util.Releaser
	Valid() bool
	Error() error
}

type Iterator interface {
	CommonIterator
	Key() []byte
	Value() []byte
}
