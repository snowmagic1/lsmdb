package table

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/errors"
	"github.com/snowmagic1/lsmdb/util"
)

type ErrCorrupted struct {
	Pos    int64
	Size   int64
	Kind   string
	Reason string
}

type block struct {
	bh            blockHandle
	data          []byte
	restartsLen   int
	restartOffset int
}

func (e *ErrCorrupted) Error() string {
	return fmt.Sprintf("corruption on %s (pos=%d): %s", e.Kind, e.Pos, e.Reason)
}

func (r *Reader) newErrCorrupted(pos, size int64, kind, reason string) error {
	return &errors.ErrCorrupted{Fd: 0, Err: &ErrCorrupted{Pos: pos, Size: size, Kind: kind, Reason: reason}}
}

func (r *Reader) newErrCorruptedBH(bh blockHandle, reason string) error {
	return r.newErrCorrupted(int64(bh.offset), int64(bh.length), "bh", reason)
}

type Reader struct {
	mu     sync.RWMutex
	reader io.ReaderAt
	err    error
	keyCmp db.Comparer

	metaBH   blockHandle
	indexBH  blockHandle
	filterBH blockHandle

	indexBlock *block
}

func (r *Reader) readRawBlock(bh blockHandle, verifyChecksum bool) ([]byte, error) {
	data := make([]byte, bh.length+blockTrailerLen)
	if _, err := r.reader.ReadAt(data, int64(bh.offset)); err != nil {
		return nil, err
	}

	if verifyChecksum {
		n := bh.length + 1
		checksumRead := binary.LittleEndian.Uint32(data[n:])
		checksumExpected := util.NewCRC(data[:n]).Value()

		if checksumRead != checksumExpected {
			return nil, r.newErrCorruptedBH(bh, fmt.Sprintf("check sum mismatch"))
		}
	}

	switch db.Compression(data[bh.length]) {
	case db.CompressionNo:
		data = data[:bh.length]
	case db.CompressionSnappy:
	default:
		return nil, r.newErrCorruptedBH(bh, fmt.Sprintf("unknown compression type"))
	}

	return data, nil
}

func (r *Reader) readBlock(bh blockHandle, verifyChecksum bool) (*block, error) {
	data, err := r.readRawBlock(bh, verifyChecksum)
	if err != nil {
		return nil, err
	}

	restartsLen := int(binary.LittleEndian.Uint32(data[len(data)-4:]))
	b := &block{
		bh:            bh,
		data:          data,
		restartsLen:   restartsLen,
		restartOffset: len(data) - (restartsLen+1)*4,
	}

	return b, nil
}

func (r *Reader) Get(key []byte) (val []byte, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.err != nil {
		return nil, r.err
	}

	return nil, nil
}

func (r *Reader) find(key []byte, useFilter bool) (rKey, value []byte, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.err != nil {
		err = r.err
		return
	}

	return
}

func NewReader(f io.ReaderAt, size int64) (*Reader, error) {
	if f == nil {
		return nil, errors.New("nil file")
	}

	r := &Reader{
		reader: f,
	}

	if size < footerLen {
		r.err = r.newErrCorrupted(0, size, "table", "too small")
		return r, nil
	}

	footerPos := size - footerLen
	var footer [footerLen]byte
	_, err := r.reader.ReadAt(footer[:], footerPos)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if string(footer[footerLen-len(magic):]) != magic {
		r.err = r.newErrCorrupted(footerPos, footerLen, "table-footer", "bad magic number")
		return r, nil
	}

	var n int

	r.metaBH, n = decodeBlockHandle(footer[:])
	if n == 0 {
		r.err = r.newErrCorrupted(footerPos, footerLen, "table-footer", "bad metaindex block handle")
	}

	r.indexBH, n = decodeBlockHandle(footer[n:])
	if n == 0 {
		r.err = r.newErrCorrupted(footerPos, footerLen, "table-footer", "bad index block handle")
		return r, nil
	}

	r.indexBlock, err = r.readBlock(r.indexBH, true)
	if err != nil {
		r.err = err
		return nil, err
	}

	return r, nil
}
