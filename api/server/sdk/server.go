/*
Package sdk is the gRPC implementation of the SDK gRPC server
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
package sdk

import (
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/libopenstorage/openstorage/alerts"

	"google.golang.org/grpc/credentials"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/spec"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/pkg/auth"
	"github.com/libopenstorage/openstorage/pkg/grpcserver"
	volumedrivers "github.com/libopenstorage/openstorage/volume/drivers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// AuthenticationType is the types of authentication
type AuthenticationType string

const (
	// Default unix domain socket location
	DefaultUnixDomainSocket = "/tmp/%s.sock"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

// ServerConfig provides the configuration to the SDK server
type ServerConfig struct {
	// Net is the transport for gRPC: unix, tcp, etc.
	// For the gRPC Server. This value goes together with `Address`.
	Net string
	// Address is the port number or the unix domain socket path.
	// For the gRPC Server. This value goes together with `Net`.
	Address string
	// RestAdress is the port number. Example: 9110
	// For the gRPC REST Gateway.
	RestPort string
	// Unix domain socket for local communication. This socket
	// will be used by the REST Gateway to communicate with the gRPC server.
	// Only set for testing. Having a '%s' can be supported to use the
	// name of the driver as the driver name.
	Socket string
	// The OpenStorage driver to use
	DriverName string
	// Cluster interface
	Cluster cluster.Cluster
	// AlertsFilterDeleter
	AlertsFilterDeleter alerts.FilterDeleter
	// Authentication configuration
	Auth *auth.JwtAuthConfig
	// Tls configuration
	Tls *TLSConfig
}

// Server is an implementation of the gRPC SDK interface
type Server struct {
	netServer   *sdkGrpcServer
	udsServer   *sdkGrpcServer
	restGateway *sdkRestGateway
}

type sdkGrpcServer struct {
	*grpcserver.GrpcServer

	name          string
	authenticator auth.Authenticator
	config        ServerConfig
	log           *logrus.Entry

	clusterServer        *ClusterServer
	nodeServer           *NodeServer
	volumeServer         *VolumeServer
	objectstoreServer    *ObjectstoreServer
	schedulePolicyServer *SchedulePolicyServer
	cloudBackupServer    *CloudBackupServer
	credentialServer     *CredentialServer
	identityServer       *IdentityServer
	alertsServer         api.OpenStorageAlertsServer
}

// Interface check
var _ grpcserver.Server = &sdkGrpcServer{}

func New(config *ServerConfig) (*Server, error) {

	// Create a gRPC server on the network
	netServer, err := newSdkGrpcServer(config)
	if err != nil {
		return nil, err
	}

	// Create a gRPC server on a unix domain socket
	udsConfig := *config
	udsConfig.Net = "unix"
	udsConfig.Address = DefaultUnixDomainSocket
	udsConfig.Tls = nil
	udsServer, err := newSdkGrpcServer(&udsConfig)
	if err != nil {
		return nil, err
	}

	// Create REST Gateway and connect it to the unix domain socket server
	restGeteway, err := newSdkRestGateway(config, udsServer)
	if err != nil {
		return nil, err
	}

	return &Server{
		netServer:   netServer,
		udsServer:   udsServer,
		restGateway: restGeteway,
	}, nil
}

// Start all servers
func (s *Server) Start() error {
	if err := s.netServer.Start(); err != nil {
		return err
	} else if err := s.udsServer.Start(); err != nil {
		return err
	} else if err := s.restGateway.Start(); err != nil {
		return err
	}

	return nil
}

// New creates a new SDK gRPC server
func newSdkGrpcServer(config *ServerConfig) (*sdkGrpcServer, error) {
	if nil == config {
		return nil, fmt.Errorf("Configuration must be provided")
	}
	if len(config.DriverName) == 0 {
		return nil, fmt.Errorf("OpenStorage Driver name must be provided")
	}

	// Create a log object for this server
	name := "SDK-" + config.Net
	log := logrus.WithFields(logrus.Fields{
		"name": name,
	})

	// Save the driver for future calls
	d, err := volumedrivers.Get(config.DriverName)
	if err != nil {
		return nil, fmt.Errorf("Unable to get driver %s info: %s", config.DriverName, err.Error())
	}

	// Setup authentication
	var authenticator auth.Authenticator
	if config.Auth != nil {
		authenticator, err = auth.New(config.Auth)
		if err != nil {
			return nil, err
		}
		log.Info(name + " authentication enabled")
	} else {
		log.Info(name + " authentication disabled")
	}

	// Setup unix domain socket name
	address := config.Address
	if config.Net == "unix" {
		address = fmt.Sprintf(config.Address, d.Name())
	}

	// Create gRPC server
	gServer, err := grpcserver.New(&grpcserver.GrpcServerConfig{
		Name:    name,
		Net:     config.Net,
		Address: address,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to setup server: %v", err)
	}

	return &sdkGrpcServer{
		GrpcServer:    gServer,
		config:        *config,
		name:          name,
		log:           log,
		authenticator: authenticator,
		identityServer: &IdentityServer{
			driver: d,
		},
		clusterServer: &ClusterServer{
			cluster: config.Cluster,
		},
		nodeServer: &NodeServer{
			cluster: config.Cluster,
		},
		volumeServer: &VolumeServer{
			driver:      d,
			cluster:     config.Cluster,
			specHandler: spec.NewSpecHandler(),
		},
		objectstoreServer: &ObjectstoreServer{
			cluster: config.Cluster,
		},
		schedulePolicyServer: &SchedulePolicyServer{
			cluster: config.Cluster,
		},
		cloudBackupServer: &CloudBackupServer{
			driver: d,
		},
		credentialServer: &CredentialServer{
			driver: d,
		},
		alertsServer: NewAlertsServer(config.AlertsFilterDeleter),
	}, nil
}

// Start is used to start the server.
// It will return an error if the server is already running.
func (s *sdkGrpcServer) Start() error {

	// Setup https if certs have been provided
	opts := make([]grpc.ServerOption, 0)
	if s.config.Tls != nil {
		creds, err := credentials.NewServerTLSFromFile(s.config.Tls.CertFile, s.config.Tls.KeyFile)
		if err != nil {
			return fmt.Errorf("Failed to create credentials from cert files: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
		s.log.Info("SDK TLS enabled")
	} else {
		s.log.Info("SDK TLS disabled")
	}

	// Setup authentication and authorization using interceptors if auth is enabled
	if s.config.Auth != nil {
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(s.auth),
				s.authorizationServerInterceptor,
				s.loggerServerInterceptor,
				grpc_recovery.UnaryServerInterceptor(),
			)))
	} else {
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				s.loggerServerInterceptor,
				grpc_recovery.UnaryServerInterceptor(),
			)))
	}

	// Start the gRPC Server
	err := s.GrpcServer.Start(func() *grpc.Server {
		grpcServer := grpc.NewServer(opts...)

		api.RegisterOpenStorageClusterServer(grpcServer, s.clusterServer)
		api.RegisterOpenStorageNodeServer(grpcServer, s.nodeServer)
		api.RegisterOpenStorageObjectstoreServer(grpcServer, s.objectstoreServer)
		api.RegisterOpenStorageVolumeServer(grpcServer, s.volumeServer)
		api.RegisterOpenStorageCredentialsServer(grpcServer, s.credentialServer)
		api.RegisterOpenStorageSchedulePolicyServer(grpcServer, s.schedulePolicyServer)
		api.RegisterOpenStorageCloudBackupServer(grpcServer, s.cloudBackupServer)
		api.RegisterOpenStorageIdentityServer(grpcServer, s.identityServer)
		api.RegisterOpenStorageMountAttachServer(grpcServer, s.volumeServer)
		api.RegisterOpenStorageAlertsServer(grpcServer, s.alertsServer)

		return grpcServer
	})
	if err != nil {
		return err
	}
	return nil
}
