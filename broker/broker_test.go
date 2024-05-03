package broker

import (
	"context"
	"github.com/google/uuid"
	"io/ioutil"
	"sync"
	"testing"
)

var broker *Broker

func init() {
	var err error
	dirName, err := ioutil.TempDir("", "data-test")
	if err != nil {
		panic(err)
	}

	broker, err = New(context.TODO(), dirName)
	if err != nil {
		panic(err)
	}
}

//func TestBroker_Chans(t *testing.T) {
//	var testChan chan string = make(chan string)
//	var wg sync.WaitGroup
//
//	wg.Add(2)
//	go func(c <-chan string) {
//	Out:
//		for {
//			select {
//			case msg := <-c:
//				t.Log(msg)
//				break Out
//			}
//		}
//		wg.Done()
//	}(testChan)
//
//	go func(c <-chan string) {
//	Out:
//		for {
//			select {
//			case msg := <-c:
//				t.Log(msg)
//				break Out
//			}
//		}
//		wg.Done()
//	}(testChan)
//	testChan <- "ok"
//	testChan <- "ok2"
//	wg.Wait()
//}

func TestBroker_New(t *testing.T) {
	tmpDir := t.TempDir()

	b, err := New(context.TODO(), tmpDir)
	if err != nil {
		t.Error(err)
	}

	if b == nil {
		t.Error("broker is nil")
	}
}

func TestBroker_Start(t *testing.T) {
	if err := broker.Start(); err != nil {
		t.Error(err)
	}
}

func TestBroker_CreateRouter(t *testing.T) {
	routerConfig := RouterConfig{"test-router", RouterDirect}
	if err := broker.CreateRouter(routerConfig); err != nil {
		t.Error(err)
	}
}

func TestBroker_CreateQueue(t *testing.T) {
	queueConfig := QueueConfig{Name: "test-queue", QueueType: QueueHdd, Bind: make([]QueueBindOptions, 0)}

	if err := broker.CreateQueue(queueConfig); err != nil {
		t.Error(err)
	}
}

func TestBroker_PublishAndGetMessageById(t *testing.T) {
	queueName := "test-queue"
	payload := "test payload"

	msgId, err := broker.Publish(queueName, []byte(payload))
	if err != nil {
		t.Error(err)
	}

	value, err := broker.GetQueueMessage(queueName, msgId)
	if err != nil {
		t.Error(err)
	}

	if string(value) != payload {
		t.Error("message not equal")
	}
}

func TestBroker_SubscribeAndPublish(t *testing.T) {
	queueName := "test-queue"
	var wg sync.WaitGroup

	wg.Add(1)

	subChan, err := broker.Subscribe(queueName)
	if err != nil {
		t.Error(err)
	}

	go func(sc <-chan []byte) {
		for msg := range sc {
			t.Log(string(msg))
			break
		}
		wg.Done()
	}(subChan)

	_, err = broker.Publish(queueName, []byte("it is working"))
	if err != nil {
		t.Log(err)
	}

	wg.Wait()
}

func BenchmarkBroker_CreateRouter_P1(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := uuid.New()
			routerConfig := RouterConfig{"test-router-" + id.String(), RouterDirect}

			if err := broker.CreateRouter(routerConfig); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkBroker_CreateQueue_P1(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := uuid.New()
			queueConfig := QueueConfig{Name: "test-queue-" + id.String(), QueueType: QueueHdd, Bind: make([]QueueBindOptions, 0)}

			if err := broker.CreateQueue(queueConfig); err != nil {
				b.Error(err)
			}
		}
	})
}
