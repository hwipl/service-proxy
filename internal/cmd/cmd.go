package cmd

import "net"

// Run is the main entry point
func Run() {
	// test service proxying
	srvAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 33000,
	}
	dstAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8000,
	}
	runTCPService(&srvAddr, &dstAddr)

	// run control server
	cntrlAddr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 32323,
	}
	runControl(&cntrlAddr)
}
