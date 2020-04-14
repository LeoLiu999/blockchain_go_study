package main

func main(){

	bc := newBlockchain()

	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()

}
