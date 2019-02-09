package main

import "crypto/rsa"

type Crypto struct {
	PublicKey

	privateKey *rsa.PrivateKey
}

func NewCrypto(pem string) {
	p := []byte(pem)
	privateKey, err := rsa.GenerateKey(p, len(p))
}
