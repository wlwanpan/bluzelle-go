package bluzelle

// Layer 2: Cryptographic Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-1-persistent-connection.md

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
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
	// EcPrivateKey pem format required
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
	// TODO: confirm why is this header needed ?
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

	digest := []byte(serializeAndConcat(binForWin))

	// Example outgoing payload (binForWin):
	// 120|MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEHLmWG+xY1lZak68iQMFCh2Vm1EfkEOkbclWWbEO1s+qpf6D6D/Yjo9CR2/zLHWkgTqo71nnjWWEU5FekTHjVmQ==
	// 2|10
	// 0|
	// 1|0
	// sig, err := ct.privKey.Sign(digest)
	ecdsaPriv := (*ecdsa.PrivateKey)(ct.privKey)
	sig, err := sign(digest, ecdsaPriv)
	if err != nil {
		log.Println("From signing digest: ", err)
		return err
	}

	ecdsaPubKey := (*ecdsa.PublicKey)(ct.privKey.PubKey())
	// !sig.Verify(digest, ct.privKey.PubKey())
	if !verify(digest, sig, ecdsaPubKey) {
		return ErrSigVerificationFailed
	}

	// blzEnvelope.Signature = sig.Serialize()
	blzEnvelope.Signature = sig
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

func serializeAndConcat(s []string) string {
	var buffer bytes.Buffer
	for _, data := range s {
		ds := deterministicSerialize(data)
		encodedDs := stringToASCII(ds)
		buffer.WriteString(string(encodedDs))
	}

	log.Println("Digest generated: ", buffer.String())
	return buffer.String()
}

func deterministicSerialize(data string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(len(data)))
	buffer.WriteString("|")
	buffer.WriteString(data)
	return buffer.String()
}

// Testing:
// Crypto for go developers (GopherCon talk) - George Tankersley
// speakerdeck.com/gtank/crypto-for-go-developers?slide=71
func sign(data []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv, data)
	if err != nil {
		return nil, err
	}

	params := priv.Curve.Params()
	curveByteSize := params.P.BitLen() / 8
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	sigSize := curveByteSize * 2
	sig := make([]byte, sigSize)
	copy(sig[curveByteSize-len(rBytes):], rBytes)
	copy(sig[sigSize-len(sBytes):], sBytes)
	return sig, nil
}

func verify(data []byte, sig []byte, pubKey *ecdsa.PublicKey) bool {
	curveByteSize := pubKey.Curve.Params().P.BitLen() / 8
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(sig[:curveByteSize])
	s.SetBytes(sig[curveByteSize:])

	return ecdsa.Verify(pubKey, data[:], r, s)
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
