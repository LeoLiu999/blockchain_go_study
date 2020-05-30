package main

import "bytes"

//TXOutput 输出
type TXOutput struct {
	//存储了多少"币" 在比特币中，value 字段存储的是 satoshi 的数量，
	//一个 satoshi 等于一百万分之一的 >BTC(0.00000001 BTC)，这也是比特币里面最小的货币单位。
	Value int

	//用一个数学难题对输出进行锁定 这个难题被存储在 ScriptPubKey 里面 输出的存储才有意义
	//比特币使用了一个叫做 Script 的脚本语言，用它来定义锁定和解锁输出的逻辑。
	PubKeyHash []byte
}

//Lock lock 输出加锁 输入解锁
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
