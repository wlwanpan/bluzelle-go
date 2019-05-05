package bluzelle

import (
	"testing"
)

const (
	TestEndPoint string = "bernoulli.bluzelle.com:51010"
	TestUUID     string = "5f493479–2447–47g6–1c36-efa5d251a25"
	TestPemfile  string = "./test.pem"
)

var blzTest *Bluzelle

func TestConnect(t *testing.T) {

	pemBytes, err := ReadPemFile(TestPemfile)
	if err != nil {
		t.Error(err)
	}
	blzTest, err = Connect(TestEndPoint, TestUUID, pemBytes)
	if err != nil {
		t.Error(err)
	}
}

func TestSize(t *testing.T) {
	if err := blzTest.Size(); err != nil {
		t.Error(err)
	}

}

// Close connection to testnet daemon.
func TestClose(t *testing.T) {
	blzTest.Close()
}
