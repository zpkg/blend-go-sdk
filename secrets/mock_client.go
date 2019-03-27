package secrets

import (
	"context"
	"fmt"
	"strings"
)

var _ Client = &MockClient{}

// NewMockClient creates a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{
		SecretValues: make(map[string]Values),
	}
}

// MockClient is a mock events client
type MockClient struct {
	SecretValues map[string]Values
}

// Put puts a value.
func (c *MockClient) Put(_ context.Context, key string, data Values, options ...Option) error {
	c.SecretValues[key] = data

	return nil
}

// Get gets a value at a given key.
func (c *MockClient) Get(_ context.Context, key string, options ...Option) (Values, error) {
	val, exists := c.SecretValues[key]
	if !exists {
		return nil, fmt.Errorf("Key not found: %s", key)
	}

	return val, nil
}

// Delete deletes a key.
func (c *MockClient) Delete(_ context.Context, key string, options ...Option) error {
	if _, exists := c.SecretValues[key]; !exists {
		return fmt.Errorf("Key not found: %s", key)
	}

	delete(c.SecretValues, key)
	return nil
}

// List lists keys on a path
func (c *MockClient) List(_ context.Context, path string, options ...Option) ([]string, error) {
	keys := make([]string, 0)
	folderSet := make(map[string]struct{})
	p := path
	if !strings.HasSuffix(path, "/") {
		p = path + "/"
	}
	for k := range c.SecretValues {
		if strings.HasPrefix(k, p) {
			s := strings.TrimPrefix(k, p)
			if strings.ContainsRune(s, '/') {
				folder := fmt.Sprintf("%s/", strings.Split(s, "/")[0])
				if _, ok := folderSet[folder]; !ok {
					folderSet[folder] = struct{}{}
					keys = append(keys, folder)
				}
			} else {
				keys = append(keys, s)
			}
		}
	}
	return keys, nil
}