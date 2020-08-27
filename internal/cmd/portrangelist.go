package cmd

import (
	"fmt"

	"github.com/hwipl/service-proxy/internal/network"
)

var (
	allowedPortRanges portRangeList
)

// portRange specifies a port range for a specific protocol
type portRange struct {
	protocol uint8
	min      uint16
	max      uint16
}

// contains checks if this port range contains the other port range
func (p *portRange) contains(other *portRange) bool {
	return p.protocol == other.protocol &&
		p.min <= other.min &&
		p.max >= other.max
}

// merge combines this port range with the other port range if they overlap
func (p *portRange) merge(other *portRange) {
	if p.protocol != other.protocol {
		return
	}

	if other.min <= p.min && other.max >= p.min {
		p.min = other.min
	}
	if other.max >= p.max && other.min <= p.max {
		p.max = other.max
	}
}

// containsPort returns if protocol and port are in the port range
func (p *portRange) containsPort(protocol uint8, port uint16) bool {
	return protocol == p.protocol && port >= p.min && port <= p.max
}

// String converts the port range to a string
func (p *portRange) String() string {
	var protocol string

	// convert protocol number to a string if possible
	switch p.protocol {
	case network.ProtocolTCP:
		protocol = "tcp"
	case network.ProtocolUDP:
		protocol = "udp"
	default:
		protocol = fmt.Sprintf("%d", p.protocol)
	}

	return fmt.Sprintf("%s:%d-%d", protocol, p.min, p.max)
}

// portRangeList is a list of portRanges
type portRangeList struct {
	l []*portRange
}

// add adds protocol and ports min and max to the list
func (p *portRangeList) add(protocol uint8, min, max uint16) {
	r := portRange{
		protocol: protocol,
		min:      min,
		max:      max,
	}

	// check if there is an existing entry that contains the new port
	// range, if the new entry can be combined with existing entries to
	// form a bigger port range, and if existing port range entries should
	// be replaced by the new one
	k := 0
	for _, existing := range p.l {
		if existing.contains(&r) {
			// existing entry already contains new one, stop here
			return
		}
		// try building a bigger port range with an existing entry
		r.merge(existing)
		if r.contains(existing) {
			// new entry contains existing one, remove existing
			// entry and add new one later
			continue
		}

		// keep entry in list
		p.l[k] = existing
		k++

	}
	// update list with remaining entries
	p.l = p.l[:k]

	// new entry, add it
	p.l = append(p.l, &r)
}

// containsPort returns if protocol and port are in any of the port ranges in
// the list
func (p *portRangeList) containsPort(protocol uint8, port uint16) bool {
	for _, r := range p.l {
		if r.containsPort(protocol, port) {
			return true
		}
	}
	return false
}

// getAll returns a list of all port ranges
func (p *portRangeList) getAll() []*portRange {
	return p.l
}
