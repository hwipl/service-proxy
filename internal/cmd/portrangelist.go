package cmd

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

// containsPort returns if protocol and port are in the port range
func (p *portRange) containsPort(protocol uint8, port uint16) bool {
	return protocol == p.protocol && port >= p.min && port <= p.max
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
	// range, or if existing port range entries should be replaced by
	// the new one
	k := 0
	for _, existing := range p.l {
		if existing.contains(&r) {
			// existing entry already contains new one, stop here
			return
		}
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
