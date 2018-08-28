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
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/libopenstorage/openstorage/alerts"

	"github.com/gobuffalo/packr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	// AuthenticationTypeUnknown is a placeholder for unknown types
	AuthenticationTypeUnknown AuthenticationType = "unknown"
	// AuthenticationTypeSharedSecret is used when using JWT tokens signed with
	// a shared secret
	AuthenticationTypeSharedSecret AuthenticationType = "shared_secret"
	// Default unix domain socket location
	DefaultUnixDomainSocket = "/run/%s.sock"
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
	SharedSecret *AuthenticationSecretsConfig
}

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
	Auth AuthenticationConfig
	// Secure Tls configuration
	Tls *TLSConfig
}

// Server is an implementation of the gRPC SDK interface
type Server struct {
	*grpcserver.GrpcServer

	socketServer         *Server
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
	alertsServer         api.OpenStorageAlertsServer
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

	// Create socket server
	if config.Net != "unix" {
		socketName := defaultUnixDomainSocket
		if (config.Socket) != 0 {
			socketName = config.Socket
		}
		socketServer, error := New(&ServerConfig{
			Net:     "unix",
			Address: fmt.Sprintf(socketName, d.Name()),
		})
	}

	// Setup authentication
	var authenticator auth.Authenticator
	if config.Auth.Enabled {
		if config.Auth.SharedSecret != nil {
			logrus.Info("SDK authentication enabled using shared secrets")
			authenticator = auth.NewSharedSecret(&auth.SharedSecretConfig{
				AdminKey: []byte(config.Auth.SharedSecret.AdminKey),
				UserKey:  []byte(config.Auth.SharedSecret.UserKey),
			})
		} else {
			return nil, fmt.Errorf("Authentication enabled without authentication configuartion provided")
		}
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

	return &Server{
		GrpcServer:    gServer,
		config:        *config,
		authenticator: authenticator,
		restPort:      config.RestPort,
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
func (s *Server) Start() error {

	opts := make([]grpc.ServerOption, 0)
	if s.config.Tls != nil {
		creds, err := credentials.NewServerTLSFromFile(s.config.Tls.CertFile, s.config.Tls.KeyFile)
		if err != nil {
			return fmt.Errorf("Failed to create credentials from cert files: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
		logrus.Info("SDK TLS enabled")
	} else {
		logrus.Info("SDK TLS disabled")
	}

	if s.config.Auth.Enabled {
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
		var err error
		address := ":" + s.restPort
		if s.config.Tls != nil {
			err = http.ListenAndServeTLS(address, s.config.Tls.CertFile, s.config.Tls.KeyFile, mux)
		} else {
			err = http.ListenAndServe(address, mux)
		}

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

	// REST Gateway Handlers
	handlers := []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) (err error){
		api.RegisterOpenStorageClusterHandler,
		/*
			api.RegisterOpenStorageNodeHandlerFromEndpoint,
			api.RegisterOpenStorageVolumeHandlerFromEndpoint,
			api.RegisterOpenStorageObjectstoreHandlerFromEndpoint,
			api.RegisterOpenStorageCredentialsHandlerFromEndpoint,
			api.RegisterOpenStorageSchedulePolicyHandlerFromEndpoint,
			api.RegisterOpenStorageCloudBackupHandlerFromEndpoint,
			api.RegisterOpenStorageIdentityHandlerFromEndpoint,
		*/
	}

	// Determine if TLS is needed for the REST Gateway to connect to the gRPC server
	/*
		var opts []grpc.DialOption
		var creds credentials.TransportCredentials
		if s.config.Tls != nil {
			var err error
			creds, err = credentials.NewClientTLSFromFile(s.config.Tls.CertFile, "")
			if err != nil {
				return nil, fmt.Errorf("Failed to setup credentials for REST gateway: %v", err)
			}
			opts = []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.FailOnNonTempDialError(true)}
			logrus.Info(">>> HERE")
		} else {
			opts = []grpc.DialOption{grpc.WithInsecure()}
			logrus.Info(">>> insecure")
		}
	*/

	creds, err := credentials.NewClientTLSFromFile(s.config.Tls.CertFile, "example.com")
	if err != nil {
		return nil, fmt.Errorf("Failed to setup credentials for REST gateway: %v", err)
	}
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		add := s.Address()
		conn, err := net.Dial("tcp", add)
		if err != nil {
			logrus.Errorf("REST Gateway failed to dial gRPC server: %v", err)
			return nil, err
		}
		conn, _, err = creds.ClientHandshake(ctx, s.Address(), conn)
		if err != nil {
			logrus.Errorf("REST Gateway failed to connect gRPC server: %v", err)
			return nil, err
		}

		return conn, nil
	}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithDialer(dialer),
	}

	conn, err := grpc.Dial(s.Address(), opts...)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup REST Gateway connection to gRPC Server: %v", err)
	}

	// Register the REST Gateway handlers
	for _, handler := range handlers {
		err := handler(context.Background(), gmux, conn)
		if err != nil {
			return nil, err
		}
	}

	/*
		err = api.RegisterOpenStorageMountAttachHandlerFromEndpoint(
			context.Background(),
			gmux,
			s.Address(),
			[]grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			return nil, err
		}

		err = api.RegisterOpenStorageAlertsHandlerFromEndpoint(
			context.Background(),
			gmux,
			s.Address(),
			[]grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			return nil, err
		}
	*/
	clustergrpc := api.NewOpenStorageClusterClient(conn)
	ci, err := clustergrpc.InspectCurrent(context.Background(), &api.SdkClusterInspectCurrentRequest{})
	fmt.Printf("%v e:%v", ci, err)

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

	tokenInfo, err := s.authenticator.AuthenticateToken(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	ctx = context.WithValue(ctx, "tokeninfo", tokenInfo)

	return ctx, nil
}

func (s *Server) loggerServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if tokenInfo, ok := ctx.Value("tokeninfo").(*auth.Token); ok {
		logrus.WithFields(logrus.Fields{
			"user":   tokenInfo.User,
			"email":  tokenInfo.Email,
			"role":   tokenInfo.Role,
			"method": info.FullMethod,
		}).Info("called")
	} else {
		logrus.WithFields(logrus.Fields{
			"method": info.FullMethod,
		}).Info("called")
	}

	return handler(ctx, req)
}

func (s *Server) authorizationServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	tokenInfo, ok := ctx.Value("tokeninfo").(*auth.Token)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Authorization called without token")
	}

	// Check user role
	if auth.RoleUser == tokenInfo.Role {
		// User is not allowed the following services and/or methods
		// Example service:
		//    openstorage.api.OpenStorageNode
		// Example method:
		//    openstorage.api.OpenStorageCluster/InspectCurrent
		//
		// TODO: Make this configurable
		blacklist := []string{
			"openstorage.api.OpenStorageCluster",
			"openstorage.api.OpenStorageNode",
		}

		for _, notallowed := range blacklist {
			if strings.Contains(info.FullMethod, notallowed) {
				return nil, status.Errorf(
					codes.PermissionDenied,
					"Not role %s is not authorized to use %s",
					tokenInfo.Role,
					info.FullMethod,
				)
			}
		}
	}

	return handler(ctx, req)
}
