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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeIsDir(t *testing.T) {
	node := &Node{
		Path:      "/var/file",
		Value:     "something",
		Directory: true,
	}
	assert.False(t, node.IsFile())
	assert.True(t, node.IsDir())
	node.Directory = false
	assert.False(t, node.IsDir())
}

func TestNodeIsFile(t *testing.T) {
	node := &Node{
		Path:      "/var/file",
		Value:     "something",
		Directory: false,
	}
	assert.True(t, node.IsFile())
	assert.False(t, node.IsDir())
}
