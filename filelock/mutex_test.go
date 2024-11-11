/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package filelock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/filelock"
	"github.com/zpkg/blend-go-sdk/uuid"
)

func Test_Mutex_RLock(t *testing.T) {
	its := assert.New(t)

	tempFilePath := filepath.Join(os.TempDir(), uuid.V4().String()+".lock")
	mu := filelock.MutexAt(tempFilePath)

	unlock, err := mu.RLock()
	its.Nil(err)

	stat, err := os.Stat(tempFilePath)
	its.Nil(err)
	its.NotNil(stat)
	its.False(stat.IsDir())

	unlock()

	stat, err = os.Stat(tempFilePath)
	its.Nil(err)
	its.NotNil(stat)
	its.False(stat.IsDir())
}

func Test_Mutex_Lock(t *testing.T) {
	its := assert.New(t)

	tempFilePath := filepath.Join(os.TempDir(), uuid.V4().String()+".lock")
	mu := filelock.MutexAt(tempFilePath)

	unlock, err := mu.Lock()
	its.Nil(err)

	stat, err := os.Stat(tempFilePath)
	its.Nil(err)
	its.NotNil(stat)
	its.False(stat.IsDir())

	unlock()

	stat, err = os.Stat(tempFilePath)
	its.Nil(err)
	its.NotNil(stat)
	its.False(stat.IsDir())
}
