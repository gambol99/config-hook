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
	dockerapi "github.com/gambol99/go-dockerclient"
	"github.com/golang/glog"
)

type ConfigHookService interface {
	/* close down any resources */
	Close()
}

const (
	DOCKER_EVENT_START   = "start"
	DOCKER_EVENT_DIE     = "die"
	DOCKER_EVENT_CREATED = "created"
	DOCKER_EVENT_DESTROY = "destroy"
	HOOK_FILE            = "FILE_"
	HOOK_CONFIG          = "CONFIG_"
	HOOK_EXEC_ONETIME    = "EXECO_"
	HOOK_EXEC            = "EXEC_"
	HOOK_REGEX           = "(FILE|CONFIG|EXECO|EXEC)_"
)

type ConfigHook struct {
	/* the agent for the k/v store */
	store store.Store
	/* the docker client */
	docker *dockerapi.Client
	/* the shutdown channel */
	shutdown_channel chan bool
}

func NewConfigHookService() (ConfigHookService, error) {
	service := new(ConfigHook)
	/* step: we need to create a store agent */
	if store, err := store.NewStore(config.Options.Store_URL); err != nil {
		glog.Errorf("Failed to create a store agent, url: %s, error: %s", config.Options.Store_URL, err)
		return nil, err
	} else {
		service.store = store
		/* step: create a docker client */
		glog.Infof("Connecting to the docker service via: %s", config.Options.Docker_Socket)
		docker, err := dockerapi.NewClient("unix://" + config.Options.Docker_Socket)
		if err != nil {
			glog.Errorf("Failed to connect to the docker service via docker, error: %s", err)
			return nil, err
		}
		service.docker = docker
		/* step: kick off the processing of events */
		if err := service.ProcessEvents(); err != nil {
			glog.Errorf("Failed to start processing events in the Hook Service, error: %s", err)
			return nil, err
		}

	}
	return service, nil
}

func (r *ConfigHook) Close() {

}

func (r *ConfigHook) ProcessEvents() error {
	/* step: we take an initial listing of all the container and process them */
	glog.Infof("Processing the container currently running")
	if containers, err := r.docker.ListContainers(dockerapi.ListContainersOptions{}); err != nil {
		glog.Errorf("Failed to processing the containers presently running, error: %s", err)
		/* CHOICE: should be continue anyhow??? */

	} else {
		/* step: iterate the containers and look for services */
		for _, container := range containers {
			go r.ProcessContainerCreation(container.ID)
		}
	}

	/* step: add ourselves as a docker event listener */
	events_channel := make(chan *dockerapi.APIEvents)
	if err := r.docker.AddEventListener(events_channel); err != nil {
		glog.Errorf("Failed to add ourselves as a docker events listen, error: %s", err)
		return err
	}

	/* step: create the go-routine and process docker events / shutdown */
	go func() {
		for {
			select {
			case event := <-events_channel:
				glog.V(4).Infof("Received docker event status: %s, id: %s", event.Status, event.ID)
				switch event.Status {
				case DOCKER_EVENT_START:
					go r.ProcessContainerCreation(event.ID)
				case DOCKER_EVENT_DESTROY, DOCKER_EVENT_DIE:
					go r.ProcessContainerDestruction(event.ID)
				}
			case <-r.shutdown_channel:
				glog.Infof("Config Hook Service recieved a shutdown signal")
				r.store.Close()
				return
			}
		}
	}()
	return nil
}

func (r *ConfigHook) ContainerEnvironment(variables []string) (map[string]string, error) {
	environment := make(map[string]string, 0)
	for _, kv := range variables {
		if found, _ := regexp.MatchString(`^(.*)=(.*)$`, kv); found {
			elements := strings.SplitN(kv, "=", 2)
			environment[elements[0]] = elements[1]
		} else {
			glog.V(3).Infof("Invalid environment variable: %s, skipping", kv)
		}
	}
	return environment, nil
}

func (r *ConfigHook) HasConfig(containerId string) (map[string]string, bool, error) {
	glog.V(5).Infof("Checking the container: %s for any config hook references", containerId)
	/* step: get the container */
	prefix := config.Options.Runtime_Prefix
	regex := fmt.Sprintf("^%s%s", prefix, HOOK_REGEX)
	glog.V(20).Infof("HasConfig() prefix: %s, regex: %s", prefix, regex)

	if container, err := r.docker.InspectContainer(containerId); err != nil {
		glog.Errorf("Failed to inspect the container: %s, error: %s", containerId, err)
		return nil, false, err
	} else {
		/* step: we get the environment variables from the container */
		if environment, err := r.ContainerEnvironment(container.Config.Env); err != nil {
			glog.Errorf("Failed to retrieve the environment variables from the container: %s, error: %s", containerId, err)
			return nil, false, err
		} else {
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
	}
}

func (r *ConfigHook) ProcessContainerCreation(containerId string) {
	glog.V(5).Infof("Processing creation of container: %s", containerId)
	if config, found, err := r.HasConfig(containerId); err != nil {
		glog.Errorf("Failed to check if the container: %s has any hooks, error: %s", containerId, err)
	} else if found {
		glog.Infof("Found config hook in container: %s, config: %v", containerId, config)

	}
}

func (r *ConfigHook) ProcessContainerDestruction(containerId string) {
	glog.V(5).Infof("Processing destruction of container: %s", containerId)

}
