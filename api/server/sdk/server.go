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
	"context"
	"fmt"
	"mime"
	"net/http"

	"github.com/gobuffalo/packr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/spec"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/pkg/auth"
	"github.com/libopenstorage/openstorage/pkg/grpcserver"
	volumedrivers "github.com/libopenstorage/openstorage/volume/drivers"
)

// AuthenticationType is the types of authentication
type AuthenticationType string

const (
	// AuthenticationTypeUnknown is a placeholder for unknown types
	AuthenticationTypeUnknown AuthenticationType = "unknown"
	// AuthenticationTypeSharedSecret is used when using JWT tokens signed with
	// a shared secret
	AuthenticationTypeSharedSecret AuthenticationType = "shared_secret"
)

// AuthenticationSecretsConfig is used when using shared secrets
type AuthenticationSecretsConfig struct {
	// Administrator role key
	AdminKey string
	// User role key
	UserKey string
}

// AuthenticationConfig provides authentication configuration for the SDK server
type AuthenticationConfig struct {
	// Determine if the authentication should be enabled
	Enabled bool
	// Type of Authentication
	Type AuthenticationType
	// Shared secret configuration
	SharedSecret AuthenticationSecretsConfig
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
	// The OpenStorage driver to use
	DriverName string
	// Cluster interface
	Cluster cluster.Cluster
	// Authentication Configuration
	Auth AuthenticationConfig
}

// Server is an implementation of the gRPC SDK interface
type Server struct {
	*grpcserver.GrpcServer

	authenticator        auth.Authenticator
	config               ServerConfig
	restPort             string
	clusterServer        *ClusterServer
	nodeServer           *NodeServer
	volumeServer         *VolumeServer
	objectstoreServer    *ObjectstoreServer
	schedulePolicyServer *SchedulePolicyServer
	cloudBackupServer    *CloudBackupServer
	credentialServer     *CredentialServer
	identityServer       *IdentityServer
}

// Interface check
var _ grpcserver.Server = &Server{}

// New creates a new SDK gRPC server
func New(config *ServerConfig) (*Server, error) {
	if nil == config {
		return nil, fmt.Errorf("Configuration must be provided")
	}
	if len(config.DriverName) == 0 {
		return nil, fmt.Errorf("OpenStorage Driver name must be provided")
	}

	// Save the driver for future calls
	d, err := volumedrivers.Get(config.DriverName)
	if err != nil {
		return nil, fmt.Errorf("Unable to get driver %s info: %s", config.DriverName, err.Error())
	}

	// Create gRPC server
	gServer, err := grpcserver.New(&grpcserver.GrpcServerConfig{
		Name:    "SDK",
		Net:     config.Net,
		Address: config.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to setup server: %v", err)
	}

	var authenticator auth.Authenticator
	authenticator = auth.NewJwtAuth(&auth.JwtAuthConfig{
		AdminKey: []byte(config.Auth.SharedSecret.AdminKey)
		UserKey: []byte(config.Auth.SharedSecret.UserKey)
	})

	return &Server{
		GrpcServer: gServer,
		config:     *config,
		authenticator: authenticator,
		restPort:   config.RestPort,
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
	}, nil
}

// Start is used to start the server.
// It will return an error if the server is already running.
func (s *Server) Start() error {

	// Start the gRPC Server
	err := s.GrpcServer.Start(func() *grpc.Server {
		opts := make([]grpc.ServerOption, 0)

		if s.config.Auth.Enabled {
			opts = append(opts, grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(s.auth)))
		}

		grpcServer := grpc.NewServer(opts...)

		api.RegisterOpenStorageClusterServer(grpcServer, s.clusterServer)
		api.RegisterOpenStorageNodeServer(grpcServer, s.nodeServer)
		api.RegisterOpenStorageObjectstoreServer(grpcServer, s.objectstoreServer)
		api.RegisterOpenStorageVolumeServer(grpcServer, s.volumeServer)
		api.RegisterOpenStorageCredentialsServer(grpcServer, s.credentialServer)
		api.RegisterOpenStorageSchedulePolicyServer(grpcServer, s.schedulePolicyServer)
		api.RegisterOpenStorageCloudBackupServer(grpcServer, s.cloudBackupServer)
		api.RegisterOpenStorageIdentityServer(grpcServer, s.identityServer)

		return grpcServer
	})
	if err != nil {
		return err
	}
	if len(s.restPort) != 0 {
		return s.startRestServer()
	}
	return nil
}

// startRestServer starts the HTTP/REST gRPC gateway.
func (s *Server) startRestServer() error {

	mux, err := s.restServerSetupHandlers()
	if err != nil {
		return err
	}

	ready := make(chan bool)
	go func() {
		ready <- true
		err := http.ListenAndServe(":"+s.restPort, mux)
		if err != nil {
			logrus.Fatalf("Unable to start SDK REST gRPC Gateway: %s\n",
				err.Error())
		}
	}()
	<-ready
	logrus.Infof("SDK gRPC REST Gateway started on port :%s", s.restPort)

	return nil
}

// restServerSetupHandlers sets up the handlers to the swagger ui and
// to the gRPC REST Gateway.
func (s *Server) restServerSetupHandlers() (*http.ServeMux, error) {

	// Create an HTTP server router
	mux := http.NewServeMux()

	// Swagger files using packr
	swaggerUIBox := packr.NewBox("./swagger-ui")
	swaggerJSONBox := packr.NewBox("./api")
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Handler to return swagger.json
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write(swaggerJSONBox.Bytes("api.swagger.json"))
	})

	// Handler to access the swagger ui. The UI pulls the swagger
	// json file from /swagger.json
	// The link below MUST have th last '/'. It is really important.
	prefix := "/swagger-ui/"
	mux.Handle(prefix,
		http.StripPrefix(prefix, http.FileServer(swaggerUIBox)))

	// Create a router just for HTTP REST gRPC Server Gateway
	gmux := runtime.NewServeMux()
	err := api.RegisterOpenStorageClusterHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageNodeHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageVolumeHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageObjectstoreHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageCredentialsHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageSchedulePolicyHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageCloudBackupHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = api.RegisterOpenStorageIdentityHandlerFromEndpoint(
		context.Background(),
		gmux,
		s.Address(),
		[]grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	// Pass all other unhandled paths to the gRPC gateway
	mux.Handle("/", gmux)

	return mux, nil
}

// Funtion defined grpc_auth.AuthFunc()
func (s *Server) auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	logrus.Infof("Token: %s", token)

	tokenInfo, err := s.authenticator(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	ctx = context.WithValue(ctx, "tokeninfo", tokenInfo)

	return ctx, nil
}
