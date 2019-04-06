package main

// Layer 4: API Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-4-api-layer.md

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

// Bluzelle represents a client connection to the bluzelle network.
type Bluzelle struct {
	// Layers
	*Metadata
	*Crypto
	*Conn

	Entry string
	UUID  string

	privPem []byte
}

// Connect initialize a new bluzelle struct.
func Connect(entry, uuid string, privPem []byte) (*Bluzelle, error) {
	blz := &Bluzelle{
		Entry:   entry,
		UUID:    uuid,
		privPem: privPem,
	}
	if err := blz.initLayers(); err != nil {
		return nil, err
	}
	return blz, nil
}

func (blz *Bluzelle) initLayers() error {
	blz.Metadata = &Metadata{blz: blz}

	crypto, err := NewCrypto(blz.privPem)
	if err != nil {
		return err
	}
	blz.Crypto = crypto

	blz.Conn = NewConn(blz.Entry)
	if err := blz.Dial(); err != nil {
		return err
	}
	return nil
}

// Adming APIs (https://docs.bluzelle.com/bluzelle-js/api)

func (blz *Bluzelle) Status() {}

func (blz *Bluzelle) Close() {}

func (blz *Bluzelle) CreateDB() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_CreateDb{
		CreateDb: &pb.DatabaseRequest{},
	}
	return blz.sendReq(blzMsg)
}

func (blz *Bluzelle) DeleteDB() {}

func (blz *Bluzelle) HasDB() {}

// PublicKey returns the corresponding public key in hex string format.
func (blz *Bluzelle) PublicKey() string {
	return blz.PPubKey()
}

func (blz *Bluzelle) GetWriters() {}

func (blz *Bluzelle) AddWriters() {}

func (blz *Bluzelle) DeleteWriters() {}

// Database APIs

func (blz *Bluzelle) Create() {}

func (blz *Bluzelle) Read() {}

func (blz *Bluzelle) Update() {}

func (blz *Bluzelle) QuickRead() {}

func (blz *Bluzelle) Delete() {}

func (blz *Bluzelle) Has() {}

func (blz *Bluzelle) Keys() {}

func (blz *Bluzelle) Size() {}

// Private

func (blz *Bluzelle) sendReq(dbMsg *pb.DatabaseMsg) error {
	signedData, err := blz.SignMsg(dbMsg)
	if err != nil {
		log.Println("Error signing data: ", err)
		return err
	}
	blz.SendMsg(signedData)

	select {
	case resp := <-blz.ReadMsg():
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
