package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	timestamp int64
	prevBlockHash []byte
	data []byte
	hash []byte
}
//使用当前时间戳  data 上一块的hash  生成当前块的hash
func (b *Block) setHash() {

	timestamp := []byte(strconv.FormatInt(b.timestamp, 10) )

	headers := bytes.Join([][]byte{timestamp, b.data, b.prevBlockHash}, []byte{} )

	hash := sha256.Sum256(headers)

	b.hash = hash[:]

}

func NewBlock(data string, prevBlockHash []byte) *Block{

	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}}

	block.setHash()

	return block
}




