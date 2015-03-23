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

package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gambol99/config-hook/pkg/utils"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

const (
	ETCD_OPTION_CERT   = "cert"
	ETCD_OPTION_CACERT = "cacert"
	ETCD_OPTION_KEY    = "key"
)

type EtcdAgent struct {
	sync.RWMutex
	// the etcd http client
	client *etcd.Client
	// a list of keys we are listing to
	watches map[string]int
	// the channel to send updates upon
	updates ServiceEventChannel
	// a stop channel for the watcher
	shutdown utils.ShutdownChannel
}

// Create a new Etcd discovery agent
//
// Agent Options:
//	cert:		the etcd certificate we should be using
//  key:	 	the certificate key is would should be using for TLS
//  cacert: 	the location of the CA certificate
func NewEtcdAgent(cfg *AgentConfig) (DiscoveryAgent, error) {
	if uri, err := url.Parse(cfg.Location); err != nil {
		return nil, ErrInvalidLocation
	} else {
		client := new(EtcdAgent)
		client.watches = make(map[string]int, 0)
		client.shutdown = make(utils.ShutdownChannel)
		// create the client
		if cfg.Get(ETCD_OPTION_CERT, "") != "" {
			etcd_client, err := etcd.NewTLSClient(client.parseHosts(uri.Host), cfg.Get(ETCD_OPTION_CERT, ""),
				cfg.Get(ETCD_OPTION_KEY, ""), cfg.Get(ETCD_OPTION_CACERT, ""))
			if err != nil {
				return nil, err
			}
			client.client = etcd_client
		} else {
			client.client = etcd.NewClient(client.parseHosts(uri.Host))
		}
		return client, nil
	}
}

func (r *EtcdAgent) Close() {
	r.shutdown <- true
}

func (r *EtcdAgent) Watch(service *Service) (ServiceEventChannel, error) {

}

func (r *EtcdAgent) Services(query string, args ...string) ([]*Service, error) {

}

func (r *EtcdAgent) UnWatch(service *Service) {

}

//
//
func (r *EtcdAgent) Endpoints(si string, args ...string) ([]*Endpoint, error) {
	list := make([]Endpoint, 0)
	glog.V(5).Infof("Listing the container nodes for service: %s, path: %s", si, si.Name)

	/* step: we get a listing of all the nodes under or branch */
	paths := make([]string, 0)
	paths, err := r.Paths(string(si.ID), &paths)
	if err != nil {
		glog.Errorf("Failed to walk the paths for service: %s, error: %s", si, err)
		return nil, err
	}
	/* step: iterate the nodes and generate the services documents */
	for _, service_path := range paths {
		glog.V(5).Infof("Retrieving service document on path: %s", service_path)
		response, err := r.client.Get(service_path, false, false)
		if err != nil {
			glog.Errorf("Failed to get service document at path: %s, error: %s", service_path, err)
			continue
		}
		/* step: convert the document into a record */
		document, err := newEtcdDocument([]byte(response.Node.Value))
		if err != nil {
			glog.Errorf("Unable to convert the response to service document, error: %s", err)
			continue
		}
		list = append(list, document.toEndpoint())
	}
	return list, nil
}

func (r *EtcdAgent) processEvents() error {
	utils.Forever(func() error {
		// we wait for any change on the root key
		response, err := r.client.Watch("/services", uint64(0), true, nil, nil)
		if err != nil {
			time.Sleep(time.Duration(5) * time.Second)
			return nil
		}

		// step: check if this key is being listened to
		if r.isWatched(response.Node.Key) && r.updates != nil {

		}
	}, r.shutdown)
	return nil
}

func (r *EtcdAgent) isWatched(key string) bool {
	r.RLock()
	defer r.Unlock()
	for watched_key, _ := range r.watches {
		if strings.HasPrefix(watched_key, key) {
			return true
		}
	}
	return false
}

func (r *EtcdAgent) parseHosts(uri string, tls bool) []string {
	hosts := make([]string, 0)
	protocol := "http"
	if tls != "" {
		protocol = "https"
	}
	for _, etcd_host := range strings.Split(uri, ",") {
		if strings.HasPrefix(etcd_host, "etcd") {
			etcd_host = strings.TrimPrefix(etcd_host, "etcd")
		}
		hosts = append(hosts, fmt.Sprintf("%s://%s", protocol, etcd_host))
	}
	return hosts
}

/* --------- Service Document Decoding ------------ */

type etcdServiceDocument struct {
	IPaddress string   `json:"ipaddress"`
	HostPort  string   `json:"host_port"`
	Port      string   `json:"port"`
	Tags      []string `json:"tags"`
}

func newEtcdDocument(data []byte) (*etcdServiceDocument, error) {
	document := &etcdServiceDocument{}
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	if err := document.isValid(); err != nil {
		return nil, err
	}
	return document, nil
}

func (e etcdServiceDocument) toEndpoint() Endpoint {
	/* check: since most registration / discovery uses port rather than host_port */
	port := ""
	if e.HostPort != "" {
		port = e.HostPort
	} else {
		port = e.Port
	}
	return Endpoint(fmt.Sprintf("%s:%s", e.IPaddress, port))
}

func (e etcdServiceDocument) isValid() error {
	if e.IPaddress == "" || e.Port == "" {
		return errors.New("Invalid service document, does not contain a ipaddress and port")
	}
	return nil
}
