package network

import (
	"log"
	"net"
)

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

// ReadFromConn reads messageLen bytes from conn
func ReadFromConn(conn net.Conn) []byte {
	buf := make([]byte, MessageLen)
	count := 0
	for count < MessageLen {
		n, err := conn.Read(buf[count:])
		if err != nil {
			log.Printf("Connection to %s: %s\n",
				conn.RemoteAddr(), err)
			return nil
		}
		count += n
	}
	return buf
}
