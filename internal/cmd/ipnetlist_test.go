package cmd

import (
	"log"
	"net"
	"reflect"
	"testing"
)

func testStringsToIPNetList(nets ...string) *ipNetList {
	var ipList ipNetList

	// add strings to ipNetList
	for _, n := range nets {
		_, ipNet, err := net.ParseCIDR(n)
		if err != nil {
			log.Fatal(err)
		}
		ipList.add(ipNet)
	}
	return &ipList
}

func testIPNetListToStrings(ipList *ipNetList) []string {
	var s []string
	for _, i := range ipList.getAll() {
		s = append(s, i.String())
	}
	return s
}

func TestIPNetListAdd(t *testing.T) {
	var addrs, want, got []string
	test := func() {
		got = testIPNetListToStrings(testStringsToIPNetList(addrs...))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, want %s", got, want)
		}
	}

	// add single ip network
	addrs = []string{"192.168.1.1/32"}
	want = addrs
	test()

	// add multiple different ip networks
	addrs = []string{
		"10.0.0.0/8",
		"172.16.0.0/16",
		"192.168.1.0/24",
		"127.0.0.0/32",
	}
	want = addrs
	test()

	// add ip networks containing each-other, network size decreasing
	addrs = []string{
		"10.0.0.0/8",
		"10.0.0.0/16",
		"10.0.0.0/32",
	}
	want = []string{"10.0.0.0/8"}
	test()

	// add ip networks containing each-other, network size increasing
	addrs = []string{
		"10.0.0.0/32",
		"10.0.0.0/16",
		"10.0.0.0/8",
		"10.0.0.0/0",
	}
	want = []string{"0.0.0.0/0"}
	test()
}
