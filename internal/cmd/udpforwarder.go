package cmd

import (
	"fmt"
	"net"
)

// udpForwarderMap maps peer addresses to forwarders
type udpForwarderMap struct {
	srvConn *net.UDPConn
	dstAddr *net.UDPAddr
	fwds    map[string]*udpForwarder
}

// get returns an udpForwarder for peer
func (u *udpForwarderMap) get(peer *net.UDPAddr) *udpForwarder {
	fwd := u.fwds[peer.String()]
	if fwd == nil {
		// create a new forwarder for this peer
		dstConn, err := net.DialUDP("udp", nil, u.dstAddr)
		if err != nil {
			fmt.Println("error creating socket for peer", peer)
			return nil
		}
		newFwd := udpForwarder{
			srvConn: u.srvConn,
			dstConn: dstConn,
			peer:    peer,
		}
		go newFwd.runForwarder()
		return &newFwd
	}

	// re-use existing forwarder
	return fwd
}

// newUDPForwarderMap creates a new udp forwarder for the udp service conn
func newUDPForwarderMap(srvConn *net.UDPConn,
	dstAddr *net.UDPAddr) *udpForwarderMap {
	u := udpForwarderMap{
		srvConn: srvConn,
		dstAddr: dstAddr,
		fwds:    make(map[string]*udpForwarder),
	}
	return &u
}

// udpForwarder forwards network traffic between a udp service proxy and
// its destination
type udpForwarder struct {
	srvConn *net.UDPConn
	dstConn *net.UDPConn
	peer    *net.UDPAddr
}

// runForwarder runs the udp forwarder
func (u *udpForwarder) runForwarder() {
	defer u.dstConn.Close()
	for {
		// read packets from proxy destination and send them to the
		// service peer
		buf := make([]byte, 2048)
		n, err := u.dstConn.Read(buf)
		if err != nil {
			return
		}
		_, err = u.srvConn.WriteToUDP(buf[:n], u.peer)
		if err != nil {
			return
		}
	}
}

// forward forwards a packet from the service peer to the proxy destination
func (u *udpForwarder) forward(b []byte) {
	_, err := u.dstConn.Write(b)
	if err != nil {
		fmt.Printf("error sending packet from %s to %s\n", u.peer,
			u.dstConn.RemoteAddr())
	}
}
