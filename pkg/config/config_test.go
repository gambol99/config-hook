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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	CONFIG = `
	{
		"prefix": "CONFIG_HOOK_",
		"docker": {
			"socket": "/var/run/docker.sock",
			"socket_proxy": "/var/run/docker.socket",
			"tcp_socket": "127.0.0.1:2739",
			"tcp_socket_proxy": "127.0.0.1:4900"
		},
		"logging": {
			"syslog_enabled": "false",
			"syslog_facility": "local5",
			"loglevel": "info"
		}
		"wait_min": "30s",
		"wait_max": "2m",
		"discovery": [
			{
				"name": "consul",
				"url": "consul.consul"
				"token": "28197w98dhw89ed79s8fsd",
				"ssl_enabled": "true",
				"ssl_verify": "true"
			},
			{
				"name": "marathon",
				"url": "http://marathon:8080,marathon1:8080",
				"username": "api",
				"password": "pass1"
			}
		],
		"templates": [
			{
				"src": "/etc/haproxy/haproxy.cfg.tmpl",
				"dest": "/etc/haproxy/haproxy.cfg",
				"wait_max": "2m",
				"exec": {
					"oncheck": "/usr/bin/haprory -V {{filename}}",
					"onchange": "/sbin/service haproxy reload",
				}
			}
		}
	}`
)
