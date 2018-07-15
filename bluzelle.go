package main

import (
	"fmt"
	"flag"
	"encoding/base64"
	"github.com/gorilla/websocket"
	// "net/http"
	// "github.com/golang/protobuf/proto"
	// "./internal/proto"
)

type Bluzelle struct {
	endpoint string
	port uint16
	uuid string
}

func wsConnect(endpoint string, msg string) <-chan string {
	c := make(chan string)

	go func() {
		ws := gowebsocket.New(endpoint)
		ws.OnConnected = func(ws gowebsocket.Socket) {
			fmt.Println("OnConnected")
			ws.SendText(msg)
		}
		ws.OnConnectError = func(err error, ws gowebsocket.Socket) {
			fmt.Println("OnConnectError")
			c <- fmt.Sprint(err)
		}
		ws.OnTextMessage = func(res string, ws gowebsocket.Socket) {
			fmt.Println("OnTextMessage")
			c <- res
		}
		ws.OnBinaryMessage = func(data []byte, ws gowebsocket.Socket) {
			fmt.Println("OnBinaryMessage")
		}
		ws.OnDisconnected = func(err error, ws gowebsocket.Socket) {
			fmt.Println("OnDisconnected")
			if err != nil {
				fmt.Println("Error")
				c <- fmt.Sprint(err)
			}
		}
		ws.Connect()
	}()

	return c
}

func (blz *Bluzelle) SendRequest(k string) string {
	// {"bzn-api": "database","msg": encoded64_msg}
	msg := fmt.Sprintf("{'bzn-api':'database', 'msg':%s}", "Uj0SKAokODA3OGUxNWMtYWM0Ny00ZGI5LTgzZGYtZGE2ZGJhNzEyMzFhEARSERIFaGVsbG8aCAEid29ybGQi")
	msgByte := []byte(msg)
	encodedBase64Msg := base64.StdEncoding.EncodeToString(msgByte)
	endpoint := fmt.Sprintf("ws://%s:%s", blz.endpoint, fmt.Sprint(blz.port))

	c := wsConnect(endpoint, encodedBase64Msg)
	fmt.Println(<-c)

	return "string to return"
}

func connect(endpoint string, port uint16, uuid string) Bluzelle {
	return Bluzelle{
		endpoint: endpoint,
		port: port,
		uuid: uuid,
	}
}

// func (blz Bluzelle) create(k string, v string) bool {}

func (blz Bluzelle) read(k string) string {
	testStr := "This is a test key value"
	return blz.SendRequest(testStr)
}

// func (blz Bluzelle) update(k string, v string) bool {}
//
// func (blz Bluzelle) delete(k string) bool {}
//
// func (blz Bluzelle) has(k string) bool {}
//
// func (blz Bluzelle) keys() []string {}
//
// func (blz Bluzelle) size() uint {}

const Endpoint string = "testnet-dev.bluzelle.com"
const Port uint16 = 51010
const Uuid string = "8c073d96-7291-11e8-adc0-fa7ae01bbebc"

func main() {

	blz := connect(Endpoint, Port, Uuid)
	value := blz.read("asdf")
	fmt.Printf(value)
}