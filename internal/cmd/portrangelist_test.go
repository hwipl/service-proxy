package cmd

import (
	"reflect"
	"testing"

	"github.com/hwipl/service-proxy/internal/network"
)

func TestPortRangeListAdd(t *testing.T) {
	var got, want []string
	var ports portRangeList
	var test = func() {
		got = []string{}
		for _, p := range ports.getAll() {
			got = append(got, p.String())
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	}

	// add single ports
	ports.add(network.ProtocolTCP, 1024, 1024)
	ports.add(network.ProtocolUDP, 1024, 1024)

	want = []string{
		"tcp:1024-1024",
		"udp:1024-1024",
	}
	test()

	// add port ranges
	ports = portRangeList{}
	ports.add(network.ProtocolTCP, 1024, 65535)
	ports.add(network.ProtocolUDP, 1024, 65535)

	want = []string{
		"tcp:1024-65535",
		"udp:1024-65535",
	}
	test()

	// add port ranges containing each other, size increasing
	ports = portRangeList{}
	ports.add(network.ProtocolTCP, 4096, 8192)
	ports.add(network.ProtocolTCP, 2048, 16384)
	ports.add(network.ProtocolTCP, 1024, 32768)

	want = []string{
		"tcp:1024-32768",
	}
	test()

	// add port ranges containing each other, size decreasing
	ports = portRangeList{}
	ports.add(network.ProtocolTCP, 1024, 32768)
	ports.add(network.ProtocolTCP, 2048, 16384)
	ports.add(network.ProtocolTCP, 4096, 8192)

	want = []string{
		"tcp:1024-32768",
	}
	test()

	// add overlapping port ranges
	ports = portRangeList{}
	ports.add(network.ProtocolTCP, 4096, 8192)
	ports.add(network.ProtocolTCP, 16384, 32768)
	ports.add(network.ProtocolTCP, 2048, 4096)
	ports.add(network.ProtocolTCP, 8192, 16384)
	ports.add(network.ProtocolTCP, 1024, 2048)

	want = []string{
		"tcp:1024-32768",
	}
	test()
}

func TestPortRangeListContainsPort(t *testing.T) {
	var want, got bool
	var ports portRangeList
	var test = func(port uint16) {
		got = ports.containsPort(network.ProtocolTCP, port)
		if got != want {
			t.Errorf("port %d: got %t, want %t", port, got, want)
		}
	}

	// prepare port range list
	ports.add(network.ProtocolTCP, 1024, 4096)

	// test port not in range
	want = false
	test(128)

	// test ports in range
	want = true
	test(1024)
	test(2048)
	test(4096)
}

func TestPortRangeListGetAll(t *testing.T) {
	var want, got string
	var ports portRangeList
	var test = func() {
		got = ""
		for i, p := range ports.getAll() {
			if i != 0 {
				got += " "
			}
			got += p.String()
		}
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}

	// test empty
	want = ""
	test()

	// test filled
	ports.add(network.ProtocolTCP, 1024, 2048)
	ports.add(network.ProtocolTCP, 4096, 8192)
	ports.add(network.ProtocolTCP, 16384, 32768)

	want = "tcp:1024-2048 tcp:4096-8192 tcp:16384-32768"
	test()
}
