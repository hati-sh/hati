# Hati

`v0.4.0-dev`

```
IN DEVELOPMENT. DO NOT USE IN PRODUCTION.
```

Hati is a simple distributed (in-memory and persistent) key/value store and a message broker.

Hati is meant to work in trusted networks - where all nodes operators know each other and can be trusted.

Once connected as a client to Hati server you can publish commands which will be processed by the server and published over the Hati network - Hati nodes can be connected to each other creating network of nodes but it is prefectly fine to run Hati as a single instance.

**Current version**
```
[req] [req]  
  ^     ^
  v     v
[Hong Kong]
```

**v2.0.0**
```
[req][req] [req][req] [req][req]
    ^          ^           ^
    v          v           v
[London]  [Hong Kong] [New York]
    ^          ^          ^
    |----------|----------|
```


## Current features

- TCP server
- JSON-RPC server
- Storing data in in-memory storage type - with sharding
- Storing data in hdd

## Configuration

While starting Hati `hati start` there are configurational flags available to be set.

- `--tcp` - indicates if should start TCP server, by default it is off
- `--tcp-host` - determines bind host for TCP server, default value is `0.0.0.0`
- `--tcp-port` - bind port for TCP server, default value is `4242`
- `--rpc` - indicates if should start JSON-RPC server, by default it is off
- `--rpc-host` - bind host for JSON-RPC server, default `0.0.0.0`
- `--rpc-port` - bind port for JSON-RPC server, default `6767`
- `--data-dir` - absolute path to directory where Hati can store files, default `/current/path/to-hati/data`
- `--cpu-num` - number of CPU cores which should be used by Hati, by default it will set for as many as available

## Protocol Commands

### Key-Value Storage

`SET <type> <ttl> <key> <value>\n` - save key with provided value to the selected storage type.

Hati offers two storage types: `memory` and `hdd` . By default `<ttl>` is set to zero `0` which means that value will be stored on the hard-drive. Ttl value is in ms, if higher than `0` Hati can guarantee that value will be stored at minimum for provided ttl value and will be removed from the storage shortly (as soon as possible) after that.

`HAS <type> <key>\n` - check if provided key exist in given storage type

`GET <type> <key>\n` - get value for provided key in given storage type

`DELETE <type> <key>\n` - get value for provided key in given storage type

`FLUSHALL <type>\n` - flush (delete) all data from given storage type

## JSON-RPC Methods

**Storage**

- `Storage.Set` - set key and value
  - `type` - storage type: `memory` / `hdd`
  - `key`
  - `value`
  - `ttl`
- `Storage.Has` - check if provided key exist
  - `type` - storage type: `memory` / `hdd`
  - `key`
- `Storage.Get` - get key
  - `type` - storage type: `memory` / `hdd`
  - `key`
- `Storage.Delete` - delete key
  - `type` - storage type: `memory` / `hdd`
  - `key`
- `Storage.FlushAll` - deletes all keys
  - `type` - storage type: `memory` / `hdd`
- `Storage.Count` - returns number of keys
  - `type` - storage type: `memory` / `hdd`

**Message broker**

- `Broker.`

## To do

Below you can find list of features planned to be released in upcoming versions.
Please keep in mind that non of these list is a set stone and items can be added/removed
but it gives overall image of what are the plans for the near future.

**v1.0.0**

- [x] Configuration options via CLI flags
- [x] Memory and persistent (LevelDB) key-value storage
- [X] TCP server
- [x] JSON-RPC server 
- [x] TCP and JSON-RPC commands:
  - [x] `SET`, `GET`, `HAS`, `DELETE`, `FLUSHALL`, `COUNT`
  - [ ] `CREATEROUTER`, `CREATEQUEUE`, `PUBLISH`, `ACK`, `GETROUTER`, `GETQUEUE`, `LISTROUTER`, `LISTQUEUE`, `FLUSH`
- [x] TCP server cancel context for graceful shutdown
- [x] RPC server cancel context for graceful shutdown
- [x] Graceful shutdown of HDD storage
- [x] Implement TTL
- [ ] CLI command to export/import backup of persistent storage
- [ ] Message broker
  - [ ] Routers management
  - [ ] Queues management
- [ ] Command to return server statistics
  - Number of keys for each storage type
  - Number of keys in each shard
  - Active TCP connections
  - Message broker statistics
- [ ] Communication via UNIX socket
- [ ] Access protected by credentials

**v1.1.0**

- [ ] Rebuilding persistent storage if number of shards has changed
- [ ] TLS support for TCP

**v2.0.0**

- [ ] Nodes clustering
- [ ] Data synchronization between nodes
  - TCP
  - P2P (? - TBD)

---

## Benchmark

### Memory Storage

```shell
goos: darwin
goarch: amd64
pkg: github.com/hati-sh/hati/storage
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkMemoryStorage_Set_Ttl0
BenchmarkMemoryStorage_Set_Ttl0-8                 292143              4055 ns/op           0.25 MB/s         480 B/op          8 allocs/op
BenchmarkMemoryStorage_Set_Ttl0_P100
BenchmarkMemoryStorage_Set_Ttl0_P100-8            751743              1553 ns/op           0.64 MB/s         520 B/op          8 allocs/op
BenchmarkMemoryStorage_Set_Ttl10
BenchmarkMemoryStorage_Set_Ttl10-8                266694              4763 ns/op           0.21 MB/s         489 B/op         10 allocs/op
BenchmarkMemoryStorage_Set_Ttl10_P100
BenchmarkMemoryStorage_Set_Ttl10_P100-8           506427              2181 ns/op           0.46 MB/s         469 B/op         10 allocs/op
BenchmarkMemoryStorage_Get
BenchmarkMemoryStorage_Get-8                     2686506               438.0 ns/op         2.28 MB/s          96 B/op          2 allocs/op
BenchmarkMemoryStorage_Get_P100
BenchmarkMemoryStorage_Get_P100-8               14420150               107.2 ns/op         9.33 MB/s          96 B/op          2 allocs/op
BenchmarkMemoryStorage_Has
BenchmarkMemoryStorage_Has-8                     2938930               408.3 ns/op         2.45 MB/s          96 B/op          2 allocs/op
BenchmarkMemoryStorage_Has_P100
BenchmarkMemoryStorage_Has_P100-8               13023224                93.40 ns/op       10.71 MB/s          96 B/op          2 allocs/op
BenchmarkMemoryStorage_CountKeys
BenchmarkMemoryStorage_CountKeys-8              11597236               106.0 ns/op         9.44 MB/s           0 B/op          0 allocs/op
BenchmarkMemoryStorage_CountKeys_P100
BenchmarkMemoryStorage_CountKeys_P100-8          6730882               170.8 ns/op         5.85 MB/s           0 B/op          0 allocs/op
BenchmarkMemoryStorage_Delete
BenchmarkMemoryStorage_Delete-8                  2391238               499.5 ns/op         2.00 MB/s          96 B/op          2 allocs/op
BenchmarkMemoryStorage_Delete_P100
BenchmarkMemoryStorage_Delete_P100-8             8513846               136.9 ns/op         7.30 MB/s          96 B/op          2 allocs/op
PASS
ok      github.com/hati-sh/hati/storage 19.217s
```