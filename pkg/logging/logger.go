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

package logging

import "github.com/gambol99/config-hook/config"

import (
	"log"
	"os"
)

type DefaultLogger struct {
	// the logger
	logger *log.Logger
}

func NewDefaultLogger(cfg *config.LoggingConfig) (Logger, error) {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", 0),
	}
	return Log, nil
}

func DefaultLogging() *config.LoggingConfig {
	return &config.LoggingConfig{
		SyslogEnabled: false,
		SyslogFacility: "local5",
		LogLevel: "info",
	}
}

func (r DefaultLogger) Info(message string, args...interface {}) {
	r.logger.Printf("[info]  %s", message, args)
}

func (r DefaultLogger) Error(message string, args...interface {}) {
	r.logger.Printf("[error] %s", message, args)
}

func (r DefaultLogger) Warn(message string, args...interface {}) {
	r.logger.Printf("[warn]  %s", message, args)
}

func (r DefaultLogger) Debug(message string, args...interface {}) {
	r.logger.Printf("[debug] %s", message, args)
}

