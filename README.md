# @hatish/hati

Hati is providing simple distributed key/value store (in-memory and persistent). In the future queues with routing.

Hati is meant to work in trusted networks - where all nodes operators know each other and can be trusted.

By default all connections are secured with TLS certificates and clients connecting to the server are required to have their connection secured with TLS certificate as well.

Hati during each startup will generate new certificate and will keep it in-memory.

You can turn off TLS by providing `--tls off` flag during startup.

## CLI

Start server

```
hati start 0.0.0.0 4242
```

Connect client to server

```
hati client 0.0.0.0 4242
```

Once connected as a client to hati server you can publish commands which will be processed by the server and published over the hati network - hati servers can be connected to each other creating networkof servers but it is prefectly fine to run hati as a single instance server.

```
  [req]      [req]      [req]
    ^          ^          ^
    v          v          v
[Londong] [Hong Kong] [New York]
    ^          ^          ^
    |----------|----------|
```

```
  [req]
    ^
    v
[Hong Kong]
```
