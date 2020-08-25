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
