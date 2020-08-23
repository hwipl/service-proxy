package cmd

import (
	"bytes"
	"testing"
)

func TestMessageSerialize(t *testing.T) {
	msg := message{
		Op:       messageAdd,
		Protocol: protocolUDP,
		Port:     65535,
		DestPort: 65535,
	}

	want := []byte{messageAdd, protocolUDP, 255, 255, 255, 255}
	got := msg.serialize()
	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMessageParse(t *testing.T) {
	msg := []byte{messageOK, protocolTCP, 255, 255, 255, 255}

	want := message{
		Op:       messageOK,
		Protocol: protocolTCP,
		Port:     65535,
		DestPort: 65535,
	}
	got := message{}
	got.parse(msg)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
