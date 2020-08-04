package cmd

import (
	"log"
	"net"
)

// tcpService stores tcp service proxy information
type tcpService struct {
	srvAddr  *net.TCPAddr
	listener *net.TCPListener
	dstAddr  *net.TCPAddr
}

// runTCPService runs a tcp service proxy that listens on srvAddr and forwards
// incoming connections to dstAddr
func runTCPService(srvAddr, dstAddr *net.TCPAddr) {
	// create service
	srv := tcpService{
		srvAddr: srvAddr,
		dstAddr: dstAddr,
	}

	// start service socket
	listener, err := net.ListenTCP("tcp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	srv.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// TODO: handle connection
		conn.Close()
	}
}
