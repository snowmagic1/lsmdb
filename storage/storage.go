package storage

import "io"

type Reader interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
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

type FileDesc struct {
	Type FileType
	Num  int64
}

type Storage interface {
	Lock() (Locker, error)

	Open(fd FileDesc) (Reader, error)

	Close() error
}
