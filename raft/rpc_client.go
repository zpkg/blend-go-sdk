package raft

import (
	"net"
	"net/rpc"
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
	// DefaultClientRedialWait is the default time to wait between rpc redial attempts.
	DefaultClientRedialWait = 5 * time.Second
	// DefaultClientConnectTimeout is the total time allowed to reach the remote.
	DefaultClientConnectTimeout = 30 * time.Second
)

// NewClient creates a new rpc client.
func NewClient(remoteAddr string) *Client {
	return &Client{
		remoteAddr:  remoteAddr,
		latch:       &worker.Latch{},
		dialTimeout: DefaultClientDialTimeout,
		redialWait:  DefaultClientRedialWait,
	}
}

// Client is an rpc peer transport over the network.
type Client struct {
	remoteAddr string
	conn       net.Conn
	client     *rpc.Client
	latch      *worker.Latch
	log        *logger.Logger

	dialTimeout time.Duration
	redialWait  time.Duration
}

// WithLogger sets the logger.
func (c *Client) WithLogger(log *logger.Logger) *Client {
	c.log = log
	return c
}

// Logger returns the logger.
func (c *Client) Logger() *logger.Logger {
	return c.log
}

// RemoteAddr returns the remote address.
func (c *Client) RemoteAddr() string { return c.remoteAddr }

// Open opens the connection.
// It waits `redialWait` time between attempts.
// It will retry indefinitely (until told to stop with `Close()`).
func (c *Client) Open() error {
	for {
		select {
		case <-c.latch.NotifyStop():
			c.latch.Stopped()
			return nil
		default:
			err := c.Dial()
			if err != nil {
				c.log.Warning(err)
				time.Sleep(c.redialWait)
				continue
			}
			return nil
		}
	}
}

// Dial dials the remote, it will only try once, and won't
// dial if the connection is already up.
func (c *Client) Dial() error {
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
func (c *Client) RequestVote(args *RequestVote) (*RequestVoteResults, error) {
	if err := c.Dial(); err != nil {
		return nil, err
	}
	var res RequestVoteResults
	err := c.client.Call(RPCMethodRequestVote, args, &res)
	if err != nil {
		c.err(c.disconnect())
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

// AppendEntries implements the append entries request handler.
func (c *Client) AppendEntries(args *AppendEntries) (*AppendEntriesResults, error) {
	if err := c.Dial(); err != nil {
		return nil, err
	}

	var res AppendEntriesResults
	err := c.client.Call(RPCMethodAppendEntries, args, &res)
	if err != nil {
		c.err(c.disconnect())
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

func (c *Client) disconnect() error {
	if c.client == nil {
		return nil
	}
	if err := c.client.Close(); err != nil {
		return exception.Wrap(err)
	}

	c.client = nil
	return nil
}

// Close closes the transport.
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	c.latch.Stop()
	err := exception.Wrap(c.client.Close())
	<-c.latch.NotifyStopped()
	return err
}

func (c *Client) err(err error) error {
	if c.log != nil && err != nil {
		c.log.Error(err)
	}
	return err
}
