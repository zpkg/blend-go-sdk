/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/vault"
)

type myConfig struct {
	Environment string `yaml:"environment" env:"SERVICE_ENV" secret:"environment"`
	Version     string `yaml:"version" env:"SERVICE_VERSION"`
	Secret      string `yaml:"secret" secret:"secret,base64"`
}

func main() {
	log := logger.All()
	client, _ := vault.New(vault.OptConfigFromEnv(), vault.OptLog(log))

	keyPath := "configTest"
	err := client.WriteInto(context.TODO(), keyPath, myConfig{
		Environment: "test",
		Version:     "wont be in output",
		Secret:      base64.StdEncoding.EncodeToString([]byte("a super secure one")),
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cfg myConfig
	err = client.ReadInto(context.TODO(), keyPath, &cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	contents, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(contents))
}
