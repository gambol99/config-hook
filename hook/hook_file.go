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
	"time"

	"github.com/golang/glog"
)

func NewHookFile(id string) *HookFile {
	h := new(HookFile)
	h.ID = id
	h.Exec = new(HookExec)
	return h
}

type HookFile struct {
	// the id which is associated to the config
	ID string `json:"id"`
	// the file which holds the content
	File string `json:"file"`
	// the key which this file should be stored at
	Key string `json:"key"`
	// the exec which should be run when content changed
	Exec *HookExec `json:"exec"`
	// the flags associated to the config
	Flags string `json:"flags"`
}

func (r HookFile) String() string {
	return fmt.Sprintf("id: %s, file: %s, key: %s, exec: (%s), flags: %s",
		r.ID, r.File, r.Key, r.Exec, r.Flags)
}

type HookExec struct {
	// the last time the exec was ran
	LastRun time.Time
	// the last exit code
	LastExitCode int
	// the exec command which should be run
	Exec string `json:"command"`
	// the check command which should be performed before hand
	Check string `json:"check"`
}

func (r HookExec) String() string {
	return fmt.Sprintf("command: %s, check: %s", r.Exec, r.Check)
}

func (r *HookFile) Set(element string, value interface{}) {
	glog.V(10).Infof("Adding the element: %s, value: %s", element, value)
	switch element {
	case "KEY":
		r.Key = value.(string)
	case "EXEC":
		r.Exec.Exec = value.(string)
	case "CHECK":
		r.Exec.Check = value.(string)
	case "FLAGS":
		r.Flags = value.(string)
	}
}

func (r HookFile) Validate() (bool, error) {
	if r.ID == "" {
		return false, errors.New("the hook config does not contain a id")
	}
	if r.File == "" {
		return false, errors.New("the hook config does not contain a file")
	}
	if r.Key == "" {
		return false, errors.New("the hook config does not contain a key")
	}
	return true, nil
}
