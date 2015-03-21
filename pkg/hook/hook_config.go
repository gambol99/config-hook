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
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gambol99/config-hook/config"

	"github.com/golang/glog"
)

func NewHooksConfig() *Hooks {
	return &Hooks{
		files: make(map[string]*HookFile, 0),
	}
}

type HookParser interface {
	// parser the hook key
	ParseKey(key string) (string, string, error)
	// is a config hook
	IsHook(key string) bool
	// check if the config has any hooks
	HasHooks() bool
	// retrieve the files
	Files(name string) *HookFile
	// validate the hooks
	Validate() error
}

type Hooks struct {
	HookParser
	// map of all the hook files
	files map[string]*HookFile
}

func (r Hooks) IsHook(key string) bool {
	return strings.HasPrefix(key, config.Options.Runtime_Prefix)
}

// Parses the key and extracts the type, the name and the element if has it
// 	key: 	the config hook key
func (r Hooks) ParseKey(key string) (string, string, string, error) {
	if strings.HasPrefix(key, hook_file_prefix) {
		matches, size := r.findMatches(key, hook_file_regex)
		if size < 1 {
			return HOOK_FILE, "", "", errors.New("Invalid config hook key: " + key + " does not match expectation")
		}
		if size == 1 {
			return HOOK_FILE, matches[0], "", nil
		}
		return HOOK_FILE, matches[0], matches[1], nil
	}

	if strings.HasPrefix(key, hook_keys_prefix) {
		matches, size := r.findMatches(key, hook_keys_regex)
		if size != 1 {
			return HOOK_KEYS, "", "", errors.New("Invalid config key for keys config: " + key)
		}
		return HOOK_KEYS, matches[0], "", nil
	}

	return "", "", "", nil
}

func (r *Hooks) Files(id string) *HookFile {
	file, found := r.files[id]
	if !found {
		file = NewHookFile(id)
		r.files[id] = file
	}
	return file
}

func (r Hooks) HasHooks() bool {
	if len(r.files) > 0 || len(r.keys) > 0 {
		return true
	}
	return false
}

func (r Hooks) Validate() error {
	// step: validate the hook files
	for id, x := range r.files {
		if err := r.Valid(); err != nil {
			glog.Errorf("invalid hook file config, error: %s", err)
			delete(r.files, id)
		}
	}
	// step: validate the keys
	for id, x := range r.keys {
		if err := r.Valid(); err != nil {
			glog.Errorf("invalid hook keys config, error: %s", err)
			delete(r.keys, id)
		}
	}
	return nil
}

func (r *Hooks) findMatches(src string, reg *regexp.Regexp) ([]string, int) {
	list := make([]string, 0)
	for index, element := range (*reg).FindStringSubmatch(src) {
		if index == 0 {
			continue
		}
		if element != "" {
			list = append(list, element)
		}
	}
	return list, len(list)
}

func (r Hooks) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Hooks config:\n")
	for _, file := range r.files {
		buffer.WriteString(fmt.Sprintf("%s\n", file))
	}
	for _, key := range r.keys {
		buffer.WriteString(fmt.Sprintf("%s\n", key))
	}
	return buffer.String()
}
