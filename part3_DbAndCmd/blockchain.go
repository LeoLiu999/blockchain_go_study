package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db   *bolt.DB
}

//获取最后一个块的hash生成新的块 并添加进bolt bucket
// bucket-> l:lastHash hash:block
func (bc *Blockchain) addBlock(data string){

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket) )

		lastHash = b.Get( []byte("l") )
		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket( []byte(blocksBucket) )

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		//跟新区块链tip tip = 最后一个块的hash
		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

//区块链迭代器 通过最后一个块的hash开始迭代
func (bc *Blockchain) Iterator() *BlockchainIterator{
	return &BlockchainIterator{bc.tip, bc.db}
}

//新建区块链 使用bolt数据库存储
//获取区块 如果没有 则新建创世块 如果存在则tip=最后一个块的hash

func newBlockchain() *Blockchain{

	var tip []byte
	db,err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}




	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket) )

		//如果不存在bucket 新建创世块并新建bucket 写入l : hash 与 hash : genesis
		if b == nil{
			genesis := newGenesisBlock()

			b, err := tx.CreateBucket( []byte(blocksBucket) )
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize() )
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip  = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}

}



