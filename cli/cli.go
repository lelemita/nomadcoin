package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/lelemita/nomadcoin/explorer"
	"github.com/lelemita/nomadcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:		Set the PORT of the server\n")
	fmt.Printf("-mode:		Choose in 'rest', 'html' and 'all'\n\n")
	os.Exit(0)
} 

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set the port of the server")
	mode := flag.String("mode", "rest", "Choose in 'rest', 'html' and 'all'")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	case "all":
		go explorer.Start(*port + 1)
		rest.Start(*port)
	default:
		usage()
	}
}
