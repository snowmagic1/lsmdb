package journal

import "io"

const (
	blockSize  = 32 * 1024
	headerSize = 7
)

type flusher interface {
	Flush() error
}

type Reader struct {
}

type Writer struct {
	w   io.Writer
	seq int
	f   flusher
	buf [blockSize]byte
}

func NewWriter(w io.Writer) *Writer {
	f, _ := w.(flusher)
	return &Writer{
		w: w,
		f: f,
	}
}
