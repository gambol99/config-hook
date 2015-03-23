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

package cmd

import (
	"github.com/alecthomas/kingpin"
)

//
// Create a new command line factory
//	name:		the name of the application
// 	usage:     	the application usage
func NewFactory(application *Application) *Factory {
	kingpin.Version(application.Version)
	factory := &Factory{commands: make(map[string]*FactoryCommand, 0)}
	return factory
}

// Parse the command line options, performing the setup on the subcommands
// 	args:		the arguments to parse
func (r *Factory) Parse(args []string) error {
	// step: we call the setup any of the sub commands
	for _, cmd := range r.commands {
		if err := cmd.Setup(cmd.cmd); err != nil {
			return err
		}
	}
	name := kingpin.Parse()

	// step: we call the implementation method
	if method, found := r.commands[name]; found {
		if err := method.Run(method.cmd); err != nil {
			return err
		}
	}
	return nil
}

// Add a flag to the default parser
//	name:		the name of the flag option
//	help:		the usage / description for the option
func (r *Factory) Flag(name, help string) *kingpin.FlagClause {
	return kingpin.Flag(name, help)
}

// Add a subcommand to the default parser
// 	name:		the command we are adding
func (r *Factory) Command(command *FactoryCommand) {
	command.cmd = kingpin.Command(command.Name, command.Usage)
	r.commands[command.Name] = command
}
