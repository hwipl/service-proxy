package pserver

import (
	"net"
	"testing"
	"time"

	"github.com/hwipl/service-proxy/internal/pclient"
)

func TestRunControlServer(t *testing.T) {
	addr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 53535,
	}
	allowedIPs := "127.0.0.1"
	allowedPorts := "tcp:53535-54545"
	go RunControlServer(&addr, nil, allowedIPs, allowedPorts)
	time.Sleep(1 * time.Second)

	// test client with not registered but already used port
	go pclient.RunControlClient(&addr, nil, "tcp:53535:53535")
	time.Sleep(1 * time.Second)

	// test client with not allowed port
	go pclient.RunControlClient(&addr, nil, "tcp:50000:50000")
	time.Sleep(1 * time.Second)

	// test client with allowed port
	go pclient.RunControlClient(&addr, nil, "tcp:53536:53536")
	time.Sleep(1 * time.Second)

	// test client with already registered port
	go pclient.RunControlClient(&addr, nil, "tcp:53536:53536")
	time.Sleep(1 * time.Second)

	// test parallel clients
	go pclient.RunControlClient(&addr, nil, "tcp:53537:53537")
	go pclient.RunControlClient(&addr, nil, "tcp:53538:53538")
	go pclient.RunControlClient(&addr, nil, "tcp:53539:53539")
	time.Sleep(1 * time.Second)
}
