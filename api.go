package main

import (
	"math/rand"
	"time"

	"github.com/wlwanpan/bluzelle-go/pb"
)

// Admin functions

func genDatabaseMsg() *pb.DatabaseMsg {
	r := rand.New(rand.NewSource(time.Now().UTC().Unix()))
	return &pb.DatabaseMsg{
		Header: &pb.DatabaseHeader{
			DbUuid:         "96764e2f-2273-4404-97c0-a05b5e36ea66",
			Nonce:          r.Uint64(),
			PointOfContact: []byte{},
		},
	}
}

func CreateDB() *pb.DatabaseMsg {
	blzMsg := genDatabaseMsg()
	blzMsg.Msg = &pb.DatabaseMsg_CreateDb{
		CreateDb: &pb.DatabaseRequest{},
	}
	return blzMsg
}

func DeleteDB() {}

func HasDB() {}

func PublicKey() {}

func GetWriters() {}

func AddWriters() {}

func DeleteWriters() {}

// Databse api function
