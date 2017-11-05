package table

import (
	"bufio"
	"errors"
	"os"
)

type Writer struct {
	writer *bufio.Writer

	prevKey []byte
}

func (w *Writer) Get(key []byte) ([]byte, error) {
	return nil, errors.New("can't get on writer")
}

func (w *Writer) Set(key, val []byte) error {

}

func (w *Writer) Close() error {

}

func NewWriter(f *os.File) *Writer {
	w := &Writer{
		prevKey: make([]byte, 0, 256),
	}

	w.writer = bufio.NewWriter(f)
}
