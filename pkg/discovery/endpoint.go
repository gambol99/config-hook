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
)

// the structure of an endpoint
type Endpoint struct {
	// the ip address for the endpoint
	Address string
	// the port the service running
	Port int
}

func (endpoint Endpoint) String() string {
	return fmt.Sprintf("endpoint: [%s:%d]", endpoint.Address, endpoint.Port)
}
