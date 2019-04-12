package bluzelle

// Layer 4: API Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-4-api-layer.md

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/wlwanpan/bluzelle-go/pb"
)

var (
	// ErrRequestTimeout is returned when a outgoing request is timed out.
	ErrRequestTimeout = errors.New("blz: request timed out")
)

const (
	// RequestTimeout time limit per db request
	RequestTimeout = 3 * time.Second
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

// Adming APIs (https://docs.bluzelle.com/bluzelle-js/api)

// Status returns the status of the daemon.
// Response: JSON-encoded swarm status.
func (blz *Bluzelle) Status() error {
	statusMsg := blz.newStatusMsg()
	statusResp, err := blz.sendStatusReq(statusMsg)
	if err != nil {
		return err
	}
	log.Println(statusResp)
	return nil
}

// Close just closes the connection to the daemon.
func (blz *Bluzelle) Close() {
	blz.close <- true
}

// CreateDB creates a new database at the given uuid.
func (blz *Bluzelle) CreateDB() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_CreateDb{
		CreateDb: &pb.DatabaseRequest{},
	}
	resp, err := blz.sendDbReq(blzMsg)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}

// DeleteDB deletes the database at the given uuid.
func (blz *Bluzelle) DeleteDB() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_DeleteDb{
		DeleteDb: &pb.DatabaseRequest{},
	}
	resp, err := blz.sendDbReq(blzMsg)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}

// HasDB queries to see if a database exists at the given uuid.
func (blz *Bluzelle) HasDB() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_HasDb{
		HasDb: &pb.DatabaseHasDb{},
	}
	resp, err := blz.sendDbReq(blzMsg)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil

}

// PublicKey returns a public key from the private pem given in the constructor.
func (blz *Bluzelle) PublicKey() string {
	return blz.PPubKey()
}

// GetWriters gets the owner and writers of the given database. The owner is the public key
// of the user that created the database. The writers array lists the public keys of users
// that are allowed to make changes to the database.
// Response: JSON {
//   owner: 'MFYwEAY...EpZop6A==',
//   writers: [
//     'MFYwEAYH...0FEoB==', ...
//   ]
// }
func (blz *Bluzelle) GetWriters() {}

// AddWriters add writers to the writers list. May only be executed by the owner of the database.
func (blz *Bluzelle) AddWriters(pubKeys ...string) {}

// DeleteWriters deletes writers from the writers list. May only be executed by the owner of the database.
func (blz *Bluzelle) DeleteWriters(pubKeys ...string) {}

// Database APIs

func (blz *Bluzelle) Create() {}

func (blz *Bluzelle) Read() {}

func (blz *Bluzelle) Update() {}

func (blz *Bluzelle) QuickRead() {}

func (blz *Bluzelle) Delete() {}

func (blz *Bluzelle) Has() {}

func (blz *Bluzelle) Keys() {}

func (blz *Bluzelle) Size() error {
	blzMsg := blz.newDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_Size{
		Size: &pb.DatabaseRequest{},
	}
	resp, err := blz.sendDbReq(blzMsg)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}

// Private methods

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

func (blz *Bluzelle) sendStatusReq(statusMsg *pb.StatusRequest) (*pb.StatusResponse, error) {
	data, err := proto.Marshal(statusMsg)
	if err != nil {
		return nil, err
	}
	if err := blz.sendMsg(data); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	select {
	case resp := <-blz.readMsg():
		statusResp := &pb.StatusResponse{}
		if err := proto.Unmarshal(resp, statusResp); err != nil {
			return nil, err
		}
		return statusResp, nil
	case <-ctx.Done():
		return nil, ErrRequestTimeout
	}
}

func (blz *Bluzelle) sendDbReq(dbMsg *pb.DatabaseMsg) (*pb.DatabaseResponse, error) {
	signedData, err := blz.SignMsg(dbMsg)
	if err != nil {
		log.Println("Error signing data: ", err)
		return nil, err
	}
	if err := blz.sendMsg(signedData); err != nil {
		log.Println("Error sending message: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	select {
	case resp := <-blz.readMsg():
		blzEnvelop := &pb.BznEnvelope{}
		if err := proto.Unmarshal(resp, blzEnvelop); err != nil {
			return nil, err
		}

		dbResp := blzEnvelop.GetDatabaseResponse()
		pbresp := &pb.DatabaseResponse{}
		if err = proto.Unmarshal(dbResp, pbresp); err != nil {
			return nil, err
		}

		dbErr := pbresp.GetHeader()
		log.Println("db uuid: ", dbErr.GetDbUuid())
		return pbresp, nil
	case <-ctx.Done():
		return nil, ErrRequestTimeout
	}
}
