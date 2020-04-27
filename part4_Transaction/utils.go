package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

//IntToHex IntToHex
func IntToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

//DbExists DbExists
func DbExists(dbFile string) bool {

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
