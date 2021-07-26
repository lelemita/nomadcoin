package main

import (
	"github.com/lelemita/nomadcoin/explorer"
	"github.com/lelemita/nomadcoin/rest"
)

func main() {
	go explorer.Start(5000)
	rest.Start(4000)
}
