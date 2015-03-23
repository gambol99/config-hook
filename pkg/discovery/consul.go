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
	"net/url"
	"sync"
	"time"

	"github.com/gambol99/config-hook/pkg/utils"

	consulapi "github.com/hashicorp/consul/api"
)

// Represent a service being watched
type ServiceWatcher struct {
	// the service being watched
	service *Service
	// the stop channel for the routine
	stop utils.ShutdownChannel
	// the number of listeners for this service
	listeners int
}

// The consul discovery agent
type ConsulAgent struct {
	sync.RWMutex
	// the consul api client
	client *consulapi.Client
	// the watch list
	watches map[string]*ServiceWatcher
}

func NewConsulAgent(cfg *AgentConfig) (DiscoveryAgent, error) {
	if uri, err := url.Parse(cfg.Location); err != nil {
		return nil, err
	} else {
		// step: create the consul client
		cfg := consulapi.DefaultConfig()
		cfg.Address = uri.Host

		client, err := consulapi.NewClient(cfg)
		if err != nil {
			return nil, err
		}

		return &ConsulAgent{
			client:  client,
			watches: make(map[string]*ServiceWatcher, 0)}, nil
	}
}

func (r *ConsulAgent) Close() {
	for _, watcher := range r.watches {
		watcher.stop <- true
	}
}

// Create a watcher for a service
// 	service:	the service up wish to keep an eye on
func (r *ConsulAgent) Watch(si *Service, updates ServiceEventChannel) error {
	// step: check if we have a watcher on this service already
	if !r.addListener(si) {
		var err error
		// step: we need to create a watcher on the service
		watcher := new(ServiceWatcher)
		watcher.listeners = 1
		watcher.service = si
		if watcher.stop, err = r.watchService(si, updates); err != nil {
			return err
		}
	}
	return nil
}

// Remove a watch on a service
// 	service:		the service you wish to remove the watcher on
func (r *ConsulAgent) UnWatch(service *Service) {
	if found := r.isWatched(service); found {
		r.removeListener(service)
	}
}

//
// Retrieve a list of the endpoint of a specific service
//	query:			the name of the service you are search for endpoints for
//	args:			a variable list of options passed to the provider
func (r *ConsulAgent) Endpoints(query string, args ...string) ([]*Endpoint, error) {
	health_checks_required := true

	// step: query for the service, along with health checks
	services, _, err := r.client.Health().Service(query, "", health_checks_required, &consulapi.QueryOptions{})
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
//	query:			the search criteria for services
//	args:			a variable list of arguments to added to the search
func (r *ConsulAgent) Services(query string, args ...string) ([]*Service, error) {
	list := make([]*Service, 0)

	return list, nil
}

// Checks to see if a service is already being watched
func (r *ConsulAgent) isWatched(si *Service) bool {
	r.RLock()
	defer r.RUnlock()
	_, found := r.watches[si.ID]
	return found
}

func (r *ConsulAgent) addListener(si *Service) bool {
	r.Lock()
	defer r.Unlock()
	watch, found := r.watches[si.ID]
	if found {
		watch.listeners += 1
	}
	return found
}

func (r *ConsulAgent) removeListener(si *Service) {
	r.Lock()
	defer r.Unlock()
	watch, found := r.watches[si.ID]
	if found {
		watch.listeners -= 1
		if watch.listeners <= 0 {
			go func() {
				watch.stop <- true
			}()
			delete(r.watches, si.ID)
		}
	}
}

func (r *ConsulAgent) watchService(si *Service, updates ServiceEventChannel) utils.ShutdownChannel {
	// our stop channel
	stop_channel := make(utils.ShutdownChannel)
	catalog := r.client.Catalog()
	wait_index := 0

	// run forever
	utils.Forever(func() error {
		if wait_index == 0 {
			// step: get the wait index for the service
			_, meta, err := catalog.Service(si.Name, "", &consulapi.QueryOptions{})
			if err != nil {
				time.Sleep(5 * time.Second)
				return nil
			} else {
				// update the wait index for this service
				wait_index = meta.LastIndex
			}
		}

		// step: build the query - make sure we have a timeout
		queryOptions := &consulapi.QueryOptions{
			WaitIndex: wait_index,
			WaitTime:  time.Duration(60) * time.Second}

		// step: making a blocking watch call for changes on the service
		_, meta, err := catalog.Service(si.Name, "", queryOptions)
		if err != nil {
			wait_index = 0
			time.Sleep(5 * time.Second)
			return nil
		}

		// step: if the wait and last index are the same, we can continue
		if r.wait_index == meta.LastIndex {
			return nil
		}

		// step: update the index
		wait_index = meta.LastIndex

		// step: construct the change event and send
		go func() {
			updates <- *si
		}()

	}, stop_channel)

	return stop_channel
}
