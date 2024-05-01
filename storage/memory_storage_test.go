package storage

import (
	"context"
	"github.com/google/uuid"
	"math/rand"
	"testing"
)

var storage *memoryStorage

var insertedKeys [][]byte

func init() {
	ctx, _ := context.WithCancel(context.Background())
	storage = NewMemoryStorage(ctx)
	storage.Start()

	i := 0
	for ; i < 10000; i++ {
		id := uuid.New()
		value := []byte(uuid.New().String())
		key := []byte(id.String())

		if ok := storage.Set(key, value, 0); !ok {
		}

		insertedKeys = append(insertedKeys, key)
	}
}

func TestNewMemoryStorageCreation(t *testing.T) {
	s := NewMemoryStorage(context.TODO())
	if s == nil {
		t.Errorf("NewMemoryStorage returned nil")
	}
}

func TestSet(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())

	if ok := storage.Set([]byte(id.String()), value, 0); !ok {
		t.Errorf("Set returned false")
	}

}

func TestGet(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := []byte(id.String())

	if ok := storage.Set(key, value, 0); !ok {
		t.Errorf("Set returned false")
	}

	resValue, err := storage.Get(key)
	if err != nil {
		t.Errorf("Get returned false")
	}

	if string(resValue) != string(value) {
		t.Errorf("Get returned wrong value")
	}
}

func BenchmarkMemoryStorage_Set_Ttl0(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		id := uuid.New()
		value := []byte(uuid.New().String())

		if ok := storage.Set([]byte(id.String()), value, 0); !ok {
			b.Errorf("Set returned false")
		}
	}
}

func BenchmarkMemoryStorage_Set_Ttl0_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := uuid.New()
			value := []byte(uuid.New().String())

			if ok := storage.Set([]byte(id.String()), value, 0); !ok {
				b.Errorf("Set returned false")
			}
		}
	})
}

func BenchmarkMemoryStorage_Set_Ttl10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		id := uuid.New()
		value := []byte(uuid.New().String())

		if ok := storage.Set([]byte(id.String()), value, 10); !ok {
			b.Errorf("Set returned false")
		}
	}
}

func BenchmarkMemoryStorage_Set_Ttl10_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := uuid.New()
			value := []byte(uuid.New().String())

			if ok := storage.Set([]byte(id.String()), value, 10); !ok {
				b.Errorf("Set returned false")
			}
		}
	})
}

func BenchmarkMemoryStorage_Get_Ttl0(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		_, err := storage.Get(insertedKeys[randInt(0, 9999)])
		if err != nil {
			b.Errorf("Get returned error %v", err)
		}
	}
}

func BenchmarkMemoryStorage_Get_Ttl0_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := storage.Get(insertedKeys[randInt(0, 9999)])
			if err != nil {
				b.Errorf("Get returned error %v", err)
			}
		}
	})
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
