package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

// controlServer stores controlServer server information
type controlServer struct {
	addr      *net.TCPAddr
	tlsConfig *tls.Config
	listener  *net.TCPListener
}

// runServer runs the control server
func (c *controlServer) runServer() {
	// start control socket
	listener, err := net.ListenTCP("tcp", c.addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	c.listener = listener
	for {
		// get new control connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// if connection is not from an allowed ip, drop it
		ip := conn.RemoteAddr().(*net.TCPAddr).IP
		if !allowedIPNets.containsIP(ip) {
			log.Printf("Dropping new connection from %s: "+
				"IP not allowed\n", conn.RemoteAddr())
			conn.Close()
			continue
		}

		// handle client connection
		handleClient(conn, c.tlsConfig, c.addr.IP)
	}
}

// RunControlServer runs the control server on addr
func RunControlServer(addr *net.TCPAddr, tlsConfig *tls.Config,
	allowedIPs, allowedPorts string) {
	// create control server
	c := controlServer{
		addr:      addr,
		tlsConfig: tlsConfig,
	}

	// parse allowed IP addresses
	if allowedIPs != "" {
		aIP := strings.Split(allowedIPs, ",")
		for _, a := range aIP {
			allowedIPNets.add(a)
		}
	}

	// parse allowed ports
	if allowedPorts != "" {
		aPorts := strings.Split(allowedPorts, ",")
		for _, a := range aPorts {
			allowedPortRanges.add(a)
		}
	}

	// output info and run control server
	ip := ""
	if addr.IP != nil {
		ip = fmt.Sprintf("%s", addr.IP)
	}
	tlsInfo := ""
	if tlsConfig != nil {
		tlsInfo = "in mTLS mode "
	}
	log.Printf("Starting server %sand listening on %s:%d\n", tlsInfo, ip,
		addr.Port)
	for _, ipNet := range allowedIPNets.getAll() {
		log.Printf("Allowing control connections from %s\n", ipNet)
	}
	for _, portRange := range allowedPortRanges.getAll() {
		log.Printf("Allowing port range %s in service registrations\n",
			portRange)
	}

	c.runServer()
}
