package main

// Layer 2: Cryptographic Layer
// (https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-1-persistent-connection.md)

import (
	"bytes"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/pem"
	"errors"
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
	privKey *secp256k1.PrivateKey

	serializedPubKey []byte
}

func NewCrypto(privPem []byte) *Crypto {
	block, _ := pem.Decode(privPem)
	if block.Type != EC_PRIVATE_KEY {
		log.Fatal("pem file loaded not ecdsa")
	}

	privKey, _ := secp256k1.PrivKeyFromBytes(block.Bytes)

	return &Crypto{
		privKey: privKey,
	}
}

func (ct *Crypto) PubKey() []byte {
	if len(ct.serializedPubKey) != 0 {
		return ct.serializedPubKey
	}
	pubKey := ct.privKey.PubKey()
	ct.serializedPubKey = elliptic.Marshal(pubKey, pubKey.X, pubKey.Y)
	return ct.serializedPubKey
}

func (ct *Crypto) PPubKey() string {
	pk := ct.PubKey()
	return base64.StdEncoding.EncodeToString(pk)
}

func (ct *Crypto) setSignature(blzEnvelope *pb.BznEnvelope, payload []byte) error {
	timeStamp := blzEnvelope.GetTimestamp()
	timeStampAsStr := strconv.Itoa(int(timeStamp))
	payloadCase := getPayloadCase(blzEnvelope)

	digest := serializeAndConcat(ct.PPubKey(), payloadCase, string(payload[:]), timeStampAsStr)

	signature, err := ct.privKey.Sign(digest)
	if err != nil {
		log.Println("From signing digest: ", err)
		return err
	}

	if !signature.Verify(digest, ct.privKey.PubKey()) {
		log.Println("From signing digest: failed to verify signature")
		return errors.New("failed to verify signature")
	}

	blzEnvelope.Signature = signature.Serialize()
	return nil
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

	if err := ct.setSignature(pbBlzEnvelop, payload); err != nil {
		return []byte{}, err
	}
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
