package main

//TXOutput 输出
type TXOutput struct {
	//存储了多少"币" 在比特币中，value 字段存储的是 satoshi 的数量，
	//一个 satoshi 等于一百万分之一的 >BTC(0.00000001 BTC)，这也是比特币里面最小的货币单位。
	Value int

	//用一个数学难题对输出进行锁定 这个难题被存储在 ScriptPubKey 里面 输出的存储才有意义
	//比特币使用了一个叫做 Script 的脚本语言，用它来定义锁定和解锁输出的逻辑。
	//由于还没有实现地址（address） ScriptPubKey 将会存储一个任意的字符串（用户定义的钱包地址）。
	ScriptPubKey string
}

//CanBeUnlockedWith 解锁该笔交易的输出
//目前还没有实现密钥，所以将会使用用户定义的地址来代替
func (txout *TXOutput) CanBeUnlockedWith(unlockedData string) bool {
	return txout.ScriptPubKey == unlockedData
}
