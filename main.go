package main

import (
	"crypto/rsa"
	"log"
	"os"
)

type Bluzelle struct {
	Entry   string
	UUID    string
	PrivPem string

	conn *Conn
}

func Connect(entry, uuid, privPem string) (*Bluzelle, error) {
	blz := &Bluzelle{
		Entry:   entry,
		UUID:    uuid,
		PrivPem: privPem,
		conn:    NewConn(entry),
	}
	if err := blz.Dial(); err != nil {
		return nil, err
	}
	return blz, nil
}

func (blz *Bluzelle) GetPublicKey() (*rsa.PublicKey, error) {
	return PublicKey(blz.PrivPem)
}

func (blz *Bluzelle) Dial() error {
	return blz.conn.Dial()
}

func (blz *Bluzelle) CreateDB() {
	createPb := CreateDB()
}

func (blz *Bluzelle) Read() string {
	return ""
}

func (blz *Bluzelle) sendReq() {

}

func main() {

	entry := "bernoulli.bluzelle.com:51010"
	uuid := "5f493479–2447–47g6–1c36-efa5d251a283"
	privPem := "MHQCAQEEIFNmJHEiGpgITlRwao/CDki4OS7BYeI7nyz+CM8NW3xToAcGBSuBBAAKoUQDQgAEndHOcS6bE1P9xjS/U+SM2a1GbQpPuH9sWNWtNYxZr0JcF+sCS2zsD+xlCcbrRXDZtfeDmgD9tHdWhcZKIy8ejQ=="

	blz, err := Connect(entry, uuid, privPem)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println(blz.GetPublicKey())

	blz.CreateDB()
}
