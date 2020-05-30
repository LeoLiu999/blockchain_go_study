package main

import "fmt"

//获取所有未花费的交易输出 计算balance
func (cli *CLI) getBalance(address string) {

	bc := newBlockchain()
	defer bc.db.Close()

	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)
	balance := 0

	for _, out := range UTXOs {

		balance += out.Value
	}

	fmt.Printf("Address:%s,Balance:%d\n", address, balance)
}
