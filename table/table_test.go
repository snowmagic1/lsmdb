package table

import (
	"os"
	"reflect"
	"testing"
)

func TestTableFooter(t *testing.T) {
	filename := "footertest"
	os.Remove(filename)
	fileForWriter, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("test: failed to open test file %v", err)
	}

	writer := NewWriter(fileForWriter, nil)
	if writer.err != nil {
		t.Errorf("test: failed to create writer %v", writer.err)
	}

	key := []byte{1, 2, 3}
	val := []byte{4, 5, 6}

	writer.Append(key, val)
	writer.Close()
	fileForWriter.Sync()
	fileForWriter.Close()

	fileForReader, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	fi, err := fileForReader.Stat()
	if err != nil {
		t.Errorf("test: cannot get file size %v", err)
	}

	reader, err := NewReader(fileForReader, fi.Size())
	if err != nil || reader.err != nil {
		t.Errorf("test: failed to create reader %v %v", err, reader.err)
	}

	retVal, err := reader.Get(key)
	if !reflect.DeepEqual(retVal, val) {
		t.Errorf("test: Get key return incorrect result, expected [%v] returned [%v]", val, retVal)
	}
}
