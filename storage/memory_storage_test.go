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

func TestNewMemoryStorage_New(t *testing.T) {
	s := NewMemoryStorage(context.TODO())
	if s == nil {
		t.Errorf("NewMemoryStorage returned nil")
	}
}

func TestMemoryStorage_Set(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())

	if ok := storage.Set([]byte(id.String()), value, 0); !ok {
		t.Errorf("Set returned false")
	}

}

func TestMemoryStorage_Get(t *testing.T) {
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

func TestMemoryStorage_Has(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := []byte(id.String())

	if ok := storage.Set(key, value, 0); !ok {
		t.Errorf("Set returned false")
	}

	if ok := storage.Has(key); !ok {
		t.Errorf("Has returned false")
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := []byte(id.String())

	if ok := storage.Set(key, value, 0); !ok {
		t.Errorf("Set returned false")
	}

	storage.Delete(key)
}

func TestMemoryStorage_Count(t *testing.T) {
	if count := storage.CountKeys(); count < 1 {
		t.Errorf("CountKeys returned 0 keys")
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

func BenchmarkMemoryStorage_Get(b *testing.B) {
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

func BenchmarkMemoryStorage_Get_P100(b *testing.B) {
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

func BenchmarkMemoryStorage_Has(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		if ok := storage.Has(insertedKeys[randInt(0, 9999)]); !ok {
			b.Errorf("Has returned false")
		}
	}
}

func BenchmarkMemoryStorage_Has_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ok := storage.Has(insertedKeys[randInt(0, 9999)]); !ok {
				b.Errorf("Has returned false")
			}
		}
	})
}

func BenchmarkMemoryStorage_CountKeys(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		if count := storage.CountKeys(); count < 1 {
			b.Errorf("CountKeys returned 0 keys")
		}
	}
}

func BenchmarkMemoryStorage_CountKeys_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if count := storage.CountKeys(); count < 1 {
				b.Errorf("CountKeys returned 0 keys")
			}
		}
	})
}

func BenchmarkMemoryStorage_Delete(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		storage.Delete(insertedKeys[randInt(0, 9999)])
	}
}

func BenchmarkMemoryStorage_Delete_P100(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	b.SetParallelism(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			storage.Delete(insertedKeys[randInt(0, 9999)])
		}
	})
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
