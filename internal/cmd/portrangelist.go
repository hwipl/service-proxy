package cmd

var (
	allowedPortRanges portRangeList
)

// portRange specifies a port range for a specific protocol
type portRange struct {
	protocol int
	min      int
	max      int
}

// containsPort returns if protocol and port are in the port range
func (p *portRange) containsPort(protocol, port int) bool {
	return protocol == p.protocol && port >= p.min && port <= p.max
}

// portRangeList is a list of portRanges
type portRangeList struct {
	l []*portRange
}

// add adds protocol and ports min and max to the list
func (p *portRangeList) add(protocol, min, max int) {
	r := portRange{
		protocol: protocol,
		min:      min,
		max:      max,
	}
	p.l = append(p.l, &r)
}

// containsPort returns if protocol and port are in any of the port ranges in
// the list
func (p *portRangeList) containsPort(protocol, port int) bool {
	for _, r := range p.l {
		if r.containsPort(protocol, port) {
			return true
		}
	}
	return false
}
