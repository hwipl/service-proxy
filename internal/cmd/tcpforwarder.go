package cmd

import "net"

// tcpForwarder forwards network traffic between a tcp service proxy and
// its destination
type tcpForwarder struct {
	srvConn net.Conn
	dstConn net.Conn
	srvData chan []byte
	dstData chan []byte
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

	for {
		select {
		case data, more := <-fwd.srvData:
			if !more {
				fwd.srvData = nil
				break
			}
			// copy data from service peer to destination
			fwd.dstConn.Write(data)
		case data, more := <-fwd.dstData:
			if !more {
				fwd.dstData = nil
				break
			}
			// copy data from destination to service peer
			fwd.srvConn.Write(data)
		}

		// if both channels are closed, stop
		if fwd.srvData == nil && fwd.dstData == nil {
			break
		}
	}

	// close everything
	srvConn.Close()
	dstConn.Close()
}
