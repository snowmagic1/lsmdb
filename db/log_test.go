package db_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/snowmagic1/lsmdb/db"
)

var logFile *os.File
var writer *db.LogWriter
var reader *db.LogReader

const logFileName = "test.log"

func init() {
	var err error

	os.Remove(logFileName)
	logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("failed to open test file ", err)
	}

	writer = db.NewLogWriter(logFile)
	reader = db.NewLogReader(logFile)
}

func testWrite(record []byte, t *testing.T) {
	if err := writer.AddRecord(record); err != nil {
		t.Errorf("failed to append, %v", err)
	}
}

func testRead(expected []byte, t *testing.T) {
	if read, err := reader.ReadRecord(); err != nil {
		t.Errorf("failed to read log, %v", err)
	} else if !reflect.DeepEqual(read, expected) {
		t.Errorf("log record mismatch, expect [%v] actual [%v]", expected, read)
	}
}
func TestSingleRecord(t *testing.T) {
	record := []byte{1, 2, 3, 4, 5}
	testWrite(record, t)
	testRead(record, t)
}

func TestMultiRecord(t *testing.T) {
	record1 := []byte{11, 2, 3, 4, 5}
	record2 := []byte{33, 21, 223, 0xa4, 5}
	testWrite(record1, t)
	testWrite(record2, t)
	testRead(record1, t)
	testRead(record2, t)
}
