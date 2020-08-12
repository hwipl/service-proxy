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
	// serverIP is the IP address the server runs services on, derived from
	// the serverAddr
	serverIP net.IP
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
	ip := ""
	if cntrlAddr.IP != nil {
		// user supplied an ip address, update serverIP with it
		serverIP = cntrlAddr.IP
		ip = fmt.Sprintf("%s", cntrlAddr.IP)
	}
	fmt.Printf("Starting server and listening on %s:%d\n", ip,
		cntrlAddr.Port)
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

	// treat remaining non-flag command line arguments as service
	// specification strings with format "<protocol>:<port>:<destPort>"
	if len(flag.Args()) == 0 {
		log.Fatal("No services specified")
	}
	var specs []*serviceSpec
	for _, arg := range flag.Args() {
		specs = append(specs, parseServiceSpec(arg))
	}
	fmt.Printf("Starting client and connecting to server %s:%d\n",
		cntrlAddr.IP, cntrlAddr.Port)

	// connect to server and configure services
	runControlClient(cntrlAddr, specs)
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
	// parse command line arguments
	parseCommandLine()
}
