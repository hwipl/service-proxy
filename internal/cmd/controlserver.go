package cmd

import (
	"crypto/tls"
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
		handleClient(conn, c.tlsConfig)
	}
}

// RunControlServer runs the control server on addr
func RunControlServer(addr *net.TCPAddr, tlsConfig *tls.Config,
	allowedIPs string) {
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

	// run control server
	c.runServer()
}
