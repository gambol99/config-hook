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
	"time"
	"fmt"
	"net/url"

	consulapi "github.com/armon/consul-api"
)

func NewConsulAgent(cfg *AgentConfig) (DiscoveryAgent, error) {
	cfg := consulapi.DefaultConfig()
	uri, err := url.Parse(cfg.Location)
	if err != nil {
		return nil, err
	}
	cfg.Address = uri.Host
	client, err := consulapi.NewClient(cfg)
	return &ConsulAgent{client, uint64(0), false}, nil
}

// The consul discovery agent
type ConsulAgent struct {
	// the consul api client
	client *consulapi.Client
	// the current wait index
	wait_index uint64
	// the kill off
	kill_off bool
}

//
//	watch for changes in the consul backend service - note, this probably isn't the best way of
//	doing it, though i've not spent much time looking at the api
//
func (r *ConsulAgent) Watch(si *Service) (EndpointEventChannel, error) {
	// channel to send back events to the endpoints store
	endpointUpdateChannel := make(EndpointEventChannel, 5)
	go func() {
		catalog := r.client.Catalog()
		for {
			if r.kill_off {
				break
			}
			if r.wait_index == 0 {
				/* step: get the wait index for the service */
				_, meta, err := catalog.Service(si.Name, "", &consulapi.QueryOptions{})
				if err != nil {
					time.Sleep(5 * time.Second)
				} else {
					/* update the wait index for this service */
					r.wait_index = meta.LastIndex
				}
			}
			/* step: build the query - make sure we have a timeout */
			queryOptions := &consulapi.QueryOptions{
			WaitIndex: r.wait_index,
			WaitTime:  DEFAULT_WAIT_TIME}

			/* step: making a blocking watch call for changes on the service */
			_, meta, err := catalog.Service(si.Name, "", queryOptions)
			if err != nil {
				r.wait_index = 0
				time.Sleep(5 * time.Second)
			} else {
				if r.kill_off {
					continue
				}
				// step: if the wait and last index are the same, we can continue
				if r.wait_index == meta.LastIndex {
					continue
				}
				// step: update the index
				r.wait_index = meta.LastIndex

				// step: construct the change event and send
				event := EndpointEvent{
					ID:     si.Name,
					Action: ENDPOINT_CHANGED,
				}
				endpointUpdateChannel <- event
			}
		}
		close(endpointUpdateChannel)
	}()
	return endpointUpdateChannel, nil
}

//
// Retrieve a list of the endpoint of a specific service
//
func (r *ConsulAgent) Endpoints(service string, args...string) ([]*Endpoint, error) {
	health_checks_required := true

	// step: query for the service, along with health checks
	services, _, err := r.client.Health().Service(service, "", health_checks_required, &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	list := make([]Endpoint, 0)
	for _, service := range services {
		service_address := service.Node.Address
		service_port := service.Service.Port
		endpoint := Endpoint(fmt.Sprintf("%s:%d", service_address, service_port))
		list = append(list, endpoint)
	}
	return list, nil
}

//
// Retrieve a list of the services available in the provider
//
func (r *ConsulAgent) Services(query string, args...string) ([]*Service, error) {
	list := make([]*Service, 0)


	return list, nil
}

func (r *ConsulAgent) Close() {
	r.kill_off = true
}
