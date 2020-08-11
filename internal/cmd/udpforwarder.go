package cmd

import "net"

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
