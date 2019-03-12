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
	"github.com/golang/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

var (
	ErrInvalidPayloadCase = errors.New("crypto: invalid payload case")

	ErrSigVerificationFailed = errors.New("crypto: failed to verify sig")
)

const (
	EcPrivateKey = "EC PRIVATE KEY"
)

type Crypto struct {
	privKey *secp256k1.PrivateKey

	serializedPubKey []byte
}

func NewCrypto(privPem []byte) *Crypto {
	block, _ := pem.Decode(privPem)
	if block.Type != EcPrivateKey {
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
	// Todo: why is this header needed ?
	return "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE" + base64.StdEncoding.EncodeToString(pk)
}

func (ct *Crypto) setSignature(blzEnvelope *pb.BznEnvelope, payload string) error {
	timeStamp := blzEnvelope.GetTimestamp()
	timeStampAsStr := strconv.Itoa(int(timeStamp))

	payloadCase, err := getPayloadCase(blzEnvelope)
	if err != nil {
		return err
	}

	binForWin := []string{
		ct.PPubKey(),
		payloadCase,
		string(blzEnvelope.GetDatabaseMsg()),
		timeStampAsStr,
	}

	digest := serializeAndConcat(binForWin)

	sig, err := ct.privKey.Sign(digest)
	if err != nil {
		log.Println("From signing digest: ", err)
		return err
	}

	if !sig.Verify(digest, ct.privKey.PubKey()) {
		return ErrSigVerificationFailed
	}

	blzEnvelope.Signature = sig.Serialize()
	return nil
}

func (ct *Crypto) SignMsg(dbMsg *pb.DatabaseMsg) ([]byte, error) {
	payload, err := proto.Marshal(dbMsg)
	if err != nil {
		return []byte{}, err
	}
	log.Println(payload)
	pbBlzEnvelop := &pb.BznEnvelope{
		Sender:    ct.PPubKey(),
		Signature: []byte{},
		Timestamp: uint64(time.Now().UTC().Unix()),
		Payload: &pb.BznEnvelope_DatabaseMsg{
			DatabaseMsg: payload,
		},
	}

	if err := ct.setSignature(pbBlzEnvelop, dbMsg.String()); err != nil {
		return []byte{}, err
	}
	return proto.Marshal(pbBlzEnvelop)
}

func serializeAndConcat(data []string) []byte {
	var buffer bytes.Buffer
	for _, d := range data {
		serData := deterministicSerialize(d)
		buffer.WriteString(serData)
	}

	str := buffer.String()
	encodedData := encodeToASCII(str)
	log.Println("Digest generated: ", encodedData)
	return []byte(encodedData)
}

func deterministicSerialize(data string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(len(data)))
	buffer.WriteString("|")
	buffer.WriteString(data)
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

func getPayloadCase(bzn *pb.BznEnvelope) (string, error) {
	// Proto index case check proto/bluzelle.proto for more details
	var output int
	if bzn.GetDatabaseMsg() != nil {
		output = 10
	}
	if bzn.GetPbftInternalRequest() != nil {
		output = 11
	}

	if output < 10 || output > 18 {
		return "", ErrInvalidPayloadCase
	}
	return strconv.Itoa(output), nil
}
