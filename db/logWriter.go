package db

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"log"
	"os"
)

type LogWriter struct {
	destfile *bufio.Writer
	crc32q   *crc32.Table
}

func NewLogWriter(file *os.File) *LogWriter {
	logWriter := &LogWriter{}
	logWriter.destfile = bufio.NewWriter(file)
	logWriter.crc32q = crc32.MakeTable(0xD5828281)

	return logWriter
}

func (logWriter *LogWriter) AddRecord(payload []byte) error {

	// TODO: block support

	return logWriter.emitPhysicalRecord(payload)
}

func check(err error) {
	if err != nil {
		log.Println("error - ", err)
	}
}

func (logWriter *LogWriter) emitPhysicalRecord(record []byte) error {

	crc := crc32.Checksum(record, logWriter.crc32q)

	// checksum | length | payload
	buf := new(bytes.Buffer)
	check(binary.Write(buf, binary.LittleEndian, crc))
	check(binary.Write(buf, binary.LittleEndian, uint16(len(record))))
	check(binary.Write(buf, binary.LittleEndian, record))

	if _, err := logWriter.destfile.Write(buf.Bytes()); err != nil {
		log.Println("failed to write log, ", err)
	} else {
		logWriter.destfile.Flush()
	}

	return nil
}
