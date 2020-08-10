package cmd

import (
	"fmt"
	"net"
)

// client stores control client information
type client struct {
	conn     net.Conn
	ip       net.IP
	tcpPorts map[int]bool
}

// addTCPService adds a tcp service to the client
func (c *client) addTCPService(port, destPort int) {
	fmt.Printf("adding new tcp service for client %s: forward port %d "+
		"to port %d\n", c.ip, port, destPort)

	// create tcp addresses and start tcp service
	srvAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	}
	dstAddr := net.TCPAddr{
		IP:   c.ip,
		Port: destPort,
	}
	if runTCPService(&srvAddr, &dstAddr) != nil {
		c.tcpPorts[port] = true
	}
}

// addService adds a service to the client
func (c *client) addService(protocol uint8, port, destPort uint16) {
	switch protocol {
	case protocolTCP:
		c.addTCPService(int(port), int(destPort))
	case protocolUDP:
		// not implemented
	default:
		// unknown protocol, stop here
		return
	}
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	defer c.conn.Close()
	defer c.stopClient()
	for {
		// read a message from the connection and parse it
		var msg message
		buf := readFromConn(c.conn)
		if buf == nil {
			fmt.Println("error reading from control client")
			return
		}
		msg.parse(buf)

		// handle message types
		switch msg.Op {
		case messageAdd:
			c.addService(msg.Protocol, msg.Port, msg.DestPort)
		case messageDel:
			// not implemented
		default:
			// unknown message, stop here
			return
		}
	}
}

// stopClient stops active client services
func (c *client) stopClient() {
	for port := range c.tcpPorts {
		tcpServices.get(port).stopService()
		tcpServices.del(port)
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
		conn:     conn,
		ip:       conn.RemoteAddr().(*net.TCPAddr).IP,
		tcpPorts: make(map[int]bool),
	}
	go c.handleClient()
}
