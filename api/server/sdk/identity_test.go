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
	"testing"

	"github.com/libopenstorage/openstorage/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func expectCapability(
	t *testing.T,
	expected api.SdkServiceCapability_OpenStorageService_Type,
	capabilities []*api.SdkServiceCapability,
) {

	for _, capOneOf := range capabilities {
		cap := capOneOf.GetService().GetType()
		if cap == expected {
			return
		}
	}

	t.Errorf("Capability %s not found in %+v", expected, capabilities)
}

func AddMetadataToContext(ctx context.Context, k, v string) context.Context {
	// You can also use the go-grpc-middleware/util/metautils utility functions
	// which give you Get, Set, Add on top of metadata.
	// But we will use the normal metadata api for simplicity:

	// Doesn't matter if there was one there or not already,
	// passing nil metadata to Join with a new one will do the
	// right thing
	md, _ := metadata.FromOutgoingContext(ctx)
	md = metadata.Join(md, metadata.New(map[string]string{
		k: v,
	}))
	return metadata.NewOutgoingContext(ctx, md)
}

func TestIdentityContextMetadata(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	version := &api.StorageVersion{
		Driver:  "mock",
		Version: "1.2.4-asdf",
	}
	s.MockDriver().EXPECT().Version().Return(version, nil).Times(1)

	// Create a connection
	c := api.NewOpenStorageIdentityClient(s.Conn())

	// setup context
	ctx := AddMetadataToContext(context.Background(), "hello", "world")
	ctx = AddMetadataToContext(ctx, "jay", "kay")
	ctx = AddMetadataToContext(ctx, "one", "two")

	r, err := c.Version(ctx, &api.SdkIdentityVersionRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r)

	details := r.GetVersion().GetDetails()
	assert.Len(t, details, 3)
	assert.Equal(t, details["hello"], "world")
	assert.Equal(t, details["jay"], "kay")
	assert.Equal(t, details["one"], "two")
}

func TestIdentityCapabilities(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	c := api.NewOpenStorageIdentityClient(s.Conn())

	// Get identities
	r, err := c.Capabilities(context.Background(), &api.SdkIdentityCapabilitiesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r)

	expectedCapabilities := []api.SdkServiceCapability_OpenStorageService_Type{
		api.SdkServiceCapability_OpenStorageService_CLUSTER,
		api.SdkServiceCapability_OpenStorageService_CLOUD_BACKUP,
		api.SdkServiceCapability_OpenStorageService_CREDENTIALS,
		api.SdkServiceCapability_OpenStorageService_NODE,
		api.SdkServiceCapability_OpenStorageService_OBJECT_STORAGE,
		api.SdkServiceCapability_OpenStorageService_SCHEDULE_POLICY,
		api.SdkServiceCapability_OpenStorageService_VOLUME,
		api.SdkServiceCapability_OpenStorageService_ALERTS,
		api.SdkServiceCapability_OpenStorageService_MOUNT_ATTACH,
	}

	for _, cap := range expectedCapabilities {
		expectCapability(t, cap, r.GetCapabilities())
	}
}

func TestIdentityVersion(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	c := api.NewOpenStorageIdentityClient(s.Conn())
	s.MockDriver().EXPECT().Version().Return(nil, fmt.Errorf("MOCK")).Times(1)
	_, err := c.Version(context.Background(), &api.SdkIdentityVersionRequest{})
	assert.Error(t, err)
	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.Internal)
	assert.Contains(t, serverError.Message(), "MOCK")

	version := &api.StorageVersion{
		Driver:  "mock",
		Version: "1.2.4-asdf",
		Details: map[string]string{
			"hello": "world",
		},
	}

	s.MockDriver().EXPECT().Version().Return(version, nil).Times(1)
	r, err := c.Version(context.Background(), &api.SdkIdentityVersionRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.NotNil(t, r.GetSdkVersion())
	assert.Equal(t, int32(api.SdkVersion_Major), r.GetSdkVersion().GetMajor())
	assert.Equal(t, int32(api.SdkVersion_Minor), r.GetSdkVersion().GetMinor())
	assert.Equal(t, int32(api.SdkVersion_Patch), r.GetSdkVersion().GetPatch())
	assert.Equal(t,
		fmt.Sprintf("%d.%d.%d",
			api.SdkVersion_Major,
			api.SdkVersion_Minor,
			api.SdkVersion_Patch,
		),
		r.GetSdkVersion().GetVersion())

	assert.NotNil(t, r.GetVersion())
	assert.Equal(t, version.GetDriver(), r.GetVersion().GetDriver())
	assert.Equal(t, version.GetVersion(), r.GetVersion().GetVersion())
	assert.Equal(t, version.GetDetails(), r.GetVersion().GetDetails())
}
