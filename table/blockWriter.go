package table

import (
	"encoding/binary"

	"github.com/snowmagic1/lsmdb/db"
)

type blockWriter struct {
	restartInternal int
	buf             []byte
	nEntries        int
	prevKey         []byte
	restartKeyLens  []uint32
	scratch         []byte
}

func (w *blockWriter) append(key, val []byte) {
	nShared := 0
	if w.nEntries%w.restartInternal == 0 {
		w.restartKeyLens = append(w.restartKeyLens, uint32(len(w.buf)))
	} else {
		nShared = db.SharedPrefixLen(w.prevKey, key)
	}

	n := binary.PutUvarint(w.scratch[0:], uint64(nShared))
	n += binary.PutUvarint(w.scratch[n:], uint64(len(key)-nShared))
	n += binary.PutUvarint(w.scratch[n:], uint64(len(val)))

	w.buf = append(w.buf, w.scratch[:n]...)
	w.buf = append(w.buf, key[nShared:]...)
	w.buf = append(w.buf, val...)

	w.prevKey = append(w.prevKey[:0], key...)

	w.nEntries++
}

func (w *blockWriter) finish() {
	if w.nEntries == 0 {
		w.restartKeyLens = append(w.restartKeyLens, 0)
	}

	restarts := append(w.restartKeyLens, uint32(len(w.restartKeyLens)))
	buf4 := w.scratch[:4]
	for _, x := range restarts {
		binary.LittleEndian.PutUint32(buf4, x)
		w.buf = append(w.buf, buf4...)
	}
}

func (w *blockWriter) reset() {
	w.buf = w.buf[:0]
	w.nEntries = 0
	w.restartKeyLens = w.restartKeyLens[:0]
}

func (w *blockWriter) len() int {
	restartLen := len(w.restartKeyLens)
	if restartLen == 0 {
		restartLen = 1
	}

	// buf len + restarts + restart lens
	return len(w.buf) + 4*restartLen + 4
}
