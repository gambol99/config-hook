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
	"fmt"
	"regexp"
	"testing"

	"github.com/gambol99/config-hook/config"
	"github.com/stretchr/testify/assert"
)

func TestNewHooksConfig(t *testing.T) {
	c := NewHooksConfig()
	assert.NotNil(t, c)
	assert.NotNil(t, c.files)
	assert.NotNil(t, c.keys)
}

func TestFiles(t *testing.T) {
	c := NewHooksConfig()
	assert.NotNil(t, c)
	c.Files("test")
	assert.Equal(t, len(c.files), 1)
}

func TestKeys(t *testing.T) {
	c := NewHooksConfig()
	assert.NotNil(t, c)
	c.Keys("test")
	assert.Equal(t, len(c.keys), 1)
}

func TestFindMatches(t *testing.T) {
	hook_file_regex = regexp.MustCompile(fmt.Sprintf("^%s%s_([[:alpha:]]+)[$_]?(KEY|CHECK|EXEC|FLAGS)?",
		"CONFIG_HOOK_", HOOK_FILE))
	hook_keys_regex = regexp.MustCompile(fmt.Sprintf("^%s%s_([[:alpha:]]+)$",
		"CONFIG_HOOK_", HOOK_KEYS))

	src := "CONFIG_HOOK_FILE_NAME"
	c := NewHooksConfig()
	assert.NotNil(t, c)
	r := regexp.MustCompile("^CONFIG_HOOK_FILE_([[:alpha:]]+)[$_]?(KEY|CHECK|EXEC|FLAGS)?")
	matches, size := c.findMatches(src, r)
	assert.NotNil(t, matches)
	assert.Equal(t, size, 1)
	assert.Equal(t, matches[0], "NAME")

	src = "CONFIG_HOOK_FILE_NAME_KEY"
	matches, size = c.findMatches(src, r)
	assert.NotNil(t, matches)
	assert.Equal(t, size, 2)
	assert.Equal(t, matches[0], "NAME")
	assert.Equal(t, matches[1], "KEY")
}

func TestIsHook(t *testing.T) {
	c := NewHooksConfig()
	assert.NotNil(t, c)
	assert.Equal(t, config.Options.Runtime_Prefix, "CONFIG_HOOK_")
	assert.False(t, c.IsHook("GONFIG_HOOK"))
	assert.True(t, c.IsHook("CONFIG_HOOK_"))
}

func TestParseKey(t *testing.T) {
	c := NewHooksConfig()
	assert.NotNil(t, c)
	assert.Equal(t, config.Options.Runtime_Prefix, "CONFIG_HOOK_")
	cfg := "CONFIG_HOO"
	hook, name, element, err := c.ParseKey(cfg)
	assert.Error(t, err, "expected to be an error")

	cfg = "CONFIG_HOOK_FILE_HAPROXY"
	hook, name, element, err = c.ParseKey(cfg)
	assert.Nil(t, err, "expected not to be an error, error: "+fmt.Sprintf("%s", err))
	assert.Equal(t, hook, "FILE")
	assert.Equal(t, name, "HAPROXY")
	assert.Equal(t, element, "")

	cfg = "CONFIG_HOOK_FILE_HAPROXY_KEY"
	hook, name, element, err = c.ParseKey(cfg)
	assert.Nil(t, err, "expected not to be an error, error: "+fmt.Sprintf("%s", err))
	assert.Equal(t, element, "KEY")

	cfg = "CONFIG_HOOK_FILE_HAPROXY_CHECK"
	hook, name, element, err = c.ParseKey(cfg)
	assert.Nil(t, err, "expected not to be an error, error: "+fmt.Sprintf("%s", err))
	assert.Equal(t, element, "CHECK")
}
