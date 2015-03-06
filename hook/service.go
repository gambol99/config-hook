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

	"github.com/gambol99/config-hook/config"
	"github.com/gambol99/config-hook/store"

	"github.com/go-fsnotify/fsnotify"
	"github.com/golang/glog"
)

type ConfigHook interface {
	// close down any resources
	Close()
}

type ConfigHookService struct {
	// the agent for the k/v store
	store store.Store
	// the docker client
	docker DockerStore
	// the shutdown channel
	shutdown ShutdownChannel
	// channel used to receive updates from the store (etcd)
	update_channel store.NodeUpdateChannel
	// a map of containerId to config hooks
	hooks map[string]*Hooks
	// the file watcher
	inotify *fsnotify.Watcher
}

const (
	HOOK_KEYS = "KEYS"
	HOOK_FILE = "FILE"
)

var (
	hook_file_regex, hook_keys_regex   *regexp.Regexp
	hook_file_prefix, hook_keys_prefix string
)

func NewConfigHook() (ConfigHook, error) {

	var err error
	service := new(ConfigHookService)
	service.update_channel = make(store.NodeUpdateChannel, 10)
	service.hooks = make(map[string]*Hooks, 0)
	service.shutdown = make(ShutdownChannel)

	// step: set the prefixes and regexes
	hook_file_regex = regexp.MustCompile(fmt.Sprintf("^%s%s_([[:alpha:]]+)[$_]?(KEY|CHECK|EXEC|FLAGS)?",
		config.Options.Runtime_Prefix, HOOK_FILE))
	hook_keys_regex = regexp.MustCompile(fmt.Sprintf("^%s%s_([[:alpha:]]+)$",
		config.Options.Runtime_Prefix, HOOK_KEYS))
	hook_file_prefix = fmt.Sprintf("%s%s", config.Options.Runtime_Prefix, HOOK_FILE)
	hook_keys_prefix = fmt.Sprintf("%s%s", config.Options.Runtime_Prefix, HOOK_KEYS)

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

	glog.V(3).Infof("%s, runtime prefix: %s", config.NAME, config.Options.Runtime_Prefix)

	// step: preprocess any container which are already running
	if err := service.preprocessContainers(); err != nil {
		glog.Errorf("Failed to preprocess the container, error: %s", err)
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
	glog.Infof("Shutting down the %s", config.NAME)
}

func (r *ConfigHookService) preprocessContainers() error {
	glog.V(6).Infof("Preprocessing any container which are already running")
	containers, err := r.docker.List()
	if err != nil {
		return err
	}
	// step: iterate the containers and pre-process them
	for _, container := range containers {
		r.processContainerCreation(container)
	}
	return nil
}

func (r *ConfigHookService) processEvents() error {
	// docker creation events
	container_created := make(DockerEvent, 10)
	container_destroyed := make(DockerEvent, 10)
	content_changes := make(chan string, 0)

	// step: add the watch
	r.docker.Watch(container_created, DOCKER_START)
	r.docker.Watch(container_destroyed, DOCKER_DESTROY)

	go func() {
		glog.Infof("Starting the event processor for config hook service")
		for {
			select {
			// a container has been created
			case id := <-container_created:
				glog.V(6).Infof("Container: %s creation event", id)
				r.processContainerCreation(id)
			// a container has been destroyed
			case id := <-container_destroyed:
				glog.V(6).Infof("Container: %s destruction event", id)
				r.processContainerDestruction(id)
			// the contents of a file has changed
			case filename := <-content_changes:
				glog.V(6).Infof("The file: %s has changed", filename)
			// we have hit a shutdown event
			case <-r.shutdown:
				glog.Infof("Request to shutdown the service")
				r.docker.Close()
			}
		}
	}()
	return nil
}

func (r *ConfigHookService) processContainerCreation(containerId string) {
	glog.V(5).Infof("Processing creation of container: %s", containerId[:12])

	// step: check if the container has any config hooks
	hooks, has_hooks, err := r.hasConfig(containerId)
	glog.V(10).Infof("Container: %s, hooks files: %V", containerId[:12], hooks.files)
	if err != nil {
		glog.Errorf("Failed to process the container: %s, error: %s", containerId[:12], err)
		return
	}
	if !has_hooks {
		glog.V(6).Infof("The container: %s has not config hooks, skipping", containerId[:12])
		return
	}

	// step: add the hooks map
	r.hooks[containerId] = hooks

	// step: process the hook files
	for _, x := range hooks.files {

	}
}

func (r *ConfigHookService) processContainerDestruction(containerId string) {
	glog.V(5).Infof("Processing destruction of container: %s", containerId)
	// step: check if the hooks config exists for this
	if _, found := r.hooks[containerId]; found {
		// step: close up any of the resources used by this

		// step: remove from the map
		delete(r.hooks, containerId)
	}
}

func (r *ConfigHookService) hasConfig(containerId string) (*Hooks, bool, error) {
	glog.V(6).Infof("Checking the container: %s for any config hook references", containerId)
	// step: get the container

	// step: get the environment of the container
	environment, err := r.docker.Environment(containerId)
	if err != nil {
		glog.Errorf("Failed to inspect the container: %s, error: %s", containerId, err)
		return nil, false, err
	}

	// step: lets attempt to find config hooks
	hooks := NewHooksConfig()

	// step: iterate the environment vars and look for hooks
	for key, value := range environment {
		if hooks.IsHook(key) {
			if hook, name, element, err := hooks.ParseKey(key); err == nil {
				switch hook {
				case HOOK_FILE:
					hooks.Files(name).Set(element, value)
				case HOOK_KEYS:
					hooks.Keys(name).File = value
				}
			}
		}
	}

	// step: we need to validate the hooks and remove anything which does satisfy
	if err := hooks.Validate(); err != nil {
		glog.Errorf("One or more configs had errors in container: %s, error: %s", containerId, err)
	}

	return hooks, hooks.HasHooks(), nil
}
