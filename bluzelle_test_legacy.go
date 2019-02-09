package main

//
// import (
// 	"encoding/binary"
// 	"testing"
// 	"time"
//
// 	"github.com/google/uuid"
// )
//
// var (
// 	testnetPort uint32 = 51010
// 	testnetURL         = "testnet.bluzelle.com"
// 	testDbUUID         = "80174b53-2dda-49f1-9d6a-6a780d"
//
// 	errMsgTemplate = "%s Error: expected %s, got %s"
// )
//
// func randUUID() string {
// 	u, _ := uuid.NewRandom()
// 	return u.String()
// }
//
// func sleep() {
// 	time.Sleep(2 * time.Second)
// }
//
// func TestSetters(t *testing.T) {
// 	blz := Connect("", 0, "")
//
// 	e := "testnet.bluzelle.ca"
// 	blz.SetEndpoint(e)
// 	if blz.Endpoint != e {
// 		t.Errorf(errMsgTemplate, "SetEndpoint", e, blz.Endpoint)
// 	}
//
// 	var p uint32 = 50000
// 	blz.SetPort(p)
// 	if blz.Port != p {
// 		t.Errorf(errMsgTemplate, "SetPort", p, blz.Port)
// 	}
//
// 	randUUID, _ := uuid.NewRandom()
// 	u := randUUID.String()
// 	blz.SetUuid(u)
// 	if blz.UUID != u {
// 		t.Errorf(errMsgTemplate, "SetUuid", u, blz.UUID)
// 	}
// }
//
// // To remove < redundant
// func TestCreate(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
//
// 	k := randUUID()
// 	v := randUUID()
//
// 	err := blz.Create(k, []byte(v))
// 	if err != nil {
// 		t.Errorf("Create Error: %s", err.Error())
// 	}
//
// 	sleep() // Allow data propagation
//
// 	readV, err := blz.Read(k)
// 	if err != nil {
// 		t.Errorf("Read Error: %s", err.Error())
// 	}
//
// 	if v != string(readV[:]) {
// 		t.Errorf(errMsgTemplate, "Create", v, readV)
// 	}
//
// 	err = blz.Create(k, []byte(v))
// 	if err != nil && err != ErrRecordExists {
// 		t.Errorf(errMsgTemplate, "Create", ErrRecordExists.Error(), err.Error())
// 	}
// }
//
// func TestUpdate(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
//
// 	k := randUUID()
// 	v := randUUID()
//
// 	err := blz.Update(k, []byte(v))
// 	if err != nil && err != ErrRecordNotFound {
// 		t.Errorf(errMsgTemplate, "Update", ErrRecordNotFound.Error(), err.Error())
// 	}
//
// 	blz.Create(k, []byte(v))
// 	sleep()
//
// 	updatedV := randUUID()
// 	err = blz.Update(k, []byte(updatedV))
// 	if err != nil {
// 		t.Errorf("Update Error: %s", err.Error())
// 	}
// 	sleep()
//
// 	readUpdatedV, err := blz.Read(k)
// 	if string(readUpdatedV[:]) != updatedV {
// 		t.Errorf(errMsgTemplate, "Update", updatedV, string(readUpdatedV[:]))
// 	}
// }
//
// func TestRemove(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
// 	k := randUUID()
//
// 	if err := blz.Create(k, []byte(randUUID())); err != nil {
// 		t.Errorf("Create Error: %s", err.Error())
// 	}
// 	sleep()
// 	if err := blz.Remove(k); err != nil {
// 		t.Errorf("Remove Error: %s", err.Error())
// 	}
// 	sleep()
// 	has, err := blz.Has(k)
// 	if err != nil {
// 		t.Errorf("Has Error: %s", err.Error())
// 	}
// 	if has {
// 		t.Errorf(errMsgTemplate, "Remove", false, has)
// 	}
//
// 	if err := blz.Remove(randUUID()); err != ErrRecordNotFound {
// 		t.Errorf(errMsgTemplate, "Remove", ErrRecordNotFound.Error(), err)
// 	}
// }
//
// func TestHas(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
// 	k := randUUID()
//
// 	has, err := blz.Has(k)
// 	if err != nil {
// 		t.Errorf("Read Error: %s", err.Error())
// 	}
// 	if has {
// 		t.Errorf(errMsgTemplate, "Has", false, has)
// 	}
//
// 	err = blz.Create(k, []byte(randUUID()))
// 	if err != nil {
// 		t.Errorf("Create Error: %s", err.Error())
// 	}
// 	sleep()
//
// 	has, err = blz.Has(k)
// 	if err != nil {
// 		t.Errorf("Read Error: %s", err.Error())
// 	}
// 	if !has {
// 		t.Errorf(errMsgTemplate, "Has", true, has)
// 	}
// }
//
// func TestSize(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
//
// 	size, err := blz.Size()
// 	if err != nil {
// 		t.Errorf("Size Error: %s", err.Error())
// 	}
// 	if size > 0 {
// 		t.Errorf(errMsgTemplate, "Size", 0, size)
// 	}
//
// 	v := []byte(randUUID())
// 	eSize := binary.Size(v)
// 	blz.Create(randUUID(), v)
//
// 	sleep()
//
// 	size, err = blz.Size()
// 	if int(size) != eSize {
// 		t.Errorf(errMsgTemplate, "Size", eSize, size)
// 	}
// }
//
// func TestKeys(t *testing.T) {
// 	blz := Connect(testnetURL, testnetPort, randUUID())
//
// 	keys, err := blz.Keys()
// 	if err != nil {
// 		t.Errorf("Keys Error: %s", err.Error())
// 	}
// 	if len(keys) > 0 {
// 		t.Errorf(errMsgTemplate, "Keys", "[]", keys)
// 	}
//
// 	n := 3
// 	ks := []string{}
// 	for i := 0; i < n; i++ {
// 		nk := randUUID()
// 		err = blz.Create(nk, []byte(""))
// 		if err != nil {
// 			t.Errorf("Create Error: %s", err.Error())
// 		}
// 		ks = append(ks, nk)
// 		sleep()
// 	}
//
// 	keys, err = blz.Keys()
// 	if err != nil {
// 		t.Errorf("Keys Error: %s", err.Error())
// 	}
// 	if len(keys) != n {
// 		t.Errorf(errMsgTemplate, "Keys", ks, keys)
// 	}
// }
