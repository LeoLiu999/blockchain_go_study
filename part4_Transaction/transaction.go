package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

//subsidy 奖励的数额。
//在比特币中，实际并没有存储这个数字，而是基于区块总数进行计算而得：
//区块总数除以 210000 就是 subsidy。挖出创世块的奖励是 50 BTC，每挖出 210000 个块后，奖励减半。
const subsidy = 10

//对于每一笔新的交易，它的输入会引用之前一笔交易的输出
//（这里有个例外，也就是coinbase 交易）。
//所谓引用之前的一个输出，也就是将之前的一个输出包含在另一笔交易的输入当中。交易的输出，也就是币实际存储的地方。
//一个输入必须引用一个输出
//一个输入可以引用之前的多笔输出
//没有被引用为输入的输出 即为余额

//Transaction 每一笔交易必须至少有一笔输入 一笔输出 输入中又必须包含其他笔交易的输出 输入来源必须包含输出（coinbase 交易例外）
type Transaction struct {
	ID   []byte
	Vin  []TXInput  //输入
	Vout []TXOutput //输出
}

//NewCoinbaseTX coinbase交易
//oinbase 交易是一种特殊的交易，它的输入不需要引用之前一笔交易的输出。
//它“凭空”产生了币（也就是产生了新币），这也是矿工获得挖出新块的奖励，可以理解为“发行新币”。
func NewCoinbaseTX(to, data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("Reward to '%s'\n", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}

	tx := Transaction{[]byte{}, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

//SetID SetID
func (tx *Transaction) SetID() {

	var hash [32]byte

	txCopy := *tx

	hash = sha256.Sum256(txCopy.Serialize())

	tx.ID = hash[:]

}

//Serialize Serialize
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)

	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

//IsCoinbase coinbase交易的输入 不需要包含前一笔的输出
//输入的Vin只有1个 Vin.txid 交易id为0 并且 Vin.Vout = -1  输入的输出索引为-1
func (tx *Transaction) IsCoinbase() bool {

	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1

}

func newUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	//找到未消费的交易输出 并保证余额足够支付
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Println("Not enough funds")
		os.Exit(1)
	}

	//生成输入列表
	for txid, outs := range validOutputs {

		txidByte, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txidByte, out, from}
			inputs = append(inputs, input)
		}

	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		//如果存在余额 把余额还回去
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx

}
