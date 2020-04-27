package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

//Block 块
type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	//Data []byte
	Transactions []*Transaction //每个块中都应包含一个交易 移除Data字段 存储交易
	Hash         []byte
	Nonce        int64
}

//Serialize 串行化
func (b *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//HashTransactions 通过仅仅一个哈希，就可以识别一个块里面的所有交易。为此，获得每笔交易的哈希，将它们关联起来，然后获得一个连接后的组合哈希。
//比特币使用了一个更加复杂的技术：
//它将一个块里面包含的所有交易表示为一个 Merkle tree ，
//然后在工作量证明系统中使用树的根哈希（root hash）。
//这个方法能够让我们快速检索一个块里面是否包含了某笔交易，即只需 root hash 而无需下载所有交易即可完成判断。
func (b *Block) HashTransactions() []byte {

	var txHashes [][]byte
	var txHash [32]byte

	for _, transaction := range b.Transactions {

		txHashes = append(txHashes, transaction.ID)

	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}

//NewBlock 生成新块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {

	block := &Block{time.Now().Unix(), prevBlockHash, txs, []byte{}, 0}
	//工作量证明
	pow := NewProofOfWork(block)

	nonce, hash := pow.run()

	block.Nonce = nonce
	block.Hash = hash[:]

	fmt.Printf("%x\n", block.Hash)

	return block
}

//newGenesisBlock 生成创世块
func newGenesisBlock(coinbase *Transaction) *Block {

	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//Deserialize 反串行化
func Deserialize(d []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
