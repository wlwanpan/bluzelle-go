package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	DefaultUuid      = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"
	DefaultEndpoint  = "127.0.0.1"
	DefaultPort      = 51010
	MaxRedirectLimit = 3
)

// ErrRedirectLimit is returned when the leader node is switched more
// than the default set limit.
var ErrRedirectLimit = errors.New("Max Leader redirect attempt reached")

// Bluzelle request api struct
type BlzReq struct {
	BznApi string `json:"bzn-api"`
	Msg    string `json:"msg"`
}

type Bluzelle struct {
	Endpoint        string
	Port            uint32
	Uuid            string
	redirectAttempt uint16
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

func (blz *Bluzelle) resetRedirectAttempt() {
	blz.redirectAttempt = 0
}

func (blz *Bluzelle) wsAddr() string {
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

func (blz *Bluzelle) sendRequest(req string) (*pb.DatabaseResponseResponse, error) {
	if blz.redirectAttempt > MaxRedirectLimit {
		blz.resetRedirectAttempt()
		return &pb.DatabaseResponseResponse{}, ErrRedirectLimit
	}

	wsAddr := blz.wsAddr()
	resp, err := wsConnect(wsAddr, req)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	dbResp := &pb.DatabaseResponse{}
	err = proto.Unmarshal(resp, dbResp)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	redirect := dbResp.GetRedirect()
	if redirect != nil {
		log.Printf("Switching to leader: %v", redirect.GetLeaderName())
		blz.SetEndpoint(redirect.GetLeaderHost())
		blz.SetPort(redirect.GetLeaderPort())
		blz.redirectAttempt++
		return blz.sendRequest(req)
	}

	return dbResp.GetResp(), nil
}

func genReq(m []byte) (string, error) {
	encodedBase64 := base64.StdEncoding.EncodeToString(m)
	blzReq := &BlzReq{
		BznApi: "database",
		Msg:    encodedBase64,
	}

	req, err := json.Marshal(blzReq)
	if err != nil {
		return "", err
	}

	return string(req), nil
}

func (blz *Bluzelle) encodeAndSendReq(msg *pb.BznMsg) (*pb.DatabaseResponseResponse, error) {
	encoded, err := proto.Marshal(msg)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	req, err := genReq(encoded)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	return blz.sendRequest(req)
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
		Endpoint:        endpoint,
		Port:            port,
		Uuid:            uuid,
		redirectAttempt: 0,
	}
}

func (blz Bluzelle) Create(k string, v []byte) error {
	createPb := &pb.DatabaseMsg_Create{Create: &pb.DatabaseCreate{Key: k, Value: v}}
	msgPb := &pb.BznMsg{
		Msg: &pb.BznMsg_Db{
			Db: &pb.DatabaseMsg{
				Header: blz.pbHeader(),
				Msg:    createPb,
			},
		},
	}
	resp, err := blz.encodeAndSendReq(msgPb)
	log.Println(resp)
	if err != nil {
		return err
	}
	return nil
}

func (blz Bluzelle) Read(k string) ([]byte, error) {
	readPb := &pb.DatabaseMsg_Read{Read: &pb.DatabaseRead{Key: k}}
	msgPb := &pb.BznMsg{
		Msg: &pb.BznMsg_Db{
			Db: &pb.DatabaseMsg{
				Header: blz.pbHeader(),
				Msg:    readPb,
			},
		},
	}
	resp, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return []byte{}, err
	}
	return resp.GetValue(), nil
}

func (blz Bluzelle) Update(k string, v string) error {
	return nil
}

func (blz Bluzelle) Remove(k string) error {
	return nil
}

func (blz Bluzelle) Has(k string) bool {
	return true
}

func (blz Bluzelle) Keys() []string {
	return []string{}
}

func (blz Bluzelle) Size() uint32 {
	return 0
}

func wsConnect(endpoint string, msg string) ([]byte, error) {
	u := url.URL{Scheme: "ws", Host: endpoint}
	log.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return []byte{}, err
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
			return resp, nil
		case err := <-errChan:
			return []byte{}, err
		}
	}
}

func main() {

	blz := Connect("testnet.bluzelle.com", 51010, "80174b53-2dda-49f1-9d6a-6a780d4")

	value, err := blz.Read("asdf")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(value)
}
