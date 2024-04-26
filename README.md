# @hatish/hati

`v0.1.0-dev`

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

## To do

- [ ] Implement TTL
- [ ] Persistent storage
- [ ] Message broker
- [ ] Nodes clustering
- [ ] Data synchronization between nodes
