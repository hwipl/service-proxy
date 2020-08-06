package cmd

import "net"

// client stores control client information
type client struct {
	conn net.Conn
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	// TODO: do something with conn
	c.conn.Close()
}

// handleClient handles the client with its control connection conn
func handleClient(conn net.Conn) {
	c := client{
		conn: conn,
	}
	go c.handleClient()
}
