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

package utils

import (
	"os"

	"errors"
	"io/ioutil"
	"sync/atomic"
)

type AtomicSwitch int64

func (p *AtomicSwitch) IsSwitched() bool {
	return atomic.LoadInt64((*int64)(p)) != 0
}

func (p *AtomicSwitch) SwitchOn() {
	atomic.StoreInt64((*int64)(p), 1)
}

func (p *AtomicSwitch) SwitchedOff() {
	atomic.StoreInt64((*int64)(p), 0)
}

type ShutdownChannel chan bool

// Check if the file is a unix socket
//	filename:		the file we should be checking
func IsValidSocket(filename string) (bool, error) {
	// step: check the docker exists
	mode, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	// step: is it socket?
	if mode.Mode()&os.ModeSocket == 0 {
		return false, errors.New("the file: " + filename + " is not a socket")
	}
	return true, nil
}

// Check to see if a file exists
//	filename:		the path of the file you are checking
func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, errors.New("The filename: " + filename + " does not exist")
		}
		return false, err
	}
	return true, nil
}

func Forever(method func() error, ch ShutdownChannel) {
	var kill_switch AtomicSwitch

	// wait for the kill switch
	go func() {
		<-ch
		kill_switch.SwitchOn()
	}()

	// keep calling the method
	go func() {
		for !kill_switch.IsSwitched() {
			if err := method(); err != nil {

			}
		}
	}()
}

// Read in the content of a file
// 	filename:		the path of the file to read
func ReadFile(filename string) (string, error) {
	found, err := FileExists(filename)
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New("the file: " + filename + " does not exist, please check")
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return content, err
}

// Reads in the environment variables and constructs a map
func GetEnvironment() map[string]string {
	env := make(map[string]string, 0)
	for _, name := range os.Environ() {
		env[name] = os.Getenv(name)
	}
	return env
}
