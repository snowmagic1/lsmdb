package db

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"os"
)

const LOG_HEADER_SIZE = 6

type LogReader struct {
	logFile *os.File
	crc32q  *crc32.Table
	offset  int64
}

func NewLogReader(file *os.File) *LogReader {
	reader := &LogReader{}
	reader.logFile = file
	reader.crc32q = crc32.MakeTable(0xD5828281)

	return reader
}

type BadRecordError struct{}

func (e *BadRecordError) Error() string {
	return fmt.Sprintf("bad log record")
}

func (logReader *LogReader) ReadRecord() (record []byte, err error) {
	return logReader.readPhysicalRecord()
}

func (logReader *LogReader) readPhysicalRecord() (record []byte, err error) {
	header, err := logReader.Read(LOG_HEADER_SIZE)
	if err != nil {
		log.Println("failed to read log header, err ", err)
		return nil, err
	}

	crc := binary.LittleEndian.Uint32(header[:4])
	l := binary.LittleEndian.Uint16(header[4:])

	record, err = logReader.Read(int(l))
	if err != nil {
		log.Println("failed to read log payload, err ", err)
		return nil, err
	}

	expected := crc32.Checksum(record, logReader.crc32q)
	if expected != crc {
		log.Printf("crc doesn't match, expect [%v] read [%v]\n", expected, crc)
		return nil, new(BadRecordError)
	}

	return
}

func (logReader *LogReader) Read(length int) (b []byte, err error) {
	b = make([]byte, length)
	read, err := logReader.logFile.ReadAt(b, logReader.offset)
	logReader.offset += int64(read)

	if read != length {
		err = new(BadRecordError)
	}

	return
}
