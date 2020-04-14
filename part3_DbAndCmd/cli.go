package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) Run(){

	cli.validateArgs()

	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addBlock":
		//解析输入参数
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.Usage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		//如果调用了addBlockCmd 并且未输入参数
		if *addBlockData == ""{
			addBlockCmd.Usage()
			os.Exit(1)
		}

		cli.AddBlock(*addBlockData)
	}

	if printchainCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *CLI) AddBlock(data string) {
	cli.bc.addBlock(data)
	fmt.Println("success")
}

func (cli *CLI) PrintChain() {
	bci := cli.bc.Iterator()

	for   {
		block := bci.Next()

		fmt.Printf("timestamp:%d\n", block.Timestamp)
		fmt.Printf("prevHash:%x\n", block.PrevBlockHash)
		fmt.Printf("data:%s\n", block.Data)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("nonce:%d\n", block.Nonce)

		pow := NewProofOfWork(block)

		fmt.Printf("validate:%s\n", strconv.FormatBool( pow.Validate() ) )
		fmt.Println()


		if  len(block.PrevBlockHash) == 0 {
			break;
		}

	}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2{
		cli.Usage()
		os.Exit(1)
	}
}

func (cli *CLI) Usage(){

	fmt.Println("Usage:")
	fmt.Println("addBlock -data DATA, add block by data")
	fmt.Println("printchain, print blockchain")

}