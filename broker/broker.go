package broker

import (
	"context"
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

func (b *Broker) Start() {}

func (b *Broker) Stop() {}

func (b *Broker) CreateRouter(config RouterConfig) error {
	if b.router[config.name] != nil {
		return ErrRouterExist
	}

	valueBytes, err := common.EncodeToBytes(config)
	if err != nil {
		return err
	}

	if err = b.routerDb.Put([]byte(config.name), valueBytes, nil); err != nil {
		return err
	}

	b.routerLock.Lock()
	defer b.routerLock.Unlock()
	b.router[config.name] = NewRouter(b.ctx, config)

	return nil
}

func (b *Broker) CreateQueue(config QueueConfig) error {
	if b.queue[config.name] != nil {
		return ErrQueueExist
	}

	valueBytes, err := common.EncodeToBytes(config)
	if err != nil {
		return err
	}

	if err = b.queueDb.Put([]byte(config.name), valueBytes, nil); err != nil {
		return err
	}

	b.queueLock.Lock()
	defer b.queueLock.Unlock()

	b.queue[config.name] = NewQueue(b.ctx, config)

	return nil
}
