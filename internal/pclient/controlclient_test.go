package pclient

import (
	"net"
	"testing"
	"time"

	"github.com/hwipl/service-proxy/internal/pserver"
)

// TestRunControlClient runs hacky control client (and control server) tests
func TestRunControlClient(t *testing.T) {
	// define server address and port
	addr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 52525,
	}

	// start a control server and give it some time to complete startup
	go pserver.RunControlServer(&addr, nil, "127.0.0.1", "tcp:52526")
	time.Sleep(1 * time.Second)

	// start a control client with a not allowed port registration and give
	// it some time to complete
	go RunControlClient(&addr, nil, "tcp:52527:52527")
	time.Sleep(1 * time.Second)

	// start a control client with an allowed port registration and give it
	// some time to complete
	go RunControlClient(&addr, nil, "tcp:52526:52526")
	time.Sleep(1 * time.Second)
}
