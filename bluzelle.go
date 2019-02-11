package main

// Layer 4: API Layer
// Doc reference (https://github.com/bluzelle/client-development-guide/blob/v0.4.x/layers/layer-4-api-layer.md)

import (
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
func Connect(entry, uuid, privPem string) (*Bluzelle, error) {
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
	return blz.crypto.GetPubKey()
}

func (blz *Bluzelle) Dial() error {
	return blz.conn.Dial()
}

func (blz *Bluzelle) sendReq(pbDbMsg *pb.DatabaseMsg) error {
	data, err := proto.Marshal(pbDbMsg)
	if err != nil {
		return err
	}
	signedData, err := blz.crypto.SignMsg(data)
	if err != nil {
		return err
	}
	blz.conn.SendMsg(signedData)

	response := <-blz.conn.ReadMsg()
	log.Printf("This is the response: %s", response)
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
	r := rand.New(rand.NewSource(time.Now().UTC().Unix()))
	return &pb.DatabaseMsg{
		Header: &pb.DatabaseHeader{
			DbUuid:         blz.UUID,
			Nonce:          r.Uint64(),
			PointOfContact: []byte{},
		},
	}
}

func main() {

	entry := "bernoulli.bluzelle.com:51010"
	uuid := "5f493479–2447–47g6–1c36-efa5d251a283"
	privPem := `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----
	` // Some Random test pem file from stackoverflow

	blz, err := Connect(entry, uuid, privPem)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println(blz.PublicKey())

	blz.CreateDB()
	select {}
}
