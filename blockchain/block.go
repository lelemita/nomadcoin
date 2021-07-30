package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/lelemita/nomadcoin/db"
	"github.com/lelemita/nomadcoin/utils"
)

type Block struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height int `json:"height"`
	Difficulty int `json:"difficulty"`
	Nonce int `json:"nonce"`
	Timestamp int `json:"timestamp`
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

func (b * Block) mine() {
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

func createBlock(data string, prevHash string, height int) *Block{
	block := Block{
		Data: data, 
		Hash: "", 
		PrevHash: prevHash,
		Height: height,
		Difficulty: Blockchain().difficulty(),
		Nonce: 0,
	}
	block.mine()
	// Block을 []byte로 바꿔서 DB에 저장
	block.persist()
	return &block
}

