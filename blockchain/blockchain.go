package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data string
	Hash string
	PrevHash string
}

type blockchain struct {
	// 너무 길어지니까 포인터의 슬라이스로 함
	blocks []*block
}

// singleton pattern: only one instance
var b *blockchain
// 딱 한번 실행되도록 하기 (goroutin, thread 가 여러개여도..)
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	// 초기화 되었는지 확인하고 딱 한번 생성
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks - 1].Hash
}

func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}