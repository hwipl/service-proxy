package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/hwipl/service-proxy/internal/pclient"
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
	// certFile is the certificate file used by this host
	certFile = ""
	// keyFile is the key file for the certificate used by this host
	keyFile = ""
	// caCertFiles is a comma-separated list of ca-certificate files
	caCertFiles = ""
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

func parseCertFiles() tls.Certificate {
	if keyFile == "" {
		log.Fatal("key file for this host's certificate must " +
			"be specified")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("cannot load certificate: ", err)
	}
	return cert
}

func parseCACertFiles() *x509.CertPool {
	files := strings.Split(caCertFiles, ",")
	caCertPool := x509.NewCertPool()
	for _, f := range files {
		caCert, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal("cannot read ca-certificate file: ", err)
		}
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Fatal("cannot parse ca-certificate file: ", f)
		}
	}
	return caCertPool
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
			allowedIPNets.add(a)
		}
	}

	// parse allowed ports
	if allowedPorts != "" {
		aPorts := strings.Split(allowedPorts, ",")
		for _, a := range aPorts {
			allowedPortRanges.add(a)
		}
	}

	// parse certificates
	var tlsConfig *tls.Config
	if certFile != "" {
		cert := parseCertFiles()
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		if caCertFiles != "" {
			clientCAs := parseCACertFiles()
			tlsConfig.ClientCAs = clientCAs
		}
	}

	// output info and start server
	tlsInfo := ""
	if tlsConfig != nil {
		tlsInfo = "in mTLS mode "
	}
	log.Printf("Starting server %sand listening on %s:%d\n", tlsInfo, ip,
		cntrlAddr.Port)
	for _, ipNet := range allowedIPNets.getAll() {
		log.Printf("Allowing control connections from %s\n", ipNet)
	}
	for _, portRange := range allowedPortRanges.getAll() {
		log.Printf("Allowing port range %s in service registrations\n",
			portRange)
	}
	RunControlServer(cntrlAddr, tlsConfig)
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

	// parse certificates
	var tlsConfig *tls.Config
	if certFile != "" {
		cert := parseCertFiles()
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   cntrlAddr.IP.String(),
		}
		if caCertFiles != "" {
			rootCAs := parseCACertFiles()
			tlsConfig.RootCAs = rootCAs
		}
	}
	// check if services are specified by user
	if registerServices == "" {
		log.Fatal("No services specified")
	}

	// print server info and run control client
	tlsInfo := ""
	if tlsConfig != nil {
		tlsInfo = "in mTLS mode "
	}
	log.Printf("Starting client %sand connecting to server %s:%d\n",
		tlsInfo, ip, cntrlAddr.Port)

	// connect to server and configure services
	pclient.RunControlClient(cntrlAddr, tlsConfig, registerServices)
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
	flag.StringVar(&certFile, "cert", certFile,
		"read this host's certificate from `file`, e.g., cert.pem")
	flag.StringVar(&keyFile, "key", keyFile,
		"read the key of this host's certificate from `file`, "+
			"e.g., key.pem")
	flag.StringVar(&caCertFiles, "ca-certs", caCertFiles,
		"read accepted ca-certificates from comma-separated list "+
			"of `files`,\ne.g., cert1.pem,cert2.pem,cert3.pem")
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
