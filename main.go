package main

import (
	"github.com/lelemita/nomadcoin/cli"
	"github.com/lelemita/nomadcoin/db"
)

func main() {
	// 함수 종료될 때 실행되는 명령
	// os.Exit(0) 이 있으면 동작 안된다. --> runtime.Goexit() 으로 대체
	defer db.DB().Close()
	cli.Start()
}