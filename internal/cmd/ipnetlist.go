package cmd

import (
	"net"
)

var (
	allowedIPNets ipNetList
)

// ipNetList is a list of IP network addresses
type ipNetList struct {
	l []*net.IPNet
}

// add adds ip network ipNet to the list
func (i *ipNetList) add(ipNet *net.IPNet) {
	// check if there is an existing entry that contains the new ip
	// network, or if an existing ip network entry should be updated with
	// the new one
	for _, existing := range i.l {
		// get netmask bits and lengths
		eOnes, eBits := existing.Mask.Size()
		nOnes, nBits := ipNet.Mask.Size()
		if eBits != nBits {
			// incompatible mask lengths, skip
			continue
		}
		if existing.Contains(ipNet.IP) && eOnes <= nOnes {
			// existing entry already contains new one, stop
			return
		}
		if ipNet.Contains(existing.IP) && eOnes > nOnes {
			// new entry contains existing one, overwrite
			// existing element and stop
			existing.IP = ipNet.IP
			existing.Mask = ipNet.Mask
			return
		}
	}

	// new element, add it
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

// getAll returns the list of ip networks
func (i *ipNetList) getAll() []*net.IPNet {
	return i.l
}
