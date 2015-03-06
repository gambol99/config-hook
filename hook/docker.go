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
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/gambol99/config-hook/config"

	dockerapi "github.com/gambol99/go-dockerclient"
	"github.com/golang/glog"
)

type DockerEvent chan string

const (
	DOCKER_START   = "start"
	DOCKER_DIE     = "die"
	DOCKER_CREATED = "created"
	DOCKER_DESTROY = "destroy"
)

// The interface to docker
type DockerStore interface {
	// retrieve the contents of a file in a container
	GetFile(containerID, filename string) (string, error)
	// Get a listing of containers
	List() ([]string, error)
	// watch for docker events
	Watch(channel DockerEvent, event_type string)
	// retrieve the environment variables for a container
	Environment(containerID string) (map[string]string, error)
	// Close down the resources
	Close()
}

// The implementation of the above
type DockerService struct {
	sync.RWMutex
	sync.Once
	// the docker client
	client *dockerapi.Client
	// the channel WE receive docker events on
	updates chan *dockerapi.APIEvents
	// a slice of those listening to creation events
	listeners map[string][]DockerEvent
	// the shutdown channel
	shutdown ShutdownChannel
}

func NewDockerStore() (DockerStore, error) {
	var err error
	// step: we have to validate the docker socket
	if valid, err := isValidSocket(config.Options.Docker_Socket); err != nil {
		glog.Errorf("Unable to validate the docker socket, error: %s", err)
		return nil, err
	} else if !valid {
		return nil, errors.New("invalid docker socket, please check")
	}
	service := new(DockerService)
	service.listeners = make(map[string][]DockerEvent, 0)
	service.shutdown = make(ShutdownChannel)
	// step: create the docker socket
	service.client, err = dockerapi.NewClient("unix://" + config.Options.Docker_Socket)
	if err != nil {
		glog.Errorf("Failed to connect to the docker service via docker, error: %s", err)
		return nil, err
	}

	return service, nil
}

func (r *DockerService) Close() {
	glog.Infof("Shutting down the docker store")
	r.shutdown <- true
}

// Retrieve a listing of the containers
func (r *DockerService) List() ([]string, error) {
	list := make([]string, 0)
	containers, err := r.client.ListContainers(dockerapi.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	// iterate the containers
	for _, container := range containers {
		list = append(list, container.ID)
	}
	return list, nil
}

func (r *DockerService) Watch(channel DockerEvent, event string) {
	r.Lock()
	defer r.Unlock()
	if r.listeners[event] == nil {
		r.listeners[event] = make([]DockerEvent, 0)
	}
	r.listeners[event] = append(r.listeners[event], channel)
	r.Once.Do(func() {
		if err := r.processEvents(); err != nil {
			glog.Fatalf("Failed to start the docker event processor, error: %s", err)
		}
	})
}

func (r *DockerService) HasFile(containerID, filename string) (bool, error) {

}

func (r *DockerService) GetFile(containerID, filename string) (string, error) {
	var buffer bytes.Buffer
	// step: construct the options
	var options dockerapi.CopyFromContainerOptions
	options.OutputStream = bufio.NewWriter(&buffer)
	options.Container = containerID
	options.Resource = filename

	// step: extract the file
	err := r.client.CopyFromContainer(options)
	if err != nil {
		glog.Errorf("Failed to copy the file: %s from the container: %s, error: %s",
			filename, containerID[:12], err)
		return "", err
	}
	return buffer.String(), nil
}

func (r *DockerService) Environment(containerId string) (map[string]string, error) {
	c, err := r.client.InspectContainer(containerId)
	if err != nil {
		return nil, err
	}
	variables := c.Config.Env
	environment := make(map[string]string, 0)
	for _, kv := range variables {
		if found, _ := regexp.MatchString(`^(.*)=(.*)$`, kv); found {
			elements := strings.SplitN(kv, "=", 2)
			environment[elements[0]] = elements[1]
		}
	}
	return environment, nil
}

func (r *DockerService) processEvents() error {
	// step: add the docker events
	updates := make(chan *dockerapi.APIEvents, 5)
	if err := r.client.AddEventListener(updates); err != nil {
		glog.Errorf("Failed to add ourselves as a docker events listen, error: %s", err)
		return err
	}

	// step: start the routine
	go func() {
		glog.Infof("Starting the docker event processor")
		for {
			select {
			case event := <-updates:
				glog.V(10).Infof("Recieved a docker event, id: %s, status: %s", event.ID[:12], event.Status)
				if listeners, found := r.listeners[event.Status]; found {
					for _, listener := range listeners {
						// DON'T BLOCK ME DUDE !!
						go func() {
							listener <- event.ID
						}()
					}
				}
			case <-r.shutdown:
				r.client.RemoveEventListener(updates)
			}
		}
	}()
	return nil
}
