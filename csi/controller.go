/*
Package csi is CSI driver interface for OSD
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

	"github.com/portworx/kvdb"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/pkg/util"

	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	volumeCapabilityMessageMultinodeVolume    = "Volume is a multinode volume"
	volumeCapabilityMessageNotMultinodeVolume = "Volume is not a multinode volume"
	volumeCapabilityMessageReadOnlyVolume     = "Volume is read only"
	volumeCapabilityMessageNotReadOnlyVolume  = "Volume is not read only"
	defaultCSIVolumeSize                      = uint64(1024 * 1024 * 1024)
)

// ControllerGetCapabilities is a CSI API functions which returns to the caller
// the capabilities of the OSD CSI driver.
func (s *OsdCsiServer) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest,
) (*csi.ControllerGetCapabilitiesResponse, error) {

	// Creating and deleting volumes supported
	capCreateDeleteVolume := &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			},
		},
	}

	// Creating and deleting snapshots
	capCreateDeleteSnapshot := &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			},
		},
	}

	// ListVolumes supported
	capListVolumes := &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
			},
		},
	}

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			capCreateDeleteVolume,
			capCreateDeleteSnapshot,
			capListVolumes,
		},
	}, nil

}

// ControllerPublishVolume is a CSI API implements the attachment of a volume
// on to a node.
func (s *OsdCsiServer) ControllerPublishVolume(
	context.Context,
	*csi.ControllerPublishVolumeRequest,
) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "This request is not supported")
}

// ControllerUnpublishVolume is a CSI API which implements the detaching of a volume
// onto a node.
func (s *OsdCsiServer) ControllerUnpublishVolume(
	context.Context,
	*csi.ControllerUnpublishVolumeRequest,
) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "This request is not supported")
}

// ValidateVolumeCapabilities is a CSI API used by container orchestration systems
// to make sure a volume specification is validiated by the CSI driver.
func (s *OsdCsiServer) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest,
) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	capabilities := req.GetVolumeCapabilities()
	if capabilities == nil || len(capabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume_capabilities must be specified")
	}
	id := req.GetVolumeId()
	if len(id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume_id must be specified")
	}

	// Log request
	logrus.Debugf("ValidateVolumeCapabilities of id %s "+
		"capabilities %#v "+
		id,
		capabilities)

	// Get grpc connection
	conn, err := s.getConn()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to connect to SDK server: %v", err)
	}

	// Get secret if any was passed
	ctx := s.setupContextWithToken(ctx, req.GetSecrets())

	// Check ID is valid with the specified volume capabilities
	volumes := api.NewOpenStorageVolumeClient(conn)
	resp, err := volumes.Inspect(ctx, &api.SdkVolumeInspectRequest{
		VolumeId: id,
	})
	if err != nil {
		return nil, err
	}
	v := resp.GetVolume()
	// Setup uninitialized response object
	result := &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeContext:      req.GetVolumeContext(),
			VolumeCapabilities: req.GetVolumeCapabilities(),
			Parameters:         req.GetParameters(),
		},
	}

	// Check capability
	for _, capability := range capabilities {
		// Currently the CSI spec defines all storage as "file systems."
		// So we do not need to check this with the volume. All we will check
		// here is the validity of the capability access type.
		if capability.GetMount() == nil && capability.GetBlock() == nil {
			return nil, status.Error(
				codes.InvalidArgument,
				"Cannot have both mount and block be undefined")
		}

		// Check access mode is setup correctly
		mode := capability.GetAccessMode()
		switch {
		case mode.Mode == csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER:
			if v.Spec.Shared {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageMultinodeVolume
				break
			}
			if v.Readonly {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY:
			if v.Spec.Shared {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageMultinodeVolume
				break
			}
			if !v.Readonly {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageNotReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY:
			if !v.Spec.Shared {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageNotMultinodeVolume
				break
			}
			if !v.Readonly {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageNotReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER ||
			mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER:
			if !v.Spec.Shared {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageNotMultinodeVolume
				break
			}
			if v.Readonly {
				result.Confirmed = nil
				result.Message = volumeCapabilityMessageReadOnlyVolume
				break
			}
		default:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"AccessMode %s is not allowed",
				mode.Mode.String())
		}

		if result.Confirmed == nil {
			return result, nil
		}
	}

	// If we passed all the checks, then it is valid.
	// result.Message needs to be empty on return
	return result, nil
}

// ListVolumes is a CSI API which returns to the caller all volume ids
// on this cluster. This includes ids created by CSI and ids created
// using other interfaces. This is important because the user could
// be requesting to mount a OSD volume created using non-CSI interfaces.
//
// This call does not yet implement tokens to due to the following
// issue: https://github.com/container-storage-interface/spec/issues/138
func (s *OsdCsiServer) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest,
) (*csi.ListVolumesResponse, error) {

	logrus.Debugf("ListVolumes req[%#v]", req)

	// Until the issue #138 on the CSI spec is resolved we will not support
	// tokenization
	if req.GetMaxEntries() != 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Driver does not support tokenization. Please see "+
				"https://github.com/container-storage-interface/spec/issues/138")
	}

	// Get grpc connection
	conn, err := s.getConn()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to connect to SDK server: %v", err)
	}

	// Get secret if any was passed
	ctx := s.setupContextWithToken(ctx, req.GetSecrets())

	// Check ID is valid with the specified volume capabilities
	volumes := api.NewOpenStorageVolumeClient(conn)

	resp, err := volumes.Enumerate(ctx, &api.SdkVolumeEnumerateRequest{})
	if err != nil {
		return err, nil
	}
	vols := resp.GetVolumeIds()

	entries := make([]*csi.ListVolumesResponse_Entry, len(vols))
	for i, v := range vols {
		// Initialize entry
		entries[i] = &csi.ListVolumesResponse_Entry{
			Volume: &csi.Volume{},
		}

		// Required
		entries[i].Volume.VolumeId = v.Id

		// This entry is optional in the API, but OSD has
		// the information available to provide it
		entries[i].Volume.CapacityBytes = int64(v.Spec.Size)

		// Attributes. We can add or remove as needed since they
		// are optional and opaque to the Container Orchestrator(CO)
		// but could be used for debugging using a csi complient client.
		entries[i].Volume.VolumeContext = osdVolumeContext(v)
	}

	return &csi.ListVolumesResponse{
		Entries: entries,
	}, nil
}

// osdVolumeContext returns the attributes of a volume as a map
// to be returned to the CSI API caller
func osdVolumeContext(v *api.Volume) map[string]string {
	return map[string]string{
		api.SpecParent: v.GetSource().GetParent(),
		api.SpecSecure: fmt.Sprintf("%v", v.GetSpec().GetEncrypted()),
		api.SpecShared: fmt.Sprintf("%v", v.GetSpec().GetShared()),
		"readonly":     fmt.Sprintf("%v", v.GetReadonly()),
		"attached":     v.AttachedState.String(),
		"state":        v.State.String(),
		"error":        v.GetError(),
	}
}

// CreateVolume is a CSI API which creates a volume on OSD
// This function supports snapshots if the parent volume id is supplied
// in the parameters.
func (s *OsdCsiServer) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest,
) (*csi.CreateVolumeResponse, error) {

	// Log request
	logrus.Debugf("CreateVolume req[%#v]", *req)

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name must be provided")
	}
	if req.GetVolumeCapabilities() == nil || len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities must be provided")
	}

	// Get parameters
	spec, locator, source, err := s.specHandler.SpecFromOpts(req.GetParameters())
	if err != nil {
		e := fmt.Sprintf("Unable to get parameters: %s\n", err.Error())
		logrus.Errorln(e)
		return nil, status.Error(codes.InvalidArgument, e)
	}

	// Get Size
	if req.GetCapacityRange() != nil && req.GetCapacityRange().GetRequiredBytes() != 0 {
		spec.Size = uint64(req.GetCapacityRange().GetRequiredBytes())
	} else {
		spec.Size = defaultCSIVolumeSize
	}

	// Get grpc connection
	conn, err := s.getConn()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to connect to SDK server: %v", err)
	}

	// Get secret if any was passed
	ctx := s.setupContextWithToken(ctx, req.GetSecrets())

	// Check ID is valid with the specified volume capabilities
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Create volume
	createResp, err := volumes.Create(ctx, &api.SdkVolumeCreateRequest{
		Name:   locator.GetName(),
		Spec:   spec,
		Labels: locator.GetVolumeLabels(),
	})
	if err != nil {
		return nil, err
	}

	// Get volume information
	inspectResp, err := volumes.Inspect(ctx, &api.SdkVolumeInspectRequest{
		VolumeId: createResp.GetVolumeId(),
	})
	if err != nil {
		return nil, err
	}

	// Create response
	volume := &csi.Volume{}
	osdToCsiVolumeInfo(volume, inspectResp.GetVolume())
	return &csi.CreateVolumeResponse{
		Volume: volume,
	}, nil
}

// DeleteVolume is a CSI API which deletes a volume
func (s *OsdCsiServer) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest,
) (*csi.DeleteVolumeResponse, error) {

	// Log request
	logrus.Debugf("DeleteVolume req[%#v]", *req)

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume id must be provided")
	}

	// Get grpc connection
	conn, err := s.getConn()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to connect to SDK server: %v", err)
	}

	// Get secret if any was passed
	ctx := s.setupContextWithToken(ctx, req.GetSecrets())

	// Check ID is valid with the specified volume capabilities
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Delete volume
	_, err := volumes.Delete(ctx, &api.SdkVolumeDeleteRequest{
		VolumeId: req.GetVolumeId(),
	})

	return &csi.DeleteVolumeResponse{}, nil
}

func osdToCsiVolumeInfo(dest *csi.Volume, src *api.Volume) {
	dest.VolumeId = src.GetId()
	dest.CapacityBytes = int64(src.Spec.GetSize())
	dest.VolumeContext = osdVolumeContext(src)
}

func csiRequestsSharedVolume(req *csi.CreateVolumeRequest) bool {
	for _, cap := range req.GetVolumeCapabilities() {
		// Check access mode is setup correctly
		mode := cap.GetAccessMode().GetMode()
		if mode == csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER ||
			mode == csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER {
			return true
		}
	}
	return false
}

// CreateSnapshot is a CSI implementation to create a snapshot from the volume
func (s *OsdCsiServer) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest,
) (*csi.CreateSnapshotResponse, error) {

	if len(req.GetSourceVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume id must be provided")
	} else if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name must be provided")
	}

	// Check if the snapshot with this name already exists
	v, err := util.VolumeFromName(s.driver, req.GetName())
	if err == nil {
		// Verify the parent is the same
		if req.GetSourceVolumeId() != v.GetSource().GetParent() {
			return nil, status.Error(codes.AlreadyExists, "Requested snapshot already exists for another source volume id")
		}

		return &csi.CreateSnapshotResponse{
			Snapshot: &csi.Snapshot{
				SnapshotId:     v.GetId(),
				SourceVolumeId: v.GetSource().GetParent(),
				CreationTime:   v.GetCtime(),
				ReadyToUse:     true,
			},
		}, nil
	}

	// Get any labels passed in by the CO
	_, locator, _, err := s.specHandler.SpecFromOpts(req.GetParameters())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to get parameters: %v", err)
	}

	// Create snapshot
	readonly := true
	snapshotID, err := s.driver.Snapshot(req.GetSourceVolumeId(), readonly, &api.VolumeLocator{
		Name:         req.GetName(),
		VolumeLabels: locator.GetVolumeLabels(),
	}, false)
	if err != nil {
		if err == kvdb.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "Volume id %s not found", req.GetSourceVolumeId())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create snapshot: %v", err)
	}

	snapInfo, err := util.VolumeFromName(s.driver, snapshotID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get information about the snapshot: %v", err)
	}

	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     snapshotID,
			SourceVolumeId: req.GetSourceVolumeId(),
			CreationTime:   snapInfo.GetCtime(),
			ReadyToUse:     true,
		},
	}, nil
}

// DeleteSnapshot is a CSI implementation to delete a snapshot
func (s *OsdCsiServer) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest,
) (*csi.DeleteSnapshotResponse, error) {

	if len(req.GetSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Snapshot id must be provided")
	}

	// If the snapshot is not found, then we can return OK
	volumes, err := s.driver.Inspect([]string{req.GetSnapshotId()})
	if (err == nil && len(volumes) == 0) ||
		(err != nil && err == kvdb.ErrNotFound) {
		return &csi.DeleteSnapshotResponse{}, nil
	} else if err != nil {
		return nil, err
	}

	err = s.driver.Delete(req.GetSnapshotId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to delete snapshot %s: %v",
			req.GetSnapshotId(),
			err)
	}

	return &csi.DeleteSnapshotResponse{}, nil
}

// ListSnapshots is not supported (we can add this later)
func (s *OsdCsiServer) ListSnapshots(
	ctx context.Context,
	req *csi.ListSnapshotsRequest,
) (*csi.ListSnapshotsResponse, error) {

	// The function ListSnapshots is also not published as
	// supported by this implementation
	return nil, status.Error(codes.Unimplemented, "ListSnapshots is not implemented")
}
