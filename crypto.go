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
	payloadCase := getPayloadCase(blzEnvelope)

	digest := serializeAndConcat(ct.PPubKey(), payloadCase, string(payload[:]), timeStampAsStr)

	signature, err := ct.privKey.Sign(rand.Reader, digest, crypto.SHA512)
	if err != nil {
		log.Fatal("From signing digest: ", err)
	}

	blzEnvelope.Signature = signature
}

func (ct *Crypto) SignMsg(payload []byte) ([]byte, error) {
	pbBlzEnvelop := &pb.BznEnvelope{
		Sender:    ct.PPubKey(),
		Signature: []byte{},
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

	str := buffer.String()
	log.Println("Digest generated: ", str)
	encodedData := encodeToASCII(str)
	return []byte(encodedData)
}

func serialize(data string) string {
	var buffer bytes.Buffer
	parsedData := encodeToASCII(data)
	buffer.WriteString(strconv.Itoa(len(parsedData)))
	buffer.WriteString("|")
	buffer.WriteString(parsedData)
	return buffer.String()
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

func getPayloadCase(bzn *pb.BznEnvelope) string {
	// Proto index case check proto/bluzelle.proto for more details
	var idx int
	if bzn.GetDatabaseMsg() != nil {
		idx = 10
	}
	if bzn.GetPbftInternalRequest() != nil {
		idx = 11
	}
	return strconv.Itoa(idx)
}
