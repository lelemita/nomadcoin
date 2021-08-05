package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"math/big"
	"os"

	"github.com/lelemita/nomadcoin/utils"
)

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string //public key
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, os.FileMode(0644))
	utils.HandleErr(err)

}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
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

func Sign(payload string, w *wallet) string {
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
	BigA, bigB := big.Int{}, big.Int{}
	BigA.SetBytes(bytes[:len(bytes)/2])
	bigB.SetBytes(bytes[len(bytes)/2:])
	return &BigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
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
		Y:     y}
	utils.HandleErr(err)
	return ecdsa.Verify(&publicKey, hash, r, s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		// has a wallet already?
		if hasWalletFile() {
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
