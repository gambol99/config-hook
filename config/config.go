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

package config

import (
	"flag"
)

const (
	DEFAULT_RUNTIME_PREFIX = "_HOOK_"
	DEFAULT_DOCKER_SOCKET  = "/var/run/docker.sock"
	DEFAULT_STORE_URL      = "etcd://localhost:4001"
)

/* the configuration options for the service */
type ConfigHookOptions struct {
	/* the docker socket file path */
	Docker_Socket string
	/* the runtime variable used to indicate configuration resolve */
	Runtime_Prefix string
	/* the url location of the store */
	Store_URL string
}

var Options ConfigHookOptions

func init() {
	flag.StringVar(&Options.Docker_Socket, "docker", DEFAULT_DOCKER_SOCKET, "the path to the docker socket file")
	flag.StringVar(&Options.Runtime_Prefix, "prefix", DEFAULT_RUNTIME_PREFIX, "the runtime prefix read from the docker env variables to indicate configs inside")
	flag.StringVar(&Options.Store_URL, "store", DEFAULT_STORE_URL, "the url for the k/v store used to push configurations")
}
