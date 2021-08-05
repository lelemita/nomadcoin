package blockchain

import (
	"errors"
	"time"

	"github.com/lelemita/nomadcoin/utils"
	"github.com/lelemita/nomadcoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

// 해당 Amount(Unspent Tx)를 생성한 TxOut을 찾는 방법
// 이 TxIn을 만든사람이 정말 그 TxOut의 주인인지를 TxIn.Signature와 TxOut.Address로 Verify
type TxIn struct {
	TxId      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxId   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

// Tx.Ins 안의 모든 TxIn에 Tx.Id로 Sign
func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	isValid := true
	for _, txIn := range tx.TxIns {
		// txIn.TxId로 txIn을 txOut한 Transaction을 찾는다
		prevTx := FindTx(Blockchain(), txIn.TxId)
		if prevTx == nil {
			isValid = false
			break
		}
		isValid = wallet.Verify(txIn.Signature, tx.Id, tx.TxOuts[txIn.Index].Address)
		if !isValid {
			break
		}
	}
	return isValid
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

// 승인할 트랜젝션들 가져오고, mempool 비우기
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

// TxIns에 해당 TxOut가 있는 Tx가 Mempool에 있는지 확인
func isOnMempool(uTxOut *UTxOut) (isExists bool) {
Outer:
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxId == uTxOut.TxId && txIn.Index == uTxOut.Index {
				isExists = true
				break Outer
			}
		}
	}
	return
}

// 돈 찍기: 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx포인터를 반환
func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	return &tx
}

var ErrorNoMoney = errors.New("not enough money")
var ErrorNotValid = errors.New("Tx Invalid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTx := range uTxOuts {
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{uTx.TxId, uTx.Index, from})
		total += uTx.Amount
	}
	// 잔돈 계산
	if change := total - amount; change != 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	if !validate(&tx) {
		return nil, ErrorNotValid
	}
	return &tx, nil
}
