package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bi *BlockchainIterator) Next() *Block{

	var block *Block

	err := bi.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket) )
		encodeBlock := b.Get(bi.currentHash)

		block = Deserialize(encodeBlock)
		return nil
	})

	if err != nil{
		log.Panic(err)
	}

	bi.currentHash = block.PrevBlockHash

	return block
}
