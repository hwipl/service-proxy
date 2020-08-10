package cmd

import (
	"flag"
	"fmt"
	"log"
	"net"
)

const (
	// defaultPort is the default port of the control server
	defaultPort = 32323
)

var (
	// serverAddr is the default listen address of the control server
	serverAddr = fmt.Sprintf(":%d", defaultPort)
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
	// parse client address
	cntrlAddr, err := net.ResolveTCPAddr("tcp", clientAddr)
	if err != nil {
		log.Fatal(err)
	}
	if cntrlAddr.IP == nil {
		log.Fatal("Invalid address to connect to as client")
	}
	if cntrlAddr.Port == 0 {
		cntrlAddr.Port = defaultPort
	}
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
