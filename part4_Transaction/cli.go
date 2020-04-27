package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//CLI CLI
type CLI struct {
}

//Run Run
func (cli *CLI) Run() {

	cli.validateArgs()

	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	getBalanceAddress := getBalanceCmd.String("address", "", "Get balance of address")
	from := sendCmd.String("from", "", "From address")
	to := sendCmd.String("to", "", "To address")
	amount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "printchain":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.Usage()
		os.Exit(1)
	}

	if printchainCmd.Parsed() {
		cli.PrintChain()
	}

	if createBlockchainCmd.Parsed() {

		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.createBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed() {

		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceAddress)

	}

	if sendCmd.Parsed() {

		if *from == "" || *to == "" || *amount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*from, *to, *amount)

	}

}

//validateArgs validateArgs
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.Usage()
		os.Exit(1)
	}
}

//Usage Usage
func (cli *CLI) Usage() {

	fmt.Println("Usage:")
	fmt.Println("createBlockchain -address ADDRESS, the address to send genesis block reward")
	fmt.Println("printchain, print blockchain")
	fmt.Println("getbalance -address ADDRESS, Get balance of address")
	fmt.Println("send -from ADDRESS -to ADDRESS -amount AMOUNT, Send AMOUNT of coins from FROM address to TO")

}
