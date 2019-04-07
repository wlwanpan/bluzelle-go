package main

// Layer 4: API Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-4-api-layer.md

import (
	"log"

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

func (blz *Bluzelle) Status() error {
	statusMsg := blz.newStatusMsg()
	return blz.sendStatusReq(statusMsg)
}

func (blz *Bluzelle) Close() {}

func (blz *Bluzelle) CreateDB() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_CreateDb{
		CreateDb: &pb.DatabaseRequest{},
	}
	return blz.sendDbReq(blzMsg)
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

func (blz *Bluzelle) sendStatusReq(statusMsg *pb.StatusRequest) error {
	data, err := proto.Marshal(statusMsg)
	if err != nil {
		return err
	}
	blz.sendMsg(data)

	select {
	case resp := <-blz.readMsg():
		log.Println(resp)
	}

	return nil
}

func (blz *Bluzelle) sendDbReq(dbMsg *pb.DatabaseMsg) error {
	signedData, err := blz.SignMsg(dbMsg)
	if err != nil {
		log.Println("Error signing data: ", err)
		return err
	}
	blz.sendMsg(signedData)

	select {
	case resp := <-blz.readMsg():
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

	if err := blz.CreateDB(); err != nil {
		log.Println(err)
	}
}
