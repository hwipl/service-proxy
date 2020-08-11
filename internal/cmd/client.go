package cmd

import (
	"fmt"
	"net"
)

// client stores control client information
type client struct {
	conn     net.Conn
	addr     *net.TCPAddr
	tcpPorts map[int]bool
	udpPorts map[int]bool
}

// addTCPService adds a tcp service to the client
func (c *client) addTCPService(port, destPort int) {
	fmt.Printf("Adding new tcp service for client %s: forward port %d "+
		"to port %d\n", c.addr, port, destPort)

	// create tcp addresses and start tcp service
	srvAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	}
	dstAddr := net.TCPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}
	if runTCPService(&srvAddr, &dstAddr) != nil {
		c.tcpPorts[port] = true
	}
}

// addUDPService adds an udp service to the client
func (c *client) addUDPService(port, destPort int) {
	fmt.Printf("Adding new udp service for client %s: forward port %d "+
		"to port %d\n", c.addr, port, destPort)

	// create udp addresses and start udp service
	srvAddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	}
	dstAddr := net.UDPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}
	if runUDPService(&srvAddr, &dstAddr) != nil {
		c.udpPorts[port] = true
	}
}

// addService adds a service to the client
func (c *client) addService(protocol uint8, port, destPort uint16) {
	switch protocol {
	case protocolTCP:
		c.addTCPService(int(port), int(destPort))
	case protocolUDP:
		c.addUDPService(int(port), int(destPort))
	default:
		// unknown protocol, stop here
		return
	}
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	defer c.conn.Close()
	defer c.stopClient()
	fmt.Println("New connection from client", c.addr)
	for {
		// read a message from the connection and parse it
		var msg message
		buf := readFromConn(c.conn)
		if buf == nil {
			fmt.Println("Closing connection to client", c.addr)
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
		s := tcpServices.get(port)
		fmt.Printf("Removing a tcp service for client %s: forward "+
			"port %d to port %d\n", c.addr, s.srvAddr.Port,
			s.dstAddr.Port)
		s.stopService()
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
		addr:     conn.RemoteAddr().(*net.TCPAddr),
		tcpPorts: make(map[int]bool),
		udpPorts: make(map[int]bool),
	}
	go c.handleClient()
}
