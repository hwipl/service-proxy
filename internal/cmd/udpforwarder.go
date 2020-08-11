package cmd

import (
	"fmt"
	"net"
)

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
