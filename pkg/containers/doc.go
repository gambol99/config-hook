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

package containers


type ContainerEvent chan string

// The interface to docker
type ContainerStore interface {
	// Injects a file into the container
	InjectFile(containerID, path, content string) error
	// retrieve the contents of a file in a container
	GetFile(containerID, path string) (string, error)
	// Get a listing of containers
	Containers() ([]string, error)
	// watch for docker events
	Watch(channel ContainerEvent, event_type string)
	// retrieve the environment variables for a container
	Environment(containerID string) (map[string]string, error)
	// Close down the resources
	Close()
}
