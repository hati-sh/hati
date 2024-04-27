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

- TCP Server
- Storing data in in-memory storage type

## Configuration

While starting Hati `hati start` there are configurational flags available to be set.

- `--host` - determines bind host for TCP server, default value is `0.0.0.0`
- `--port` - bind port for TCP server, default value is `4242`
- `--rpc` - indicates if should start JSON-RPC server, default `off`
- `--rpc-host` - bind host for JSON-RPC server, default `0.0.0.0`
- `--rpc-port` - bind port for JSON-RPC server, default `6666`
- `--data-dir` - absolute path to directory where Hati can store files, default `/current/path/to-hati/data`
- `--cpu-num` - number of CPU cores which should be used by Hati, by default it will set for as many as available

## Commands

### Key-Value Storage

`SET <type> <ttl> <key> <value>` - save key with provided value to the selected storage type.

Hati offers two storage types: `memory` and `hdd` . By default `<ttl>` is set to zero `0` which means that value will be stored on the hard-drive. Ttl value is in ms, if higher than `0` Hati can guarantee that value will be stored at minimum for provided ttl value and will be removed from the storage shortly (as soon as possible) after that.

**^^ TTL IS NOT IMPLEMENTED YET ^^**

`HAS <type> <key>` - check if provided key exist in given storage type

`GET <type> <key>` - get value for provided key in given storage type

`DELETE <type> <key>` - get value for provided key in given storage type

`FLUSHALL <type>` - flush (delete) all data from given storage type

## To do

- [ ] Implement TTL
- [ ] Persistent storage
- [ ] Message broker
- [ ] Nodes clustering
- [ ] Data synchronization between nodes
