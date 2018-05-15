package secrets

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/exception"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
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
	DefaultBufferPoolSize = 1024

	// ReflectTagName is a reflect tag name.
	ReflectTagName = "secret"

	// Version1 is a constant.
	Version1 = "1"
	// Version2 is a constant.
	Version2 = "2"
)

// New returns a new client.
func New() (*Client, error) {
	return NewFromConfig(&Config{})
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
	var certPool *CertPool
	if caPaths := cfg.GetRootCAs(); len(caPaths) > 0 {
		certPool, err = NewCertPool()
		if err != nil {
			return nil, err
		}
		err = certPool.AddPaths(caPaths...)
		if err != nil {
			return nil, err
		}
		xport.TLSClientConfig = &tls.Config{
			RootCAs: certPool.Pool(),
		}
	}
	client := &Client{
		remote:     remote,
		bufferPool: NewBufferPool(DefaultBufferPoolSize),
		token:      cfg.GetToken(),
		certPool:   certPool,
		client: &http.Client{
			Timeout:   cfg.GetTimeout(),
			Transport: xport,
		},
	}

	client.kv1 = &kv1{client: client}
	client.kv2 = &kv2{client: client}
	return client, nil
}

// NewFromEnv is a helper to create a client from a config read from the environment.
func NewFromEnv() (*Client, error) {
	return NewFromConfig(NewConfigFromEnv())
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

	kv1 *kv1
	kv2 *kv2

	bufferPool *BufferPool
	client     HTTPClient
	certPool   *CertPool
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

// WithHTTPClient sets the http client.
func (c *Client) WithHTTPClient(hc HTTPClient) *Client {
	c.client = hc
	return c
}

// HTTPClient sets the http client.
func (c *Client) HTTPClient() HTTPClient {
	return c.client
}

// CertPool returns the cert pool.
func (c *Client) CertPool() *CertPool {
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
func (c *Client) Put(key string, data Values, options ...Option) error {
	backend, err := c.backend()
	if err != nil {
		return err
	}

	return backend.Put(key, data, options...)
}

// Get gets a value at a given key.
func (c *Client) Get(key string, options ...Option) (Values, error) {
	backend, err := c.backend()
	if err != nil {
		return nil, err
	}

	return backend.Get(key, options...)
}

// Delete puts a key.
func (c *Client) Delete(key string, options ...Option) error {
	backend, err := c.backend()
	if err != nil {
		return err
	}
	return backend.Delete(key, options...)
}

// ReadInto reads a secret into an object.
func (c *Client) ReadInto(key string, obj interface{}, options ...Option) error {
	response, err := c.Get(key, options...)
	if err != nil {
		return err
	}
	return util.Reflection.PatchStrings(ReflectTagName, response, obj)
}

// WriteInto writes an object into a secret at a given key.
func (c *Client) WriteInto(key string, obj interface{}, options ...Option) error {
	return c.Put(key, util.Reflection.DecomposeStrings(ReflectTagName, obj), options...)
}

// --------------------------------------------------------------------------------
// utility methods
// --------------------------------------------------------------------------------

func (c *Client) backend() (KV, error) {
	version, err := c.getVersion()
	if err != nil {
		return nil, err
	}
	switch version {
	case Version1:
		return c.kv1, nil
	case Version2:
		return c.kv2, nil
	}
	return nil, exception.New("invalid kv version").WithMessagef("version: %s", version)
}

func (c *Client) getVersion() (string, error) {
	meta, err := c.getMountMeta()
	if err != nil {
		return "", err
	}
	return meta.Data.Options["version"], nil
}

func (c *Client) getMountMeta() (*MountResponse, error) {
	req := c.createRequest(MethodGet, "/v1/sys/internal/ui/mounts/secret/")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response MountResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
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

// applyOptions applies options to a request.
func (c *Client) applyOptions(req *http.Request, options ...Option) {
	for _, opt := range options {
		opt(req)
	}
}

func (c *Client) createRequest(method, path string, options ...Option) *http.Request {
	remote := c.copyRemote()
	remote.Path = path
	req := &http.Request{
		Method: method,
		URL:    remote,
		Header: http.Header{
			HeaderVaultToken: []string{c.Token()},
		},
	}
	c.applyOptions(req, options...)
	return req
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
