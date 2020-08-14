package cmd

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
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
	// registerServices is a comma-separated list of services to register
	// on the server
	registerServices = ""
)

func parseTCPAddr(addr string) *net.TCPAddr {
	// parse server address, check if it's a valid tcp address
	cntrlAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		// it's not a valid tcp address, check if it's an ip address
		cntrlIP, err := net.ResolveIPAddr("ip", addr)
		if err != nil {
			log.Fatal("cannot parse server address: ", addr)
		}
		cntrlAddr = &net.TCPAddr{
			IP:   cntrlIP.IP,
			Port: defaultPort,
			Zone: cntrlIP.Zone,
		}
	}
	return cntrlAddr
}

// run in server mode
func runServer() {
	ip := ""
	cntrlAddr := parseTCPAddr(serverAddr)
	if cntrlAddr.IP != nil {
		// user supplied an ip address, update serverIP with it
		serverIP = cntrlAddr.IP
		ip = fmt.Sprintf("%s", cntrlAddr.IP)
	}
	log.Printf("Starting server and listening on %s:%d\n", ip,
		cntrlAddr.Port)
	runControl(cntrlAddr)
}

// run in client mode
func runClient() {
	ip := ""
	cntrlAddr := parseTCPAddr(clientAddr)
	if cntrlAddr.IP != nil {
		ip = fmt.Sprintf("%s", cntrlAddr.IP)
	}
	if cntrlAddr.Port == 0 {
		cntrlAddr.Port = defaultPort
	}

	// parse service specifications in registerServices (format
	// "<protocol>:<port>:<destPort>")
	if registerServices == "" {
		log.Fatal("No services specified")
	}
	services := strings.Split(registerServices, ",")
	var specs []*serviceSpec
	for _, s := range services {
		specs = append(specs, parseServiceSpec(s))
	}
	log.Printf("Starting client and connecting to server %s:%d\n",
		ip, cntrlAddr.Port)

	// connect to server and configure services
	runControlClient(cntrlAddr, specs)
}

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// set command line arguments
	flag.StringVar(&serverAddr, "s", serverAddr,
		"start server (default) and listen on `address`")
	flag.StringVar(&clientAddr, "c", clientAddr,
		"start client and connect to `address`; requires -r")
	flag.StringVar(&registerServices, "r", registerServices,
		"register comma-separated list of `services` on server,\n"+
			"e.g., tcp:8000:80,udp:53000:53000")
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
