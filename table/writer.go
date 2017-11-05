package table

import (
	"errors"
	"io"
)

type Writer struct {
	writer io.Writer
}

func (w *Writer) Get(key []byte) ([]byte, error) {
	return nil, errors.New("can't get on writer")
}

func (w *Writer) Set(key, val []byte) error {

}
