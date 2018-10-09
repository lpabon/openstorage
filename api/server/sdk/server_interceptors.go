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
	"encoding/json"

	sdk_auth "github.com/libopenstorage/openstorage-sdk-auth/pkg/auth"
	"github.com/sirupsen/logrus"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorContextkey string

const (
	InterceptorContextTokenKey InterceptorContextkey = "tokenclaims"
)

// Authenticate user and add authorization information back in the context
func (s *sdkGrpcServer) auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	// Authenticate user
	claims, err := s.authenticator.AuthenticateToken(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	// Add authorization information back into the context so that other
	// functions can get access to this information
	ctx = context.WithValue(ctx, InterceptorContextTokenKey, claims)

	return ctx, nil
}

func (s *sdkGrpcServer) loggerServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if claims, ok := ctx.Value(InterceptorContextTokenKey).(*sdk_auth.Claims); ok {
		// Change claims to JSON string to print into log
		claimsJSON, _ := json.Marshal(claims)
		logrus.WithFields(logrus.Fields{
			"name":   claims.Name,
			"email":  claims.Email,
			"role":   claims.Role,
			"claims": string(claimsJSON),
			"method": info.FullMethod,
		}).Info("audit")
	} else {
		logrus.WithFields(logrus.Fields{
			"method": info.FullMethod,
		}).Info("audit without authentication")
	}

	return handler(ctx, req)
}

func (s *sdkGrpcServer) authorizationServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	/*
		tokenInfo, ok := ctx.Value(InterceptorContextTokenKey).(*auth.Token)
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
						"Role %s is not authorized to use %s",
						tokenInfo.Role,
						info.FullMethod,
					)
				}
			}
		}
	*/

	return handler(ctx, req)
}
