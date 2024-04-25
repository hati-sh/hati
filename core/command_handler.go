package core

import (
	"context"
	"errors"
)

type CommandHandler struct {
	ctx context.Context
}

func (ch *CommandHandler) processPayload(payload []byte) ([]byte, error) {
	// fmt.Println("command handler processing payload ProcessPayload: ")
	// fmt.Println(string(payload))

	if payload != nil {
		response := []byte("+OK\n")

		return response, nil
	}

	return nil, errors.New("+ERR\n")
}
