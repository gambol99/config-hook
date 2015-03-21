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
	"errors"
)

var (
	ErrInvalidConfiguration = errors.New("Invalid configuration, please check config")
)

const (
	AUTHOR                 = "Rohith <gambol99@gmail.com>"
	NAME                   = "Config Hook Service"
	DEFAULT_RUNTIME_PREFIX = "HOOK_"
	DEFAULT_DOCKER_SOCKET  = "/var/run/docker.sock"
	DEFAULT_STORE_URL      = "etcd://127.0.0.1:4001"
)

type Configuration struct {
	// the prefix in environment variables
	Prefix string `json:"prefix"`
	// The docker configuration, assuming we are acting as a proxy
	Docker DockerConfig `json:"docker,omitempty"`
	// Should we be logging to syslog
	Logging *LoggingConfig `json:"logging"`
	// A map of discovery providers for services
	Discovery []*map[string]string `json:"discovery"`
	// A map of templates to
	Templates []* `json:"templates"`
}

//
// Logging config
//
type LoggingConfig struct {
	// whether or not we should log to syslog
	SyslogEnabled bool `json:"syslog_enabled"`
	// the facility for the above if enabled
	SyslogFacility string `json:"syslog_facility"`
	// the log level
	LogLevel string `json:"loglevel"`
}

//
// The configuration for docker
//
type DockerConfig struct {
	// the location of the read docker socket
	Socket string `json:"socket"`
	// the socket we shoukd create as a proxy
	SocketProxy string `json:"socket_proxy"`
	// the real tcp socket used by docker
	TCPSocket string `json:"tcp_socket"`
	// the
	TCPSocket_Proxy string `json:"tcp_socket_proxy"`
}
