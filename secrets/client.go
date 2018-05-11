package secrets

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/logger"
)

const (
	// MethodGet is a request method.
	MethodGet = "GET"
	// MethodPost is a request method.
	MethodPost = "POST"
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

	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultBufferPoolSize = 256
)

// New returns a new client.
func New() (*Client, error) {
	xport := &http.Transport{}
	err := http2.ConfigureTransport(xport)
	if err != nil {
		return nil, err
	}
	remote, err := url.ParseRequestURI(DefaultAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		remote:     remote,
		bufferPool: NewBufferPool(DefaultBufferPoolSize),
		client: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: xport,
		},
	}, nil
}

// NewFromConfig returns a new client from a config.
func NewFromConfig(cfg *Config) (*Client, error) {
	xport := &http.Transport{}
	err := http2.ConfigureTransport(xport)
	if err != nil {
		return nil, err
	}
	remote, err := url.ParseRequestURI(cfg.GetAddr())
	if err != nil {
		return nil, err
	}
	return &Client{
		remote:     remote,
		bufferPool: NewBufferPool(DefaultBufferPoolSize),
		token:      cfg.GetToken(),
		client: &http.Client{
			Timeout:   cfg.GetTimeout(),
			Transport: xport,
		},
	}, nil
}

// Must does things with the error such as panic.
func Must(c *Client, err error) *Client {
	if err != nil {
		panic(err)
	}
	return c
}

// Client is a client to talk to the secrets store.
type Client struct {
	remote *url.URL
	token  string
	log    *logger.Logger

	bufferPool *BufferPool
	client     *http.Client
	certPool   *x509.CertPool
}

// WithRemote set the client remote url.
func (c *Client) WithRemote(remote *url.URL) *Client {
	c.remote = remote
	return c
}

// Remote returns the client remote addr.
func (c *Client) Remote() *url.URL {
	return c.remote
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
func (c *Client) Put(key string, data Values) error {
	contents, err := c.jsonBody(data)
	if err != nil {
		return err
	}
	req := c.createRequest(MethodPut, filepath.Join("/v1/", key))
	req.Body = contents
	res, err := c.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// Get gets a value at a given key.
func (c *Client) Get(key string) (Values, error) {
	response, err := c.Meta(key)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// Meta gets the metadata for a key.
func (c *Client) Meta(key string) (*Secret, error) {
	req := c.createRequest(MethodGet, filepath.Join("/v1/", key))
	res, err := c.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response Secret
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Delete puts a key.
func (c *Client) Delete(key string) error {
	req := c.createRequest(MethodDelete, filepath.Join("/v1/", key))
	res, err := c.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

// ListMounts lists mounts.
func (c *Client) ListMounts() (map[string]Mount, error) {
	req := c.createRequest(MethodGet, "/v1/sys/mounts")
	res, err := c.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var result map[string]Mount
	err = c.readJSON(res, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Mount creates a new mount.
func (c *Client) Mount(path string, mountInfo *MountInput) error {
	r := c.createRequest(MethodPost, fmt.Sprintf("/v1/sys/mounts/%s", path))
	body, err := c.jsonBody(mountInfo)
	if err != nil {
		return err
	}
	r.Body = body
	return c.discard(c.send(r))
}

func (c *Client) jsonBody(input interface{}) (io.ReadCloser, error) {
	buf := c.bufferPool.Get()
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) readJSON(r io.Reader, output interface{}) error {
	return json.NewDecoder(r).Decode(output)
}

// copyRemote returns a copy of our remote.
func (c *Client) copyRemote() *url.URL {
	remoteCopy := *c.remote
	return &remoteCopy
}

func (c *Client) createRequest(method, pathFormat string, args ...interface{}) *http.Request {
	remote := c.copyRemote()
	remote.Path = fmt.Sprintf(pathFormat, args...)
	return &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token()},
		},
	}
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
		buf := c.bufferPool.Get()
		defer buf.Close()
		io.Copy(buf, res.Body)
		return nil, fmt.Errorf("non-2xx returned from remote: %d; %v", res.StatusCode, buf.String())
	}
	return res.Body, nil
}

func (c *Client) discard(res io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer res.Close()
	if _, err := io.Copy(ioutil.Discard, res); err != nil {
		return err
	}
	return nil
}
