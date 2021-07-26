package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lelemita/nomadcoin/blockchain"
	"github.com/lelemita/nomadcoin/utils"
)

const port string = ":4000"

type URL string

type URLDescription struct {
	URL URL `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

type AddBlockBody struct {
	Message string
}

// type TextMarshaler interface
func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL: URL("/"),
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: URL("/blocks"),
			Method: "POST",
			Description: "Add A Block",
			Payload: "message:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks (rw http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var addBlockBody AddBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func main() {
	http.HandleFunc("/", documentation )
	http.HandleFunc("/blocks", blocks )
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
