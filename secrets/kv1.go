package secrets

import (
	"encoding/json"
	"path/filepath"
)

// assert kv1 implements kv.
var (
	_ KV = &kv1{}
)

// kv1 defines key value version 1 interactions
type kv1 struct {
	client *VaultClient
}

func (kv1 kv1) Put(key string, data Values, options ...Option) error {
	contents, err := kv1.client.jsonBody(data)
	if err != nil {
		return err
	}
	req := kv1.client.createRequest(MethodPut, filepath.Join("/v1/", key), options...)
	req.Body = contents
	res, err := kv1.client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

func (kv1 kv1) Get(key string, options ...Option) (Values, error) {
	req := kv1.client.createRequest(MethodGet, filepath.Join("/v1/", key), options...)
	res, err := kv1.client.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretV1
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// Delete puts a key.
func (kv1 kv1) Delete(key string, options ...Option) error {
	req := kv1.client.createRequest(MethodDelete, filepath.Join("/v1/", key), options...)
	res, err := kv1.client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}
