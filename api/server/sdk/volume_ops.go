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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/pkg/util"
	"github.com/portworx/kvdb"
)

// Create creates a volume
func (s *VolumeServer) Create(
	ctx context.Context,
	req *api.VolumeCreateRequest,
) (*api.VolumeCreateResponse, error) {

	if req.GetLocator() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"Must supply a locator object")
	} else if len(req.GetLocator().GetName()) == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"Must supply a unique name in locator")
	} else if req.GetSource() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"Must supply source object")
	} else if req.GetSpec() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"Must supply spec object")
	}

	locator := req.GetLocator()
	spec := req.GetSpec()
	source := req.GetSource()
	volName := locator.GetName()

	// Check if the volume has already been created or is in process of creation
	v, err := util.VolumeFromName(s.driver, volName)
	if err == nil {
		// Check the requested arguments match that of the existing volume
		if spec.Size != v.GetSpec().GetSize() {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"Existing volume has a size of %v which differs from requested size of %v",
				v.GetSpec().GetSize(),
				spec.Size)
		}
		if v.GetSpec().GetShared() != req.GetSpec().GetShared() {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"Existing volume has shared=%v while request is asking for shared=%v",
				v.GetSpec().GetShared(),
				req.GetSpec().GetShared())
		}
		if v.GetSource().GetParent() != source.GetParent() {
			return nil, status.Error(codes.AlreadyExists, "Existing volume has conflicting parent value")
		}

		// Return information on existing volume
		return &api.VolumeCreateResponse{
			Id: v.GetId(),
		}, nil
	}

	// Check if the caller is asking to create a snapshot or for a new volume
	var id string
	if source != nil && len(source.GetParent()) != 0 {
		// Get parent volume information
		parent, err := util.VolumeFromName(s.driver, source.Parent)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"unable to get parent volume information: %s",
				err.Error())
		}

		// Create a snapshot from the parent
		id, err = s.driver.Snapshot(parent.GetId(), false, &api.VolumeLocator{
			Name: volName,
		})
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"unable to create snapshot: %s\n",
				err.Error())
		}
	} else {
		// Create the volume
		id, err = s.driver.Create(locator, source, spec)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"Failed to create volume: %v",
				err.Error())
		}
	}

	return &api.VolumeCreateResponse{
		Id: id,
	}, nil
}

// CreateSimpleVolume provides a simple API to create a volume
func (s *VolumeServer) CreateSimpleVolume(
	ctx context.Context,
	req *api.VolumeCreateSimpleVolumeRequest,
) (*api.VolumeCreateResponse, error) {

	// Generate objects from the parameters passed in
	spec, locator, source, err := s.specHandler.SpecFromOpts(req.GetParameters())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Unable to get parameters: %s\n",
			err.Error())
	}
	locator.Name = req.GetName()
	spec.Size = uint64(req.GetSize())

	// Create volume
	return s.Create(ctx, &api.VolumeCreateRequest{
		Locator: locator,
		Source:  source,
		Spec:    spec,
	})
}

// Delete deletes a volume
func (s *VolumeServer) Delete(
	ctx context.Context,
	req *api.VolumeDeleteRequest,
) (*api.VolumeDeleteResponse, error) {

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply volume id")
	}

	// If the volume is not found, return OK to be idempotent
	volumes, err := s.driver.Inspect([]string{req.GetVolumeId()})
	if (err == nil && len(volumes) == 0) ||
		(err != nil && err == kvdb.ErrNotFound) {
		return &api.VolumeDeleteResponse{}, nil
	} else if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Failed to determine if volume id %s exists: %v",
			req.GetVolumeId(),
			err.Error())
	}

	err = s.driver.Delete(req.GetVolumeId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Failed to delete volume %s: %v",
			req.GetVolumeId(),
			err.Error())
	}

	return &api.VolumeDeleteResponse{}, nil
}

// Inspect returns information about a volume
func (s *VolumeServer) Inspect(
	ctx context.Context,
	req *api.VolumeInspectRequest,
) (*api.VolumeInspectResponse, error) {

	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply volume id")
	}

	vols, err := s.driver.Inspect([]string{req.GetVolumeId()})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Failed to inspect volume %s: %v",
			req.GetVolumeId(),
			err.Error())
	}

	return &api.VolumeInspectResponse{
		Volume: vols[0],
	}, nil
}

// Enumerate returns a list of volumes
func (s *VolumeServer) Enumerate(
	ctx context.Context,
	req *api.VolumeEnumerateRequest,
) (*api.VolumeEnumerateResponse, error) {

	vols, err := s.driver.Enumerate(req.GetLocator(), nil)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Failed to enumerate volumes: %v",
			err.Error())
	}

	return &api.VolumeEnumerateResponse{
		Volumes: vols,
	}, nil
}
