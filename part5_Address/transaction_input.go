package main

import "bytes"

//TXInput 输入
type TXInput struct {
	//存储这笔交易的 ID
	Txid []byte
	//存储这笔交易中所有输出的索引 一笔输入必须有之前的输出 这里有个例外，也就是coinbase 交易
	Vout      int
	Signature []byte
	PubKey    []byte
}

//UsesKey 解锁  输出加锁 输入解锁
//unspent transactions outputs, UTXO 未花费的交易输出
//当检查余额时，并不需要知道整个区块链上所有的 UTXO，
//只需要关注那些能够解锁的那些 UTXO
//解锁该笔输入的输出
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
