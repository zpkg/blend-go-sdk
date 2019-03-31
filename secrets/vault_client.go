package secrets

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// assert VaultClient implements Client
var _ Client = &VaultClient{}

// NewVaultClient returns a new client.
func NewVaultClient() (*VaultClient, error) {
	return NewVaultClientFromConfig(&Config{})
}

// NewVaultClientFromConfig returns a new client from a config.
func NewVaultClientFromConfig(cfg *Config) (*VaultClient, error) {
	xport := &http.Transport{}
	err := http2.ConfigureTransport(xport)
	if err != nil {
		return nil, err
	}
	remote, err := url.ParseRequestURI(cfg.AddrOrDefault())
	if err != nil {
		return nil, err
	}
	var certPool *CertPool
	if caPaths := cfg.RootCAs; len(caPaths) > 0 {
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
	client := &VaultClient{
		remote:     remote,
		mount:      cfg.MountOrDefault(),
		bufferPool: NewBufferPool(DefaultBufferPoolSize),
		token:      cfg.Token,
		certPool:   certPool,
		client: &http.Client{
			Timeout:   cfg.TimeoutOrDefault(),
			Transport: xport,
		},
	}

	client.kv1 = &kv1{client: client}
	client.kv2 = &kv2{client: client}
	return client, nil
}

// NewVaultClientFromEnv is a helper to create a client from a config read from the environment.
func NewVaultClientFromEnv() (*VaultClient, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewVaultClientFromConfig(cfg)
}

// Must does things with the error such as panic.
func Must(c *VaultClient, err error) *VaultClient {
	if err != nil {
		panic(err)
	}
	return c
}

// VaultClient is a client to talk to the secrets store.
type VaultClient struct {
	Remote *url.URL
	Token  string
	Mount  string
	Log    logger.Log

	KV1 *kv1
	KV2 *kv2

	Client   HTTPClient
	CertPool *CertPool
}

// Put puts a value.
func (c *VaultClient) Put(ctx context.Context, key string, data Values, options ...RequestOption) error {
	backend, err := c.backend(ctx, key)
	if err != nil {
		return err
	}

	return backend.Put(ctx, key, data, options...)
}

// Get gets a value at a given key.
func (c *VaultClient) Get(ctx context.Context, key string, options ...RequestOption) (Values, error) {
	backend, err := c.backend(ctx, key)
	if err != nil {
		return nil, err
	}

	return backend.Get(ctx, key, options...)
}

// Delete puts a key.
func (c *VaultClient) Delete(ctx context.Context, key string, options ...RequestOption) error {
	backend, err := c.backend(ctx, key)
	if err != nil {
		return err
	}
	return backend.Delete(ctx, key, options...)
}

// ReadInto reads a secret into an object.
func (c *VaultClient) ReadInto(ctx context.Context, key string, obj interface{}, options ...RequestOption) error {
	response, err := c.Get(ctx, key, options...)
	if err != nil {
		return err
	}
	return RestoreJSON(response, obj)
}

// WriteInto writes an object into a secret at a given key.
func (c *VaultClient) WriteInto(ctx context.Context, key string, obj interface{}, options ...RequestOption) error {
	data, err := DecomposeJSON(obj)
	if err != nil {
		return err
	}
	return c.Put(ctx, key, data, options...)
}

// --------------------------------------------------------------------------------
// utility methods
// --------------------------------------------------------------------------------

func (c *VaultClient) backend(ctx context.Context, key string) (KV, error) {
	version, err := c.getVersion(ctx, key)
	if err != nil {
		return nil, err
	}
	switch version {
	case Version1:
		return c.KV1, nil
	case Version2:
		return c.KV2, nil
	default:
		return c.KV1, nil
	}
}

func (c *VaultClient) getVersion(ctx context.Context, key string) (string, error) {
	meta, err := c.getMountMeta(ctx, filepath.Join(c.mount, key))
	if err != nil {
		return "", err
	}
	return meta.Data.Options["version"], nil
}

func (c *VaultClient) getMountMeta(ctx context.Context, key string) (*MountResponse, error) {
	req := c.createRequest(MethodGet, filepath.Join("/v1/sys/internal/ui/mounts/", key))
	req = req.WithContext(ctx)

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

func (c *VaultClient) jsonBody(input interface{}) (io.ReadCloser, error) {
	buf := c.bufferPool.Get()
	err := json.NewEncoder(buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *VaultClient) readJSON(r io.Reader, output interface{}) error {
	return json.NewDecoder(r).Decode(output)
}

// copyRemote returns a copy of our remote.
func (c *VaultClient) copyRemote() *url.URL {
	remoteCopy := *c.remote
	return &remoteCopy
}

// applyOptions applies options to a request.
func (c *VaultClient) applyOptions(req *http.Request, options ...RequestOption) {
	for _, opt := range options {
		opt(req)
	}
}

func (c *VaultClient) createRequest(method, path string, options ...RequestOption) *http.Request {
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

func (c *VaultClient) send(req *http.Request) (io.ReadCloser, error) {
	if c.log != nil {
		c.log.Trigger(req.Context(), NewEvent(req))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		buf := c.bufferPool.Get()
		defer buf.Close()
		io.Copy(buf, res.Body)
		return nil, exception.New(ExceptionClassForStatus(res.StatusCode), exception.OptMessagef("status: %d; %v", res.StatusCode, buf.String()))
	}
	return res.Body, nil
}

func (c *VaultClient) discard(res io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer res.Close()
	if _, err := io.Copy(ioutil.Discard, res); err != nil {
		return err
	}
	return nil
}
