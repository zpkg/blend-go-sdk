
export RAFT_PEERS=127.0.0.1:6060,127.0.0.1:6061,127.0.0.1:6062,127.0.0.1:6063,127.0.0.1:6064

run:
	@recover --delay 0 --wait 10000 --verbose -- go run main.go --id=worker-0 --bind-addr=:6060 & \
	recover --delay 10 --wait 10000 --verbose -- go run main.go --id=worker-1 --bind-addr=:6061 & \
	recover --delay 20 --wait 10000 --verbose -- go run main.go --id=worker-2 --bind-addr=:6062 & \
	recover --delay 30 --wait 10000 --verbose -- go run main.go --id=worker-3 --bind-addr=:6063 & \
	recover --delay 40 --wait 10000 --verbose -- go run main.go --id=worker-4 --bind-addr=:6064 && fg

run-node0:
	@RAFT_ID=worker-0 RAFT_BIND_ADDR=:6060 recover --wait 5000 --verbose -- go run main.go

run-node1:
	@RAFT_ID=worker-1 RAFT_BIND_ADDR=:6061 recover --wait 5000 --verbose -- go run main.go

run-node2:
	@RAFT_ID=worker-2 RAFT_BIND_ADDR=:6062 recover --wait 5000 --verbose -- go run main.go

run-node3:
	@RAFT_ID=worker-3 RAFT_BIND_ADDR=:6063 recover --wait 5000 --verbose -- go run main.go

run-node4:
	@RAFT_ID=worker-4 RAFT_BIND_ADDR=:6064 recover --wait 5000 --verbose -- go run main.go
