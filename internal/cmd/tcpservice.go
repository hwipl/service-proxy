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

// runService runs the tcp service proxy
func (t *tcpService) runService() {
	// start service socket
	listener, err := net.ListenTCP("tcp", t.srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	t.listener = listener
	for {
		// get new service connection
		srvConn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// open connection to proxy destination
		dstConn, err := net.DialTCP("tcp", nil, t.dstAddr)

		// start forwarding traffic between connections
		go runTCPForwarder(srvConn, dstConn)
	}
}

// runTCPService runs a tcp service proxy that listens on srvAddr and forwards
// incoming connections to dstAddr
func runTCPService(srvAddr, dstAddr *net.TCPAddr) {
	// create service
	srv := tcpService{
		srvAddr: srvAddr,
		dstAddr: dstAddr,
	}

	// run service
	go srv.runService()
}
