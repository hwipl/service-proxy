package cmd

import (
	"log"
	"net"
	"sync"
)

var (
	// tcpServices stores all active tcp services identified by port
	tcpServices tcpServiceMap
)

// tcpServiceMap stores active tcp services identified by port
type tcpServiceMap struct {
	m sync.Mutex
	s map[int]*tcpService
}

// add adds the service entry identified by port to the tcpServiceMap and
// returns true if successful
func (t tcpServiceMap) add(port int, service *tcpService) bool {
	t.m.Lock()
	defer t.m.Unlock()

	if t.s == nil {
		t.s = make(map[int]*tcpService)
	}
	if t.s[port] == nil {
		t.s[port] = service
		return true
	}
	return false
}

// del removes the service identified by port from the tcpServiceMap
func (t tcpServiceMap) del(port int) {
	t.m.Lock()
	defer t.m.Unlock()

	delete(t.s, port)
}

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
		if err != nil {
			srvConn.Close()
			continue
		}

		// start forwarding traffic between connections
		runTCPForwarder(srvConn, dstConn)
	}
}

// runTCPService runs a tcp service proxy that listens on srvAddr and forwards
// incoming connections to dstAddr
func runTCPService(srvAddr, dstAddr *net.TCPAddr) *tcpService {
	// create service
	srv := tcpService{
		srvAddr: srvAddr,
		dstAddr: dstAddr,
	}

	if tcpServices.add(srvAddr.Port, &srv) {
		// run service
		go srv.runService()
		return &srv
	}

	return nil
}
