package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

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

func (b * Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		blockAsString := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			fmt.Printf("Block as String: %s\nHash: %s", blockAsString, hash)
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
		Difficulty: difficulty,
		Nonce: 0,
	}
	block.mine()
	// Block을 []byte로 바꿔서 DB에 저장
	block.persist()
	return &block
}

