package blockchain

import (
	"sync"
	"testing"

	"github.com/lelemita/nomadcoin/utils"
)

type fakeDB struct {
	// 테스트 해야할 결과문이 두개 이상이면 fake func 필요
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (fakeDB) SaveBlock(hash string, data []byte) {}
func (fakeDB) SaveChain(data []byte)              {}
func (fakeDB) DeleteAllBlocks()                   {}

func TestBlockchain(t *testing.T) {
	t.Run("Should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}
		bc := Blockchain()
		if bc.Height != 1 {
			t.Errorf("Blockchain() should create a blockchain with a height of %d, got %d.", 1, bc.Height)
		}
	})
	t.Run("Should restore blockchain", func(t *testing.T) {
		// once.Do() 를 다시 한번 하기 위해 새로 만들어줌
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{
					Height:            8,
					NewestHash:        "xxxx",
					CurrentDifficulty: 2,
				}
				return utils.ToBytes(bc)
			},
		}
		bc := Blockchain()
		if bc.Height != 8 {
			t.Errorf("Blockchain() should retore a blockchain with a height of %d, got %d.", 8, bc.Height)
		}
	})
}
