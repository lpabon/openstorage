/*
Package grpc is gRPC server implementation to OpenStorage
Copyright 2018 Portworx

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
package grpc

import (
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/spec"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/volume"
	volumedrivers "github.com/libopenstorage/openstorage/volume/drivers"
)

// ServerConfig provides the configuration to the
// the gRPC CSI server created by NewServer()
type ServerConfig struct {
	Net        string
	Address    string
	DriverName string
	Cluster    cluster.Cluster
}

// Server is a OSD CSI compliant server which
// proxies CSI requests for a single specific driver
type Server struct {
	listener    net.Listener
	server      *grpc.Server
	driver      volume.VolumeDriver
	cluster     cluster.Cluster
	wg          sync.WaitGroup
	running     bool
	lock        sync.Mutex
	specHandler spec.SpecHandler
}

// NewServer creates a gRPC CSI complient server on the
// specified port and transport.
func NewServer(config *ServerConfig) (*Server, error) {
	if nil == config {
		return nil, fmt.Errorf("Configuration must be provided")
	}
	if len(config.Address) == 0 {
		return nil, fmt.Errorf("Address must be provided")
	}
	if len(config.Net) == 0 {
		return nil, fmt.Errorf("Net must be provided")
	}
	if len(config.DriverName) == 0 {
		return nil, fmt.Errorf("OpenStorage Driver name must be provided")
	}

	// Save the driver for future calls
	d, err := volumedrivers.Get(config.DriverName)
	if err != nil {
		return nil, fmt.Errorf("Unable to get driver %s info: %s", config.DriverName, err.Error())
	}

	l, err := net.Listen(config.Net, config.Address)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup server: %s", err.Error())
	}

	return &Server{
		listener:    l,
		driver:      d,
		cluster:     config.Cluster,
		specHandler: spec.NewSpecHandler(),
	}, nil
}

// Start is used to start the server.
// It will return an error if the server is already running.
func (s *Server) Start() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.running {
		return fmt.Errorf("Server already running")
	}

	s.server = grpc.NewServer()
	api.RegisterOpenStorageClusterServer(s.server, s)
	reflection.Register(s.server)

	// Start listening for requests
	logrus.Infof("[%s] gRPC SDK API Server ready on %s", s.driver.Name(), s.Address())
	waitForServer := make(chan bool)
	s.goServe(waitForServer)
	<-waitForServer

	s.running = true
	return nil
}

// Stop is used to stop the gRPC CSI complient server.
// It can be called multiple times. It does nothing if the server
// has already been stopped.
func (s *Server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.running {
		return
	}

	s.server.Stop()
	s.wg.Wait()
	s.running = false
}

// Address returns the address of the server which can be
// used by clients to connect.
func (s *Server) Address() string {
	return s.listener.Addr().String()
}

// IsRunning returns true if the server is currently running
func (s *Server) IsRunning() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.running
}

func (s *Server) goServe(started chan<- bool) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		started <- true
		err := s.server.Serve(s.listener)
		if err != nil {
			logrus.Fatalf("ERROR: Unable to start gRPC server: %s\n", err.Error())
		}
	}()
}
