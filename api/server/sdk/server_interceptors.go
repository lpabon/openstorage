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
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/libopenstorage/openstorage/pkg/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Funtion defined grpc_auth.AuthFunc()
func (s *sdkGrpcServer) auth(ctx context.Context) (context.Context, error) {
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

func (s *sdkGrpcServer) loggerServerInterceptor(
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

func (s *sdkGrpcServer) authorizationServerInterceptor(
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
	if "user" == tokenInfo.Role {
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
