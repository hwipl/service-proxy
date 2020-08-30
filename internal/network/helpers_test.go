package network

import (
	"bytes"
	"log"
	"net"
	"testing"
)

func TestWriteToReadFromConn(t *testing.T) {
	in, out := net.Pipe()
	data := []byte{1, 2, 3, 4, 5, 6}
	go func() {
		if !WriteToConn(in, data) {
			log.Fatal("error writing to conn")
		}
	}()
	want := data
	got := ReadFromConn(out)
	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}
