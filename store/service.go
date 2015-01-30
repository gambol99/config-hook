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

package store

import (
	"errors"
)

const VERBOSE_LEVEL = 5

type Store interface {
	/* retrieve a key from the store */
	Get(key string) (*Node, error)
	/* List all the keys under a path */
	Paths(path string, paths *[]string) ([]string, error)
	/* Get a list of all the nodes under the path */
	List(path string) ([]*Node, error)
	/* set a key in the store */
	Set(key string, value string) error
	/* delete a key from the store */
	Delete(key string) error
	/* recursively delete a path */
	RemovePath(path string) error
	/* release all the resources */
	Close()
}

var (
	InvalidUrlErr       = errors.New("Invalid URI error, please check backend url")
	InvalidDirectoryErr = errors.New("Invalid directory specified")
)

func NewStore(location string) (Store, error) {



	return nil, nil
}
