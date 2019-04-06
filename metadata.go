package main

import (
	"math/rand"
	"time"

	"github.com/wlwanpan/bluzelle-go/pb"
)

// Layer 3: Metadata Layer
// https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-3-metadata.md

type Metadata struct {
	blz *Bluzelle
}

func (md *Metadata) newDatabaseMsg() *pb.DatabaseMsg {
	return &pb.DatabaseMsg{
		Header: &pb.DatabaseHeader{
			DbUuid:         md.blz.UUID,
			Nonce:          randNonce(),
			PointOfContact: []byte{},
		},
	}
}

func randNonce() uint64 {
	now := time.Now().UTC().Unix()
	r := rand.New(rand.NewSource(now))
	return r.Uint64()
}
