package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/lelemita/nomadcoin/db"
	"github.com/lelemita/nomadcoin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
}

func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

// singleton pattern: only one instance
var b *blockchain
// 딱 한번 실행되도록 하기 (goroutin, thread 가 여러개여도..)
var once sync.Once

func (b *blockchain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(b))
}

func (b *blockchain) AddBlock (data string){
	block := createBlock(data, b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	// 초기화 되었는지 확인하고 딱 한번 생성
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("Height: %d\nNewest Hash: %s\n", b.Height, b.NewestHash)
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("Height: %d\nNewest Hash: %s\n", b.Height, b.NewestHash)
	return b
}
