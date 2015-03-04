[![Build Status](https://drone.io/github.com/gambol99/config-hook/status.png)](https://drone.io/github.com/gambol99/config-hook/latest)

### **Config Hook**
-----

The Config Hook service is a agent service that works together with [config-fs](http://github.com/gambol99/config-fs). The use case is to provide a convenient method of distributing configuration files and templates via the containers themselves. For example, you have a container for haproxy which requires a templated config and is using [config-fs](http://github.com/gambol99/config-fs) to generate the content. Instead of separating the config from the service, we can place the config template *INSIDE* of the container and tag (during runtime environment variables or dockerfile) with config hook prefix. Once the container is started, the service see's it, jumps inside to grab the files and injects them into the K/V store. We can also extend upon this and perform requested actions when content is changed, i.e. call some script inside the container when a config is changed.

#### **Configuration**
---


	[jest@starfury config-hook]$ stage/config-hook --help
	Usage of stage/config-hook:
	  -docker="/var/run/docker.sock": the path to the docker socket file
	  -etcd-cacert="": the etcd ca certificate file (optional)
	  -etcd-cert="": the etcd certificate file (optional)
	  -etcd-keycert="": the etcd key certificate file (optional)
	  -prefix="CONFIG_HOOK_": the runtime prefix read from the docker env variables to indicate configs inside
	  -stderrthreshold=0: logs at or above this threshold go to stderr
	  -store="etcd://127.0.0.1:4001": the url for the k/v store used to push configurations
	  -v=0: log level for V logs
	  -vmodule=: comma-separated list of pattern=N settings for file-filtered logging

#### **Building**
----
Assuming the following GO environment
  
    [jest@starfury config-hook]$ set | grep GO
    GOPATH=/home/jest/go
    GOROOT=/usr/lib/golang
    
    cd $GOPATH && mkdir -p src/github.com/config-hook 
    cd src/github.com/gambol99 && git clone https://github.com/gambol99/config-hook.git
    cd config-hook && make

An alternative would be to build inside a golang container
  
    cd /tmp && git clone https://github.com/gambol99/config-hook.git 
    cd config-hook
    docker run --rm -v "$PWD":/go/src/github.com/gambol99/config-hook \
      -w /go/src/github.com/gambol99/config-hook -e GOOS=linux golang:1.3.3 make
    stage/config-hook --help

#### **Runtime Hooks**
---
All the runtime variables are read from the environment variables of the container. The default prefix is CONFIG_HOOK_[TYPE]. We have the following hook types

 >  * FILE:       a file / template to injected into the store
 >  * KEYS:     a file containing a series of KEY=VALUE pairs which are injected into the k/v store

#### **File Types**
**Format**: PREFIX_FILE_[NAME]=[PATH];[KEY];[EXEC];[FLAGS]

> - NAME: the name is an arbitrary identifier for the type and should be unique, additions with simply override the former
> - PATH: the path of the file | template INSIDE the container
> - KEY:  the path in the K/V store the config should be stored

**Optional**:
> - EXEC:  a command line execute when the content of PATH has changed
> - FLAGS: a comma separated list of options i.e. OT (onetime)

**Examples**:

A HAProxy example

>  CONFIG_HOOK_HAPROXY=/configs/haproxy.cfg;/env/prod/configs/haproxy.cfg;/usr/bin/ha_restart

The above statement will extract the *'/config/haproxy.cfg'* file from within the container and push into the K/V store on the key *'/env/prod/configs/haproxy.cfg'*. Whenever the content of the *'/env/prod/configs/haproxy.cfg'* file changes, the *'/usr/bin/ha_restart'* will be executed within container as a docker exec.

**Additional**

Note, if you don't like the compact format above you can spread the above sections into multiple environment variables i.e.

    HK_FILE_<NAME>=/config/haproxy.cfg
    HK_FILE_<NAME>_KEY=/env/%ENVIRONMENT%/configs/haproxy.cfg
    HK_FILE_<NAME>_EXEC=/config/haproxy.cfg
    HK_FILE_<NAME>_FLAGS=/config/haproxy.cfg

#### **Keys Types**

**Format**: PREFIX_FILE_[NAME]=[PATH];[FLAGS}

> - NAME: the name is an arbitrary identifier for the type 
> - PATH: the path of the file within the container which has the key pairs

**Optional**:

>  - FLAGS: a comma separated list of options i.e. OT (onetime)

**Content**

The contents of the keys file is simple newline separated list of KEY=VALUE 

	KEY_ONE=VALUE_ONE
	KEY_TWO=VALUE_TWO
	...

