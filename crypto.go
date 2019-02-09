package main

import (
	"crypto/x509"
	"crypto/rsa"
	"encoding/pem"
)

type Crypto struct {
	// PublicKey

	privateKey *rsa.PrivateKey
}

func GenBlz() {}

func PublicKey(privatePem string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(privatePem))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}