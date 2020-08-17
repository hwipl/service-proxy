package cmd

import (
	"log"
	"net"
)

// client stores control client information
type client struct {
	conn     net.Conn
	addr     *net.TCPAddr
	laddr    *net.TCPAddr
	tcpPorts map[int]bool
	udpPorts map[int]bool
}

// addTCPService adds a tcp service to the client
func (c *client) addTCPService(port, destPort int) bool {
	log.Printf("Adding new service for client %s: forward tcp port %d "+
		"to port %d\n", c.addr, port, destPort)

	// create tcp addresses and start tcp service
	srvAddr := net.TCPAddr{
		IP:   serverIP,
		Port: port,
	}
	srcAddr := net.TCPAddr{
		IP: c.laddr.IP,
	}
	dstAddr := net.TCPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}
	if runTCPService(&srvAddr, &srcAddr, &dstAddr) == nil {
		return false
	}
	c.tcpPorts[port] = true
	return true
}

// addUDPService adds an udp service to the client
func (c *client) addUDPService(port, destPort int) bool {
	log.Printf("Adding new service for client %s: forward udp port %d "+
		"to port %d\n", c.addr, port, destPort)

	// create udp addresses and start udp service
	srvAddr := net.UDPAddr{
		IP:   serverIP,
		Port: port,
	}
	srcAddr := net.UDPAddr{
		IP: c.laddr.IP,
	}
	dstAddr := net.UDPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}
	if runUDPService(&srvAddr, &srcAddr, &dstAddr) == nil {
		return false
	}
	c.udpPorts[port] = true
	return true
}

// addService adds a service to the client
func (c *client) addService(protocol uint8, port, destPort uint16) bool {
	// check if port is allowed
	if !allowedPortRanges.containsPort(protocol, port) {
		return false
	}

	// start service
	switch protocol {
	case protocolTCP:
		return c.addTCPService(int(port), int(destPort))
	case protocolUDP:
		return c.addUDPService(int(port), int(destPort))
	default:
		// unknown protocol, stop here
		return false
	}
}

// handleAddMsg handles the client's add message
func (c *client) handleAddMsg(msg *message) bool {
	// try to add service
	if ok := c.addService(msg.Protocol, msg.Port, msg.DestPort); ok {
		msg.Op = messageOK
	} else {
		msg.Op = messageErr
	}

	// send result back to client
	if !writeToConn(c.conn, msg.serialize()) {
		return false
	}
	return true
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	defer c.conn.Close()
	defer c.stopClient()
	log.Println("New connection from client", c.addr)
	for {
		// read a message from the connection and parse it
		var msg message
		buf := readFromConn(c.conn)
		if buf == nil {
			log.Println("Closing connection to client", c.addr)
			return
		}
		msg.parse(buf)

		// handle message types
		switch msg.Op {
		case messageAdd:
			if !c.handleAddMsg(&msg) {
				return
			}
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
		log.Printf("Removing a service for client %s: forward tcp "+
			"port %d to port %d\n", c.addr, s.srvAddr.Port,
			s.dstAddr.Port)
		s.stopService()
		tcpServices.del(port)
	}
	for port := range c.udpPorts {
		s := udpServices.get(port)
		log.Printf("Removing a service for client %s: forward udp "+
			"port %d to port %d\n", c.addr, s.srvAddr.Port,
			s.dstAddr.Port)
		s.stopService()
		udpServices.del(port)
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

// writeToConn writes buf to conn
func writeToConn(conn net.Conn, buf []byte) bool {
	count := 0
	for count < len(buf) {
		n, err := conn.Write(buf[count:])
		if err != nil {
			return false
		}
		count += n
	}
	return true
}

// handleClient handles the client with its control connection conn
func handleClient(conn net.Conn) {
	c := client{
		conn:     conn,
		addr:     conn.RemoteAddr().(*net.TCPAddr),
		laddr:    conn.LocalAddr().(*net.TCPAddr),
		tcpPorts: make(map[int]bool),
		udpPorts: make(map[int]bool),
	}
	go c.handleClient()
}
