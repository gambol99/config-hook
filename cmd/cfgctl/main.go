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
	"fmt"
	"os"

	"github.com/gambol99/config-hook/pkg/cmd"

	"github.com/alecthomas/kingpin"
)

func main() {
	factory := cmd.NewFactory(&cmd.Application{
		Name: "cfgctl",
		Usage: "",
		Version: "0.0.0",
	})

	factory.Flag("config", "the path to a configuration file").Short('c').File()
	factory.Flag("test", "a test flag").String()

	factory.Command(&cmd.FactoryCommand{
		Name: "render",
		Usage: "usage for the render",
		Setup: func(c *kingpin.CmdClause) error {
			c.Flag("template", "add a template to the required resources").Strings()
			c.Flag("store", "a store provider to add the resources").Short('s').Strings()
			return nil
		},
		Run: func(c *kingpin.CmdClause) error {
			fmt.Println("Running the command")
			return nil
		},
	})
	factory.Parse(os.Args)
}
