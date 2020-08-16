package cmd

import (
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
	log.Println("Connected to server", c.serverAddr)

	// send service specs to server
	for _, spec := range c.specs {
		log.Printf("Sending service registration %s to server", spec)
		m := spec.toMessage()
		tcpWriteToConn(c.conn, m.serialize())
	}

	// keep connection open
	buf := make([]byte, messageLen)
	for {
		// ignore messages from server for now
		_, err := c.conn.Read(buf)
		if err != nil {
			log.Println("Closing connection to server:", err)
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
