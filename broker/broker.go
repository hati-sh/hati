package broker

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/hati-sh/hati/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var ErrRouterExist = errors.New("router exist")
var ErrQueueExist = errors.New("queue exist")

var routerDbWriteOptions = &opt.Options{
	WriteBuffer: common.BROKER_ROUTER_HDD_WRITE_BUFFER,
}

var queueDbWriteOptions = &opt.Options{
	WriteBuffer: common.BROKER_QUEUE_HDD_WRITE_BUFFER,
}

type Broker struct {
	ctx        context.Context
	dataDir    string
	router     map[string]*Router
	routerDb   *leveldb.DB
	routerLock sync.RWMutex
	queue      map[string]*Queue
	queueDb    *leveldb.DB
	queueLock  sync.RWMutex
}

func New(ctx context.Context, dataDir string) (*Broker, error) {
	routerDb, err := common.OpenDatabase(dataDir, "router", routerDbWriteOptions)
	if err != nil {
		return nil, err
	}

	queueDb, err := common.OpenDatabase(dataDir, "queue", routerDbWriteOptions)
	if err != nil {
		return nil, err
	}

	brokerInstance := &Broker{
		ctx:      ctx,
		dataDir:  dataDir,
		router:   make(map[string]*Router),
		routerDb: routerDb,
		queue:    make(map[string]*Queue),
		queueDb:  queueDb,
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
	_ = b.routerDb.Close()
	_ = b.queueDb.Close()
}

func (b *Broker) CreateRouter(config RouterConfig) error {
	if b.router[config.Name] != nil {
		return ErrRouterExist
	}

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
	if b.queue[config.Name] != nil {
		return ErrQueueExist
	}

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
