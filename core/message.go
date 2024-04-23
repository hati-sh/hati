package core

import (
	"errors"
	"fmt"
)

// Message
type Message struct {
	payload    []byte
	extraSpace [4]byte
}

func NewMessage() Message {
	return Message{}
}

func (m *Message) Bytes() []byte {
	m.payload = append(m.payload, []byte("EOP\n")...)

	msgBytes := append(MESSAGE_HEADER[:], m.extraSpace[:]...)
	msgBytes = append(msgBytes, []byte("CL:"+fmt.Sprint(len(m.payload))+"\n")...)
	msgBytes = append(msgBytes, m.payload...)
	msgBytes = append(msgBytes, MESSAGE_EOF[:]...)

	return msgBytes
}

func (m *Message) SetPayload(payload []byte) {
	m.payload = payload
}

func (m *Message) SetExtraSpace(extraSpace [4]byte) {
	m.extraSpace = extraSpace
}

func ParseBytesToMessage(in []byte) (*Message, error) {
	// incomingHeader := in[0:7]

	return nil, nil
}

var ErrInvalidHeader = errors.New("invalid header")
