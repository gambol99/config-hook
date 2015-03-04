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
	"regexp"
	"strings"

	"github.com/gambol99/config-hook/config"
	"github.com/gambol99/config-hook/store"

	"github.com/golang/glog"
)

type ConfigHook interface {
	// close down any resources
	Close()
}

const (
	DOCKER_EVENT_START   = "start"
	DOCKER_EVENT_DIE     = "die"
	DOCKER_EVENT_CREATED = "created"
	DOCKER_EVENT_DESTROY = "destroy"

	HOOK_FILE         = "FILE_"
	HOOK_CONFIG       = "CONFIG_"
	HOOK_EXEC_ONETIME = "EXECO_"
	HOOK_EXEC         = "EXEC_"
	HOOK_REGEX        = "(FILE|CONFIG|EXECO|EXEC)_"
)

type ConfigHookService struct {
	// the agent for the k/v store
	store store.Store
	// the docker client
	docker DockerStore
	// the shutdown channel
	shutdown ShutdownChannel
	// channel used to receive updates from the store (etcd)
	update_channel store.NodeUpdateChannel
}

func NewConfigHook() (ConfigHook, error) {

	var err error
	service := new(ConfigHookService)
	service.update_channel = make(store.NodeUpdateChannel, 10)
	service.shutdown = make(ShutdownChannel)

	// step: we need to create a store agent
	service.store, err = store.NewStore(config.Options.Store_URL, service.update_channel)
	if err != nil {
		glog.Errorf("Failed to create a store agent, url: %s, error: %s", config.Options.Store_URL, err)
		return nil, err
	}

	// step: create the docker store
	service.docker, err = NewDockerStore()
	if err != nil {
		glog.Errorf("Failed to create a docker store agent, error: %s", err)
		return nil, err
	}

	// step: kick off the processing of events
	if err := service.processEvents(); err != nil {
		glog.Errorf("Failed to start processing events in the Hook Service, error: %s", err)
		return nil, err
	}

	return service, nil
}

func (r *ConfigHookService) Close() {

}

func (r *ConfigHookService) processEvents() error {
	// docker creation events
	dockers_creates := make(DockerEvent, 10)
	dockers_destroys := make(DockerEvent, 10)

	for {
		select {
		case id := <-dockers_creates:
			var _ = id

		case id := <-dockers_destroys:
			var _ = id

		case <-r.shutdown:

		}
	}
}

func (r *ConfigHookService) hasConfig(containerId string) (map[string]string, bool, error) {
	glog.V(5).Infof("Checking the container: %s for any config hook references", containerId)
	/* step: get the container */
	prefix := config.Options.Runtime_Prefix
	regex := fmt.Sprintf("^%s%s", prefix, HOOK_REGEX)

	// step: get the environment of the container
	environment, err := r.docker.Environment(containerId)
	if err != nil {
		glog.Errorf("Failed to inspect the container: %s, error: %s", containerId, err)
		return nil, false, err
	}

	// step: lets attempt to find config hooks
	list := make(map[string]string, 0)
	found := false
	/* step: iterate the environment vars and look for hooks */
	for name, value := range environment {
		glog.V(20).Infof("HasConfig() checking key: %s, value: %s", name, value)
		/* check: does it start with the prefix? */
		if strings.HasPrefix(name, prefix) {
			/* check: does the left over start with a hook type? */
			if matched, _ := regexp.MatchString(regex, name); matched {
				glog.V(20).Infof("HasConfig() found matching regex key: %s, value: %s", name, value)
				list[name] = value
				found = true
			}
		}
	}
	glog.V(20).Infof("HasConfig() list: %v, found: %s", list, found)
	return list, found, nil
}

func (r *ConfigHookService) processContainerCreation(containerId string) {
	glog.V(5).Infof("Processing creation of container: %s", containerId)
	if config, found, err := r.hasConfig(containerId); err != nil {
		glog.Errorf("Failed to check if the container: %s has any hooks, error: %s", containerId, err)

	} else if found {
		glog.Infof("Found config hook in container: %s, config: %v", containerId, config)

	}
}

func (r *ConfigHookService) processContainerDestruction(containerId string) {
	glog.V(5).Infof("Processing destruction of container: %s", containerId)

}
