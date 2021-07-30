package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {
	data := "hello"
	difficulty := 2
	target := strings.Repeat("0", difficulty)
	nonce := 1
	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data + fmt.Sprint(nonce))))
		if strings.HasPrefix(hash, target) {
			fmt.Printf("Data:%s\nHash:%s\nTarget:%s\nNonce:%d\n", data, hash, target, nonce)
			return
		}
		nonce += 1
	}
}