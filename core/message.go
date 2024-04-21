package core

// Message
type Message struct {
	payload    []byte
	extraSpace [4]byte
}

func NewMessage() Message {
	return Message{}
}

func (m *Message) Bytes() []byte {
	msgBytes := append(MESSAGE_HEADER[:], m.payload...)
	msgBytes = append(msgBytes, MESSAGE_EOF[:]...)
	msgBytes = append(msgBytes, m.extraSpace[:]...)

	return msgBytes
}

func (m *Message) SetPayload(payload []byte) {
	m.payload = payload
}

func (m *Message) SetExtraSpace(extraSpace [4]byte) {
	m.extraSpace = extraSpace
}
