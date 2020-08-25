package cmd

import "testing"

func TestParseTCPAddr(t *testing.T) {
	parseAllowedIP("127.0.0.1")
	parseAllowedIP("192.168.1.0/24")
	parseAllowedIP("2000::1")
	parseAllowedIP("fe80::1/64")

	want := "127.0.0.1/32\n" +
		"192.168.1.0/24\n" +
		"2000::1/128\n" +
		"fe80::/64"
	got := ""
	for i, n := range allowedIPNets.getAll() {
		if i > 0 {
			got += "\n"
		}
		got += n.String()
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestParseAllowedPort(t *testing.T) {
	parseAllowedPort("udp:1024")
	parseAllowedPort("tcp:1024-2048")
	parseAllowedPort("tcp:8192-4096")

	want := "udp:1024-1024\n" +
		"tcp:1024-2048\n" +
		"tcp:4096-8192"
	got := ""
	for i, p := range allowedPortRanges.getAll() {
		if i > 0 {
			got += "\n"
		}
		got += p.String()
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
