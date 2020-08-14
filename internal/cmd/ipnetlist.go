package cmd

import "net"

var (
	allowedIPNets ipNetList
)

// ipNetList is a list of IP network addresses
type ipNetList struct {
	l []*net.IPNet
}

// add adds ip network ipNet to the list
func (i *ipNetList) add(ipNet *net.IPNet) {
	i.l = append(i.l, ipNet)
}

// containsIP checks if any ip network in the list contains ip
func (i *ipNetList) containsIP(ip net.IP) bool {
	for _, ipNet := range i.l {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
