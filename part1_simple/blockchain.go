package main

type Blockchain struct {

	blocks []*Block

}
//获取最后一个块的hash生成新的块
func (bc *Blockchain) addBlock(data string){

	prevBlock := bc.blocks[len(bc.blocks) -1]

	newBlock := NewBlock(data, prevBlock.hash)

	bc.blocks = append(bc.blocks, newBlock)

}
//生成创世块
func createGenesisBlock() *Block{

	return NewBlock("Genesis Block", []byte{})
}

//新建区块链
func newBlockchain() *Blockchain{

	return &Blockchain{[]*Block{createGenesisBlock()}}

}



