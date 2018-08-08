package bluzelle

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/wlwanpan/bluzelle-go/proto"
)

type Bluzelle struct {
	endpoint string
	port     uint16
	uuid     string
}

func (blz Bluzelle) genAddr() string {
	portStr := fmt.Sprint(blz.port)
	strArr := []string{blz.endpoint, ":", portStr}
	return strings.Join(strArr, "")
}

func (blz Bluzelle) genHeaderPb() *pb.DatabaseHeader {
	return &pb.DatabaseHeader{
		DbUuid:        blz.uuid,
		TransactionId: rand.Uint64(),
	}
}

func Connect(endpoint string, port uint16, uuid string) *Bluzelle {
	return &Bluzelle{
		endpoint: endpoint,
		port:     port,
		uuid:     uuid,
	}
}

func wsConnect(endpoint string, message string) <-chan string {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: endpoint}
	fmt.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	res := make(chan string)

	go func() {
		defer close(res)
		for {
			msgType, msg, err := c.ReadMessage()
			fmt.Println(msgType)
			if err != nil {
				fmt.Println("Error ReadMessage:")
				fmt.Println(fmt.Sprint(err))
				return
			}
			fmt.Println(msg)
			res <- string(msg[:])
			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}()

	return res
}

func (blz *Bluzelle) sendRequest(k string) string {
	prePackProto := "Uj0SKAokODA3OGUxNWMtYWM0Ny00ZGI5LTgzZGYtZGE2ZGJhNzEyMzFhEARSERIFaGVsbG8aCAEid29ybGQi"
	msg := fmt.Sprintf("{'bzn-api':'database', 'msg':%s}", prePackProto)
	msgByte := []byte(msg)
	encodedBase64Msg := base64.StdEncoding.EncodeToString(msgByte)

	wsAddr := blz.genAddr()

	c := wsConnect(wsAddr, encodedBase64Msg)

	return <-c
}

func (blz Bluzelle) Read(k string) string {

	readPb := &pb.DatabaseMsg_Read{
		Read: &pb.DatabaseRead{
			Key: k,
		},
	}

	blzMsgPb := &pb.BznMsg{
		Msg: pb.isBznMsg_Msg{
			Db: pb.DatabaseMsg{
				Header: blz.genHeaderPb(),
				Msg:    readPb,
			},
		},
	}
	encodedBlzMsgPb, err := proto.Marshal(blzMsgPb)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(encodedBlzMsgPb)
	return blz.sendRequest(encodedBlzMsgPb)
}

func main() {

	const Endpoint string = "testnet-dev.bluzelle.com"
	const Port uint16 = 51010
	const Uuid string = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"

	blz := connect(Endpoint, Port, Uuid)
	value := blz.read("asdf")
	fmt.Printf(value)

}
