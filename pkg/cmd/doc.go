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

import "github.com/alecthomas/kingpin"

type Application struct {
	// the factory name
	Name string
	// the usage for the factory
	Usage string
	// the version
	Version string
}

//
// A command factory
//
type Factory struct {
	// a list of subcommands
	commands map[string]*FactoryCommand
}

//
// The structure of a command
//
type FactoryCommand struct {
	Name string
	// the alias for the command
	Alias string
	// the usage of the for the command
	Usage string
	// an example usage
	Example string
	// the setup code
	Setup func(c *kingpin.CmdClause) error
	// the method to run
	Run func(c *kingpin.CmdClause) error
	// the command struct
	cmd *kingpin.CmdClause
}

