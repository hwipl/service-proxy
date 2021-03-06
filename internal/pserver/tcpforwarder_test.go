package pserver

import (
	"bytes"
	"log"
	"net"
	"testing"

	"github.com/hwipl/service-proxy/internal/network"
)

func TestTCPForwarder(t *testing.T) {
	// create listeners
	tcpAddr := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}
	srvListener, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer srvListener.Close()
	dstListener, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer dstListener.Close()

	// create connections
	srvConn, err := net.DialTCP("tcp", &tcpAddr,
		srvListener.Addr().(*net.TCPAddr))
	if err != nil {
		log.Fatal(err)
	}
	defer srvConn.Close()
	dstConn, err := net.DialTCP("tcp", &tcpAddr,
		dstListener.Addr().(*net.TCPAddr))
	if err != nil {
		log.Fatal(err)
	}
	defer dstConn.Close()

	// accept connections
	srvClient, err := srvListener.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}
	defer srvClient.Close()
	dstClient, err := dstListener.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}
	defer dstClient.Close()

	// create forwarder between connections
	forwarder := tcpForwarder{
		srvConn: srvClient,
		dstConn: dstConn,
		srvData: make(chan []byte),
		dstData: make(chan []byte),
	}
	go forwarder.runForwarder()

	// test writing data to service connection and reading from destination
	var want, got []byte
	want = []byte{1, 2, 3, 4, 5, 6}
	network.WriteToConn(srvConn, want)
	got = network.ReadFromConn(dstClient)

	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}

	// test writing to destination connection and reading from service
	want = []byte{6, 5, 4, 3, 2, 1}
	network.WriteToConn(dstClient, want)
	got = network.ReadFromConn(srvConn)

	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}
