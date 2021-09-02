package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	// //toGetTestKey() 참조
	testKey string = "30770201010420722f026193ca462a71ef1bbca650f4c93e2920297310fef8ab03656934473a87a00a06082a8648ce3d030107a144034200041ea0b7a88eb16ecce8cac6e99c1ca7fbe87fa909708cbea6f591690b5efa1a41dadfea9638baf06343918fc521ec5513771bea492a81036672e9c570e156b4e6"
	testSig string = "12232ff554b277720ad5fe38a5785e1872a6b4523814a71ef4ade3170647af3e1e3b59a796a970bc7cb46971cda372a6685b82c134391121c3316bd6dea26d67"
	// 유효한 hexadecimal 아무거나 (이거는 block hash 가져옴)
	testPayload  string = "0001961555cb88268649070d839cc7bb8951f4512ad1faa024cc6dd19bafb193"
	wrongPayload string = "c2bb8314fa7640b970a2a5ca49b529c4d09129d949119b1bcb1388dd9ad419dc"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	// return utils.ToBytes(makeTestWallet().privateKey), nil
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestWallet(t *testing.T) {
	t.Run("New wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return false },
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})

	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return true },
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	bytes, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(bytes)
	return &wallet{key, aFromK(key)}
}

func TestSign(t *testing.T) {
	w = makeTestWallet()
	sig := W.Sign(testPayload)
	// 서명은 일치여부로 검사불가(랜덤). hexadecimal 인지를 검사
	_, err := hex.DecodeString(sig)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s\n", sig)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		ok    bool
		input string
	}
	tests := []test{
		{true, testPayload},
		{false, wrongPayload},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := W.Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xxxxxxx")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex")
	}
}

// func toGetTestKey(t *testing.T) {
// 	// createPrivKey()
// 	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	// persistKey
// 	bytes, _ := x509.MarshalECPrivateKey(privKey)
// 	// see hexadecimal key
// 	t.Logf("%x\n", bytes)
// }
