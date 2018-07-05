package secrets

import (
	"encoding/json"
	"path/filepath"
)

// assert kv1 implements kv.
var (
	_ KV = &kv2{}
)

// kv2 defines key value version 2 interactions
type kv2 struct {
	client *VaultClient
}

func (kv2 kv2) Put(key string, data Values, options ...Option) error {
	contents, err := kv2.client.jsonBody(SecretData{Data: data})
	if err != nil {
		return err
	}
	req := kv2.client.createRequest(MethodPut, filepath.Join("/v1/", key), options...)
	req.Body = contents
	res, err := kv2.client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}

func (kv2 kv2) Get(key string, options ...Option) (Values, error) {
	req := kv2.client.createRequest(MethodGet, filepath.Join("/v1/", key), options...)
	res, err := kv2.client.send(req)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var response SecretV2
	if err := json.NewDecoder(res).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data.Data, nil
}

// Delete puts a key.
func (kv2 kv2) Delete(key string, options ...Option) error {
	req := kv2.client.createRequest(MethodDelete, filepath.Join("/v1/", key), options...)
	res, err := kv2.client.send(req)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}
