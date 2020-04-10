package main

import (
	"time"
)

type Block struct {
	timestamp int64
	prevBlockHash []byte
	data []byte
	hash []byte
	nonce int64
}

func NewBlock(data string, prevBlockHash []byte) *Block{

	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data), []byte{}, 0}
	//工作量证明
	pow := NewProofOfWork(block)

	nonce, hash := pow.run()

	block.nonce = nonce
	block.hash = hash[:]

	return block
}




