package network

import "net"

// TCPWriteToConn writes data to conn
func TCPWriteToConn(conn net.Conn, data []byte) bool {
	count := 0
	for count < len(data) {
		n, err := conn.Write(data[count:])
		if err != nil {
			// do more in this case? abort connection?
			return false
		}
		count += n
	}
	return true
}
