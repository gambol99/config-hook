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

package config

import (
	"errors"
	"time"
)

type HookTemplate struct {
	// the source template which the config is being generated from
	Source string `json:"src"`
	// the destination where the content should be placed
	Dest string `json:"dest"`
	// the exec which should be run in order to check the container
	Exec *TemplateExec `json:"exec"`
	// the minimum amount of time to wait
	WaitMin time.Duration `json:"wait_min"`
	// the max amount of time to wait for this template
	WaitMax time.Duration `json:"wait_max"`
}

func (r *HookTemplate) Set(element string, value interface{}) {
	switch element {
	case "KEY":
		r.Key = value.(string)
	case "EXEC":
		r.Exec.Exec = value.(string)
	case "CHECK":
		r.Exec.Check = value.(string)
	case "FLAGS":
		r.Flags = value.(string)
	case "":
		r.File = value.(string)
	}
}

func (r HookTemplate) Valid() error {
	if r.ID == "" {
		return errors.New("the hook config does not contain a id")
	}
	if r.File == "" {
		return errors.New("the hook config does not contain a file")
	}
	if r.Key == "" {
		return errors.New("the hook config does not contain a key")
	}
	return nil
}
