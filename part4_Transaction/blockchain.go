package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

//Blockchain 区块链对象
type Blockchain struct {
	tip []byte //最后一个块的hash值
	db  *bolt.DB
}

//FindUnspentTransations 通过地址找到区块链中所有未花费的交易输出
func (bc *Blockchain) FindUnspentTransations(address string) []Transaction {

	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int)

	bci := bc.Iterator()

	for {

		block := bci.Next()
		for _, tx := range block.Transactions {

			txID := hex.EncodeToString(tx.ID)
		Outputs:
			//outIdx Vout 存储的输出索引
			for outIdx, out := range tx.Vout {

				//查询该交易是否已被花费
				if spentTXOs[txID] != nil {

					for _, spentOut := range spentTXOs[txID] {
						//核对该笔交易是否被花费
						if spentOut == outIdx {
							continue Outputs
						}

					}

				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}

				if tx.IsCoinbase() == false {

					for _, in := range tx.Vin {

						if in.CanUnlockOutputWith(address) {

							inTxID := hex.EncodeToString(in.Txid)
							//只要引用到输入中的输出 都已经被消费 将当前transaction id作为索引 这个输入中的输出存入spentTXOs
							spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)

						}

					}

				}

			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	return unspentTXs
}

//FindUTXO 返回未花费的output
func (bc *Blockchain) FindUTXO(address string) []TXOutput {

	var UTXOs []TXOutput
	//找到未花费的交易
	unspentTransactions := bc.FindUnspentTransations(address)

	for _, unTx := range unspentTransactions {

		for _, out := range unTx.Vout {

			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}

	}

	return UTXOs
}

//FindSpendableOutputs 找到所有的未花费输出，并且确保它们存储了足够的值
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {

	unspentOutputs := make(map[string][]int)

	unspentTXs := bc.FindUnspentTransations(address)

	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txid := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {

			if out.CanBeUnlockedWith(address) && accumulated < amount {

				accumulated += out.Value
				unspentOutputs[txid] = append(unspentOutputs[txid], outIdx)

				if accumulated >= amount {

					break Work
				}
			}

		}
	}

	return accumulated, unspentOutputs

}

//MineBlock 挖出新块 发送币意味着创建新的交易 需要通过挖出新块的方式将交易打包到区块链中
//完成交易后需要添加一个新块到区块链中
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	block := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(block.Hash, block.Serialize())

		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)

		if err != nil {
			log.Panic(err)
		}

		bc.tip = block.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block

}

//FindTransaction 从链中找到transaction
func (bc *Blockchain) FindTransaction(txid []byte) (Transaction, error) {

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {

			if bytes.Compare(txid, tx.ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	return Transaction{}, errors.New("Transaction is not found")
}

//Iterator 区块链迭代器 通过最后一个块的hash开始迭代
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

//CreateBlockchain 新建区块链
func CreateBlockchain(address string) *Blockchain {

	if DbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesisBlock := newGenesisBlock(cbtx)
	//放入数据库
	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesisBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

//newBlockchain 返回区块链对象
func newBlockchain() *Blockchain {

	if !DbExists(dbFile) {
		fmt.Println("Blockchain is not exists,create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}

}
