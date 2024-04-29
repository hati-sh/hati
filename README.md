# @hatish/hati

`v0.2.0-dev`

```
IN DEVELOPMENT. DO NOT USE IN PRODUCTION.
```

Hati is a simple distributed (in-memory and persistent) key/value store and a message broker.

Hati is meant to work in trusted networks - where all nodes operators know each other and can be trusted.

Once connected as a client to Hati server you can publish commands which will be processed by the server and published over the Hati network - Hati nodes can be connected to each other creating network of nodes but it is prefectly fine to run Hati as a single instance.

```
  [req]      [req]      [req]
    ^          ^          ^
    v          v          v
[London] [Hong Kong] [New York]
    ^          ^          ^
    |----------|----------|
```

```
  [req]
    ^
    v
[Hong Kong]
```

## Current features

- TCP server
- JSON-RPC server
- Storing data in in-memory storage type - with sharding

## Configuration

While starting Hati `hati start` there are configurational flags available to be set.

- `--tcp` - indicates if should start TCP server, by default it is off
- `--tcp-host` - determines bind host for TCP server, default value is `0.0.0.0`
- `--tcp-port` - bind port for TCP server, default value is `4242`
- `--rpc` - indicates if should start JSON-RPC server, by default it is off 
- `--rpc-host` - bind host for JSON-RPC server, default `0.0.0.0`
- `--rpc-port` - bind port for JSON-RPC server, default `6666`
- `--data-dir` - absolute path to directory where Hati can store files, default `/current/path/to-hati/data`
- `--cpu-num` - number of CPU cores which should be used by Hati, by default it will set for as many as available

## Protocol Commands

### Key-Value Storage

`SET <type> <ttl> <key> <value>\n` - save key with provided value to the selected storage type.

Hati offers two storage types: `memory` and `hdd` . By default `<ttl>` is set to zero `0` which means that value will be stored on the hard-drive. Ttl value is in ms, if higher than `0` Hati can guarantee that value will be stored at minimum for provided ttl value and will be removed from the storage shortly (as soon as possible) after that.

**^^ TTL IS NOT IMPLEMENTED YET ^^**

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
**A**
- [ ] Persistent storage
- [ ] Implement TTL
- [ ] Rpc server cancel context for graceful shutdown
- [ ] CLI command to export/import backup of persistent storage 
- [ ] Message broker

**B**
- [ ] Rebuilding persistent database if number of shards has changed
- [ ] Nodes clustering
- [ ] Data synchronization between nodes
  - TCP
  - P2P (? - TBD)
