package util

import "hash/crc32"

var table = crc32.MakeTable(crc32.Castagnoli)

type CRC32 uint32

func NewCRC(b []byte) CRC32 {
	return CRC32(0).Update(b)
}

func (c CRC32) Update(b []byte) CRC32 {
	return CRC32(crc32.Update(uint32(c), table, b))
}

func (c CRC32) Value() uint32 {
	return uint32(c>>15|c<<17) + 0xa282ead8
}
