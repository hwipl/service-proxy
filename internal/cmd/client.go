package cmd

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/hwipl/service-proxy/internal/network"
)

// client stores control client information
type client struct {
	conn  net.Conn
	addr  *net.TCPAddr
	laddr *net.TCPAddr
	// serverIP is the IP address the server runs services on
	serverIP net.IP
	tcpPorts map[int]bool
	udpPorts map[int]bool
}

// addTCPService adds a tcp service to the client
func (c *client) addTCPService(port, destPort int) bool {
	log.Printf("Adding new service for client %s: forward tcp port %d "+
		"to port %d\n", c.addr, port, destPort)

	// create tcp addresses
	srvAddr := net.TCPAddr{
		IP:   c.serverIP,
		Port: port,
	}
	srcAddr := net.TCPAddr{
		IP: c.laddr.IP,
	}
	dstAddr := net.TCPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}

	// check if port is allowed
	if !allowedPortRanges.containsPort(network.ProtocolTCP, uint16(port)) {
		log.Printf("Could not create tcp service %s<->%s: "+
			"port not allowed\n", &srvAddr, &dstAddr)
		return false
	}

	// start tcp service
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

	// create udp addresses
	srvAddr := net.UDPAddr{
		IP:   c.serverIP,
		Port: port,
	}
	srcAddr := net.UDPAddr{
		IP: c.laddr.IP,
	}
	dstAddr := net.UDPAddr{
		IP:   c.addr.IP,
		Port: destPort,
	}

	// check if port is allowed
	if !allowedPortRanges.containsPort(network.ProtocolUDP, uint16(port)) {
		log.Printf("Could not create udp service %s<->%s: "+
			"port not allowed\n", &srvAddr, &dstAddr)
		return false
	}

	// start udp service
	if runUDPService(&srvAddr, &srcAddr, &dstAddr) == nil {
		return false
	}
	c.udpPorts[port] = true
	return true
}

// addService adds a service to the client
func (c *client) addService(protocol uint8, port, destPort uint16) bool {
	// start service
	switch protocol {
	case network.ProtocolTCP:
		return c.addTCPService(int(port), int(destPort))
	case network.ProtocolUDP:
		return c.addUDPService(int(port), int(destPort))
	default:
		// unknown protocol, stop here
		return false
	}
}

// handleAddMsg handles the client's add message
func (c *client) handleAddMsg(msg *network.Message) bool {
	// try to add service
	if ok := c.addService(msg.Protocol, msg.Port, msg.DestPort); ok {
		msg.Op = network.MessageOK
	} else {
		msg.Op = network.MessageErr
	}

	// send result back to client
	if !network.WriteToConn(c.conn, msg.Serialize()) {
		return false
	}
	return true
}

// handleClient handles the client and its control connection
func (c *client) handleClient() {
	defer c.conn.Close()
	defer c.stopClient()
	for {
		// read a message from the connection and parse it; if there is
		// no message within 30s, assume client is dead and stop
		var msg network.Message
		c.conn.SetDeadline(time.Now().Add(30 * time.Second))
		buf := network.ReadFromConn(c.conn)
		if buf == nil {
			log.Println("Closing connection to client", c.addr)
			return
		}
		msg.Parse(buf)

		// handle message types
		switch msg.Op {
		case network.MessageAdd:
			if !c.handleAddMsg(&msg) {
				return
			}
		case network.MessageDel:
			// not implemented
		case network.MessageNop:
			// just ignore NOP
		default:
			// unknown message, stop here
			log.Println("Unknown message from client", c.addr)
			log.Println("Closing connection to client", c.addr)
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

// handleClient handles the client with its control connection conn
func handleClient(conn net.Conn, tlsConfig *tls.Config, serverIP net.IP) {
	c := client{
		conn:     conn,
		addr:     conn.RemoteAddr().(*net.TCPAddr),
		laddr:    conn.LocalAddr().(*net.TCPAddr),
		serverIP: serverIP,
		tcpPorts: make(map[int]bool),
		udpPorts: make(map[int]bool),
	}
	tlsInfo := ""
	if tlsConfig != nil {
		tlsConn := tls.Server(conn, tlsConfig)
		tlsConn.SetDeadline(time.Now().Add(15 * time.Second))
		if err := tlsConn.Handshake(); err != nil {
			log.Println("TLS handshake with client", c.addr,
				"failed:", err)
			tlsConn.Close()
			return
		}
		clientCert := tlsConn.ConnectionState().PeerCertificates[0]
		commonName := clientCert.Subject.CommonName
		tlsInfo = " (CN=" + commonName + ")"
		c.conn = tlsConn
	}
	log.Printf("New connection from client %s%s\n", c.addr, tlsInfo)
	go c.handleClient()
}
