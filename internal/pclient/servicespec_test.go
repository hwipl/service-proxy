package pclient

import (
	"testing"

	"github.com/hwipl/service-proxy/internal/network"
)

func TestServiceSpecToMessage(t *testing.T) {
	s := ServiceSpec{
		Protocol: "tcp",
		Port:     1024,
		DestPort: 1024,
	}
	want := network.Message{
		Op:       network.MessageAdd,
		Protocol: network.ProtocolTCP,
		Port:     1024,
		DestPort: 1024,
	}
	got := s.ToMessage()
	if *got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServiceSpecFromMessage(t *testing.T) {
	m := network.Message{
		Op:       network.MessageAdd,
		Protocol: network.ProtocolTCP,
		Port:     1024,
		DestPort: 1024,
	}
	want := ServiceSpec{
		Protocol: "tcp",
		Port:     1024,
		DestPort: 1024,
	}
	got := ServiceSpec{}
	got.FromMessage(&m)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServiceSpecString(t *testing.T) {
	s := ServiceSpec{
		Protocol: "tcp",
		Port:     1024,
		DestPort: 1024,
	}
	want := "tcp:1024:1024"
	got := s.String()
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestParseServiceSpec(t *testing.T) {
	s := "tcp:1024:1024"
	want := ServiceSpec{
		Protocol: "tcp",
		Port:     1024,
		DestPort: 1024,
	}
	got := ParseServiceSpec(s)
	if *got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
