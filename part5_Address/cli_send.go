package main

import "fmt"

func (cli *CLI) send(from, to string, amount int) {

	bc := newBlockchain()

	defer bc.db.Close()

	//新建交易
	tx := newUTXOTransaction(from, to, amount, bc)

	bc.MineBlock([]*Transaction{tx})

	fmt.Println("success")
}
