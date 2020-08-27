package cmd

import (
	"bytes"
	"encoding/binary"
	"log"
)

const (
	MessageLen = 6

	// message types
	MessageOK  = 1
	MessageAdd = 2
	MessageDel = 3
	MessageErr = 4
	MessageNop = 5

	// protocol numbers
	ProtocolTCP = 6
	ProtocolUDP = 17
)

// Message stores a control Message
type Message struct {
	Op       uint8
	Protocol uint8
	Port     uint16
	DestPort uint16
}

// Serialize writes message to a byte slice
func (m *Message) Serialize() []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, m)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// Parse reads message from byte slice b
func (m *Message) Parse(b []byte) {
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, m)
	if err != nil {
		log.Println("error reading message:", err)
	}
}
