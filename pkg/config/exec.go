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
	"fmt"
)

// This represents a exec
type TemplateExec struct {
	// the exec command which should be run
	Exec string `json:"onchange"`
	// the check command which should be performed before hand
	Check string `json:"oncheck"`
}

func (r TemplateExec) String() string {
	return fmt.Sprintf("command: %s, check: %s", r.Exec, r.Check)
}
