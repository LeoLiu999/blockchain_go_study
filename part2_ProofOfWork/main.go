package main

import (
	"fmt"
	"strconv"
)

func main(){

	bc := newBlockchain()
	bc.addBlock("send 1 BTC to A")
	bc.addBlock("send 2 BTC to B")
	bc.addBlock("send 3 BTC to C")

	for _,block:= range(bc.blocks) {

		fmt.Printf("timestamp:%d\n", block.timestamp)
		fmt.Printf("data:%s\n", block.data)
		fmt.Printf("hash:%x\n", block.hash)
		fmt.Printf("prevBlockHash:%x\n", block.prevBlockHash)
		fmt.Printf("nonce:%d\n", block.nonce)

		pow := NewProofOfWork(block)
		validate := pow.Validate()
		fmt.Printf("validate:%s\n", strconv.FormatBool(validate) )

		fmt.Println()

	}


}
