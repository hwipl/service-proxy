package cmd

import (
	"fmt"
	"log"
	"net"
)

// controlClient stores control client information
type controlClient struct {
	serverAddr *net.TCPAddr
	specs      []*serviceSpec
	conn       net.Conn
}

// runClient runs the control client
func (c *controlClient) runClient() {
	// connect to server
	conn, err := net.DialTCP("tcp", nil, c.serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	defer c.conn.Close()

	// send service specs to server
	for _, spec := range c.specs {
		m := spec.toMessage()
		tcpWriteToConn(c.conn, m.serialize())
	}

	// keep connection open
	buf := make([]byte, messageLen)
	for {
		// ignore messages from server for now
		_, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Closing connection to server:", err)
			return
		}
	}
}

// runControlClient runs the control client
func runControlClient(cntrlAddr *net.TCPAddr, specs []*serviceSpec) {
	c := controlClient{
		serverAddr: cntrlAddr,
		specs:      specs,
	}
	c.runClient()
}
