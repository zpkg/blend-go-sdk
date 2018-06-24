package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

var (
	// Assert RPCClient is a client.
	_ Client = &RPCClient{}
)

// NewRPCClient creates a new rpc client.
func NewRPCClient(remoteAddr string) *RPCClient {
	return &RPCClient{
		remoteAddr: remoteAddr,
		transport:  &http.Transport{},
		client:     &http.Client{},
		timeout:    DefaultClientTimeout,
	}
}

// RPCClient is the net/rpc client to talk to other nodes.
type RPCClient struct {
	sync.Mutex

	remoteAddr string
	log        *logger.Logger

	transport *http.Transport
	client    *http.Client

	timeout time.Duration
}

// Timeout is the timeout for dialing new connections
func (c *RPCClient) Timeout() time.Duration {
	return c.timeout
}

// WithTimeout sets the DialTimeout
func (c *RPCClient) WithTimeout(d time.Duration) *RPCClient {
	c.timeout = d
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
func (c *RPCClient) RemoteAddr() string {
	return c.remoteAddr
}

// Open opens the connection.
func (c *RPCClient) Open() error {
	c.client = &http.Client{
		Timeout:   c.timeout,
		Transport: c.transport,
	}
	return nil
}

// Close is a nop right now.
func (c *RPCClient) Close() error {
	return nil
}

// RequestVote implements the request vote handler.
func (c *RPCClient) RequestVote(args *RequestVote) (*RequestVoteResults, error) {
	var res RequestVoteResults
	err := c.callWithTimeout(RPCMethodRequestVote, args, &res)
	if err != nil {
		return nil, exception.New(err)
	}
	return &res, nil
}

// AppendEntries implements the append entries request handler.
func (c *RPCClient) AppendEntries(args *AppendEntries) (*AppendEntriesResults, error) {
	var res AppendEntriesResults
	err := c.callWithTimeout(RPCMethodAppendEntries, args, &res)
	if err != nil {
		return nil, exception.New(err)
	}
	return &res, nil
}

// call invokes a method with the default call timeout.
func (c *RPCClient) callWithTimeout(method string, args interface{}, reply interface{}) error {

	reqURL, err := url.Parse(fmt.Sprintf("http://%s/%s", c.remoteAddr, method))
	if err != nil {
		return exception.Wrap(err)
	}

	body, err := c.encode(args)
	if err != nil {
		return exception.Wrap(err)
	}

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Body:   body,
	}

	if c.log != nil {
		defer func() {
			c.log.Trigger(logger.NewHTTPRequestEvent(req).WithFlag("rpc.call"))
		}()
	}

	res, err := c.client.Do(req)
	if err != nil {
		return exception.Wrap(err)
	}
	if res.StatusCode > 299 {
		return exception.New("non-2xx returned from rpc server").WithMessagef("status code returned: %d", res.StatusCode)
	}

	if err := c.decode(reply, res.Body); err != nil {
		return err
	}
	return nil
}

func (c *RPCClient) encode(obj interface{}) (io.ReadCloser, error) {
	buffer := new(bytes.Buffer)

	if err := json.NewEncoder(buffer).Encode(obj); err != nil {
		return nil, exception.New(err)
	}
	return ioutil.NopCloser(buffer), nil
}

func (c *RPCClient) decode(obj interface{}, contents io.ReadCloser) error {
	if contents == nil {
		return exception.New("response body unset; cannot continue")
	}
	defer contents.Close()
	return exception.New(json.NewDecoder(contents).Decode(&obj))
}

func (c *RPCClient) err(err error) error {
	if c.log != nil && err != nil {
		c.log.Error(err)
	}
	return err
}
