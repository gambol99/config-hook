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

package main

import (
	"flag"

	"github.com/gambol99/config-hook/store"
	"github.com/golang/glog"
	"os"
)

const (
	DEFAULT_RUNTIME_PREFIX = "CONFIG_FS_"
	DEFAULT_DOCKER_SOCKET  = "/var/run/docker.sock"
	DEFAULT_STORE_URL      = "etcd://localhost:4001"
)

/* the configuration options for the service */
var Options struct {
	/* the docker socket file path */
	docker_socket string
	/* the runtime variable used to indicate configuration resolve */
	runtime_prefix string
	/* the url location of the store */
	store_url string
}

func init() {
	flag.StringVar(&Options.docker_socket, "docker", DEFAULT_DOCKER_SOCKET, "the path to the docker socket file")
	flag.StringVar(&Options.runtime_prefix, "prefix", DEFAULT_RUNTIME_PREFIX, "the runtime prefix read from the docker env variables to indicate configs inside")
	flag.StringVar(&Options.store_url, "store", DEFAULT_STORE_URL, "the url for the k/v store used to push configurations")
}

func main() {
	/* step: parse the command line options */
	flag.Parse()
	/* step: print the banner */
	glog.Infof("Starting the Config Hook Service, version: %s (author: %s)", VERSION, AUTHOR)

	/* step: we need to create a store agent */
	if store, err := store.NewStore(Options.store_url); err != nil {
		glog.Errorf("Failed to create a store agent, url: %s, error: %s", Options.store_url, err)
		os.Exit(1)
	} else {
		var _ = store

		/* step: create a docker client */


		/* step: we take an initial listing of all the container and process them */


		/* step: we wait for docker events and process the container */


	}
}
