package core

import (
	"bytes"
	"context"
	"errors"
	"github.com/hati-sh/hati/broker"
	"github.com/hati-sh/hati/storage"
)

type CommandHandler struct {
	ctx            context.Context
	storageManager *storage.Manager
	broker         *broker.Broker
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
	} else if bytes.Equal(payloadArr[0], CmdHas) {
		if ch.has(payloadArr) {
			return CmdOk, nil
		}
		return nil, errors.New(string(CmdErr))
	} else if bytes.Equal(payloadArr[0], CmdDelete) {
		ch.delete(payloadArr)
		return CmdOk, nil
	} else if bytes.Equal(payloadArr[0], CmdCount) {
		return ch.count(payloadArr)
	} else if bytes.Equal(payloadArr[0], CmdFlushAll) {
		ch.flushAll(payloadArr)
		return CmdOk, nil
	}

	return nil, errors.New(string(CmdErr))
}

func (ch *CommandHandler) set(payloadArr [][]byte) ([]byte, error) {
	storageType := storage.Type(payloadArr[1])

	if len(payloadArr) < 5 {
		return nil, errors.New(string(CmdErr))
	}

	key := payloadArr[3]
	ttl := payloadArr[2]
	value := bytes.Join(payloadArr[4:], []byte(" "))

	if err := ch.storageManager.Set(storageType, key, value, ttl); err != nil {
		return nil, err
	}

	return CmdOk, nil
}

func (ch *CommandHandler) get(payloadArr [][]byte) ([]byte, error) {
	storageType := storage.Type(payloadArr[1])
	key := payloadArr[2]

	value, err := ch.storageManager.Get(storageType, key)
	if err != nil {
		return nil, err
	}
	value = append(value, byte('\n'))

	return value, nil
}

func (ch *CommandHandler) has(payloadArr [][]byte) bool {
	storageType := storage.Type(payloadArr[1])
	key := payloadArr[2]

	has := ch.storageManager.Has(storageType, key)

	return has
}

func (ch *CommandHandler) delete(payloadArr [][]byte) {
	storageType := storage.Type(payloadArr[1])
	key := payloadArr[2]

	ch.storageManager.Delete(storageType, key)
}

func (ch *CommandHandler) count(payloadArr [][]byte) ([]byte, error) {
	storageType := storage.Type(payloadArr[1])

	keysCount, err := ch.storageManager.Count(storageType)
	return []byte(string(rune(keysCount))), err
}

func (ch *CommandHandler) flushAll(payloadArr [][]byte) {
	storageType := storage.Type(payloadArr[1])
	ch.storageManager.FlushAll(storageType)
}
