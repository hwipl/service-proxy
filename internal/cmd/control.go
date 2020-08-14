package cmd

import (
	"log"
	"net"
)

// control stores control server information
type control struct {
	addr     *net.TCPAddr
	listener *net.TCPListener
}

// runControl runs the control server
func (c *control) runControl() {
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
			conn.Close()
			continue
		}

		// handle client connection
		handleClient(conn)
	}
}

// runControl runs the control server on addr
func runControl(addr *net.TCPAddr) {
	// create control server
	c := control{
		addr: addr,
	}

	// run control server
	c.runControl()
}
