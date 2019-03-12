package main

// TODO: Seperate 3 and 4.
// Layer 3 + 4: Metadata and API Layer
// (https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-4-api-layer.md)

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

// Bluzelle represents a client connection to the bluzelle network.
type Bluzelle struct {
	Entry string
	UUID  string

	// Bluzelle layers
	conn   *Conn
	crypto *Crypto
}

// Connect initialize a new bluzelle struct.
func Connect(entry, uuid string, privPem []byte) (*Bluzelle, error) {
	blz := &Bluzelle{
		Entry:  entry,
		UUID:   uuid,
		conn:   NewConn(entry),
		crypto: NewCrypto(privPem),
	}
	if err := blz.Dial(); err != nil {
		return nil, err
	}
	return blz, nil
}

// PublicKey returns the corresponding public key from the bluzelle private pem.
func (blz *Bluzelle) PublicKey() []byte {
	return blz.crypto.PubKey()
}

// PPublicKey returns the corresponding public key in hex string format.
func (blz *Bluzelle) PPublicKey() string {
	return blz.crypto.PPubKey()
}

func (blz *Bluzelle) Dial() error {
	return blz.conn.Dial()
}

func (blz *Bluzelle) sendReq(dbMsg *pb.DatabaseMsg) error {
	signedData, err := blz.crypto.SignMsg(dbMsg)
	if err != nil {
		log.Println("Error signing data: ", err)
		return err
	}
	blz.conn.SendMsg(signedData)

	select {
	case resp := <-blz.conn.ReadMsg():
		blzEnvelop := &pb.BznEnvelope{}
		if err := proto.Unmarshal(resp, blzEnvelop); err != nil {
			log.Fatal(err)
		}
		dbResp := blzEnvelop.GetDatabaseResponse()
		pbresp := &pb.DatabaseResponse{}
		if err = proto.Unmarshal(dbResp, pbresp); err != nil {
			log.Fatal(err)
		}
		dbErr := pbresp.GetHeader()
		log.Println("db uuid: ", dbErr.GetDbUuid())
	}
	return nil
}

func (blz *Bluzelle) CreateDB() {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_CreateDb{
		CreateDb: &pb.DatabaseRequest{},
	}
	blz.sendReq(blzMsg)
}

func (blz *Bluzelle) newDatabaseMsg() *pb.DatabaseMsg {
	return &pb.DatabaseMsg{
		Header: &pb.DatabaseHeader{
			DbUuid:         blz.UUID,
			Nonce:          blz.randNonce(),
			PointOfContact: []byte{},
		},
	}
}

func (blz *Bluzelle) randNonce() uint64 {
	now := time.Now().UTC().Unix()
	r := rand.New(rand.NewSource(now))
	return r.Uint64()
}

// ReadPemFile reads a private pem file.
func ReadPemFile(path string) ([]byte, error) {
	privKeyFile, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer privKeyFile.Close()
	pemStats, err := privKeyFile.Stat()
	if err != nil {
		return []byte{}, err
	}
	log.Println("Loaded pem file", pemStats.Name())

	return ioutil.ReadAll(privKeyFile)
}

func main() {
	entry := "bernoulli.bluzelle.com:51010"
	uuid := "5f493479–2447–47g6–1c36-efa5d251a283"

	pemBytes, err := ReadPemFile("./test.pem")
	if err != nil {
		log.Fatal(err)
	}

	blz, err := Connect(entry, uuid, pemBytes)
	if err != nil {
		log.Fatal(err)
	}

	blz.CreateDB()
}
