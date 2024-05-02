package broker

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"github.com/hati-sh/hati/common"
	"github.com/hati-sh/hati/common/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var ErrRouterExist = errors.New("router exist")
var ErrQueueExist = errors.New("queue exist")
var ErrQueueNotExist = errors.New("queue not exist")
var ErrQueueMessageIdEmpty = errors.New("message id empty")

var routerDbWriteOptions = &opt.Options{
	WriteBuffer: common.BROKER_ROUTER_HDD_WRITE_BUFFER,
}

var queueDbWriteOptions = &opt.Options{
	WriteBuffer: common.BROKER_QUEUE_HDD_WRITE_BUFFER,
}

var queueMessageDbWriteOptions = &opt.Options{
	WriteBuffer: common.BROKER_QUEUE_HDD_WRITE_BUFFER,
}

type Broker struct {
	ctx                 context.Context
	dataDir             string
	router              map[string]*Router
	routerDb            *leveldb.DB
	routerLock          sync.RWMutex
	queue               map[string]*Queue
	queueDb             *leveldb.DB
	queueMessageDb      *leveldb.DB
	queueMessageChanMap map[string]chan []byte
	queueLock           sync.RWMutex
}

func New(ctx context.Context, dataDir string) (*Broker, error) {
	routerDb, err := common.OpenDatabase(dataDir, "router", routerDbWriteOptions)
	if err != nil {
		return nil, err
	}

	queueDb, err := common.OpenDatabase(dataDir, "queue", queueDbWriteOptions)
	if err != nil {
		return nil, err
	}

	queueMessageDb, err := common.OpenDatabase(dataDir, "queue_message", queueMessageDbWriteOptions)
	if err != nil {
		return nil, err
	}

	brokerInstance := &Broker{
		ctx:                 ctx,
		dataDir:             dataDir,
		router:              make(map[string]*Router),
		routerDb:            routerDb,
		queue:               make(map[string]*Queue),
		queueDb:             queueDb,
		queueMessageDb:      queueMessageDb,
		queueMessageChanMap: make(map[string]chan []byte, common.BROKER_QUEUE_MESSAGE_CHAN_SIZE),
	}

	return brokerInstance, nil
}

func (b *Broker) Start() error {
	if err := b.loadRoutersFromDatabase(); err != nil {
		return err
	}

	if err := b.loadQueuesFromDatabase(); err != nil {
		return err
	}

	return nil
}

func (b *Broker) Stop() {
	if err := b.routerDb.Close(); err != nil {
		logger.Error(err.Error())
	}

	if err := b.queueDb.Close(); err != nil {
		logger.Error(err.Error())
	}

	if err := b.queueMessageDb.Close(); err != nil {
		logger.Error(err.Error())
	}

	b.router = nil
	b.queue = nil
	b.queueMessageChanMap = nil
}

func (b *Broker) CreateRouter(config RouterConfig) error {
	b.routerLock.RLock()
	if b.router[config.Name] != nil {
		b.routerLock.RUnlock()
		return ErrRouterExist
	}
	b.routerLock.RUnlock()

	valueBytes, err := common.EncodeToBytes(config)
	if err != nil {
		return err
	}

	if err = b.routerDb.Put([]byte(config.Name), valueBytes, nil); err != nil {
		return err
	}

	b.routerLock.Lock()
	defer b.routerLock.Unlock()
	b.router[config.Name] = NewRouter(b.ctx, config)

	return nil
}

func (b *Broker) CreateQueue(config QueueConfig) error {
	b.queueLock.RLock()
	if b.queue[config.Name] != nil {
		b.queueLock.RUnlock()
		return ErrQueueExist
	}
	b.queueLock.RUnlock()

	valueBytes, err := common.EncodeToBytes(config)
	if err != nil {
		return err
	}

	if err = b.queueDb.Put([]byte(config.Name), valueBytes, nil); err != nil {
		return err
	}

	b.queueLock.Lock()
	defer b.queueLock.Unlock()

	b.queue[config.Name] = NewQueue(b.ctx, config)

	return nil
}

func (b *Broker) Publish(queueName string, payload []byte) (string, error) {
	if b.queue[queueName] == nil {
		return "", ErrQueueNotExist
	}

	msgIdUuid := uuid.New()
	msgId := msgIdUuid.String()

	if err := b.queueMessageDb.Put([]byte(queueName+"-"+msgId), payload, nil); err != nil {
		return "", err
	}

	if b.queueMessageChanMap[queueName] != nil {
		b.queueMessageChanMap[queueName] <- payload
	}

	return msgId, nil
}

func (b *Broker) GetQueueMessage(queueName string, messageId string) ([]byte, error) {
	if queueName == "" || b.queue[queueName] == nil {
		return nil, ErrQueueNotExist
	}

	if messageId == "" {
		return nil, ErrQueueMessageIdEmpty
	}

	value, err := b.queueMessageDb.Get([]byte(queueName+"-"+messageId), nil)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (b *Broker) Subscribe(queueName string) (<-chan []byte, error) {
	if queueName == "" || b.queue[queueName] == nil {
		return nil, ErrQueueNotExist
	}

	if b.queueMessageChanMap[queueName] == nil {
		b.queueMessageChanMap[queueName] = make(chan []byte, common.BROKER_QUEUE_MESSAGE_CHAN_SIZE)
	}

	return b.queueMessageChanMap[queueName], nil
}

func (b *Broker) loadRoutersFromDatabase() error {
	iter := b.routerDb.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		routerConfig := RouterConfig{}
		dec := gob.NewDecoder(bytes.NewReader(value))
		err := dec.Decode(&routerConfig)
		if err != nil {
			return err
		}

		b.router[string(key)] = NewRouter(b.ctx, routerConfig)
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}

func (b *Broker) loadQueuesFromDatabase() error {
	iter := b.queueDb.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		queueConfig := QueueConfig{}
		dec := gob.NewDecoder(bytes.NewReader(value))
		err := dec.Decode(&queueConfig)
		if err != nil {
			return err
		}

		b.queue[string(key)] = NewQueue(b.ctx, queueConfig)
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}
