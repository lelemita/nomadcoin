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
	return "fakeAddress"
}

func TestAddTx(t *testing.T) {
	tx := Tx{TxIns: []*TxIn{{TxId: "Tx01", Index: 0}}}
	mem = &mempool{
		Txs: map[string]*Tx{"fakeTx": &tx},
	}
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			b := &Block{
				Transactions: []*Tx{{
					Id: "Tx01",
					TxOuts: []*TxOut{
						{Address: "fakeAddress", Amount: 50}, {Address: "fakeAddress", Amount: 50},
						{Address: "fakeAddress", Amount: 50}, {Address: "fakeAddress", Amount: 50}},
				}}}
			return utils.ToBytes(b)
		},
	}
	b = &blockchain{NewestHash: "xxxx"}

	t.Run("not enough money in account", func(t *testing.T) {
		_, err := mem.AddTx("fakeAddress", 200)
		if err != ErrorNoMoney {
			t.Error("AddTx() should raise Error when money is not enough.")
		}

	})

	myWallet = fakeWallet{
		fakeVerify: func() bool { return false },
	}
	t.Run("When wrong signature", func(t *testing.T) {
		_, err := mem.AddTx("fakeAddress", 30)
		if err != ErrorNotValid {
			t.Error("AddTx() should raise ErrorNotValid for wrong signature.")
		}
	})

	myWallet = fakeWallet{
		fakeVerify: func() bool { return true },
	}
	t.Run("When there is changes", func(t *testing.T) {
		tx, _ := mem.AddTx("fakeAddress", 30)
		if tx.TxOuts[0].Address != "fakeAddress" && tx.TxOuts[0].Amount != 70 {
			t.Errorf("the change should 70, but got %d", tx.TxOuts[0].Amount)
		}
	})
}

func TestAddPeerTx(t *testing.T) {
	tx := Tx{Id: "Tx01"}
	mem = &mempool{Txs: make(map[string]*Tx)}
	mem.AddPeerTx(&tx)
	if mem.Txs["Tx01"] != &tx {
		t.Error("AddPeerTx() should append transaction to mempool.")
	}
}

// AddTx() test에 통합하려다가 실패. 이게 발생할 수가 있나?
// UTxOutsByAddress() 에서 나온 TxOuts로 TxIns를 만든건데, 재검사를 또 해야 하나?
func TestValidate(t *testing.T) {
	tx := Tx{TxIns: []*TxIn{{TxId: "never", Index: 0}}}
	isValid := validate(&tx)
	if isValid {
		t.Error("validate() should return false for no history txIn")
	}
}
