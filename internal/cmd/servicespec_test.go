package cmd

import "testing"

func TestServiceSpecToMessage(t *testing.T) {
	s := serviceSpec{
		protocol: "tcp",
		port:     1024,
		destPort: 1024,
	}
	want := Message{
		Op:       MessageAdd,
		Protocol: ProtocolTCP,
		Port:     1024,
		DestPort: 1024,
	}
	got := s.toMessage()
	if *got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServiceSpecFromMessage(t *testing.T) {
	m := Message{
		Op:       MessageAdd,
		Protocol: ProtocolTCP,
		Port:     1024,
		DestPort: 1024,
	}
	want := serviceSpec{
		protocol: "tcp",
		port:     1024,
		destPort: 1024,
	}
	got := serviceSpec{}
	got.fromMessage(&m)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestServiceSpecString(t *testing.T) {
	s := serviceSpec{
		protocol: "tcp",
		port:     1024,
		destPort: 1024,
	}
	want := "tcp:1024:1024"
	got := s.String()
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestParseServiceSpec(t *testing.T) {
	s := "tcp:1024:1024"
	want := serviceSpec{
		protocol: "tcp",
		port:     1024,
		destPort: 1024,
	}
	got := parseServiceSpec(s)
	if *got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
