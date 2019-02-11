package main

// Layer 2: Cryptographic Layer
// Doc reference (https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-1-persistent-connection.md)

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

type Crypto struct {
	pubKey  *rsa.PublicKey
	privKey []byte
}

func NewCrypto(privPem string) *Crypto {
	block, _ := pem.Decode([]byte(privPem))
	cert, _ := x509.ParseCertificate(block.Bytes)
	pubKey := cert.PublicKey.(*rsa.PublicKey)

	return &Crypto{
		pubKey:  pubKey,
		privKey: block.Bytes,
	}
}

func (ct *Crypto) GetPubKey() []byte {
	block := &pem.Block{
		Bytes: x509.MarshalPKCS1PublicKey(ct.pubKey),
	}
	return pem.EncodeToMemory(block)
}

func (ct *Crypto) SignMsg(payload []byte) ([]byte, error) {
	pbBlzEnvelop := &pb.BznEnvelope{
		Sender:    string(ct.GetPubKey()),
		Timestamp: uint64(time.Now().UTC().Unix()),
		Signature: []byte{}, // calc signature
		Payload: &pb.BznEnvelope_DatabaseMsg{
			DatabaseMsg: payload,
		},
	}
	data, err := proto.Marshal(pbBlzEnvelop)
	return data, err
}
