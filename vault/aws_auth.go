/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/blend/go-sdk/ex"
)

// AWSAuth defines vault aws auth methods
type AWSAuth struct {
	CredentialProvider CredentialProvider
}

// NewAWSAuth creates a new AWS struct
func NewAWSAuth(opts ...AWSAuthOption) (*AWSAuth, error) {
	auth := &AWSAuth{
		CredentialProvider: GetIAMAuthCredentials,
	}
	var err error
	for _, opt := range opts {
		if err = opt(auth); err != nil {
			return nil, err
		}
	}
	return auth, nil
}

// AWSAuthOption mutates an AWSAuth instance
type AWSAuthOption func(*AWSAuth) error

// CredentialProvider defines the credential provider func interface
type CredentialProvider func(roleARN string) (*credentials.Credentials, error)

// OptAWSAuthCredentialProvider sets the credential provider
func OptAWSAuthCredentialProvider(cp CredentialProvider) AWSAuthOption {
	return func(a *AWSAuth) error {
		a.CredentialProvider = cp
		return nil
	}
}

// AWSIAMLogin returns a vault token given the instance role which invokes this function
func (a *AWSAuth) AWSIAMLogin(ctx context.Context, client HTTPClient, baseURL url.URL, roleName, roleARN, service, region string) (string, error) {
	headers, err := a.GetCallerIdentitySignedHeaders(roleARN, service, region)
	if err != nil {
		return "", ex.New(err)
	}

	request, err := createVaultLoginRequest(roleName, baseURL, headers)
	if err != nil {
		return "", ex.New(err)
	}

	res, err := client.Do(request)
	if err != nil {
		return "", ex.New(err)
	}
	defer res.Body.Close()

	var response AWSAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", ex.New(err)
	}
	return response.Auth.ClientToken, nil
}

// GetCallerIdentitySignedHeaders gets signed get caller identity request headers
func (a *AWSAuth) GetCallerIdentitySignedHeaders(roleARN, service, region string) (http.Header, error) {
	body := strings.NewReader(STSGetIdentityBody)
	req, err := http.NewRequest(MethodPost, STSURL, body)
	if err != nil {
		return nil, ex.New(err)
	}

	credentials, err := a.CredentialProvider(roleARN)
	if err != nil {
		return nil, ex.New(err)
	}

	signer := v4.NewSigner(credentials)
	return signer.Sign(req, body, service, region, time.Now())
}

// GetIAMAuthCredentials is a credential provider to be passed in as input into the AWSAuth struct
func GetIAMAuthCredentials(roleARN string) (*credentials.Credentials, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, ex.New(err)
	}
	credentials := stscreds.NewCredentials(session, roleARN)
	return credentials, nil
}

func createVaultLoginRequest(roleName string, baseURL url.URL, header http.Header) (*http.Request, error) {
	baseURL.Path = AWSAuthLoginPath
	headers, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	body := map[string]string{
		"role":                    roleName,
		"iam_http_request_method": MethodPost,
		"iam_request_url":         base64.StdEncoding.EncodeToString([]byte(STSURL)),
		"iam_request_body":        base64.StdEncoding.EncodeToString([]byte(STSGetIdentityBody)),
		"iam_request_headers":     base64.StdEncoding.EncodeToString(headers),
	}

	contents, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &http.Request{
		URL:    &baseURL,
		Method: MethodPost,
		Body:   ioutil.NopCloser(bytes.NewReader(contents)),
	}, nil
}
