package table

import (
	"os"
)

type blockHandle struct {
	offset, length uint64
}

type block []byte

type Reader struct {
	file *os.File
}

func (r *Reader) Get(key []byte) (val []byte, err error) {

}

func (r *Reader) Set(key, val []byte) error {
	return errors.New("can't set on reader")
}

func (r *Reader) Delete(key, val []byte) error {
	return errors.New("can't set on reader")
}

func NewReader(f *os.File) *Reader {

}