package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/wlwanpan/bluzelle-go/proto"
)

const (
	DefaultUuid     = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"
	DefaultEndpoint = "127.0.0.1"
	DefaultPort     = 51010
)

type BlzApi struct {
	BznApi string `json:"bzn-api"`
	Msg    string `json:"msg"`
}

type Bluzelle struct {
	Endpoint string
	Port     uint32
	Uuid     string
}

func (blz *Bluzelle) SetEndpoint(endpoint string) {
	blz.Endpoint = endpoint
}

func (blz *Bluzelle) SetPort(port uint32) {
	blz.Port = port
}

func (blz *Bluzelle) SetUuid(uuid string) {
	blz.Uuid = uuid
}

func (blz *Bluzelle) getWsAddr() string {
	p := fmt.Sprint(blz.Port)
	strArr := []string{blz.Endpoint, ":", p}
	return strings.Join(strArr, "")
}

func (blz *Bluzelle) pbHeader() *pb.DatabaseHeader {
	return &pb.DatabaseHeader{
		DbUuid:        blz.Uuid,
		TransactionId: rand.Uint64(),
	}
}

func (blz *Bluzelle) sendRequest(m []byte) (*pb.DatabaseResponseResponse, error) {
	encodedBase64 := base64.StdEncoding.EncodeToString(m)
	api := &BlzApi{
		BznApi: "database",
		Msg:    encodedBase64,
	}

	apiJson, err := json.Marshal(api)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	wsAddr := blz.getWsAddr()
	msg := string(apiJson)
	dbResp, err := wsConnect(wsAddr, msg)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	redirect := dbResp.GetRedirect()
	if redirect != nil {
		log.Printf("Switching to leader: %v", redirect.GetLeaderName())
		blz.SetEndpoint(redirect.GetLeaderHost())
		blz.SetPort(redirect.GetLeaderPort())
		return blz.sendRequest(m)
	}

	return dbResp.GetResp(), nil
}

func wsConnect(endpoint string, msg string) (*pb.DatabaseResponse, error) {
	u := url.URL{Scheme: "ws", Host: endpoint}
	log.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &pb.DatabaseResponse{}, err
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
				errChan <- err
				return
			}
			respChan <- r
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(msg))
	for {
		select {
		case resp := <-respChan:
			dbResp := &pb.DatabaseResponse{}
			err := proto.Unmarshal(resp, dbResp)
			if err != nil {
				return dbResp, err
			}
			return dbResp, nil
		case err := <-errChan:
			log.Println(err.Error())
			return &pb.DatabaseResponse{}, err
		}
	}
}

func Connect(endpoint string, port uint32, uuid string) *Bluzelle {
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
	var resp *pb.DatabaseResponseResponse

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
	return string(resp.GetValue()[:]), nil
}

func main() {

	blz := Connect("testnet.bluzelle.com", 51010, "80174b53-2dda-49f1-9d6a-6a780d4")

	value, err := blz.Read("asdf")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(value)
}
