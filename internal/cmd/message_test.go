package cmd

import (
	"bytes"
	"testing"
)

func TestMessageSerialize(t *testing.T) {
	msg := Message{
		Op:       MessageAdd,
		Protocol: ProtocolUDP,
		Port:     65535,
		DestPort: 65535,
	}

	want := []byte{MessageAdd, ProtocolUDP, 255, 255, 255, 255}
	got := msg.Serialize()
	if bytes.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMessageParse(t *testing.T) {
	msg := []byte{MessageOK, ProtocolTCP, 255, 255, 255, 255}

	want := Message{
		Op:       MessageOK,
		Protocol: ProtocolTCP,
		Port:     65535,
		DestPort: 65535,
	}
	got := Message{}
	got.Parse(msg)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
