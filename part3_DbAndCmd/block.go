package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	PrevBlockHash []byte
	Data []byte
	Hash []byte
	Nonce int64
}

func (b *Block) Serialize() []byte{

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func NewBlock(data string, prevBlockHash []byte) *Block{

	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}, 0}
	//工作量证明
	pow := NewProofOfWork(block)

	nonce, hash := pow.run()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}
//生成创世块
func newGenesisBlock() *Block{

	return NewBlock("Genesis Block", []byte{})
}
//很奇怪这里为什么不能使用Unserialize作为函数名
func Deserialize(d []byte) *Block{

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}