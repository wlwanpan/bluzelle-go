package main

// Layer 1: Persistent Connection
// https://devel-docs.bluzelle.com/client-development-guide/layers/layer-1-persistent-connection

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Conn represents the persistent layer for Bluzelle.
type Conn struct {
	// Endpoint represents the entry point for the bluzelle network.
	Endpoint string

	// IncomingMsg
	IncomingMsg chan []byte

	webConn *websocket.Conn
}

// NewConn creates a new conn
func NewConn(endpoint string) *Conn {
	return &Conn{
		Endpoint:    endpoint,
		IncomingMsg: make(chan []byte),
		webConn:     nil,
	}
}

func (conn *Conn) Dial() error {
	u := url.URL{Scheme: "ws", Host: conn.Endpoint}
	log.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	conn.webConn = c
	conn.webConn.SetPongHandler(func(msg string) error {
		log.Printf("From pong handler: %s", msg)
		return nil
	})
	go func() {
		for {
			messageType, r, err := c.ReadMessage()
			if err != nil {
				log.Println("Error from conn:", err)
				log.Println(messageType)
				log.Println(r)
			}
			log.Println(messageType)
			log.Println(r)
			conn.IncomingMsg <- r
		}
	}()

	conn.sendPingMsg()
	return nil
}

func (conn *Conn) readMsg() <-chan []byte {
	return conn.IncomingMsg
}

func (conn *Conn) sendMsg(data []byte) error {
	return conn.webConn.WriteMessage(websocket.TextMessage, data)
}

func (conn *Conn) sendPingMsg() error {
	return conn.webConn.WriteMessage(websocket.PingMessage, []byte{})
}
