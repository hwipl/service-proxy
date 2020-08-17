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
	for {
		// read reply messages from server
		var msg message
		buf := readFromConn(c.conn)
		if buf == nil {
			log.Println("Closing connection to server")
			return
		}
		msg.parse(buf)

		// handle message types
		var spec serviceSpec
		switch msg.Op {
		case messageOK:
			spec.fromMessage(&msg)
			log.Println("Server replied OK for service registration", &spec)
		case messageErr:
			spec.fromMessage(&msg)
			log.Println("Server replied Error for service registration", &spec)
		default:
			// unknown message, stop here
			log.Println("Unknown reply from server, " +
				"closing connection")
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
