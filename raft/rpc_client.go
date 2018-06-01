package raft

import (
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

const (
	// RPCMethodRequestVote is an rpc method.
	RPCMethodRequestVote = "ServerMethods.RequestVote"
	// RPCMethodAppendEntries is an rpc method.
	RPCMethodAppendEntries = "ServerMethods.AppendEntries"

	// DefaultClientDialTimeout is a default.
	DefaultClientDialTimeout = 500 * time.Millisecond
	// DefaultClientCallTimeout is the default time to wait for a call result.
	DefaultClientCallTimeout = 500 * time.Millisecond
	// DefaultClientRedialWait is the default time to wait between rpc redial attempts.
	DefaultClientRedialWait = 5 * time.Second
	// DefaultClientConnectTimeout is the total time allowed to reach the remote.
	DefaultClientConnectTimeout = 30 * time.Second
)

// NewRPCClient creates a new rpc client.
func NewRPCClient(remoteAddr string) *RPCClient {
	return &RPCClient{
		remoteAddr:  remoteAddr,
		latch:       &worker.Latch{},
		dialTimeout: DefaultClientDialTimeout,
		callTimeout: DefaultClientCallTimeout,
		redialWait:  DefaultClientRedialWait,
	}
}

// RPCClient is the net/rpc client to talk to other nodes.
type RPCClient struct {
	sync.Mutex

	remoteAddr string
	conn       net.Conn
	client     *rpc.Client
	latch      *worker.Latch
	log        *logger.Logger

	dialTimeout time.Duration
	callTimeout time.Duration
	redialWait  time.Duration
}

// DialTimeout is the timeout for dialing new connections
func (c *RPCClient) DialTimeout() time.Duration {
	return c.dialTimeout
}

// WithDialTimeout sets the DialTimeout
func (c *RPCClient) WithDialTimeout(d time.Duration) *RPCClient {
	c.dialTimeout = d
	return c
}

// CallTimeout is the timeout for individual rpc calls
func (c *RPCClient) CallTimeout() time.Duration {
	return c.callTimeout
}

// WithCallTimeout is the timeout for individual rpc calls
func (c *RPCClient) WithCallTimeout(d time.Duration) *RPCClient {
	c.callTimeout = d
	return c
}

// WithLogger sets the logger.
func (c *RPCClient) WithLogger(log *logger.Logger) *RPCClient {
	c.log = log
	return c
}

// Logger returns the logger.
func (c *RPCClient) Logger() *logger.Logger {
	return c.log
}

// WithRemoteAddr sets the remote addr.
func (c *RPCClient) WithRemoteAddr(addr string) *RPCClient {
	c.remoteAddr = addr
	return c
}

// RemoteAddr returns the remote address.
func (c *RPCClient) RemoteAddr() string { return c.remoteAddr }

// Open dials the remote, it will only try once, and won't
// dial if the connection is already up.
func (c *RPCClient) Open() error {
	c.Lock()
	defer c.Unlock()

	if c.client != nil {
		return nil
	}

	var err error
	c.conn, err = net.DialTimeout("tcp", c.remoteAddr, c.dialTimeout)
	if err != nil {
		return exception.Wrap(err)
	}
	c.client = rpc.NewClient(c.conn)
	return nil
}

// RequestVote implements the request vote handler.
func (c *RPCClient) RequestVote(args *RequestVote) (*RequestVoteResults, error) {
	if err := c.Open(); err != nil {
		return nil, err
	}
	var res RequestVoteResults
	err := c.callWithTimeout(RPCMethodRequestVote, args, &res)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

// AppendEntries implements the append entries request handler.
func (c *RPCClient) AppendEntries(args *AppendEntries) (*AppendEntriesResults, error) {
	if err := c.Open(); err != nil {
		return nil, err
	}

	var res AppendEntriesResults
	err := c.callWithTimeout(RPCMethodAppendEntries, args, &res)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

// call invokes a method with the default call timeout.
func (c *RPCClient) callWithTimeout(method string, args interface{}, reply interface{}) error {
	timeout := time.NewTimer(c.callTimeout)

	c.Lock()
	defer c.Unlock()

	result := c.client.Go(method, args, reply, nil)
	select {
	case <-timeout.C:
		c.client.Close()
		c.client = nil
		return exception.New("rpc call timeout").WithMessagef("method: %s", method)
	case <-result.Done:
		if result.Error != nil {
			return exception.Wrap(result.Error)
		}
		return nil
	}
}

// Close closes the transport.
func (c *RPCClient) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.client == nil {
		return nil
	}
	c.latch.Stop()
	err := exception.Wrap(c.client.Close())
	<-c.latch.NotifyStopped()
	return err
}

func (c *RPCClient) err(err error) error {
	if c.log != nil && err != nil {
		c.log.Error(err)
	}
	return err
}
