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
func (t *tcpServiceMap) add(port int, service *tcpService) bool {
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
func (t *tcpServiceMap) del(port int) {
	t.m.Lock()
	defer t.m.Unlock()

	delete(t.s, port)
}

// get gets the service identified by port from the tcpServiceMap
func (t *tcpServiceMap) get(port int) *tcpService {
	t.m.Lock()
	defer t.m.Unlock()

	return t.s[port]
}

// tcpService stores tcp service proxy information
type tcpService struct {
	srvAddr  *net.TCPAddr
	listener *net.TCPListener
	srcAddr  *net.TCPAddr
	dstAddr  *net.TCPAddr
	mutex    *sync.Mutex
	done     bool
}

// runService runs the tcp service proxy
func (t *tcpService) runService() {
	defer t.listener.Close()
	for {
		// get new service connection
		srvConn, err := t.listener.Accept()
		if err != nil {
			if t.getDone() {
				// service is shutting down, ignore errors
				return
			}
			log.Fatal(err)
		}

		// open connection to proxy destination
		dstConn, err := net.DialTCP("tcp", t.srcAddr, t.dstAddr)
		if err != nil {
			srvConn.Close()
			continue
		}

		// start forwarding traffic between connections
		runTCPForwarder(srvConn, dstConn)
	}
}

// setDone marks the service as done
func (t *tcpService) setDone() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.done = true
}

// getDone checks if the service is done
func (t *tcpService) getDone() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.done
}

// stopService stops the tcp service proxy
func (t *tcpService) stopService() {
	// set service to done and close its listener
	t.setDone()
	t.listener.Close()
	// this is pretty gracefully. active forwarder connections will remain
	// open until they are done. also close active forwarders?
}

// runTCPService runs a tcp service proxy that listens on srvAddr and forwards
// incoming connections to dstAddr from srcAddr
func runTCPService(srvAddr, srcAddr, dstAddr *net.TCPAddr) *tcpService {
	// create service
	listener, err := net.ListenTCP("tcp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	srv := tcpService{
		srvAddr:  srvAddr,
		srcAddr:  srcAddr,
		dstAddr:  dstAddr,
		listener: listener,
		mutex:    &sync.Mutex{},
	}

	if tcpServices.add(srvAddr.Port, &srv) {
		// run service
		go srv.runService()
		return &srv
	}

	return nil
}
