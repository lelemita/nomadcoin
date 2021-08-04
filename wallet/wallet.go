package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/lelemita/nomadcoin/utils"
)

/*
[Sending a signed message]
1) we hash the msg
	"I did" -> hash(x) -> "hased_message"
2) generate key pair
	KeyPair(privateK, publicK)
	(save priv to a file: wallet)
3) sign the hash
	("hashed_message" + privateK) -> "signature"

[verifing the message]
4) verify
	("hashed_message" + "signature" + publicK) -> true/false
*/

func Start() {
	//1. generate key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	//2. hash massage
	message := "I love you"
	hashedMsg := utils.Hash(message)
	//3. sign to hashed message
	hashAsBytes, err := hex.DecodeString(hashedMsg)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleErr(err)

	fmt.Printf("R: %d\nS: %d\n", r, s)
}
