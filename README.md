[![Build Status](https://drone.io/github.com/gambol99/config-hook/status.png)](https://drone.io/github.com/gambol99/config-hook/latest)

#### **Config Hook**
-----

The Config Hook service is a agent service that works together with [config-fs](http://github.com/gambol99/config-fs). The use case is to provide a convenient method of distributing configuration files and templates via the containers themselves. For example, you have a container for haproxy which requires a templated config and is using [config-fs](http://github.com/gambol99/config-fs) to generate the content. Instead of separating the config from the service, we can place the config template *INSIDE* of the container and tag (during runtime environment variables or dockerfile) with config hook prefix. Once the container is started, the service see's it, jumps inside to grab the files and injects them into the K/V store. 

##### **Runtime Hooks**

All the runtime variables are read from the environment variables of the container. The default prefix is CONFIG_HOOK_[ATTR]. To expose one or more files from the container you can;

	CONFIG_HOOK_FILE_<NAME>=<path in container>
	e.g.
	CONFIG_HOOK_FILE_HAPROXY=/etc/haproxy/haproxy.cfg
	CONFIG_HOOK_FILE_APP=/app/config/database.yml

Note: you can also get the config-hook service to jump inside and perform an exec. Thus, you can expose a config file, the hook will wait for the content to be generated and then issue an exec inside the running container. Example;

	CONFIG_HOOK_FILE_HAPROXY=/etc/haproxy/haproxy.cfg;/env/%ENV/config/haproxy.cfg
	CONFIG_HOOK_EXEC_ONETIME_HAPROXY=/usr/local/bin/reload-haproxy

##### **Example**

You have a haproxy container which is using config-fs to generate the configuration file. Lets assuming the proxy container is mapping the file /config/env/prod/configs/haproxy.cfg into the /etc/haproxy/haproxy.cfg file. The file which is generating this can be found in the container at /config/haproxy.cfg

	CONFIG_HOOK_FILE_HAPROXY=/etc/haproxy/haproxy.cfg;/env/%ENV/config/haproxy.cfg
	CONFIG_HOOK_EXEC-ONETIME_HAPROXY=/usr/local/bin/reload-haproxy

