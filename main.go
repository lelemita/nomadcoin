package main

import (
	"github.com/lelemita/nomadcoin/blockchain"
	"github.com/lelemita/nomadcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}