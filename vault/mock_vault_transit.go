/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"encoding/base64"
)

// Assert MockTransitClient implements TransitClient
var (
	_ Client        = MockTransitClient{}
	_ TransitClient = MockTransitClient{}
)

// MockTransitClient skips interactions with the vault for encryption/decryption
type MockTransitClient struct {
	Client
}

// Encrypt just returns the input in the mock
func (m MockTransitClient) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt just returns the input in the mock
func (m MockTransitClient) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(ciphertext)
}
