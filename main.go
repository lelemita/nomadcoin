package main

import (
	"github.com/lelemita/nomadcoin/cli"
	"github.com/lelemita/nomadcoin/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
