package main

import (
	"encoding/binary"
	"unsafe"
)

func getEndian() binary.ByteOrder {
	var i int32 = 0x1
	v := (*[4]byte)(unsafe.Pointer(&i))
	if v[0] == 0 {
		return binary.BigEndian
	}

	return binary.LittleEndian
}
