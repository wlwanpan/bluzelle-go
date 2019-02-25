package main

// Layer 2: Cryptographic Layer
// (https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-1-persistent-connection.md)

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"log"
	"strconv"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/gogo/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

const (
	EC_PRIVATE_KEY = "EC PRIVATE KEY"
)

type Crypto struct {
	privKey *ecdsa.PrivateKey

	serializedPubKey []byte
}

func NewCrypto(privPem []byte) *Crypto {
	block, _ := pem.Decode(privPem)
	if block.Type != EC_PRIVATE_KEY {
		log.Fatal("pem file loaded not ecdsa")
	}

	privKey, _ := secp256k1.PrivKeyFromBytes(block.Bytes)

	return &Crypto{
		privKey: privKey.ToECDSA(),
	}
}

func (ct *Crypto) PubKey() []byte {
	if len(ct.serializedPubKey) != 0 {
		return ct.serializedPubKey
	}
	pubKey := ct.privKey.PublicKey
	ct.serializedPubKey = elliptic.Marshal(pubKey, pubKey.X, pubKey.Y)
	return ct.serializedPubKey
}

func (ct *Crypto) PPubKey() string {
	pk := ct.PubKey()
	return base64.StdEncoding.EncodeToString(pk)
}

func (ct *Crypto) setSignature(blzEnvelope *pb.BznEnvelope, payload []byte) {
	timeStamp := blzEnvelope.GetTimestamp()
	timeStampAsStr := strconv.Itoa(int(timeStamp))

	payloadCase := blzEnvelope.GetDatabaseMsg()

	digest := serializeAndConcat(ct.PPubKey(), string(payloadCase), string(payload), timeStampAsStr)

	signature, err := ct.privKey.Sign(rand.Reader, digest, crypto.SHA512)
	if err != nil {
		log.Fatal("From signing digest: ", err)
	}

	blzEnvelope.Signature = signature
}

func (ct *Crypto) SignMsg(payload []byte) ([]byte, error) {
	pbBlzEnvelop := &pb.BznEnvelope{
		Sender:    ct.PPubKey(),
		Timestamp: uint64(time.Now().UTC().Unix()),
		Payload: &pb.BznEnvelope_DatabaseMsg{
			DatabaseMsg: payload,
		},
	}

	ct.setSignature(pbBlzEnvelop, payload)
	return proto.Marshal(pbBlzEnvelop)
}

func serializeAndConcat(data ...string) []byte {
	var buffer bytes.Buffer
	for _, d := range data {
		serData := serialize(d)
		buffer.WriteString(serData)
	}

	log.Println("Digest generated: ", buffer.String())
	return buffer.Bytes()
}

func serialize(data string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(len(data)))
	buffer.WriteString("|")
	buffer.WriteString(data)
	return encodeToASCII(buffer.String())
}

func encodeToASCII(str string) string {
	rs := make([]rune, 0, len(str))
	for _, r := range str {
		if r <= 127 {
			rs = append(rs, r)
		}
	}
	return string(rs)
}
