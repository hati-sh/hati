package storage

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

var storage *memoryStorage

func init() {
	ctx, _ := context.WithCancel(context.Background())
	storage = NewMemoryStorage(ctx)
	storage.Start()
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
