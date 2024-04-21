package core

import "fmt"

var MESSAGE_HEADER [8]byte = [8]byte{'+', 'h', 'a', 't', 'i', '+'}
var MESSAGE_EOF [8]byte = [8]byte{'-', '-', 'h', 'a', 't', 'i', '\n', '\r'}
var MESSAGE_HATI_SPACE [4]byte = [4]byte{}
var COMMAND_DELIMITER [2]byte = [2]byte{'+', '\n'}

type Hati struct{}

func NewHati() Hati {
	return Hati{}
}

func (h *Hati) Start() error {

	fmt.Println(MESSAGE_HEADER)

	return nil
}
