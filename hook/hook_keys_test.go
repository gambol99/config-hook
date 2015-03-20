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

package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHookKeyValid(t *testing.T) {
	key := NewHookKeys("test_key")
	assert.NotNil(t, key)
	key.File = "/opt/file/keys"
	key.Flags = ""
	valid, err := key.Valid()
	assert.Nil(t, err)
	assert.True(t, valid)
	key.File = ""
	valid, err = key.Valid()
	assert.NotNil(t, err)
	assert.False(t, valid)
}
