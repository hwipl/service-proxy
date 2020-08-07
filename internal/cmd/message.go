package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	messageLen = 6
)

// message stores a control message
type message struct {
	op       uint8
	protocol uint8
	port     uint16
	destPort uint16
}

// serialize writes message to a byte slice
func (m *message) serialize() []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, m)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// parse reads message from byte slice b
func (m *message) parse(b []byte) {
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, m)
	if err != nil {
		fmt.Println("error reading message:", err)
	}
}
