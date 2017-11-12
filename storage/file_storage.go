package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

var (
	errFileOpen = errors.New("leveldb/storage: file still open")
	errReadOnly = errors.New("leveldb/storage: storage is read-only")
)

type fileLock interface {
	release() error
}

type unixFileLock struct {
	f *os.File
}

func (fl *unixFileLock) release() error {
	return nil
}

type fileStorage struct {
	rootpath string
	readOnly bool
	mu       sync.Mutex
	flock    fileLock
	logw     *os.File
	logSize  int64
	open     int
}

func OpenDiskStorage(path string, readOnly bool) (Storage, error) {
	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return nil, fmt.Errorf("leveldb/storage: %s is not a directory", path)
		}
	} else if os.IsNotExist(err) && !readOnly {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	flock, err := newFileLock(filepath.Join(path, "LOCK"), readOnly)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			flock.release()
		}
	}()

	var (
		logw    *os.File
		logSize int64
	)

	if !readOnly {
		logw, err = os.OpenFile(filepath.Join(path, "LOG"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		logSize, err = logw.Seek(0, os.SEEK_END)
		if err != nil {
			logw.Close()
			return nil, err
		}
	}

	fs := &fileStorage{
		rootpath: path,
		readOnly: readOnly,
		flock:    flock,
		logw:     logw,
		logSize:  logSize,
	}

	runtime.SetFinalizer(fs, (*fileStorage).Close)

	return fs, nil
}

func (fs *fileStorage) Open(fd FileDesc) (Reader, error) {
	return nil, nil
}

func (fs *fileStorage) Close() error {
	return nil
}

func (fs *fileStorage) Lock() (Locker, error) {
	return nil, nil
}

func newFileLock(path string, readonly bool) (fl fileLock, err error) {
	var flag int
	if readonly {
		flag = os.O_RDONLY
	} else {
		flag = os.O_RDWR
	}

	f, err := os.OpenFile(path, flag, 0)
	if os.IsNotExist(err) {
		f, err = os.OpenFile(path, flag|os.O_CREATE, 0644)
	}

	if err != nil {
		return
	}

	if err = setFileLock(f, readonly, true); err != nil {
		f.Close()
		return
	}

	fl = &unixFileLock{f: f}

	return
}

func setFileLock(f *os.File, readonly, lock bool) error {
	how := syscall.LOCK_UN
	if lock {
		if readonly {
			how = syscall.LOCK_SH
		} else {
			how = syscall.LOCK_EX
		}
	}

	return syscall.Flock(int(f.Fd()), how|syscall.LOCK_NB)
}

func fsGenName(fd FileDesc) string {
	switch fd.Type {
	case TypeManifest:
		return fmt.Sprintf("MANIFEST-%06d", fd.Num)
	case TypeJournal:
		return fmt.Sprintf("%06d.log", fd.Num)
	case TypeTable:
		return fmt.Sprintf("%06d.ldb", fd.Num)
	case TypeTemp:
		return fmt.Sprintf("%06d.tmp", fd.Num)
	default:
		panic("invalid file type")
	}
}

func (fs *fileStorage) Create(fd FileDesc) (Writer, error) {
	if !FileDescOk(fd) {
		return nil, ErrInvalidFile
	}

	if fs.readOnly {
		return nil, errReadOnly
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.open < 0 {
		return nil, ErrClosed
	}

	fileName := filepath.Join(fs.rootpath, fsGenName(fd))
	of, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	fs.open++

	return &fileWrap{File: of, fs: fs, fd: fd}, nil
}

type fileWrap struct {
	*os.File
	fs     *fileStorage
	fd     FileDesc
	closed bool
}
