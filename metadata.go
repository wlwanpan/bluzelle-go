package main

import (
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
			DbUuid: md.blz.UUID,
			Nonce:  randUint64(),
		},
	}
}

func (md *Metadata) newStatusMsg() *pb.StatusRequest {
	return &pb.StatusRequest{}
}
