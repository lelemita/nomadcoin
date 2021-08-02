package blockchain

import (
	"errors"
	"time"

	"github.com/lelemita/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id string `json:"id"`
	Timestamp int `json:"timestamp"`
	TxIns []*TxIn `json:"txIns"`
	TxOuts []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
} 

type TxIn struct {
	Owner string `json:"owner"`
	Amount int `json:"amount"`
}

type TxOut struct {
	Owner string `json:"owner"`
	Amount int `json:"amount"`
}

// 돈 찍기: 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx포인터를 반환
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn {
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut {
		{address, minerReward},
	}
	tx := Tx {
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	beforeFrom := Blockchain().BalanceByAddress(from)
	if beforeFrom < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	change := total - amount
	if change != 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}

	txOuts = append(txOuts, &TxOut{to, amount})
	tx := Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("nico", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}