package bluzelle

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	testnetUrl  string = "testnet.bluzelle.com"
	testnetPort uint32 = 51010
	testDbUuid  string = "80174b53-2dda-49f1-9d6a-6a780d"

	errMsgTemplate = "%s Error: expected %s, got %s"
)

func randUuid() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

func sleep() {
	time.Sleep(2 * time.Second)
}

func TestSetters(t *testing.T) {
	blz := Connect("", 0, "")

	e := "testnet.bluzelle.ca"
	blz.SetEndpoint(e)
	if blz.Endpoint != e {
		t.Errorf(errMsgTemplate, "SetEndpoint", e, blz.Endpoint)
	}

	var p uint32 = 50000
	blz.SetPort(p)
	if blz.Port != p {
		t.Errorf(errMsgTemplate, "SetPort", p, blz.Port)
	}

	randUuid, _ := uuid.NewRandom()
	u := randUuid.String()
	blz.SetUuid(u)
	if blz.Uuid != u {
		t.Errorf(errMsgTemplate, "SetUuid", u, blz.Uuid)
	}
}

// To remove < redundant
func TestCreate(t *testing.T) {
	blz := Connect(testnetUrl, testnetPort, randUuid())

	k := randUuid()
	v := randUuid()

	err := blz.Create(k, []byte(v))
	if err != nil {
		t.Errorf("Create Error: %s", err.Error())
	}

	sleep() // Allow data propagation

	readV, err := blz.Read(k)
	if err != nil {
		t.Errorf("Read Error: %s", err.Error())
	}

	if v != string(readV[:]) {
		t.Errorf(errMsgTemplate, "Create", v, readV)
	}

	err = blz.Create(k, []byte(v))
	if err != nil && err != ErrRecordExists {
		t.Errorf(errMsgTemplate, "Create", ErrRecordExists.Error(), err.Error())
	}
}

func TestUpdate(t *testing.T) {
	blz := Connect(testnetUrl, testnetPort, randUuid())

	k := randUuid()
	v := randUuid()

	err := blz.Update(k, []byte(v))
	if err != nil && err != ErrRecordNotFound {
		t.Errorf(errMsgTemplate, "Update", ErrRecordNotFound.Error(), err.Error())
	}

	blz.Create(k, []byte(v))
	sleep()

	updatedV := randUuid()
	err = blz.Update(k, []byte(updatedV))
	if err != nil {
		t.Errorf("Update Error: %s", err.Error())
	}
	sleep()

	readUpdatedV, err := blz.Read(k)
	if string(readUpdatedV[:]) != updatedV {
		t.Errorf(errMsgTemplate, "Update", updatedV, string(readUpdatedV[:]))
	}
}

func TestHas(t *testing.T) {
	blz := Connect(testnetUrl, testnetPort, randUuid())
	k := randUuid()

	has, err := blz.Has(k)
	if err != nil {
		t.Errorf("Read Error: %s", err.Error())
	}
	if has {
		t.Errorf(errMsgTemplate, "Has", false, has)
	}

	err = blz.Create(k, []byte(randUuid()))
	if err != nil {
		t.Errorf("Create Error: %s", err.Error())
	}
	sleep()

	has, err = blz.Has(k)
	if err != nil {
		t.Errorf("Read Error: %s", err.Error())
	}
	if !has {
		t.Errorf(errMsgTemplate, "Has", true, has)
	}
}

func TestSize(t *testing.T) {
	blz := Connect(testnetUrl, testnetPort, randUuid())

	size, err := blz.Size()
	if err != nil {
		t.Errorf("Size Error: %s", err.Error())
	}
	if size > 0 {
		t.Errorf(errMsgTemplate, "Size", 0, size)
	}

	v := []byte(randUuid())
	eSize := binary.Size(v)
	blz.Create(randUuid(), v)

	sleep()

	size, err = blz.Size()
	if int(size) != eSize {
		t.Errorf(errMsgTemplate, "Size", eSize, size)
	}
}

func TestKeys(t *testing.T) {
	blz := Connect(testnetUrl, testnetPort, randUuid())

	keys, err := blz.Keys()
	if err != nil {
		t.Errorf("Keys Error: %s", err.Error())
	}
	if len(keys) > 0 {
		t.Errorf(errMsgTemplate, "Keys", "[]", keys)
	}

	n := 3
	ks := []string{}
	for i := 0; i < n; i++ {
		nk := randUuid()
		err := blz.Create(nk, []byte(""))
		if err != nil {
			t.Errorf("Create Error: %s", err.Error())
		}
		ks = append(ks, nk)
		sleep()
	}

	keys, err = blz.Keys()
	if err != nil {
		t.Errorf("Keys Error: %s", err.Error())
	}
	if len(keys) != n {
		t.Errorf(errMsgTemplate, "Keys", ks, keys)
	}

}
