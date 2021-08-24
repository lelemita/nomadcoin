package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	// //toGetTestKey() 참조
	testKey string = "30770201010420722f026193ca462a71ef1bbca650f4c93e2920297310fef8ab03656934473a87a00a06082a8648ce3d030107a144034200041ea0b7a88eb16ecce8cac6e99c1ca7fbe87fa909708cbea6f591690b5efa1a41dadfea9638baf06343918fc521ec5513771bea492a81036672e9c570e156b4e6"
	// 유효한 hexadecimal 아무거나 (이거는 block hash 가져옴)
	testPayload string = "0001961555cb88268649070d839cc7bb8951f4512ad1faa024cc6dd19bafb193"
	testSig     string = "12232ff554b277720ad5fe38a5785e1872a6b4523814a71ef4ade3170647af3e1e3b59a796a970bc7cb46971cda372a6685b82c134391121c3316bd6dea26d67"
)

func makeTestWallet() *wallet {
	bytes, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(bytes)
	return &wallet{key, aFromK(key)}
}

func TestSign(t *testing.T) {
	sig := Sign(testPayload, makeTestWallet())
	// 서명은 일치여부로 검사불가(랜덤). hexadecimal 인지를 검사
	_, err := hex.DecodeString(sig)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s\n", sig)
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
