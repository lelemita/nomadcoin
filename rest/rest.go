package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lelemita/nomadcoin/blockchain"
	"github.com/lelemita/nomadcoin/p2p"
	"github.com/lelemita/nomadcoin/utils"
	"github.com/lelemita/nomadcoin/wallet"
)

var port string

type url string

// type TextMarshaler interface
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type addPeerPayload struct {
	Address, Port string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "message:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgread to WebSockets",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)

	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	block, err := blockchain.FindBlock(vars["hash"])
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// Middleware: All request에 Content-type을 설정
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Conetext-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

// ?total=true 이면 총액 반환 / 그 외는 목록
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	isTotal := r.URL.Query().Get("total")
	encoder := json.NewEncoder(rw)
	switch isTotal {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		utils.HandleErr(encoder.Encode(balanceResponse{address, amount}))
	default:
		utils.HandleErr(encoder.Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	// Decoder가 r.Body의 내용을 paload struct로 변환
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	go p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
	// 익명 구조체
	// json.NewEncoder(rw).Encode(struct {
	// 	Address string `json:"address"`
	// }{Address: address})
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	ip_port := fmt.Sprintf("127.0.0.1%s", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://%s\n", ip_port)
	log.Fatal(http.ListenAndServe(ip_port, router))
}
