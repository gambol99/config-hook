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

func TestNewHookFile(t *testing.T) {
	config := NewHookFile("test")
	assert.NotNil(t, config, "the hook file should have been set")
	assert.NotNil(t, config.Exec, "the hook file exec should have been set")
	assert.Equal(t, config.ID, "test")
}

func TestSet(t *testing.T) {
	config := NewHookFile("test")
	config.Set("KEY", "key")
	assert.Equal(t, config.Key, "key")
	config.Set("FLAGS", "flags")
	assert.Equal(t, config.Flags, "flags")
	config.Set("EXEC", "exec")
	assert.Equal(t, config.Exec.Exec, "exec")
	config.Set("CHECK", "check")
	assert.Equal(t, config.Exec.Check, "check")
}

func TestValidate(t *testing.T) {
	config := NewHookFile("test")
	valid, err := config.Validate()
	assert.NotNil(t, err, "the error should have been raised")
	assert.Equal(t, valid, false, "the valid should have been false")
	config.File = "/usr/hello"
	config.Key = "/usr/key"
	valid, err = config.Validate()
	assert.Nil(t, err, "the error should not have been raised")
	assert.Equal(t, valid, true, "the valid flag should have been true")
}
