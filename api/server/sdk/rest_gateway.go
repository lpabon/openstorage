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

type SdkRestGateway struct {}

func NewSdkRestGateway(server *Server)  *SdkRestGateway{
	return nil
	
}

func (*s SdkRestGateway) Start() error {
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
func (s *SdkRestGateway) restServerSetupHandlers() (*http.ServeMux, error) {

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

	clustergrpc := api.NewOpenStorageClusterClient(conn)
	ci, err := clustergrpc.InspectCurrent(context.Background(), &api.SdkClusterInspectCurrentRequest{})
	fmt.Printf("%v e:%v", ci, err)

	// Pass all other unhandled paths to the gRPC gateway
	mux.Handle("/", gmux)

	return mux, nil
}