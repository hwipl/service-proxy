package pclient

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hwipl/service-proxy/internal/network"
)

// ServiceSpec stores the specification of a service
type ServiceSpec struct {
	Protocol string
	Port     uint16
	DestPort uint16
}

// ToMessage converts a service specification to a message
func (s *ServiceSpec) ToMessage() *network.Message {
	m := network.Message{
		Op:       network.MessageAdd,
		Port:     s.Port,
		DestPort: s.DestPort,
	}
	switch s.Protocol {
	case "tcp":
		m.Protocol = network.ProtocolTCP
	case "udp":
		m.Protocol = network.ProtocolUDP
	default:
		log.Fatalf("unknown protocol \"%s\" in service "+
			"specification\n", s.Protocol)
	}
	return &m
}

// FromMessage fills this service specification from a message
func (s *ServiceSpec) FromMessage(msg *network.Message) {
	s.Port = msg.Port
	s.DestPort = msg.DestPort
	switch msg.Protocol {
	case network.ProtocolTCP:
		s.Protocol = "tcp"
	case network.ProtocolUDP:
		s.Protocol = "udp"
	default:
		s.Protocol = "unknown"
	}
}

// String converts the service spec to a string
func (s *ServiceSpec) String() string {
	return fmt.Sprintf("%s:%d:%d", s.Protocol, s.Port, s.DestPort)
}

// ParseServiceSpec parses spec as a service specification with the format
// "<protocol>:<port>:<destPort>"
func ParseServiceSpec(spec string) *ServiceSpec {
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
	s := ServiceSpec{
		Protocol: protocol,
		Port:     uint16(port),
		DestPort: uint16(destPort),
	}
	return &s
}
