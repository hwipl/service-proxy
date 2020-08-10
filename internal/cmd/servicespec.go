package cmd

import (
	"log"
	"strconv"
	"strings"
)

// serviceSpec stores the specification of a service
type serviceSpec struct {
	protocol string
	port     uint16
	destPort uint16
}

// toMessage converts a service specification to a message
func (s *serviceSpec) toMessage() *message {
	m := message{
		op:       messageAdd,
		port:     s.port,
		destPort: s.destPort,
	}
	switch s.protocol {
	case "tcp":
		m.protocol = protocolTCP
	case "udp":
		m.protocol = protocolUDP
	default:
		log.Fatalf("unknown protocol \"%s\" in service "+
			"specification\n", s.protocol)
	}
	return &m
}

// parseServiceSpec parses spec as a service specification with the format
// "<protocol>:<port>:<destPort>"
func parseServiceSpec(spec string) *serviceSpec {
	errFmt := "Error parsing service specification %s"
	parts := strings.Split(spec, ":")
	if len(parts) != 3 {
		log.Fatalf(errFmt, spec)
	}

	// parse protocol
	protocol := parts[0]

	// parse port
	port, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		log.Fatalf(errFmt, spec)
	}

	// parse destination port
	destPort, err := strconv.ParseUint(parts[2], 10, 16)
	if err != nil {
		log.Fatalf(errFmt, spec)
	}

	// return as serviceSpec
	s := serviceSpec{
		protocol: protocol,
		port:     uint16(port),
		destPort: uint16(destPort),
	}
	return &s
}
