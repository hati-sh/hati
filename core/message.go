package core

import (
	"errors"
	"fmt"
)

var ErrInvalidHeader = errors.New("invalid header")
var ErrInvalidPayload = errors.New("invalid payload")
var ErrInvalidEOF = errors.New("invalid message EOF")

// Message
type Message struct {
	header     [len(MESSAGE_HEADER)]byte
	payload    []byte
	extraSpace [4]byte
}

func NewMessage() Message {
	return Message{}
}

func (m *Message) Bytes() []byte {
	payloadLength := len(m.payload)

	m.payload = append(m.payload, []byte("EOP\n")...)

	msgBytes := append(MESSAGE_HEADER[:], m.extraSpace[:]...)
	msgBytes = append(msgBytes, []byte("CL:"+fmt.Sprint(payloadLength)+"\n")...)
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

// func ParseBytesToMessage(in []byte) (*Message, error) {
// 	receivedHeaderBytes := in[0:8]

// 	if bytes.Compare(receivedHeaderBytes, MESSAGE_HEADER[:]) != 0 {
// 		return nil, ErrInvalidHeader
// 	}

// 	firstNewLineIdx := bytes.IndexByte(in, '\n')
// 	if firstNewLineIdx <= 15 {
// 		return nil, ErrInvalidHeader
// 	}

// 	var extraSpace [4]byte = [4]byte{in[9], in[10], in[11], in[12]}
// 	clHeader := in[12:firstNewLineIdx]

// 	if bytes.Compare(clHeader[0:3], []byte{'C', 'L', ':'}) != 0 {
// 		return nil, ErrInvalidHeader
// 	}

// 	var contentLength int

// 	contentLength, err := strconv.Atoi(string(in[15:firstNewLineIdx]))
// 	if err != nil {
// 		return nil, ErrInvalidPayload
// 	}

// 	payloadStartIdx := firstNewLineIdx + 1
// 	payloadEndIdx := payloadStartIdx + contentLength

// 	payload := in[payloadStartIdx:payloadEndIdx]

// 	eopStartIdx := payloadEndIdx
// 	eopEndIdx := eopStartIdx + 4

// 	if bytes.Compare(in[eopStartIdx:eopEndIdx], []byte{'E', 'O', 'P', '\n'}) != 0 {
// 		fmt.Println("x <")

// 		return nil, ErrInvalidPayload
// 	}

// 	if bytes.Compare(in[eopStartIdx:eopEndIdx], []byte{'E', 'O', 'P', '\n'}) != 0 {
// 		fmt.Println("y <")

// 		return nil, ErrInvalidEOF
// 	}

// 	return &Message{
// 		header:     [8]byte(receivedHeaderBytes),
// 		payload:    payload,
// 		extraSpace: extraSpace,
// 	}, nil
// }
