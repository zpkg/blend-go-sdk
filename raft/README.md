raft
----

This is a very basic implementation of the raft consensus protocol. It focuses exclusively on leader elections.

# Design 

- Nodes should be configured ahead of time, the node configuration should be static. This prevents dynamic resizing of nodes, but also means that we don't have to worry about malicious nodes being added without our control.
- We use `net/rpc` to communicate with other nodes. This is a very bare bones and fast rpc system used by a number of smaller services and is preferrable to something more heavy like grpc.
- RPC calls are encoded using `encoding/gob` (the default for `net/rpc`.
- The defaults for the config are optimized for a three node cluster on a relatively low latency network. You'd want to tune these yourself if you're running in a different scenario.

# Example

The example found in `_example/node/main.go` is the absolute barest bones implementation.

It creates a new config from environment variables, creates a new raft node, then adds configured peers as rpc endpoints. 