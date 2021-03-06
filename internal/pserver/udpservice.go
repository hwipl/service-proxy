package pserver

import (
	"log"
	"net"
	"sync"
)

var (
	// udpServices stores all active udp services identified by port
	udpServices udpServiceMap
)

// udpServiceMap stores active udp services identified by port
type udpServiceMap struct {
	m sync.Mutex
	u map[int]*udpService
}

// add adds the service entry identified by port to the udpServiceMap and
// returns true if successful
func (u *udpServiceMap) add(port int, service *udpService) bool {
	u.m.Lock()
	defer u.m.Unlock()

	if u.u == nil {
		u.u = make(map[int]*udpService)
	}
	if u.u[port] == nil {
		u.u[port] = service
		return true
	}
	return false
}

// del removes the service identified by port from the udpServiceMap
func (u *udpServiceMap) del(port int) {
	u.m.Lock()
	defer u.m.Unlock()

	delete(u.u, port)
}

// get gets the service identified by port from the udpServiceMap
func (u *udpServiceMap) get(port int) *udpService {
	u.m.Lock()
	defer u.m.Unlock()

	return u.u[port]
}

// udpService stores udp service proxy information
type udpService struct {
	srvAddr *net.UDPAddr
	conn    *net.UDPConn
	srcAddr *net.UDPAddr
	dstAddr *net.UDPAddr
	fwds    *udpForwarderMap
}

// runService runs the udp service proxy
func (u *udpService) runService() {
	defer u.conn.Close()
	u.fwds = newUDPForwarderMap(u.conn, u.srcAddr, u.dstAddr)
	for {
		// read packet from socket
		buf := make([]byte, 2048)
		n, addr, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		// get forwarder for peer address and forward packet
		fwd := u.fwds.get(addr)
		fwd.forward(buf[:n])
	}
}

// stopService stops the udp service proxy
func (u *udpService) stopService() {
	u.conn.Close()
	u.fwds.stopAll()
}

// runUDPService runs an udp service proxy that listens on srvAddr and forwards
// incomming packets to dstAddr
func runUDPService(srvAddr, srcAddr, dstAddr *net.UDPAddr) *udpService {
	// create service
	srv := udpService{
		srvAddr: srvAddr,
		srcAddr: srcAddr,
		dstAddr: dstAddr,
	}

	if udpServices.add(srvAddr.Port, &srv) {
		// create udp listener/udp conn
		conn, err := net.ListenUDP("udp", srvAddr)
		if err != nil {
			log.Printf("Could not create udp service %s<->%s: %s\n",
				srvAddr, dstAddr, err)
			udpServices.del(srvAddr.Port)
			return nil
		}
		srv.conn = conn

		// run service
		go srv.runService()
		return &srv
	}

	log.Printf("Could not create udp service %s<->%s: service already "+
		"active\n", srvAddr, dstAddr)
	return nil
}
