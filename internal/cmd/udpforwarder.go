package cmd

import "net"

// udpForwarder forwards network traffic between a udp service proxy and
// its destination
type udpForwarder struct {
	srvConn *net.UDPConn
	dstConn *net.UDPConn
}

// runForwarder runs the udp forwarder
func (u *udpForwarder) runForwarder() {
	// not implemented
}
