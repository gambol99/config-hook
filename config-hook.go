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
	"os"
	"os/signal"
	"syscall"

	"github.com/gambol99/config-hook/hook"
	"github.com/golang/glog"
)

func main() {
	/* step: parse the command line options */
	flag.Parse()
	/* step: print the banner */
	glog.Infof("Starting the Config Hook Service, version: %s (author: %s)", VERSION, AUTHOR)
	/* step: we create the hook service and wait */
	if service, err := hook.NewConfigHookService(); err != nil {
		glog.Errorf("Failed to create the hook service, error: %s", err )
		os.Exit(1)
	} else {
		/* step: we wait for any kill signals */
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		/* step: wait on the signal */
		<-signalChannel
		/* step: shutdown the service */
		service.Close()
	}
}
