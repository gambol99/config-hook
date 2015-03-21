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

// The channel used to indicated when changes has occurred on a service
type ServiceEventChannel chan Service

// The interface for a discovery service provider
type DiscoveryAgent interface {
	// Shutdown the discovery agent
	Close() error
	// Query the service for a list of endpoints
	Endpoints(service string, args...string) ([]*Endpoint, error)
	// Query for a list of services
	Services(query string, args...string) ([]*Service, error)
	// Watch a service and fire back to use when changes occur
	Watch(service *Service, updates ServiceEventChannel) error
}

// The configuration passed to the discovery agent
type AgentConfig struct {
	// the url / location of the service
	Location string
	// the configuration passed to the agent
	Config map[string]string
}

func (r AgentConfig) Get(name, dft string) string {
	if item, found := r.Config[name]; found {
		return item
	}
	return dft
}

