package main

import "fmt"
//新建区块链
func (cli *CLI) createBlockchain(address string) {

	bc := CreateBlockchain(address)
	defer bc.db.Close()

	fmt.Println("success")


}
