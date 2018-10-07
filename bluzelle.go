package bluzelle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

// Default const (Might change to swarmdb specs, check open source gitter channel)
const (
	DefaultUuid     = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"
	DefaultEndpoint = "127.0.0.1"
	DefaultPort     = 51010

	ConnTimeout      = 5 * time.Second
	MaxRedirectLimit = 3
)

var (
	// ErrRedirectLimit is returned when the leader node is switched more
	// than the default set limit.
	ErrRedirectLimit = errors.New("Max Leader redirect attempt reached")

	// ErrConnTimeout is returned when websocket connection to the swarm is
	// longer than the default set limit.
	ErrConnTimeout = errors.New("Connection timeout")

	// ErrRecordExists is returned when created an record that already exist
	// on the bluzelle db.
	ErrRecordExists = errors.New("Bluzelle: Record exists")

	// ErrRecordNotFound is returned when updating, removing or reading
	// a record that does not exist on the bluzelle db.
	ErrRecordNotFound = errors.New("Bluzelle: Record not found")

	// ErrValueSizeTooLarge is returned when creating or updating with a value
	// over max limit of 307200 characters.
	ErrValueSizeTooLarge = errors.New("Bluzelle: Value size too large")
)

// The Bluzelle type represents a connection to bluzelle swarmdb.
// CRUD operations are called this connection.
type Bluzelle struct {
	// Websocket addr of leaderhost
	Endpoint string

	// Port of leaderhost
	Port uint32

	// Uuid of db reference on the swarm
	Uuid string

	// redirectAttempt tracks number of leaderhost redirected
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

// Generate websocket addr from endpoint and port
func (blz *Bluzelle) wsAddr() string {
	p := fmt.Sprint(blz.Port)
	strArr := []string{blz.Endpoint, ":", p}
	return strings.Join(strArr, "")
}

// Generate protobuf bluzelle msg from bluzelle struct
func (blz *Bluzelle) pbBznMsg() *pb.BznMsg {
	return &pb.BznMsg{
		Msg: &pb.BznMsg_Db{
			Db: &pb.DatabaseMsg{
				Header: &pb.DatabaseHeader{
					DbUuid:        blz.Uuid,
					TransactionId: rand.Uint64(),
				},
			},
		},
	}
}

func (blz *Bluzelle) sendRequest(req string) (*pb.DatabaseResponseResponse, error) {
	if blz.redirectAttempt > MaxRedirectLimit {
		blz.redirectAttempt = 0
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

	dbRespResp := dbResp.GetResp()
	respErr := dbRespResp.GetError()
	err = parseBlzErr(respErr)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	return dbRespResp, nil
}

func (blz *Bluzelle) encodeAndSendReq(msg *pb.BznMsg) (*pb.DatabaseResponseResponse, error) {
	encoded, err := proto.Marshal(msg)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	encodedBase64 := base64.StdEncoding.EncodeToString(encoded)
	blzReq := &struct {
		BznApi string `json:"bzn-api"`
		Msg    string `json:"msg"`
	}{
		BznApi: "database",
		Msg:    encodedBase64,
	}

	req, err := json.Marshal(blzReq)
	if err != nil {
		return &pb.DatabaseResponseResponse{}, err
	}

	return blz.sendRequest(string(req))
}

// Connect creates a new client connection using the given endpoint, port and db uuid.
// If zero value args are passed, the connection will default back to localhost with
// port 51010.
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
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Create{
		Create: &pb.DatabaseCreate{Key: k, Value: v},
	}

	_, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return err
	}
	return nil
}

func (blz Bluzelle) Read(k string) ([]byte, error) {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Read{
		Read: &pb.DatabaseRead{Key: k},
	}

	resp, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return []byte{}, err
	}
	return resp.GetValue(), nil
}

func (blz Bluzelle) Update(k string, v []byte) error {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Update{
		Update: &pb.DatabaseUpdate{Key: k, Value: v},
	}

	_, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return err
	}
	return nil
}

func (blz Bluzelle) Remove(k string) error {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Delete{
		Delete: &pb.DatabaseDelete{Key: k},
	}

	_, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return err
	}
	return nil
}

func (blz Bluzelle) Has(k string) (bool, error) {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Has{
		Has: &pb.DatabaseHas{Key: k},
	}

	resp, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return false, err
	}
	return resp.GetHas(), nil
}

// Keys returns a array of all the keys saved on the current db uuid.
func (blz Bluzelle) Keys() ([]string, error) {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Keys{
		Keys: &pb.DatabaseEmpty{},
	}

	resp, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return []string{}, err
	}
	return resp.GetKeys(), nil
}

func (blz Bluzelle) Size() (int32, error) {
	msgPb := blz.pbBznMsg()
	msgPb.GetDb().Msg = &pb.DatabaseMsg_Size{
		Size: &pb.DatabaseEmpty{},
	}

	resp, err := blz.encodeAndSendReq(msgPb)
	if err != nil {
		return 0, err
	}
	return resp.GetSize(), nil
}

func parseBlzErr(e string) error {
	switch e {
	case "RECORD_EXISTS":
		return ErrRecordExists
	case "RECORD_NOT_FOUND":
		return ErrRecordNotFound
	case "VALUE_SIZE_TOO_LARGE":
		return ErrValueSizeTooLarge
	}
	return nil
}

func wsConnect(endpoint string, msg string) ([]byte, error) {
	s := time.Now()
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

		diff := time.Now().Sub(s)
		if diff >= ConnTimeout {
			return []byte{}, ErrConnTimeout
		}
	}
}
