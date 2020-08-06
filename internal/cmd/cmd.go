package cmd

import "net"

// Run is the main entry point
func Run() {
	// run control server
	cntrlAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 32323,
	}
	runControl(&cntrlAddr)
}
