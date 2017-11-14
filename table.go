package lsmdb

import (
	"github.com/snowmagic1/lsmdb/iterator"
	"github.com/snowmagic1/lsmdb/storage"
	"github.com/snowmagic1/lsmdb/table"
)

type tFile struct {
	fd         storage.FileDesc
	seekLeft   int32
	size       int64
	imin, imax internalKey
}

func newTableFile(fd storage.FileDesc, size int64, imin, imax internalKey) *tFile {
	f := &tFile{
		fd:   fd,
		size: size,
		imin: imin,
		imax: imax,
	}

	f.seekLeft = int32(size / 16384)
	if f.seekLeft < 100 {
		f.seekLeft = 100
	}

	return f
}

type tOps struct {
	s      *session
	noSync bool
}

func (t *tOps) create() (*tWriter, error) {
	fd := storage.FileDesc{storage.TypeTable, t.s.allocFileNum()}
	fw, err := t.s.stor.Create(fd)
	if err != nil {
		return nil, err
	}

	return &tWriter{
		t:  t,
		fd: fd,
		w:  fw,
	}, nil
}

func (t *tOps) createFrom(src iterator.Iterator) (f *tFile, n int, err error) {
	w, err := t.create()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			w.drop()
		}
	}()

	for src.Next() {
		err = w.append(src.Key(), src.Value())
		if err != nil {
			return
		}
	}

	err = src.Error()
	if err != nil {
		return
	}

	n = w.tw.EntriesLen()
	f, err = w.finish()

	return
}

type tWriter struct {
	t           *tOps
	fd          storage.FileDesc
	w           storage.Writer
	tw          *table.Writer
	first, last []byte
}

func (w *tWriter) append(key, val []byte) error {
	if w.first == nil {
		w.first = append([]byte{}, key...)
	}
	w.last = append(w.last[:0], key...)

	return w.tw.Append(key, val)
}

func (w *tWriter) drop() {

}

func (w *tWriter) finish() (f *tFile, err error) {
	defer w.close()
	err = w.tw.Close()
	if err != nil {
		return
	}

	if !w.t.noSync {
		err = w.w.Sync()
		if err != nil {
			return
		}
	}

	f = newTableFile(w.fd, int64(w.tw.BytesLen()), internalKey(w.first), internalKey(w.last))

	return
}

func (w *tWriter) close() {

}
