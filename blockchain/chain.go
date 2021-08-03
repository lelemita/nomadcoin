package blockchain

import (
	"sync"

	"github.com/lelemita/nomadcoin/db"
	"github.com/lelemita/nomadcoin/utils"
)

const (
	defaultDifficulty int = 2
	difficultyIntterval int = 5
	// 2분/1블록 목표.
	blockInterval int = 2
	allowedRange int = 2
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
}

// singleton pattern: only one instance
var b *blockchain
// 딱 한번 실행되도록 하기 (goroutin, thread 가 여러개여도..)
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock(){
	block := createBlock(b.NewestHash, b.Height + 1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

// struct를 mutate시키지 않으므로 단순 function으로 구현
func Blocks(b *blockchain) []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor) 
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newest := allBlocks[0]
	lastCalculated := allBlocks[difficultyIntterval - 1]
	actualTime := (newest.Timestamp - lastCalculated.Timestamp) / 60
	expectedTime := blockInterval * difficultyIntterval 
	if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	} else if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	}
	return b.CurrentDifficulty
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % difficultyIntterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

// 한 블록안에서, 특정인의, TxOuts에는 있고 TxIns에는 없는 TxOuts 목록
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	// Unspent Transaction
	var uTxOuts []*UTxOut
	// spent transaction outputs
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			// input으로 사용된 tx들 찾기
			for _, txIn := range tx.TxIns {
				if txIn.Owner == address {
					creatorTxs[txIn.TxId] = true
				}
			}
			// outs 중에서 input에 없는 것들 찾기
			for idx, txOut := range tx.TxOuts {
				if txOut.Owner == address {
					// _: value / ok: 해당 키 값의 존재여부
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, idx, txOut.Amount}
						// 이미 mempool에 있는 Tx 값은 쓰면 안된다.
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

// 총 자산
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int = 0
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	// 딱 한번 생성
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
	})
	// fmt.Printf("Height: %d\nNewest Hash: %s\n", b.Height, b.NewestHash)
	return b
}
