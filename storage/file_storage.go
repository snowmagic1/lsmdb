package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
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
	path     string
	readOnly bool
	flock    fileLock
	logw     *os.File
	logSize  int64
}

func OpenFile(path string, readOnly bool) (Storage, error) {
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
		path:     path,
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
