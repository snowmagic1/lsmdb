package wallog

import (
	"bufio"
	"log"
	"os"
)

type LogRW struct {
	currIndex  uint64
	fileWriter *bufio.Writer
	fileReader *bufio.Reader
}

func NewLogRW(file *os.File) *LogRW {
	logRW := &LogRW{}
	logRW.fileWriter = bufio.NewWriter(file)

	return logRW
}

func (_logRW *LogRW) Write(key, val string) {

	_logRW.currIndex++

	entry := &logRecord{
		index: _logRW.currIndex,
		key:   key,
		val:   val}

	buf := entry.toBytes()
	_, err := _logRW.fileWriter.Write(buf)

	if err != nil {
		log.Println("failed to write ", err)
		return
	}

	newEntry := fromBytes(buf)
	log.Println("newEntry ", newEntry)
}
