package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"math/big"
	"os"

	"github.com/lelemita/nomadcoin/utils"
)

const (
	fileName string = "nomadcoin.wallet"
)

type MyWallet struct{}

func (MyWallet) Sign(payload string) string {
	return sign(payload, w)
}
func (MyWallet) Verify(signature, payload, address string) bool {
	return verify(signature, payload, address)
}
func (MyWallet) GetAddress() string {
	return w.Address
}

var W *MyWallet = &MyWallet{}

type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

// test할 때는, 이 인터페이스를 원하는대로 구현하면 된다.
// 함수들을 외부 파일에 독립적으로 테스트 하기 위함.
func (layer) hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

var files fileLayer = layer{}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string //public key
}

var w *wallet

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileName, bytes, os.FileMode(0644))
	utils.HandleErr(err)

}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := files.readFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInts(a, b *big.Int) string {
	z := append(a.Bytes(), b.Bytes()...)
	return hex.EncodeToString(z)
}

// Address(public) from (private) Key
func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X, key.Y)
}

func sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r, s)
}

func restoreBigInts(signature string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(bytes[:len(bytes)/2])
	bigB.SetBytes(bytes[len(bytes)/2:])
	return &bigA, &bigB, nil
}

func verify(signature, payload, address string) bool {
	// signature
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	// hash
	hash, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	// public key
	x, y, err := restoreBigInts(address)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	utils.HandleErr(err)
	return ecdsa.Verify(&publicKey, hash, r, s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		// has a wallet already?
		if files.hasWalletFile() {
			// yes -> restore from file
			w.privateKey = restoreKey()
		} else {
			// no -> create prv key, save to file
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
