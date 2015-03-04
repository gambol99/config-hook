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
)

const (
	HOOK_TYPE_FILE   = 0
	HOOK_TYPE_CONFIG = 1
	EXEC_FOREVER     = 0
	EXEC_ONETIME     = 1
)

type HookConfig struct {
	// the type of hook
	hook int
	// the path to the file
	file_path string
	// the action to perform if any on a change of content
	action *HookExec
}

func (r HookConfig) String() string {
	hook_type := "file"
	if r.hook == HOOK_TYPE_CONFIG {
		hook_type = "config"
	}
	return fmt.Sprintf("hook: type: %s, resource: %s, action: %s", hook_type, r.file_path, *r.action)
}

type HookExec struct {
	/* the command to execute on the change of a config */
	execute_path string
	/* the command to execute to validate the config, if any */
	check_path string
	/* the type of exec - onetime or forever */
	exec_type int
}

func (r HookExec) String() string {
	exec_type := "one-time"
	if r.exec_type == EXEC_FOREVER {
		exec_type = "forever"
	}
	return fmt.Sprintf("exec: execute: %s, check: %s, type: %s", r.execute_path, r.check_path, exec_type)
}
