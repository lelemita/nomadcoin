package blockchain

import (
	"testing"

	"github.com/lelemita/nomadcoin/utils"
)

type fakeWallet struct {
	fakeVerify func() bool
}

func (f fakeWallet) Verify(signature, payload, address string) bool {
	return f.fakeVerify()
}
func (f fakeWallet) Sign(payload string) string {
	return "fakeSign"
}
func (f fakeWallet) GetAddress() string {
	return "xxxx"
}

func TestMakeTx(t *testing.T) {
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			b := &Block{
				Transactions: []*Tx{{
					TxOuts: []*TxOut{{Address: "from", Amount: 50}, {Address: "from", Amount: 50}},
				}}}
			return utils.ToBytes(b)
		},
	}
	b = &blockchain{NewestHash: "xxxx"}

	t.Run("not enough money in account", func(t *testing.T) {
		_, err := makeTx("from", "to", 200)
		if err != ErrorNoMoney {
			t.Error("makeTx() should raise Error when money is not enough.")
		}

	})

	myWallet = fakeWallet{
		fakeVerify: func() bool { return false },
	}
	t.Run("When wrong signature", func(t *testing.T) {
		_, err := makeTx("from", "to", 30)
		if err != ErrorNotValid {
			t.Error("makeTx() should raise ErrorNotValid for wrong signature.")
		}

	})

	myWallet = fakeWallet{
		fakeVerify: func() bool { return true },
	}
	t.Run("When there is changes", func(t *testing.T) {
		tx, _ := makeTx("from", "to", 30)
		if tx.TxOuts[0].Address != "from" && tx.TxOuts[0].Amount != 70 {
			t.Errorf("the change should 70, but got %d", tx.TxOuts[0].Amount)
		}
	})
	// t.Run("When there is no changes", func(t *testing.T) {
	// 	tx, _ := makeTx("from", "to", 100)
	// 	if tx.TxOuts[0].Address != "from" && tx.TxOuts[0].Amount != 100 {
	// 		t.Errorf("the change should 100, but got %d", tx.TxOuts[0].Amount)
	// 	}
	// })

}
