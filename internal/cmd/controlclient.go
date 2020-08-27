package cmd

import (
	"crypto/tls"
	"log"
	"net"
	"time"
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
	if tlsConfig != nil {
		c.conn = tls.Client(conn, tlsConfig)
	}
	defer c.conn.Close()
	log.Println("Connected to server", c.serverAddr)

	// send service specs to server
	active := 0
	for _, spec := range c.specs {
		log.Printf("Sending service registration %s to server", spec)
		m := spec.toMessage()
		tcpWriteToConn(c.conn, m.Serialize())

		// read reply messages from server
		var msg Message
		buf := readFromConn(c.conn)
		if buf == nil {
			log.Println("Closing connection to server")
			return
		}
		msg.Parse(buf)

		// handle message types
		var spec serviceSpec
		replyFmt := "Server reply: service registration %s %s\n"
		switch msg.Op {
		case MessageOK:
			spec.fromMessage(&msg)
			log.Printf(replyFmt, &spec, "OK")
			active++
		case MessageErr:
			spec.fromMessage(&msg)
			log.Printf(replyFmt, &spec, "ERROR")
		default:
			// unknown message, stop here
			log.Println("Unknown reply from server, " +
				"closing connection")
			return
		}
	}

	// are any services active on the server?
	if active == 0 {
		log.Println("Could not register any service on the server, " +
			"closing connection")
		return
	}
	log.Printf("Registered %d service(s) on the server, "+
		"keeping connection open", active)

	// keep connection open
	go func() {
		for {
			// send a keep-alive/NOP message every 15 seconds
			time.Sleep(15 * time.Second)
			keepAlive := Message{Op: MessageNop}
			if !tcpWriteToConn(c.conn, keepAlive.Serialize()) {
				return
			}
		}
	}()
	for {
		if readFromConn(c.conn) == nil {
			log.Println("Closing connection to server")
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
