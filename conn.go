package bluzelle

// Layer 1: Persistent Connection
// https://devel-docs.bluzelle.com/client-development-guide/layers/layer-1-persistent-connection

import (
	"errors"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	// ErrWsConnNotInitialized is returned when writing to a nil websocket.
	ErrWsConnNotInitialized = errors.New("conn: websocket connection not initialized")
)

// Conn represents the persistent layer for Bluzelle.
type Conn struct {
	// Endpoint represents the entry point for the bluzelle network.
	Endpoint string

	// IncomingMsg
	IncomingMsg chan []byte

	wsConn *websocket.Conn

	close chan bool
}

// NewConn creates a new conn
func NewConn(endpoint string) *Conn {
	return &Conn{
		Endpoint:    endpoint,
		IncomingMsg: make(chan []byte),
		wsConn:      nil,
		close:       make(chan bool),
	}
}

// EndpointURL formats the endpoint to a ws protocol.
func (conn *Conn) EndpointURL() string {
	url := url.URL{Scheme: "ws", Host: conn.Endpoint}
	return url.String()
}

// Dial initiates the websocket connection to blz endpoint.
func (conn *Conn) Dial() error {
	url := conn.EndpointURL()
	log.Println("Connecting to: ", url)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	conn.wsConn = c
	conn.setPingPongHandlers()

	go func() {
		for {
			select {
			case <-conn.close:
				conn.closeConn()
				return
			default:
				// TODO: Manage incoming message to match outgoing request.
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
		}
	}()

	log.Printf("Pinging: %s", url)
	conn.sendPingMsg()
	return nil
}

func (conn *Conn) setPingPongHandlers() {
	conn.wsConn.SetPongHandler(func(msg string) error {
		log.Printf("Received pong from: %s", conn.EndpointURL())
		return nil
	})
	conn.wsConn.SetPingHandler(func(msg string) error {
		log.Printf("Received ping from: %s, responding pong", conn.EndpointURL())
		return conn.sendPongMsg()
	})
}

// Close sends a socket close message and closes the connection.
func (conn *Conn) closeConn() {
	// TODO: To remove error log before release.
	if conn.wsConn == nil {
		return
	}
	if err := conn.sendCloseMsg(); err != nil {
		log.Printf("conn: err sending close message: %s", err)
	}
	if err := conn.wsConn.Close(); err != nil {
		log.Printf("conn: err closing connection: %s", err)
	}
	conn.wsConn = nil
}

func (conn *Conn) readMsg() <-chan []byte {
	return conn.IncomingMsg
}

func (conn *Conn) sendMsg(data []byte) error {
	if conn.wsConn == nil {
		return ErrWsConnNotInitialized
	}
	return conn.wsConn.WriteMessage(websocket.TextMessage, data)
}

func (conn *Conn) sendPingMsg() error {
	return conn.wsConn.WriteMessage(websocket.PingMessage, nil)
}

func (conn *Conn) sendPongMsg() error {
	return conn.wsConn.WriteMessage(websocket.PongMessage, nil)
}

func (conn *Conn) sendCloseMsg() error {
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	return conn.wsConn.WriteMessage(websocket.CloseMessage, closeMsg)
}
