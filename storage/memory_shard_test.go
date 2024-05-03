package storage

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/hati-sh/hati/common"
	"testing"
)

var shardMap MemoryShardMap
var insertedShardMapKeys []string

func init() {
	shardMap = newMemoryShardMap(common.STORAGE_DEFAULT_NUMBER_OF_SHARDS)

	i := 0
	for ; i < 10000; i++ {
		id := uuid.New()
		value := []byte(uuid.New().String())
		key := id.String()

		shardMap.Set(key, value)

		insertedShardMapKeys = append(insertedShardMapKeys, key)
	}
}

func TestMemoryShardMap_Set(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := id.String()

	shardMap.Set(key, value)
}

func TestMemoryShardMap_Get(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := id.String()

	shardMap.Set(key, value)

	val, ok := shardMap.Get(key)
	if !ok || !bytes.Equal(val, value) {
		t.Error("ShardMap.Get fail")
	}
}

func TestMemoryShardMap_Has(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := id.String()

	shardMap.Set(key, value)

	ok := shardMap.Has(key)
	if !ok {
		t.Error("ShardMap.Has fail")
	}
}

func TestMemoryShardMap_Delete(t *testing.T) {
	id := uuid.New()
	value := []byte(uuid.New().String())
	key := id.String()

	shardMap.Set(key, value)
	shardMap.Delete(key)

	ok := shardMap.Has(key)
	if ok {
		t.Error("ShardMap.Delete fail")
	}
}

func TestMemoryShardMap_Count(t *testing.T) {
	if count := shardMap.CountKeys(); count < 1 {
		t.Errorf("CountKeys returned 0 keys")
	}
}

func BenchmarkMemoryShardMap_Set(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		id := uuid.New()
		value := []byte(uuid.New().String())

		shardMap.Set(id.String(), value)
	}
}

func BenchmarkMemoryShardMap_Set_P10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := uuid.New()
			value := []byte(uuid.New().String())

			shardMap.Set(id.String(), value)
		}
	})
}

func BenchmarkMemoryShardMap_Get(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()

	i := 0
	for ; i < b.N; i++ {
		_, ok := shardMap.Get(string(insertedShardMapKeys[randInt(0, 9999)]))
		if !ok {
			b.Error("ShardMap.Get fail")
		}
	}
}

func BenchmarkMemoryShardMap_Get_P10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, ok := shardMap.Get(string(insertedShardMapKeys[randInt(0, 9999)]))
			if !ok {
				b.Error("ShardMap.Get fail")
			}
		}
	})
}

func BenchmarkMemoryShardMap_Has_P1(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ok := shardMap.Has(string(insertedShardMapKeys[randInt(0, 9999)]))
			if !ok {
				b.Error("ShardMap.Has fail")
			}
		}
	})
}

func BenchmarkMemoryShardMap_Has_P10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ok := shardMap.Has(string(insertedShardMapKeys[randInt(0, 9999)]))
			if !ok {
				b.Error("ShardMap.Has fail")
			}
		}
	})
}

func BenchmarkMemoryShardMap_CountKeys_P1(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if count := shardMap.CountKeys(); count < 1 {
				b.Errorf("CountKeys returned 0 keys")
			}
		}
	})
}

func BenchmarkMemoryShardMap_CountKeys_P10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if count := shardMap.CountKeys(); count < 1 {
				b.Errorf("CountKeys returned 0 keys")
			}
		}
	})
}

func BenchmarkMemoryShardMap_Delete_P1(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			shardMap.Delete(string(insertedShardMapKeys[randInt(0, 9999)]))
		}
	})
}

func BenchmarkMemoryShardMap_Delete_P10(b *testing.B) {
	b.SetBytes(1)
	b.ReportAllocs()
	b.SetParallelism(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			shardMap.Delete(string(insertedShardMapKeys[randInt(0, 9999)]))
		}
	})
}
