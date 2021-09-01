package blockchain

import (
	"reflect"
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

func TestBlocks(t *testing.T) {
	fakeBlock := 0
	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: ""},
	}
	t.Run("", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				defer func() {
					fakeBlock++
				}()
				return utils.ToBytes(blocks[fakeBlock])
			},
		}
		bc := &blockchain{}
		blocksResult := Blocks(bc)
		if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
			t.Error("Blocks() should return a slice of blocks.")
		}
	})
}

func TestFindTx(t *testing.T) {
	t.Run("Tx not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height:       2,
					Transactions: []*Tx{},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{
			NewestHash: "xxxx",
		}, "test")
		if tx != nil {
			t.Error("Tx should be not found.")
		}
	})
	t.Run("Tx should be found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height:       2,
					Transactions: []*Tx{{Id: "test"}},
				}
				return utils.ToBytes(b)
			},
		}
		tx := FindTx(&blockchain{
			NewestHash: "xxxx",
		}, "test")
		if tx == nil {
			t.Error("Tx should be found.")
		}
	})
}

// table test
func TestGetDifficulty(t *testing.T) {
	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func() { fakeBlock++ }()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}
	type test struct {
		height int
		want   int
	}
	tests := []test{
		{height: 0, want: defaultDifficulty},
		{height: 2, want: defaultDifficulty},
		{height: difficultyIntterval, want: defaultDifficulty + 1},
	}
	for _, tc := range tests {
		bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDifficulty}
		di := getDifficulty(bc)
		if di != tc.want {
			t.Errorf("getDifficulty{} should %d got %d.", tc.want, di)
		}
	}
}

func TestAddPeerBlock(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	mem.Txs["test_tx"] = &Tx{}
	nb := &Block{
		Difficulty:   2,
		Hash:         "test",
		Transactions: []*Tx{{Id: "test_tx"}},
	}
	bc.AddPeerBlock(nb)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock() should mutate the blockchain.")
	}
}

func TestReplace(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	blocks := []*Block{{Difficulty: 3, Hash: "nb_hash"}, {}, {}}
	bc.Replace(blocks)
	if bc.CurrentDifficulty != 3 || bc.Height != 3 || bc.NewestHash != "nb_hash" {
		t.Error("Replace() should mutate the blockchain.")
	}
}
