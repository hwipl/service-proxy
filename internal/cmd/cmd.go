package cmd

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
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
	// allowedIPs is a comma-separated list of all IPs allowed to connect
	// to the server
	allowedIPs = "0.0.0.0/0"
	// allowedPorts is a comma-separated list of protocol and port (range)
	// pairs, that are allowed as services on the server
	allowedPorts = "udp:1024-65535,tcp:1024-65535"
	// certFiles are the comma-separated certificate and key files used by
	// this host
	certFiles = ""
	// tlsConfig contains the tls config
	tlsConfig *tls.Config
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

func parseAllowedIP(addr string) {
	// check if it is a cidr address
	ip, ipNet, err := net.ParseCIDR(addr)
	if err != nil {
		// not cidr, check if we can parse it as regular ip
		ip = net.ParseIP(addr)
		if ip == nil {
			log.Fatal("cannot parse allowed IP: ", addr)
		}
		// create ip net
		netmask := net.CIDRMask(32, 32)
		if ip.To4() == nil { // ipv6 address
			netmask = net.CIDRMask(128, 128)
		}
		ipNet = &net.IPNet{
			IP:   ip,
			Mask: netmask,
		}
	}
	allowedIPNets.add(ipNet)
}

func parseAllowedPort(port string) {
	// get protocol and port range
	protPorts := strings.Split(port, ":")
	if len(protPorts) != 2 {
		log.Fatal("cannot parse allowed port: ", port)
	}

	// parse protocol
	protocol := uint8(0)
	switch protPorts[0] {
	case "tcp":
		protocol = protocolTCP
	case "udp":
		protocol = protocolUDP
	default:
		log.Fatal("unknown protocol in allowed port: ", port)
	}

	// get min and max port from port range
	minmax := strings.Split(protPorts[1], "-")
	if len(minmax) < 1 || len(minmax) > 2 {
		log.Fatal("cannot parse allowed port: ", port)
	}
	min, err := strconv.ParseUint(minmax[0], 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	getMax := func() string {
		if len(minmax) == 2 {
			return minmax[1]
		} else {
			return minmax[0]
		}
	}
	max, err := strconv.ParseUint(getMax(), 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	if min > max {
		min, max = max, min
	}

	// add port range to allowed port ranges
	allowedPortRanges.add(protocol, uint16(min), uint16(max))
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

	// parse allowed IP addresses
	if allowedIPs != "" {
		aIP := strings.Split(allowedIPs, ",")
		for _, a := range aIP {
			parseAllowedIP(a)
		}
	}

	// parse allowed ports
	if allowedPorts != "" {
		aPorts := strings.Split(allowedPorts, ",")
		for _, a := range aPorts {
			parseAllowedPort(a)
		}
	}

	// output info and start server
	log.Printf("Starting server and listening on %s:%d\n", ip,
		cntrlAddr.Port)
	for _, ipNet := range allowedIPNets.getAll() {
		log.Printf("Allowing control connections from %s\n", ipNet)
	}
	for _, portRange := range allowedPortRanges.getAll() {
		log.Printf("Allowing port range %s in service registrations\n",
			portRange)
	}
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
	flag.StringVar(&allowedIPs, "allowed-ips", allowedIPs,
		"set comma-separated list of `IPs` the server accepts\n"+
			"service registrations from, e.g.:\n"+
			"127.0.0.1,192.168.1.0/24")
	flag.StringVar(&allowedPorts, "allowed-ports", allowedPorts,
		"set comma-separated list of `ports` the server accepts\n"+
			"in service registrations, e.g.:\n"+
			"udp:2048-65000,tcp:8000")
	flag.StringVar(&certFiles, "cert", certFiles,
		"read this host's certificate and key from comma-separated "+
			"`files`,\ne.g., cert.pem,key.pem")
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
