package blockchain

import (
	"reflect"
	"testing"

	"github.com/lelemita/nomadcoin/utils"
)

func TestCreateBlock(t *testing.T) {
	myWallet = fakeWallet{}
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block NOT found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("xxxx")
		if err == nil {
			t.Error("The block should not be found.")
		}
	})
	t.Run("Block IS found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{Height: 2}
				return utils.ToBytes(b)
			},
		}
		block, _ := FindBlock("xxxx")
		if block.Height != 2 {
			t.Error("The block should be found.")
		}
	})
}
