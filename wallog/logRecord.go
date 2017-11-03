package wallog

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"log"
	"reflect"
)

type logRecord struct {
	index uint64
	key   string
	val   string
}

func check(err error) {
	if err != nil {
		log.Println("error - ", err)
	}
}

func getRecordSize(record *logRecord) uint16 {

	// size
	size := uint16(reflect.TypeOf(uint16(0)).Size())

	// index
	size += uint16(reflect.TypeOf(record.index).Size())

	// key
	keyBytes := []byte(record.key)
	size += uint16(reflect.TypeOf(uint16(0)).Size())
	size += uint16(len(keyBytes))

	// val
	valBytes := []byte(record.val)
	size += uint16(reflect.TypeOf(uint16(0)).Size())
	size += uint16(len(valBytes))

	// checksum
	size += sha1.Size

	return size
}

func (record *logRecord) toBytes() []byte {

	// key
	keyBytes := []byte(record.key)
	// val
	valBytes := []byte(record.val)

	size := getRecordSize(record)

	h := sha1.New()

	h.Write([]byte(record.key))
	bs := h.Sum(nil)

	buf := new(bytes.Buffer)
	check(binary.Write(buf, binary.LittleEndian, size))
	check(binary.Write(buf, binary.LittleEndian, uint16(len(keyBytes))))
	check(binary.Write(buf, binary.LittleEndian, keyBytes))
	check(binary.Write(buf, binary.LittleEndian, uint16(len(valBytes))))
	check(binary.Write(buf, binary.LittleEndian, valBytes))
	check(binary.Write(buf, binary.LittleEndian, bs))

	return buf.Bytes()
}

func fromBytes(b []byte) *logRecord {

	buf := bytes.NewBuffer(b)

	record := &logRecord{}

	r := bufio.NewReader(bytes.NewReader(buf.Bytes()))
	var err error

	// size, err := readUInt16(r)
	// check(err)

	record.key, err = readString(r)
	check(err)

	record.val, err = readString(r)
	check(err)

	return record
}

func readUInt16(r *bufio.Reader) (uint16, error) {
	b := make([]byte, 2)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func readString(r *bufio.Reader) (string, error) {
	len, err := readUInt16(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, len)
	r.Read(buf)

	return string(buf), nil
}
