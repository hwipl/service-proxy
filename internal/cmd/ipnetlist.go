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
	// network, or if existing ip network entries should be replaced by
	// the new one
	k := 0
	for _, existing := range i.l {
		// get netmask bits and lengths
		eOnes, eBits := existing.Mask.Size()
		nOnes, nBits := ipNet.Mask.Size()
		if eBits == nBits {
			if existing.Contains(ipNet.IP) && eOnes <= nOnes {
				// existing entry already contains new one,
				// stop here
				return
			}
			if ipNet.Contains(existing.IP) && eOnes > nOnes {
				// new entry contains existing one, remove
				// existing entry and add new one later
				continue
			}
		}

		// keep entry in list
		i.l[k] = existing
		k++

	}
	// update list with remaining entries
	i.l = i.l[:k]

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
