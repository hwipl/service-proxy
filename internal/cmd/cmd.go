package cmd

import (
	"flag"
	"log"
	"net"
)

var (
	serverAddr = ":32323"
)

// run in server mode
func runServer(cntrlAddr *net.TCPAddr) {
	runControl(cntrlAddr)
}

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// set command line arguments
	flag.StringVar(&serverAddr, "s", serverAddr,
		"start server (default) and listen on `address`")
	flag.Parse()

	// parse server address
	cntrlAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	runServer(cntrlAddr)
}

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

	// parse command line arguments
	parseCommandLine()
}
