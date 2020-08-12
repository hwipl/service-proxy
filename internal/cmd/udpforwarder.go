package cmd

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// udpForwarderMap maps peer addresses to forwarders
type udpForwarderMap struct {
	mutex   sync.Mutex
	srvConn *net.UDPConn
	dstAddr *net.UDPAddr
	fwds    map[string]*udpForwarder
}

// get returns an udpForwarder for peer
func (u *udpForwarderMap) get(peer *net.UDPAddr) *udpForwarder {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	fwd := u.fwds[peer.String()]
	if fwd == nil {
		// create a new forwarder for this peer
		dstConn, err := net.DialUDP("udp", nil, u.dstAddr)
		if err != nil {
			fmt.Println("error creating socket for peer", peer)
			return nil
		}
		newFwd := udpForwarder{
			fwdMap:  u,
			srvConn: u.srvConn,
			dstConn: dstConn,
			peer:    peer,
			srvData: make(chan []byte),
			dstData: make(chan []byte),
		}
		u.fwds[peer.String()] = &newFwd
		go newFwd.runForwarder()
		return &newFwd
	}

	// re-use existing forwarder
	return fwd
}

// del removes the udpForwarder for peer
func (u *udpForwarderMap) del(peer *net.UDPAddr) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	delete(u.fwds, peer.String())
}

// stopAll stops all udpForwarders in the map
func (u *udpForwarderMap) stopAll() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for peer, fwd := range u.fwds {
		// close destination socket and remove element from map
		fwd.dstConn.Close()
		delete(u.fwds, peer)
	}
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
	fwdMap  *udpForwarderMap
	srvConn *net.UDPConn
	dstConn *net.UDPConn
	peer    *net.UDPAddr
	srvData chan []byte
	dstData chan []byte
}

// runForwarder runs the udp forwarder
func (u *udpForwarder) runForwarder() {
	defer u.dstConn.Close()
	defer u.fwdMap.del(u.peer)

	// read data from destination conn to channel
	go udpReadToChannel(u.dstConn, u.dstData)

	// create ticker for detecting dead connection
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// start forwarding traffic
	last := -1
	pkts := 0
	for {
		select {
		case data, more := <-u.srvData:
			if !more {
				// no more data from service connection,
				// close destination connection and stop
				return
			}
			_, err := u.dstConn.Write(data)
			if err != nil {
				fmt.Printf("error sending packet from %s "+
					"to %s\n", u.peer,
					u.dstConn.RemoteAddr())
				return
			}
			pkts++
		case data, more := <-u.dstData:
			if !more {
				// no more data from destination connection,
				// stop here
				return
			}
			_, err := u.srvConn.WriteToUDP(data, u.peer)
			if err != nil {
				fmt.Printf("error sending packet from %s "+
					"to %s\n", u.dstConn.RemoteAddr(),
					u.peer)
				return
			}
			pkts++
		case <-ticker.C:
			if last == pkts {
				// no packets forwarded since last timer tick,
				// assume connection is dead and stop here
				fmt.Printf("Cleaning up udp forwarder "+
					"between peer %s and %s\n", u.peer,
					u.dstConn.RemoteAddr())
				return
			}
			last = pkts
		}
	}
}

// forward forwards a packet from the service peer to the proxy destination
func (u *udpForwarder) forward(b []byte) {
	u.srvData <- b
}

// udpReadToChannel reads data from conn and writes it to channel
// TODO: use/merge with tcpReadToChannel?
func udpReadToChannel(conn net.Conn, channel chan<- []byte) {
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			close(channel)
			return
		}
		channel <- buf[:n]
	}
}
