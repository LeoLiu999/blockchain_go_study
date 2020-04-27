package main

//TXInput 输入
type TXInput struct {
	//存储这笔交易的 ID
	Txid []byte
	//存储这笔交易中所有输出的索引 一笔输入必须有之前的输出 这里有个例外，也就是coinbase 交易
	Vout int
	//signature
	//ScriptSig 是一个脚本，提供了可作用于一个输出的 ScriptPubKey 的数据。
	//如果 ScriptSig 提供的数据是正确的，那么输出就会被解锁，然后被解锁的值就可以被用于产生新的输出；
	//如果数据不正确，输出就无法被引用在输入中，或者说，也就是无法使用这个输出。这种机制，保证了用户无法花费属于其他人的币。
	//由于还没有实现地址（address），所以 ScriptSig 将仅仅存储一个任意用户定义的钱包地址。将在后续实现公钥（public key）和签名（signature）。
	ScriptSig string
}

//CanUnlockOutputWith 解锁
//unspent transactions outputs, UTXO 未花费的交易输出
//当检查余额时，并不需要知道整个区块链上所有的 UTXO，
//只需要关注那些能够解锁的那些 UTXO
//目前还没有实现密钥，所以将会使用用户定义的地址来代替。
//解锁该笔输入的输出
func (txin *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return txin.ScriptSig == unlockingData
}
