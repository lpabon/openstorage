/*
CSI Interface for OSD
Copyright 2017 Portworx

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
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/config"
	"github.com/stretchr/testify/assert"
)

func TestClusterNodeStatus(t *testing.T) {
	capi := &clusterApi{}
	ts := httptest.NewServer(http.HandlerFunc(capi.nodeStatus))
	defer ts.Close()

	// Initialize cluster
	err := cluster.Init(config.ClusterConfig{})
	assert.NotNil(t, err)

	// no routing is needed because the server will invoke the function
	_, err = http.Get(ts.URL)
	assert.NotNil(t, err)
}
