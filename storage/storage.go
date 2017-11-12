package storage

import (
	"errors"
	"io"
)

type Syncer interface {
	Sync() error
}

type Reader interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
}

type Writer interface {
	io.WriteCloser
	Syncer
}

type Locker interface {
	Unlock()
}

type FileType int

const (
	TypeManifest FileType = 1 << iota
	TypeJournal
	TypeTable
	TypeTemp

	TypeAll = TypeManifest | TypeJournal | TypeTable | TypeTemp
)

var (
	ErrInvalidFile = errors.New("leveldb/storage: invalid file for argument")
	ErrLocked      = errors.New("leveldb/storage: already locked")
	ErrClosed      = errors.New("leveldb/storage: closed")
)

type FileDesc struct {
	Type FileType
	Num  int64
}

func FileDescOk(fd FileDesc) bool {
	switch fd.Type {
	case TypeManifest:
	case TypeJournal:
	case TypeTable:
	case TypeTemp:
	default:
		return false
	}

	return fd.Num > 0
}

type Storage interface {
	Lock() (Locker, error)

	Open(fd FileDesc) (Reader, error)

	Create(fd FileDesc) (Writer, error)

	Close() error
}
