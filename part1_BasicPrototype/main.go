package main

import "fmt"

func main(){

	bc := newBlockchain()
	bc.addBlock("send 1 BTC to A")
	bc.addBlock("send 2 BTC to B")

	for _, block := range bc.blocks{

		fmt.Printf("timestamp:%d\n", block.timestamp)
		fmt.Printf("prevHash:%x\n", block.prevBlockHash)
		fmt.Printf("data:%s\n", block.data)
		fmt.Printf("hash:%x\n", block.hash)
		fmt.Println()
	}

}

