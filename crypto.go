package main

// Layer 2: Cryptographic Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-1-persistent-connection.md

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
	// ErrPemfileNotECDSA is returned when a wrong pem file format is provided.
	ErrPemfileNotECDSA = errors.New("crypto: pem file loaded not ecdsa")

	// ErrInvalidPayloadCase is returned when could not find a payload case.
	ErrInvalidPayloadCase = errors.New("crypto: invalid payload case")

	// ErrSigVerificationFailed is returned when a the priv/pub key generation failed.
	ErrSigVerificationFailed = errors.New("crypto: fail to verify sig")
)

const (
	EcPrivateKey string = "EC PRIVATE KEY"
)

type Crypto struct {
	privKey *secp256k1.PrivateKey

	serializedPubKey []byte
}

func NewCrypto(privPem []byte) (*Crypto, error) {
	block, _ := pem.Decode(privPem)
	if block.Type != EcPrivateKey {
		return nil, ErrPemfileNotECDSA
	}

	privKey, _ := secp256k1.PrivKeyFromBytes(block.Bytes)
	return &Crypto{privKey: privKey}, nil
}

func (ct *Crypto) GenPubKey() []byte {
	if len(ct.serializedPubKey) != 0 {
		return ct.serializedPubKey
	}
	pubKey := ct.privKey.PubKey()
	ct.serializedPubKey = elliptic.Marshal(pubKey, pubKey.X, pubKey.Y)
	return ct.serializedPubKey
}

func (ct *Crypto) PPubKey() string {
	pk := ct.GenPubKey()
	// Todo: confirm why is this header needed ?
	return "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE" + base64.StdEncoding.EncodeToString(pk)
}

func (ct *Crypto) setMsgSig(blzEnvelope *pb.BznEnvelope) error {
	timeStamp := blzEnvelope.GetTimestamp()

	payloadCase, err := getPayloadCase(blzEnvelope)
	if err != nil {
		return err
	}

	binForWin := []string{
		ct.PPubKey(),
		strconv.Itoa(payloadCase),
		string(blzEnvelope.GetDatabaseMsg()),
		strconv.Itoa(int(timeStamp)),
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
	pbBlzEnvelop := &pb.BznEnvelope{
		Sender:    ct.PPubKey(),
		Signature: []byte{},
		Timestamp: uint64(time.Now().UTC().Unix()),
		Payload: &pb.BznEnvelope_DatabaseMsg{
			DatabaseMsg: payload,
		},
	}

	if err := ct.setMsgSig(pbBlzEnvelop); err != nil {
		return []byte{}, err
	}
	return proto.Marshal(pbBlzEnvelop)
}

func serializeAndConcat(s []string) []byte {
	var buffer bytes.Buffer
	for _, data := range s {
		ds := deterministicSerialize(data)
		encodedDs := stringToASCII(ds)
		buffer.WriteString(string(encodedDs))
	}

	result := buffer.Bytes()
	log.Println("Digest generated: ", result)
	return result
}

func deterministicSerialize(data string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(len(data)))
	buffer.WriteString("|")
	buffer.WriteString(data)
	return buffer.String()
}

func getPayloadCase(bzn *pb.BznEnvelope) (int, error) {
	// Proto index case check proto/proto/bluzelle.proto for more details
	var output int
	var err error
	switch bzn.GetPayload().(type) {
	case *pb.BznEnvelope_DatabaseMsg:
		output = 10
	case *pb.BznEnvelope_PbftInternalRequest:
		output = 11
	case *pb.BznEnvelope_DatabaseResponse:
		output = 12
	case *pb.BznEnvelope_Json:
		output = 13
	case *pb.BznEnvelope_Audit:
		output = 14
	case *pb.BznEnvelope_Pbft:
		output = 15
	case *pb.BznEnvelope_PbftMembership:
		output = 16
	case *pb.BznEnvelope_StatusRequest:
		output = 17
	case *pb.BznEnvelope_StatusResponse:
		output = 18
	default:
		err = ErrInvalidPayloadCase
	}
	return output, err
}
