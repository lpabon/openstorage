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
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/libopenstorage/openstorage/api"
)

// Enumerate is an api
func (s *Server) Enumerate(ctx context.Context, req *api.ClusterEnumerateRequest) (*api.ClusterEnumerateResponse, error) {
	c, err := s.cluster.Enumerate()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ClusterEnumerateResponse{
		Cluster: c.ToOSCluster(),
	}, nil
}

// Inspect is an API
func (s *Server) Inspect(ctx context.Context, req *api.ClusterInspectRequest) (*api.ClusterInspectResponse, error) {
	if len(req.GetNodeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Node id must be provided")
	}

	node, err := s.cluster.Inspect(req.GetNodeId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ClusterInspectResponse{
		Node: node.ToOSNode(),
	}, nil
}
