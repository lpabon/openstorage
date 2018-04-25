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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/libopenstorage/openstorage/api"

	"github.com/portworx/kvdb"
)

func TestSdkVolumeCreate(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	name := "myvol"
	size := uint64(1234)
	req := &api.VolumeCreateRequest{
		Locator: &api.VolumeLocator{
			Name: name,
		},
		Source: &api.Source{},
		Spec: &api.VolumeSpec{
			Size: size,
		},
	}

	// Create response
	id := "myid"
	gomock.InOrder(
		s.MockDriver().
			EXPECT().
			Inspect([]string{name}).
			Return(nil, fmt.Errorf("not found")).
			Times(1),

		s.MockDriver().
			EXPECT().
			Enumerate(&api.VolumeLocator{Name: name}, nil).
			Return(nil, fmt.Errorf("not found")).
			Times(1),

		s.MockDriver().
			EXPECT().
			Create(req.GetLocator(), req.GetSource(), req.GetSpec()).
			Return(id, nil).
			Times(1),
	)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	r, err := c.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, r.GetId(), "myid")
}

func TestSdkVolumeCreateSimpleVolume(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	name := "myvol"
	size := int64(1234)
	req := &api.VolumeCreateSimpleVolumeRequest{
		Name: name,
		Size: size,
	}

	// Create response
	id := "myid"
	gomock.InOrder(
		s.MockDriver().
			EXPECT().
			Inspect([]string{name}).
			Return(nil, fmt.Errorf("not found")).
			Times(1),

		s.MockDriver().
			EXPECT().
			Enumerate(&api.VolumeLocator{Name: name}, nil).
			Return(nil, fmt.Errorf("not found")).
			Times(1),

		s.MockDriver().
			EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(id, nil).
			Times(1),
	)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	r, err := c.CreateSimpleVolume(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, r.GetId(), "myid")
}

func TestSdkVolumeDelete(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	id := "myvol"
	req := &api.VolumeDeleteRequest{
		VolumeId: id,
	}

	// Create response
	gomock.InOrder(
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{},
			}, nil).
			Times(1),

		s.MockDriver().
			EXPECT().
			Delete(id).
			Return(nil).
			Times(1),
	)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	_, err := c.Delete(context.Background(), req)
	assert.NoError(t, err)
}

func TestSdkVolumeDeleteReturnOkWhenVolumeNotFound(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	id := "myvol"
	req := &api.VolumeDeleteRequest{
		VolumeId: id,
	}

	// Create response
	s.MockDriver().
		EXPECT().
		Inspect([]string{id}).
		Return(nil, kvdb.ErrNotFound).
		Times(1)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	_, err := c.Delete(context.Background(), req)
	assert.NoError(t, err)
}

func TestSdkVolumeDeleteBadArguments(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	req := &api.VolumeDeleteRequest{}

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	_, err := c.Delete(context.Background(), req)
	assert.Error(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "volume id")
}

func TestSdkVolumeInspect(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	id := "myid"
	req := &api.VolumeInspectRequest{
		VolumeId: id,
	}

	s.MockDriver().
		EXPECT().
		Inspect([]string{id}).
		Return([]*api.Volume{
			&api.Volume{
				Id: id,
			},
		}, nil).
		Times(1)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	r, err := c.Inspect(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, r.GetVolume())
	assert.Equal(t, r.GetVolume().GetId(), id)
}

func TestSdkVolumeEnumerate(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	id := "myid"
	s.MockDriver().
		EXPECT().
		Enumerate(nil, nil).
		Return([]*api.Volume{
			&api.Volume{
				Id: id,
			},
		}, nil).
		Times(1)

	// Setup client
	c := api.NewOpenStorageVolumeClient(s.Conn())

	// Get info
	r, err := c.Enumerate(context.Background(), &api.VolumeEnumerateRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r.GetVolumes())
	assert.Len(t, r.GetVolumes(), 1)
	assert.Equal(t, r.GetVolumes()[0].GetId(), id)
}
