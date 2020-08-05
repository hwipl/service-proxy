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

// tcpWriteToConn writes data to conn
func tcpWriteToConn(conn net.Conn, data []byte) {
	count := 0
	for count < len(data) {
		n, err := conn.Write(data[count:])
		if err != nil {
			// do more in this case? abort connection?
			return
		}
		count += n
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

	// read data from connections to channels
	go tcpReadToChannel(fwd.srvConn, fwd.srvData)
	go tcpReadToChannel(fwd.dstConn, fwd.dstData)

	// start forwarding traffic
	for {
		select {
		case data, more := <-fwd.srvData:
			if !more {
				// no more data from service connection,
				// disable channel, close reading side of
				// service connection and close writing side of
				// destination connection
				fwd.srvData = nil
				fwd.srvConn.(*net.TCPConn).CloseRead()
				fwd.dstConn.(*net.TCPConn).CloseWrite()
				break
			}
			// copy data from service peer to destination
			tcpWriteToConn(fwd.dstConn, data)
		case data, more := <-fwd.dstData:
			if !more {
				// no more data from destination connection,
				// disable channel, close reading side of
				// destination connection and close writing
				// side of service connection
				fwd.dstData = nil
				fwd.dstConn.(*net.TCPConn).CloseRead()
				fwd.srvConn.(*net.TCPConn).CloseWrite()
				break
			}
			// copy data from destination to service peer
			tcpWriteToConn(fwd.srvConn, data)
		}

		// if both channels are closed, stop
		if fwd.srvData == nil && fwd.dstData == nil {
			break
		}
	}
}
