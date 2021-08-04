package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
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

// const (
// 	signature  string = "9a424d9f2e2d711d0b497940dc8de26e678860fb7282f394b42fa535129ec109914146328c552475c626848e2e3fd65e32204ee6ef78d55185238a8572ca2204"
// 	privateKey string = "30770201010420700a712c77abf88e962d20ca99fe287bbb24e5151b6d1dc1aa3bb5f327e0c315a00a06082a8648ce3d030107a144034200049696a1bbae8d4d22d08606de9a7fdb2ba4198683da770c5c6c691758a2c5583a41c6051e1c6b2b45f97a625f8038c6ef7fed472b5cc2d9dc3dc40d4646d30a39"
// 	hashedMsg  string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
// )

func Start() {
	//1. generate key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	// 1-1. privateKey를 저장하기 위해 Bytes로 변환
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	fmt.Printf("%x\n\n", keyAsBytes)

	// 2. hash massage
	message := "i love you"
	hashedMsg := utils.Hash(message)
	// 3. sign to hashed message
	hashAsBytes, err := hex.DecodeString(hashedMsg)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleErr(err)

	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("%x\n\n", signature)

	//4. Verify
	isOk := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
	fmt.Println(isOk)
}
