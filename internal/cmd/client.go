package cmd

import (
	"fmt"
	"net"
)

// client stores control client information
type client struct {
	conn net.Conn
	ip   net.IP
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	defer c.conn.Close()
	for {
		// read a message from the connection and parse it
		var msg message
		buf := readFromConn(c.conn)
		if buf == nil {
			fmt.Println("error reading from control client")
			return
		}
		msg.parse(buf)
	}
}

// readFromConn reads messageLen bytes from conn
func readFromConn(conn net.Conn) []byte {
	buf := make([]byte, messageLen)
	count := 0
	for count < messageLen {
		n, err := conn.Read(buf[count:])
		if err != nil {
			return nil
		}
		count += n
	}
	return buf
}

// handleClient handles the client with its control connection conn
func handleClient(conn net.Conn) {
	c := client{
		conn: conn,
		ip:   conn.RemoteAddr().(*net.TCPAddr).IP,
	}
	go c.handleClient()
}
