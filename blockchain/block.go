package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/lelemita/nomadcoin/db"
	"github.com/lelemita/nomadcoin/utils"
)

// the number of zero in hash's front
// for this, change the value of Nonce
const difficulty int = 2


type Block struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height int `json:"height"`
	Difficulty int `json:"difficulty"`
	Nonce int `json:"nonce"`
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

func createBlock(data string, prevHash string, height int) *Block{
	block := Block{
		Data: data, 
		Hash: "", 
		PrevHash: prevHash,
		Height: height,
	}
	payload := block.Data + block.PrevHash +  fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	// Block을 []byte로 바꿔서 DB에 저장
	block.persist()
	return &block
}

