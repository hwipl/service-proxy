package cmd

import (
	"log"
	"net"
)

// udpService stores udp service proxy information
type udpService struct {
	srvAddr *net.UDPAddr
	conn    *net.UDPConn
	dstAddr *net.UDPAddr
	fwds    *udpForwarderMap
}

// runService runs the udp service proxy
func (u *udpService) runService() {
	// start service socket
	conn, err := net.ListenUDP("udp", u.srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	u.conn = conn
	u.fwds = newUDPForwarderMap(u.conn, u.dstAddr)
	for {
		// read packet from socket
		buf := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		// get forwarder for peer address and forward packet
		fwd := u.fwds.get(addr)
		fwd.forward(buf[:n])
	}
}

// runUDPService runs an udp service proxy that listens on srvAddr and forwards
// incomming packets to dstAddr
func runUDPService(srvAddr, dstAddr *net.UDPAddr) *udpService {
	// create service
	srv := udpService{
		srvAddr: srvAddr,
		dstAddr: dstAddr,
	}

	// run service
	go srv.runService()
	return &srv
}