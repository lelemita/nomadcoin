package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/lelemita/nomadcoin/db"
	"github.com/lelemita/nomadcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transaction"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("Block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockByte := db.Block(hash)
	if blockByte == nil {
		return nil, ErrNotFound
	} else {
		block := &Block{}
		block.restore(blockByte)
		return block, nil
	}
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce += 1
		}
	}
}

func createBlock(prevHash string, height int, diff int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	block.mine()
	// mining이 오래걸리므로, 트랜젝션은 그 뒤에 더함
	block.Transactions = Mempool.TxToConfirm()
	// Block을 []byte로 바꿔서 DB에 저장
	block.persist()
	return &block
}
