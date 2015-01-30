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
	"flag"
	"fmt"
	"net/url"
	"strings"
	"sync"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type EtcdStoreClient struct {
	/* a lock for the watcher map */
	sync.RWMutex
	/* a list of etcd hosts */
	hosts []string
	/* the etcd client - under the hood is http client which should be pooled i believe */
	client *etcd.Client
	/* stop channel for the client */
	stop_channel chan bool
}

/* etcd options for TLS support */
var EtcdOptions struct {
	cert_file, key_file, cacert_file string
}

func init() {
	flag.StringVar(&EtcdOptions.cert_file, "etcd-cert", "", "the etcd certificate file (optional)")
	flag.StringVar(&EtcdOptions.key_file, "etcd-keycert", "", "the etcd key certificate file (optional)")
	flag.StringVar(&EtcdOptions.cacert_file, "etcd-cacert", "", "the etcd ca certificate file (optional)")
}

const (
	ETCD_PREFIX = "etcd://"
)

func NewEtcdStoreClient(location *url.URL) (Store, error) {
	/* step: create the client */
	store := new(EtcdStoreClient)
	store.hosts = store.ParseHostsURL(location)
	store.stop_channel = make(chan bool)

	glog.Infof("Creating a Etcd Agent for K/V Store, hosts: %s", store.hosts)

	/* step: create the etcd client */
	if EtcdOptions.cacert_file != "" {
		client, err := etcd.NewTLSClient(store.hosts, EtcdOptions.cert_file, EtcdOptions.key_file, EtcdOptions.cacert_file)
		if err != nil {
			glog.Errorf("Failed to create a TLS connection to etcd: %s, error: %s", *location, err)
			return nil, err
		}
		store.client = client
	} else {
		store.client = etcd.NewClient(store.hosts)
	}
	/* step: are we using tls or not? */
	return store, nil
}

func (r EtcdStoreClient) ParseHostsURL(location *url.URL) []string {
	hosts := make([]string, 0)
	/* step: determine the protocol */
	protocol := "http"
	if EtcdOptions.cacert_file != "" {
		protocol = "https"
	}
	for _, host := range strings.Split(location.Host, ",") {
		hosts = append(hosts, fmt.Sprintf("%s://%s", protocol, host))
	}
	return hosts
}

func (r *EtcdStoreClient) Close() {
	glog.Infof("Shutting down the etcd client")
	r.stop_channel <- true
}

func (r *EtcdStoreClient) ValidateKey(key string) string {
	/* step: if it doesnt start with a / - add it */
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	/* step: if it ends with a slash, remove it */
	if len(key) > 1 && strings.HasSuffix(key, "/") {
		key = key[:len(key)-1]
	}
	return key
}

func (r *EtcdStoreClient) Get(key string) (*Node, error) {
	lookup := r.ValidateKey(key)
	/* step: lets check the cache */
	if response, err := r.GetRaw(lookup); err != nil {
		glog.Errorf("Failed to get the key: %s, error: %s", lookup, err)
		return nil, err
	} else {
		return r.CreateNode(response.Node), nil
	}
}

func (r *EtcdStoreClient) GetRaw(key string) (*etcd.Response, error) {
	glog.V(VERBOSE_LEVEL).Infof("GetRaw() key: %s", key)
	response, err := r.client.Get(key, false, true)
	if err != nil {
		glog.Errorf("Failed to get the key: %s, error: %s", key, err)
		return nil, err
	}
	return response, nil
}

func (r *EtcdStoreClient) Set(key string, value string) error {
	glog.V(VERBOSE_LEVEL).Infof("Set() key: %s, value: %s", key, value)
	_, err := r.client.Set(key, value, uint64(0))
	if err != nil {
		glog.Errorf("Failed to set the key: %s, error: %s", key, err)
		return err
	}
	return nil
}

func (r *EtcdStoreClient) Delete(key string) error {
	glog.V(VERBOSE_LEVEL).Infof("Delete() deleting the key: %s", key)
	if _, err := r.client.Delete(key, false); err != nil {
		glog.Errorf("Delete() failed to delete key: %s, error: %s", key, err)
		return err
	}
	return nil
}

func (r *EtcdStoreClient) RemovePath(path string) error {
	glog.V(VERBOSE_LEVEL).Infof("RemovePath() deleting the path: %s", path)
	if _, err := r.client.Delete(path, true); err != nil {
		glog.Errorf("RemovePath() failed to delete key: %s, error: %s", path, err)
		return err
	}
	return nil
}

func (r *EtcdStoreClient) List(path string) ([]*Node, error) {
	key := r.ValidateKey(path)
	glog.V(VERBOSE_LEVEL).Infof("List() path: %s", key)
	if response, err := r.GetRaw(path); err != nil {
		glog.Errorf("List() failed to get path: %s, error: %s", key, err)
		return nil, err
	} else {
		list := make([]*Node, 0)
		if response.Node.Dir == false {
			glog.Errorf("List() path: %s is not a directory node", key)
			return nil, InvalidDirectoryErr
		}
		for _, item := range response.Node.Nodes {
			list = append(list, r.CreateNode(item))
		}
		return list, nil
	}
}

func (e *EtcdStoreClient) Paths(path string, paths *[]string) ([]string, error) {
	response, err := e.client.Get(path, false, true)
	if err != nil {
		return nil, errors.New("Unable to complete walking the tree" + err.Error())
	}
	for _, node := range response.Node.Nodes {
		if node.Dir {
			e.Paths(node.Key, paths)
		} else {
			glog.Infof("Found service container: %s appending now", node.Key)
			*paths = append(*paths, node.Key)
		}
	}
	return *paths, nil
}

func (r *EtcdStoreClient) CreateNode(response *etcd.Node) *Node {
	node := &Node{}
	node.Path = response.Key
	if response.Dir == false {
		node.Directory = false
		node.Value = response.Value
	} else {
		node.Directory = true
	}
	return node
}
