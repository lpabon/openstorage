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

	"github.com/libopenstorage/openstorage/api/client"
)

func (s *Server) Enumerate(context.Context, *client.ClusterEnumerateRequest) (*client.ClusterEnumerateResponse, error) {
	c, err := s.cluster.Enumerate()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cluster := client.Cluster{
		Status: c.Status,
		Id:     c.Id,
		NodeId: c.NodeId,
		Nodes:  c.Nodes,
	}

	return &client.ClusterEnumerateResponse{
		Cluster: c,
	}, nil
}

func (s *Server) Inspect(context.Context, *client.ClusterInspectRequest) (*client.ClusterInspectResponse, error) {
	panic("not implemented")
}

func (s *Server) SetSize(context.Context, *client.ClusterSetSizeRequest) (*client.ClusterSetSizeResponse, error) {
	panic("not implemented")
}

func (s *Server) Remove(context.Context, *client.ClusterRemoveRequest) (*client.ClusterRemoveResponse, error) {
	panic("not implemented")
}

func (s *Server) NodeStatus(context.Context, *client.ClusterNodeStatusRequest) (*client.ClusterNodeStatusResponse, error) {
	panic("not implemented")
}

func (s *Server) PeerStatus(context.Context, *client.ClusterPeerStatusRequest) (*client.ClusterPeerStatusResponse, error) {
	panic("not implemented")
}

func (s *Server) GetNodeIdFromIp(context.Context, *client.ClusterGetNodeIdFromIpRequest) (*client.ClusterGetNodeIdFromIpResponse, error) {
	panic("not implemented")
}

func (s *Server) EnumerateAlerts(context.Context, *client.ClusterEnumerateAlertsRequest) (*client.ClusterEnumerateAlertsResponse, error) {
	panic("not implemented")
}

func (s *Server) ClearAlert(context.Context, *client.ClusterClearAlertRequest) (*client.ClusterClearAlertResponse, error) {
	panic("not implemented")
}

func (s *Server) EraseAlert(context.Context, *client.ClusterEraseAlertRequest) (*client.ClusterEraseAlertResponse, error) {
	panic("not implemented")
}
