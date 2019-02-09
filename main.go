package main

import "log"

func main() {

	entry := "bernoulli.bluzelle.com:51010"
	// uuid := "5f493479–2447–47g6–1c36-efa5d251a283"
	// private_pem := "MHQCAQEEIFNmJHEiGpgITlRwao/CDki4OS7BYeI7nyz+CM8NW3xToAcGBSuBBAAKoUQDQgAEndHOcS6bE1P9xjS/U+SM2a1GbQpPuH9sWNWtNYxZr0JcF+sCS2zsD+xlCcbrRXDZtfeDmgD9tHdWhcZKIy8ejQ=="

	conn := NewConn(entry)

	if err := conn.Dial(); err != nil {
		log.Println(err)
	}

	conn.SendMsg([]byte("asdf"))

	select {
	case msg := <-conn.ReadMsg():
		log.Println(msg)
	}

}
