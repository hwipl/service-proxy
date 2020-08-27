package cmd

import (
	"net"

	"github.com/hwipl/service-proxy/internal/network"
)

// tcpForwarder forwards network traffic between a tcp service proxy and
// its destination
type tcpForwarder struct {
	srvConn net.Conn
	dstConn net.Conn
	srvData chan []byte
	dstData chan []byte
}

// runForwarder runs the tcp forwarder
func (t *tcpForwarder) runForwarder() {
	// read data from connections to channels
	go tcpReadToChannel(t.srvConn, t.srvData)
	go tcpReadToChannel(t.dstConn, t.dstData)

	// start forwarding traffic
	for {
		select {
		case data, more := <-t.srvData:
			if !more {
				// no more data from service connection,
				// disable channel, close reading side of
				// service connection and close writing side of
				// destination connection
				t.srvData = nil
				t.srvConn.(*net.TCPConn).CloseRead()
				t.dstConn.(*net.TCPConn).CloseWrite()
				break
			}
			// copy data from service peer to destination
			network.TCPWriteToConn(t.dstConn, data)
		case data, more := <-t.dstData:
			if !more {
				// no more data from destination connection,
				// disable channel, close reading side of
				// destination connection and close writing
				// side of service connection
				t.dstData = nil
				t.dstConn.(*net.TCPConn).CloseRead()
				t.srvConn.(*net.TCPConn).CloseWrite()
				break
			}
			// copy data from destination to service peer
			network.TCPWriteToConn(t.srvConn, data)
		}

		// if both channels are closed, stop
		if t.srvData == nil && t.dstData == nil {
			break
		}
	}

}

// tcpReadToChannel reads data from conn and writes it to channel
func tcpReadToChannel(conn net.Conn, channel chan<- []byte) {
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			channel <- data
		}
		if err != nil {
			close(channel)
			break
		}
	}
}

// runTCPForwarder starts forwarding traffic between a connection to the
// service proxy and a connection to the destination
func runTCPForwarder(srvConn, dstConn net.Conn) {
	fwd := tcpForwarder{
		srvConn: srvConn,
		dstConn: dstConn,
		srvData: make(chan []byte),
		dstData: make(chan []byte),
	}
	go fwd.runForwarder()
}
