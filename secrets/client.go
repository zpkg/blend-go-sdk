package secrets

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/logger"
)

const (
	// MethodGet is a request method.
	MethodGet = "GET"
	// MethodPut is a request method.
	MethodPut = "PUT"
	// MethodDelete is a request method.
	MethodDelete = "DELETE"

	// HeaderVaultToken is the vault token header.
	HeaderVaultToken = "X-Vault-Token"
	// HeaderContentType is the content type header.
	HeaderContentType = "Content-Type"
	// ContentTypeApplicationJSON is a content type.
	ContentTypeApplicationJSON = "application/json"
)

// New returns a new client.
func New() *Client {
	return &Client{
		addr: DefaultAddr,
		client: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: &http.Transport{},
		},
	}
}

// NewFromConfig returns a new client from a config.
func NewFromConfig(cfg *Config) *Client {
	return New().WithAddr(cfg.GetAddr()).
		WithToken(cfg.GetToken()).
		WithTimeout(cfg.GetTimeout())
}

// Client is a client to talk to the secrets store.
type Client struct {
	addr  string
	token string
	log   *logger.Logger

	client   *http.Client
	certPool *x509.CertPool
}

// WithAddr set the client remote addr.
func (c *Client) WithAddr(addr string) *Client {
	c.addr = addr
	return c
}

// Addr returns the client addr.
func (c *Client) Addr() string {
	return c.addr
}

// WithToken sets the token.
func (c *Client) WithToken(token string) *Client {
	c.token = token
	return c
}

// Token returns the token.
func (c *Client) Token() string {
	return c.token
}

// WithTimeout sets the client timeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.Timeout = timeout
	return c
}

// Timeout returns the timeout.
func (c *Client) Timeout() time.Duration {
	return c.client.Timeout
}

// WithCertPool returns the cert pool.
func (c *Client) WithCertPool(certPool *x509.CertPool) *Client {
	c.certPool = certPool
	return c
}

// CertPool returns the cert pool.
func (c *Client) CertPool() *x509.CertPool {
	return c.certPool
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

// Put puts a value.
func (c *Client) Put(key string, value Values) error {
	contents, err := json.Marshal(StorageData{Value: value})
	if err != nil {
		return err
	}

	req, err := c.createRequest(MethodPut)
	if err != nil {
		return err
	}
	req.Header.Add(HeaderContentType, ContentTypeApplicationJSON)
	req.URL.Path = filepath.Join("/v1/", key)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
	res, err := c.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// Get gets a value at a given key.
func (c *Client) Get(key string) (Values, error) {
	req, err := c.createRequest(MethodGet)
	if err != nil {
		return nil, err
	}
	req.URL.Path = filepath.Join("/v1/", key)
	res, err := c.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response StorageResponse
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data.Value, nil
}

// Delete puts a key.
func (c *Client) Delete(key string) error {
	req, err := c.createRequest(MethodDelete)
	if err != nil {
		return err
	}
	req.URL.Path = filepath.Join("/v1/", key)
	res, err := c.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

func (c *Client) createRequest(method string) (*http.Request, error) {
	remote, err := url.ParseRequestURI(c.addr)
	if err != nil {
		return nil, err
	}
	return &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token()},
		},
	}, nil
}

func (c *Client) send(req *http.Request) (io.ReadCloser, error) {
	if c.log != nil {
		c.log.Trigger(NewEvent(req))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		contents, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("non-2xx returned from remote: %d; %v", res.StatusCode, string(contents))
	}
	return res.Body, nil
}
