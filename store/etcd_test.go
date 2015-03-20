/*
Copyright 2014 Rohith All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"time"
)

const (
	ETCD_URL = "etcd://127.0.0.1:4001"
	ETCD_KEY = "/TEST_KEY"
	ETCD_VAL = "VALUE"
)

var (
	client Store
	updates NodeUpdateChannel
)

func TestSetup(t *testing.T) {
	location, err := url.Parse(ETCD_URL)
	assert.Nil(t, err)
	assert.NotNil(t, location)
	updates = make(NodeUpdateChannel, 10)
	client, err = NewEtcdStoreClient(location, updates)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestSetGet(t *testing.T) {
	err := client.Set(ETCD_KEY, ETCD_VAL)
	assert.Nil(t, err)
	node, err := client.Get(ETCD_KEY)
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, ETCD_KEY, node.Path)
	assert.Equal(t, ETCD_VAL, node.Value)
	assert.False(t, node.IsDir())
	assert.True(t, node.IsFile())
}

func TestDelete(t *testing.T) {
	err := client.Set(ETCD_KEY, ETCD_VAL)
	assert.Nil(t, err)
	assert.Nil(t, client.Delete(ETCD_KEY))
}

func TestWatch(t *testing.T) {
	err := client.Set(ETCD_KEY, ETCD_VAL)
	assert.Nil(t, err)
	// check it is there
	node, err := client.Get(ETCD_KEY)
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, ETCD_KEY, node.Path)
	assert.Equal(t, ETCD_VAL, node.Value)

	// add a watch
	client.Watch(ETCD_KEY)
}

func TestWatchNotification(t *testing.T) {
	// update the key
	err := client.Set(ETCD_KEY, ETCD_VAL)
	assert.Nil(t, err)
	timeout := time.After(time.Duration(5) * time.Second)

	// wait for a change
	select {
	case event := <-updates:
		timeout = nil
		assert.NotNil(t, event.Node)
		assert.Equal(t, 1, event.Operation)
		assert.Equal(t, ETCD_KEY, event.Node.Path)
	case <- timeout:
		assert.Fail(t, "we timed out waiting for an event")
	}
}

func TestUnWatch(t *testing.T) {
	client.Unwatch(ETCD_KEY)
	timeout := time.After(time.Duration(100) * time.Millisecond)

	// wait for a change
	select {
	case <-updates:
		assert.Fail(t, "we should not have recieved an event here")
	case <- timeout:
	}
}

func TestList(t *testing.T) {
	err := client.Set("/test/one", "1")
	assert.Nil(t, err)
	err = client.Set("/test/two", "2")
	assert.Nil(t, err)
	err = client.Set("/test/three", "3")
	assert.Nil(t, err)

	list, err := client.List("/test")
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.NotEmpty(t, list)
	assert.Equal(t, 3, len(list))
}

