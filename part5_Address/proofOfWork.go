package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24 //定义挖矿的难度 hash前24位为0
const maxNonce = math.MaxInt64

//ProofOfWork 工作量证明 工作量证明的核心是计算区块链中块的hash值 有了工作量证明之后让挖矿成为可能
type ProofOfWork struct {
	block  *Block
	target *big.Int //hashInt 为big.Int 需要和big.Int比较 这里需要设置为big.Int
}

//计算区块的hash值 返回hash值、计数器
//使用块中的数据+计数器 获取hash值
func (pow *ProofOfWork) run() (int64, []byte) {

	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

	fmt.Printf("开始挖矿\n")

	for nonce <= maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		//将hash转化为int
		hashInt.SetBytes(hash[:])

		//比较hash与target的值 比target小 hash才有效
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}

	}
	return nonce, hash[:]
}

//链接块中的数据 将块中的数据转换成byte切片数组
func (pow *ProofOfWork) prepareData(nonce int64) []byte {

	return bytes.Join(
		[][]byte{
			IntToHex(pow.block.Timestamp),
			//块里存储了交易信息 不是Data了
			//pow.block.Data,
			pow.block.HashTransactions(),
			pow.block.PrevBlockHash,
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce))}, []byte{})

}

//Validate 对工作量证明进行验证
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	//计算hash
	hash := sha256.Sum256(data)
	//比较hash
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

//NewProofOfWork NewProofOfWork
func NewProofOfWork(block *Block) *ProofOfWork {

	target := big.NewInt(1)
	//左移256-targetBits位 hash为sum256算法生产的256位
	//最终target变为0000010000000000000000000000000000000000000000000000000000000000
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block, target}

	return pow
}
