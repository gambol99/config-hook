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

package discovery

import (
	"fmt"
	"strings"
)

// The structure for a service
type Service struct {
	// a unique id for the service
	ID string
	// a name for the service
	Name string
	// Tags / meta data associated
	Tags []string
}

func (service Service) String() string {
	return fmt.Sprintf("service: id: %s, name: %s, tags: (%s)",
		service.ID, service.Name, strings.Join(service.Tags, ","))
}
