package table

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/snowmagic1/lsmdb/db"
	"github.com/snowmagic1/lsmdb/util"
)

type Writer struct {
	writer io.Writer
	err    error

	// options
	keyCmp               db.Comparer
	blockRestartInternal int
	blockSize            int

	dataBlock  blockWriter
	indexBlock blockWriter
	// filter

	pendingDataBH blockHandle
	offset        uint64
	nEntries      int

	scratch   [50]byte
	bHScratch []byte
}

func (w *Writer) Append(key, val []byte) error {
	if w.err != nil {
		return w.err
	}

	if w.keyCmp.Compare(w.dataBlock.prevKey, key) >= 0 {
		w.err = fmt.Errorf("Set called in non-increasing key, prev %v curr %v", w.dataBlock.prevKey, key)
		return w.err
	}

	// filter

	// flush it if thereis bh pending
	w.addIndexEntryForPendingDataBlock(key)
	w.dataBlock.append(key, val)

	if w.dataBlock.len() >= w.blockSize {
		if err := w.writeCurrDataBlock(); err != nil {
			w.err = err
			return w.err
		}
	}

	return nil
}

func (w *Writer) addIndexEntryForPendingDataBlock(key []byte) {
	if w.pendingDataBH.length == 0 {
		return
	}

	seperator := w.keyCmp.Separator(w.scratch[:0], w.dataBlock.prevKey, key)

	n := encodeBlockHandle(w.bHScratch, w.pendingDataBH)
	w.indexBlock.append(seperator, w.bHScratch[:n])

	w.pendingDataBH = blockHandle{}
}

func (w *Writer) writeCurrDataBlock() error {
	w.dataBlock.finish()
	bh, err := w.writeRawBlock(w.dataBlock.buf, db.CompressionNo)
	if err != nil {
		w.err = err
		log.Println("failed to write block, ", err)
		return err
	}

	// filter

	// reset the per-block state
	w.dataBlock.reset()
	w.pendingDataBH = bh

	return err
}

func (w *Writer) writeRawBlock(b []byte, compression db.Compression) (blockHandle, error) {
	w.scratch[0] = uint8(compression)
	b = append(b, w.scratch[:1]...)
	crc := util.NewCRC(b).Value()

	binary.LittleEndian.PutUint32(w.scratch[0:4], crc)
	b = append(b, w.scratch[:4]...)

	if _, err := w.writer.Write(b); err != nil {
		log.Println("failed to write data block")
		return blockHandle{}, err
	}

	bh := blockHandle{w.offset, uint64(len(b) - blockTrailerLen)}
	w.offset += uint64(len(b))

	return bh, nil
}

func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}

	if err := w.writeCurrDataBlock(); err != nil {
		log.Println("failed to flush block, ", err)
		w.err = err
		return w.err
	}

	// filter
	metaIndexBH := blockHandle{}

	// index

	// append for the last data block
	w.addIndexEntryForPendingDataBlock(nil)

	w.indexBlock.finish()
	indexBH, err := w.writeRawBlock(w.indexBlock.buf, db.CompressionNo)
	if err != nil {
		w.err = err
		return w.err
	}

	// write table footer
	footer := w.scratch[:footerLen]
	for i := range footer {
		footer[i] = 0
	}
	n := encodeBlockHandle(footer, metaIndexBH)
	encodeBlockHandle(footer[n:], indexBH)
	copy(footer[footerLen-len(magic):], magic)
	if _, err := w.writer.Write(footer); err != nil {
		w.err = err
		return w.err
	}

	w.err = errors.New("writer is closed")

	return nil
}

func NewWriter(f *os.File, o *db.Options) *Writer {
	bhMaxSize := 2 * binary.MaxVarintLen64

	w := &Writer{
		keyCmp:               o.GetComparer(),
		blockRestartInternal: o.GetBlockRestartInterval(),
		blockSize:            o.GetBlockSize(),
		bHScratch:            make([]byte, bhMaxSize),
	}

	w.dataBlock.restartInternal = o.GetBlockRestartInterval()
	w.dataBlock.scratch = w.scratch[0:]

	w.indexBlock.restartInternal = 1
	w.indexBlock.scratch = w.scratch[0:]

	w.writer = f

	return w
}
