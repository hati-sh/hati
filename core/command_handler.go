package core

import (
	"bytes"
	"context"
	"errors"

	"github.com/hati-sh/hati/broker"
	"github.com/hati-sh/hati/storage"
)

type CommandHandler struct {
	ctx     context.Context
	storage *storage.Storage
	broker  *broker.Broker
}

func (ch *CommandHandler) processPayload(payload []byte) ([]byte, error) {
	if payload == nil {
		return nil, errors.New(string(CmdErr))
	}

	payloadArr := bytes.Split(bytes.Trim(payload, " "), []byte(" "))

	if bytes.Equal(payloadArr[0], CmdSet) {
		return ch.set(payloadArr)
	} else if bytes.Equal(payloadArr[0], CmdGet) {
		return ch.get(payloadArr)
	}

	return nil, errors.New(string(CmdErr))
}

func (ch *CommandHandler) set(payloadArr [][]byte) ([]byte, error) {
	key := payloadArr[2]
	value := bytes.Join(payloadArr[3:], []byte(" "))

	if err := ch.storage.Set(storage.Memory, key, value); err != nil {
		return nil, err
	}

	return CmdOk, nil
}

func (ch *CommandHandler) get(payloadArr [][]byte) ([]byte, error) {
	key := payloadArr[2]

	value, err := ch.storage.Memory.Get(key)
	if err != nil {
		return nil, err
	}
	value = append(value, byte('\n'))

	return value, nil
}
