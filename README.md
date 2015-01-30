#### **Config Hook** 
-----

The Config Hook service is a agent service that works together with [config-fs](http://github.com/gambol99/config-fs). The use case is to provide a convenient method of distributing configuration files and templates via the containers themselves. For example, you have a container for haproxy which requires a templated config and is using [config-fs](http://github.com/gambol99/config-fs) to generate the content. Instead of separating the config from the service, we can place the config template *INSIDE* of the container and tag (during runtime environment variables or dockerfile) with config hook prefix. Once the container is started, the service see's it, jumps inside to grab the files and injects them into the K/V store. 

