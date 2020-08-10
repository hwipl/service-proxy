package cmd

import (
	"flag"
	"log"
	"net"
)

var (
	// serverAddr is the default listen address of the control server
	serverAddr = ":32323"
	// clientAddr is the address of a control server the client connects to
	clientAddr = ""
)

// run in server mode
func runServer() {
	// parse server address
	cntrlAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	runControl(cntrlAddr)
}

// run in client mode
func runClient() {
	// not implemented
}

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// set command line arguments
	flag.StringVar(&serverAddr, "s", serverAddr,
		"start server (default) and listen on `address`")
	flag.StringVar(&clientAddr, "c", clientAddr,
		"start client and connect to `address`")
	flag.Parse()

	// if client address is specified on the command line, run as client
	if clientAddr != "" {
		runClient()
		return
	}
	// otherwise run as server
	runServer()
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
