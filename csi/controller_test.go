/*
CSI Interface for OSD
Copyright 2017 Portworx

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
package csi

import (
	"fmt"
	"testing"

	"github.com/libopenstorage/openstorage/api"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestControllerGetCapabilities(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Make a call
	c := csi.NewControllerClient(s.Conn())
	r, err := c.ControllerGetCapabilities(context.Background(), &csi.ControllerGetCapabilitiesRequest{})
	assert.Nil(t, err)

	// Verify
	expectedValues := []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	}
	caps := r.GetResult().GetCapabilities()
	assert.Len(t, caps, len(expectedValues))
	found := 0
	for _, expectedCap := range expectedValues {
		for _, cap := range caps {
			if cap.GetRpc().GetType() == expectedCap {
				found++
				break
			}
		}
	}
	assert.Equal(t, found, len(expectedValues))
}

func TestControllerPublishVolume(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Make a call
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ControllerPublishVolume(context.Background(), &csi.ControllerPublishVolumeRequest{})
	assert.NotNil(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.Unimplemented)
	assert.Contains(t, serverError.Message(), "not supported")
}

func TestControllerUnPublishVolume(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Make a call
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ControllerUnpublishVolume(context.Background(), &csi.ControllerUnpublishVolumeRequest{})
	assert.NotNil(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.Unimplemented)
	assert.Contains(t, serverError.Message(), "not supported")
}

func TestControllerValidateVolumeCapabilitiesBadArguments(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	req := &csi.ValidateVolumeCapabilitiesRequest{}

	// Missing everything
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "Version")

	// Miss capabilities and id
	req.Version = &csi.Version{}
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "capabilities")

	// Miss id and capabilities len is 0
	req.Version = &csi.Version{}
	req.VolumeCapabilities = []*csi.VolumeCapability{}
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "capabilities")

	// Miss id
	req.Version = &csi.Version{}
	req.VolumeCapabilities = []*csi.VolumeCapability{
		&csi.VolumeCapability{},
	}
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "volume_id")
}

func TestControllerValidateVolumeInvalidId(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"
	gomock.InOrder(
		// First time called it will say it is not there
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return(nil, fmt.Errorf("Id not found")),

		// Second time called it will not return an error,
		// but return an empty list
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{}, nil),

		// Third time it is called, it will return
		// a good list with a volume with an id that does
		// not match (even if this probably never could happen)
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id: "bad volume id",
				},
			}, nil),

		// Fourth time driver will return a list with more than
		// one volume, which should be unexpected since it only
		// asked for one volume.
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id: "bad volume id 1",
				},
				&api.Volume{
					Id: "bad volume id 2",
				},
			}, nil),
	)

	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{},
		},
		VolumeId: id,
	}

	// Missing everything
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.NotFound)
	assert.Contains(t, serverError.Message(), "ID not found")

	// Send again, same result
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.NotFound)
	assert.Contains(t, serverError.Message(), "ID not found")

	// Now it should be an internal id error
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.Internal)
	assert.Contains(t, serverError.Message(), "Driver volume id")

	// Now driver should have returned an unexpected number of volumes
	_, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok = status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.Internal)
	assert.Contains(t, serverError.Message(), "unexpected number of volumes")
}

func TestControllerValidateVolumeInvalidCapabilities(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"
	s.MockDriver().
		EXPECT().
		Inspect([]string{id}).
		Return([]*api.Volume{
			&api.Volume{
				Id: id,
			},
		}, nil).
		Times(1)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{},
		},
		VolumeId: id,
	}

	// Make request
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)

	serverError, ok := status.FromError(err)
	assert.True(t, ok)

	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "Cannot have both")
}

func TestControllerValidateVolumeAccessModeSNWR(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"

	// RO SH
	// x   x
	// *   x
	// x   *
	// *   *
	gomock.InOrder(
		// not-RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// not-RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),

		// RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),
	)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
		},
		VolumeId: id,
	}

	// Expect non-RO and non-SH
	c := csi.NewControllerClient(s.Conn())
	r, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, r.GetResult().Supported)

	// Expect RO and non-SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect non-RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)
}

func TestControllerValidateVolumeAccessModeSNRO(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"

	// RO SH
	// x   x
	// *   x
	// x   *
	// *   *
	gomock.InOrder(
		// not-RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// not-RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),

		// RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),
	)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
				},
			},
		},
		VolumeId: id,
	}

	// Expect non-RO and non-SH
	c := csi.NewControllerClient(s.Conn())
	r, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and non-SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, r.GetResult().Supported)

	// Expect non-RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)
}

func TestControllerValidateVolumeAccessModeMNRO(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"

	// RO SH
	// x   x
	// *   x
	// x   *
	// *   *
	gomock.InOrder(
		// not-RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// not-RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),

		// RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),
	)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
				},
			},
		},
		VolumeId: id,
	}

	// Expect non-RO and non-SH
	c := csi.NewControllerClient(s.Conn())
	r, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and non-SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect non-RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, r.GetResult().Supported)
}

func TestControllerValidateVolumeAccessModeMNWR(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"

	// RO SH
	// x   x
	// *   x
	// x   *
	// *   *
	gomock.InOrder(
		// not-RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// RO and not-SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: false,
					},
				},
			}, nil),

		// not-RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: false,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),

		// RO and SH
		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id:       id,
					Readonly: true,
					Spec: &api.VolumeSpec{
						Shared: true,
					},
				},
			}, nil),
	)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
				},
			},
		},
		VolumeId: id,
	}

	// Expect non-RO and non-SH
	c := csi.NewControllerClient(s.Conn())
	r, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect RO and non-SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)

	// Expect non-RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, r.GetResult().Supported)

	// Expect RO and SH
	r, err = c.ValidateVolumeCapabilities(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, r.GetResult().Supported)
}

func TestControllerValidateVolumeAccessModeUnknown(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup mock
	id := "testvolumeid"
	s.MockDriver().
		EXPECT().
		Inspect([]string{id}).
		Return([]*api.Volume{
			&api.Volume{
				Id:       id,
				Readonly: false,
				Spec: &api.VolumeSpec{
					Shared: false,
				},
			},
		}, nil).
		Times(1)

	// Setup request
	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version: &csi.Version{},
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_UNKNOWN,
				},
			},
		},
		VolumeId: id,
	}

	// Expect non-RO and non-SH
	c := csi.NewControllerClient(s.Conn())
	_, err := c.ValidateVolumeCapabilities(context.Background(), req)
	assert.NotNil(t, err)
}

func TestControllerListVolumesInvalidArguments(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// Setup request
	req := &csi.ListVolumesRequest{}

	// Expect error without version
	_, err := c.ListVolumes(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "Version")

	// Expect error with maxentries set
	// To be removed once CSI Spec issue #138 is resolved
	req.Version = &csi.Version{}
	req.MaxEntries = 1
	_, err = c.ListVolumes(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok = status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.Unimplemented)
	assert.Contains(t, serverError.Message(), "token")
}

func TestControllerListVolumesEnumerateError(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// Setup mock
	s.MockDriver().
		EXPECT().
		Enumerate(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("TEST")).
		Times(1)

	// Setup request
	req := &csi.ListVolumesRequest{
		Version: &csi.Version{},
	}

	// Expect error without version
	_, err := c.ListVolumes(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.Internal)
	assert.Contains(t, serverError.Message(), "TEST")
}

func TestControllerListVolumes(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// Setup mock
	mockVolumeList := []*api.Volume{
		&api.Volume{
			Id:            "one",
			AttachedState: api.AttachState_ATTACH_STATE_INTERNAL,
			Readonly:      false,
			State:         api.VolumeState_VOLUME_STATE_ERROR,
			Spec: &api.VolumeSpec{
				Shared: true,
				Size:   uint64(11),
			},
			Error: "TEST",
		},
		&api.Volume{
			Id:            "two",
			AttachedState: api.AttachState_ATTACH_STATE_EXTERNAL,
			Readonly:      true,
			State:         api.VolumeState_VOLUME_STATE_ATTACHED,
			Spec: &api.VolumeSpec{
				Shared: true,
				Size:   uint64(22),
			},
		},
		&api.Volume{
			Id:            "three",
			AttachedState: api.AttachState_ATTACH_STATE_INTERNAL_SWITCH,
			Readonly:      true,
			State:         api.VolumeState_VOLUME_STATE_TRY_DETACHING,
			Spec: &api.VolumeSpec{
				Shared: false,
				Size:   uint64(33),
			},
		},
	}
	s.MockDriver().
		EXPECT().
		Enumerate(gomock.Any(), gomock.Any()).
		Return(mockVolumeList, nil).
		Times(1)

	// Setup request
	req := &csi.ListVolumesRequest{
		Version: &csi.Version{},
	}

	// Expect error without version
	r, err := c.ListVolumes(context.Background(), req)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	volumes := r.GetResult().GetEntries()
	assert.Equal(t, len(mockVolumeList), len(volumes))
	assert.Equal(t, len(r.GetResult().GetNextToken()), 0)

	found := 0
	for _, mv := range mockVolumeList {
		for _, v := range volumes {
			info := v.GetVolumeInfo()
			assert.NotNil(t, info)

			if mv.GetId() == info.GetId() {
				found++
				assert.Equal(t, info.GetCapacityBytes(), mv.GetSpec().GetSize())

				attributes := info.GetAttributes()
				assert.Equal(t, attributes["readonly"], fmt.Sprintf("%v", mv.GetReadonly()))
				assert.Equal(t, attributes["shared"], fmt.Sprintf("%v", mv.GetSpec().GetShared()))
				assert.Equal(t, attributes["state"], mv.GetState().String())
				assert.Equal(t, attributes["attached"], mv.GetAttachedState().String())
				assert.Equal(t, attributes["error"], mv.GetError())
				break
			}
		}
	}
	assert.Equal(t, found, len(mockVolumeList))
}

func TestControllerCreateVolume(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// Setup request
	name := "myvol"
	size := uint64(1234)
	req := &csi.CreateVolumeRequest{
		Version: &csi.Version{},
		Name:    name,
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{},
		},
		CapacityRange: &csi.CapacityRange{
			RequiredBytes: size,
		},
	}

	// Setup mock functions
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

		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id: id,
					Locator: &api.VolumeLocator{
						Name: name,
					},
					Spec: &api.VolumeSpec{
						Size: size,
					},
				},
			}, nil).
			Times(1),
	)

	r, err := c.CreateVolume(context.Background(), req)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	volumeInfo := r.GetResult().GetVolumeInfo()

	assert.Equal(t, id, volumeInfo.GetId())
	assert.Equal(t, size, volumeInfo.GetCapacityBytes())
}

func TestControllerCreateVolumeSnapshot(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// Setup request
	mockParentID := "parendId"
	name := "myvol"
	size := uint64(1234)
	req := &csi.CreateVolumeRequest{
		Version: &csi.Version{},
		Name:    name,
		VolumeCapabilities: []*csi.VolumeCapability{
			&csi.VolumeCapability{},
		},
		CapacityRange: &csi.CapacityRange{
			RequiredBytes: size,
		},
		Parameters: map[string]string{
			api.SpecParent: mockParentID,
		},
	}

	// Setup mock functions
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
			Inspect([]string{mockParentID}).
			Return([]*api.Volume{
				&api.Volume{
					Id: mockParentID,
				},
			}, nil).
			Times(1),

		s.MockDriver().
			EXPECT().
			Snapshot(mockParentID, false, &api.VolumeLocator{
				Name: name,
			}).
			Return(id, nil).
			Times(1),

		s.MockDriver().
			EXPECT().
			Inspect([]string{id}).
			Return([]*api.Volume{
				&api.Volume{
					Id: id,
					Locator: &api.VolumeLocator{
						Name: name,
					},
					Spec: &api.VolumeSpec{
						Size: size,
					},
				},
			}, nil).
			Times(1),
	)

	r, err := c.CreateVolume(context.Background(), req)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	volumeInfo := r.GetResult().GetVolumeInfo()

	assert.Equal(t, id, volumeInfo.GetId())
	assert.Equal(t, size, volumeInfo.GetCapacityBytes())
}

func TestControllerDeleteVolumeInvalidArguments(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// No version
	req := &csi.DeleteVolumeRequest{}
	_, err := c.DeleteVolume(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "Version")

	// No id
	req.Version = &csi.Version{}
	_, err = c.DeleteVolume(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok = status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.InvalidArgument)
	assert.Contains(t, serverError.Message(), "Volume id")
}

func TestControllerDeleteVolumeError(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// No version
	myid := "myid"
	req := &csi.DeleteVolumeRequest{
		Version:  &csi.Version{},
		VolumeId: myid,
	}

	// Setup mock
	s.MockDriver().EXPECT().Delete(myid).Return(fmt.Errorf("TEST")).Times(1)

	_, err := c.DeleteVolume(context.Background(), req)
	assert.NotNil(t, err)
	serverError, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, serverError.Code(), codes.Internal)
	assert.Contains(t, serverError.Message(), "Unable to delete")
}

func TestControllerDeleteVolume(t *testing.T) {
	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()
	c := csi.NewControllerClient(s.Conn())

	// No version
	myid := "myid"
	req := &csi.DeleteVolumeRequest{
		Version:  &csi.Version{},
		VolumeId: myid,
	}

	// Setup mock
	s.MockDriver().EXPECT().Delete(myid).Return(nil).Times(1)

	_, err := c.DeleteVolume(context.Background(), req)
	assert.Nil(t, err)
}
