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
	"errors"
	"fmt"
)

func NewHookKeys(id string) *HookKeys {
	return &HookKeys{
		ID:    id,
		File:  "",
		Flags: "",
	}
}

type HookKeys struct {
	// the id which is associated to the keys
	ID string `json:"id"`
	// the file which holds the content
	File string `json:"file"`
	// the flags associated to the config
	Flags string `json:"flags"`
}

func (r HookKeys) String() string {
	return fmt.Sprintf("id: %s, file: %s, flags: %s", r.ID, r.File, r.Flags)
}

func (r HookKeys) Valid() (bool, error) {
	if r.ID == "" {
		return false, errors.New("the hook config does not contain a id")
	}
	if r.File == "" {
		return false, errors.New("the hook config does not contain a file")
	}
	return true, nil
}
