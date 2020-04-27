package main

import (
	"fmt"
	"strconv"
)

//PrintChain 遍历区块链
func (cli *CLI) PrintChain() {
	bc := newBlockchain()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("timestamp:%d\n", block.Timestamp)
		fmt.Printf("prevHash:%x\n", block.PrevBlockHash)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("nonce:%d\n", block.Nonce)

		pow := NewProofOfWork(block)

		fmt.Printf("validate:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}
}
