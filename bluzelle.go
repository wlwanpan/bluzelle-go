package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/wlwanpan/bluzelle-go/proto"
)

const (
	DefaultUuid     = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"
	DefaultEndpoint = "127.0.0.1"
	DefaultPort     = 51010
)

type Bluzelle struct {
	Endpoint string
	Port     uint16
	Uuid     string
}

func (blz Bluzelle) getWsAddr() string {
	p := fmt.Sprint(blz.Port)
	strArr := []string{blz.Endpoint, ":", p}
	return strings.Join(strArr, "")
}

func (blz Bluzelle) pbHeader() *pb.DatabaseHeader {
	return &pb.DatabaseHeader{
		DbUuid:        blz.Uuid,
		TransactionId: rand.Uint64(),
	}
}

func (blz *Bluzelle) sendRequest(m []byte) (string, error) {
	encodedBase64 := base64.StdEncoding.EncodeToString(m)
	// From rb encoded -> "Ui4SJAofODAxNzRiNTMtMmRkYS00OWYxLTlkNmEtNmE3ODBkNBDkBloGEgRhc2Rm"
	msg := fmt.Sprintf("{'bzn-api':'database', 'msg':'%s'}", encodedBase64)
	wsAddr := blz.getWsAddr()

	return wsConnect(wsAddr, msg)
}

func wsConnect(endpoint string, msg string) (string, error) {
	u := url.URL{Scheme: "ws", Host: endpoint}
	log.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", err
	}
	defer c.Close()

	respChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		defer close(respChan)
		defer close(errChan)
		for {
			_, r, err := c.ReadMessage()
			if err != nil {
				log.Println("Got Error")
				errChan <- err
				return
			}
			log.Println("Got msg")
			respChan <- r
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			log.Println(t.String())
			c.WriteMessage(websocket.TextMessage, []byte(msg))
		case resp := <-respChan:
			log.Println(resp)
		case err := <-errChan:
			log.Println(err.Error())
			return "", err
		}
	}
}

func Connect(endpoint string, port uint16, uuid string) *Bluzelle {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if port == 0 {
		port = DefaultPort
	}
	if uuid == "" {
		uuid = DefaultUuid
	}

	return &Bluzelle{
		Endpoint: endpoint,
		Port:     port,
		Uuid:     uuid,
	}
}

func (blz Bluzelle) Read(k string) (string, error) {
	var err error
	var encodedBlzMsgPb []byte
	var resp string

	read := &pb.DatabaseMsg_Read{
		Read: &pb.DatabaseRead{
			Key: k,
		},
	}
	blzMsgPb := &pb.BznMsg{
		Msg: &pb.BznMsg_Db{
			Db: &pb.DatabaseMsg{
				Header: blz.pbHeader(),
				Msg:    read,
			},
		},
	}

	encodedBlzMsgPb, err = proto.Marshal(blzMsgPb)
	if err != nil {
		return "", err
	}

	resp, err = blz.sendRequest(encodedBlzMsgPb)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func main() {

	blz := Connect("testnet.bluzelle.com", 51010, "80174b53-2dda-49f1-9d6a-6a780d4")

	value, err := blz.Read("asdf")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(value)
}
