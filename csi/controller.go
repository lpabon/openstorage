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

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/pkg/util"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.pedge.io/dlog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	volumeCapabilityMessageMultinodeVolume    = "Volume is a multinode volume"
	volumeCapabilityMessageNotMultinodeVolume = "Volume is not a multinode volume"
	volumeCapabilityMessageReadOnlyVolume     = "Volume is read only"
	volumeCapabilityMessageNotReadOnlyVolume  = "Volume is not read only"
)

// ControllerGetCapabilities is a CSI API functions which returns to the caller
// the capabilities of the OSD CSI driver.
func (s *OsdCsiServer) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest,
) (*csi.ControllerGetCapabilitiesResponse, error) {

	// Volume capabilities
	capCreateDeleteVolume := &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			},
		},
	}
	capListVolumes := &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
			},
		},
	}

	//
	// Here we will be adding the following in future patches:
	// List Volumes Capabability
	// Get Storage Capacity
	//

	return &csi.ControllerGetCapabilitiesResponse{
		Reply: &csi.ControllerGetCapabilitiesResponse_Result_{
			Result: &csi.ControllerGetCapabilitiesResponse_Result{
				Capabilities: []*csi.ControllerServiceCapability{
					capCreateDeleteVolume,
					capListVolumes,
				},
			},
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
// Note: The method used here to return errors is still not part of the spec.
// See: https://github.com/container-storage-interface/spec/pull/115
// Discussion:  https://groups.google.com/forum/#!topic/kubernetes-sig-storage-wg-csi/TpTrNFbRa1I
//
func (s *OsdCsiServer) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest,
) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	// Probably we may use version in the future, but for now, let's just log it
	version := req.GetVersion()
	if version == nil {
		return nil, status.Error(codes.InvalidArgument, "Version must be specified")
	}
	capabilities := req.GetVolumeCapabilities()
	if capabilities == nil || len(capabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume_capabilities must be specified")
	}
	id := req.GetVolumeId()
	if len(id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume_id must be specified")
	}
	attributes := req.GetVolumeAttributes()

	// Log request
	dlog.Debugf("ValidateVolumeCapabilities of id %s "+
		"capabilities %#v "+
		"version %#v "+
		"attributes %#v ",
		id,
		capabilities,
		version,
		attributes)

	// Check ID is valid with the specified volume capabilities
	volumes, err := s.driver.Inspect([]string{id})
	if err != nil || len(volumes) == 0 {
		return nil, status.Error(codes.NotFound, "ID not found")
	}
	if len(volumes) != 1 {
		errs := fmt.Sprintf(
			"Driver returned an unexpected number of volumes when one was expected: %d",
			len(volumes))
		dlog.Errorln(errs)
		return nil, status.Error(codes.Internal, errs)
	}
	v := volumes[0]
	if v.Id != id {
		errs := fmt.Sprintf(
			"Driver volume id [%s] does not equal requested id of: %s",
			v.Id,
			id)
		dlog.Errorln(errs)
		return nil, status.Error(codes.Internal, errs)
	}

	// Setup uninitialized response object
	result := &csi.ValidateVolumeCapabilitiesResponse_Result{
		Supported: true,
	}
	resp := &csi.ValidateVolumeCapabilitiesResponse{
		Reply: &csi.ValidateVolumeCapabilitiesResponse_Result_{
			Result: result,
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
				result.Supported = false
				result.Message = volumeCapabilityMessageMultinodeVolume
				break
			}
			if v.Readonly {
				result.Supported = false
				result.Message = volumeCapabilityMessageReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY:
			if v.Spec.Shared {
				result.Supported = false
				result.Message = volumeCapabilityMessageMultinodeVolume
				break
			}
			if !v.Readonly {
				result.Supported = false
				result.Message = volumeCapabilityMessageNotReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY:
			if !v.Spec.Shared {
				result.Supported = false
				result.Message = volumeCapabilityMessageNotMultinodeVolume
				break
			}
			if !v.Readonly {
				result.Supported = false
				result.Message = volumeCapabilityMessageNotReadOnlyVolume
				break
			}
		case mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER ||
			mode.Mode == csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER:
			if !v.Spec.Shared {
				result.Supported = false
				result.Message = volumeCapabilityMessageNotMultinodeVolume
				break
			}
			if v.Readonly {
				result.Supported = false
				result.Message = volumeCapabilityMessageReadOnlyVolume
				break
			}
		default:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"AccessMode %s is not allowed",
				mode.Mode.String())
		}

		if !result.Supported {
			return resp, nil
		}
	}

	// If we passed all the checks, then it is valid
	result.Message = "Volume is supported"
	return resp, nil
}

// ListVolumes is a CSI API which returns to the caller all volume ids
// on this cluster. This includes ids created by CSI and ids created
// using other interfaces. This is important because the user could
// be requesting to mount a OSD volume created using non-CSI interfaces.
//
// This call does not yet implement tokens to due the following
// issue: https://github.com/container-storage-interface/spec/issues/138
func (s *OsdCsiServer) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest,
) (*csi.ListVolumesResponse, error) {

	// Future: Once CSI is released, check version
	// for now, just log it.
	dlog.Debugf("ListVolumes req[%#v]", req)

	// REQUIRE: Must at least that version is provided
	if req.GetVersion() == nil {
		return nil, status.Error(codes.InvalidArgument, "Version must be provided")
	}

	// Until the issue #138 on the CSI spec is resolved we will not support
	// tokenization
	if req.GetMaxEntries() != 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Driver does not support tokenization. Please see "+
				"https://github.com/container-storage-interface/spec/issues/138")
	}

	volumes, err := s.driver.Enumerate(&api.VolumeLocator{}, nil)
	if err != nil {
		errs := fmt.Sprintf("Unable to get list of volumes: %s", err.Error())
		dlog.Errorln(errs)
		return nil, status.Error(codes.Internal, errs)
	}
	entries := make([]*csi.ListVolumesResponse_Result_Entry, len(volumes))
	for i, v := range volumes {
		// Initialize entry
		entries[i] = &csi.ListVolumesResponse_Result_Entry{
			VolumeInfo: &csi.VolumeInfo{},
		}

		// Required
		entries[i].VolumeInfo.Id = v.Id

		// This entry is optional in the API, but OSD has
		// the information available to provide it
		entries[i].VolumeInfo.CapacityBytes = v.Spec.Size

		// Attributes. We can add or remove as needed since they
		// are optional and opaque to the CO but could be used
		// for debuggig using a csi complient client.
		entries[i].VolumeInfo.Attributes = map[string]string{
			"readonly": fmt.Sprintf("%v", v.Readonly),
			"shared":   fmt.Sprintf("%v", v.Spec.Shared),
			"state":    v.State.String(),
			"attached": v.AttachedState.String(),
			"error":    v.Error,
		}
	}

	return &csi.ListVolumesResponse{
		Reply: &csi.ListVolumesResponse_Result_{
			Result: &csi.ListVolumesResponse_Result{
				Entries: entries,
			},
		},
	}, nil

}

// CreateVolume is a CSI API which creates a volume on OSD
// This function supports snapshots if the parent volume id is supplied
// in the parameters.
func (s *OsdCsiServer) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest,
) (*csi.CreateVolumeResponse, error) {

	// Log request
	dlog.Debugf("CreateVolume req[%#v]", *req)

	// Check arguments
	if req.GetVersion() == nil {
		return nil, status.Error(codes.InvalidArgument, "Version must be provided")
	}
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name must be provided")
	}
	if req.GetVolumeCapabilities() == nil || len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities must be provided")
	}
	if req.GetCapacityRange() == nil {
		return nil, status.Error(codes.InvalidArgument, "CapacityRange must be specified")
	}
	if req.GetCapacityRange().GetRequiredBytes() == 0 {
		return nil, status.Error(codes.Internal, "Capacity range required bytes cannot be zero")
	}

	// Create response
	volume := &csi.VolumeInfo{}
	resp := &csi.CreateVolumeResponse{
		Reply: &csi.CreateVolumeResponse_Result_{
			Result: &csi.CreateVolumeResponse_Result{
				VolumeInfo: volume,
			},
		},
	}

	// Check if the volume has already been created or is in process of creation
	v, err := util.VolumeFromName(s.driver, req.GetName())
	if err == nil {
		osdToCsiVolumeInfo(volume, v)
		return resp, nil
	}

	// Get parameters
	spec, locator, source, err := s.specHandler.SpecFromOpts(req.GetParameters())
	if err != nil {
		e := fmt.Sprintf("Unable to get parameters: %s\n", err.Error())
		dlog.Errorln(e)
		return nil, status.Error(codes.Internal, e)
	}

	// Check if the caller is asking to create a snapshot or for a new volume
	var id string
	if source != nil && len(source.GetParent()) != 0 {
		// Get parent volume information
		parent, err := util.VolumeFromName(s.driver, source.Parent)
		if err != nil {
			e := fmt.Sprintf("unable to get parent volume information: %s\n", err.Error())
			dlog.Errorln(e)
			return nil, status.Error(codes.InvalidArgument, e)
		}

		// Create a snapshot from the parent
		id, err = s.driver.Snapshot(parent.GetId(), false, &api.VolumeLocator{
			Name: req.GetName(),
		})
		if err != nil {
			e := fmt.Sprintf("unable to create snapshot: %s\n", err.Error())
			dlog.Errorln(e)
			return nil, status.Error(codes.Internal, e)
		}
	} else {
		// Get Capabilities and Size
		spec.Size = req.GetCapacityRange().GetRequiredBytes()
		for _, cap := range req.GetVolumeCapabilities() {
			// Check access mode is setup correctly
			mode := cap.GetAccessMode().GetMode()
			if mode == csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER ||
				mode == csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER {
				spec.Shared = true
			}
		}

		// Create the volume
		locator.Name = req.GetName()
		id, err = s.driver.Create(locator, source, spec)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// id must have been set
	v, err = util.VolumeFromName(s.driver, id)
	if err != nil {
		e := fmt.Sprintf("Unable to find newly created volume: %s", err.Error())
		dlog.Errorln(e)
		return nil, status.Error(codes.Internal, e)
	}
	osdToCsiVolumeInfo(volume, v)
	return resp, nil

}

func osdToCsiVolumeInfo(dest *csi.VolumeInfo, src *api.Volume) {
	dest.Id = src.GetId()
	dest.CapacityBytes = src.Spec.GetSize()
}

/*
For next patches what still needs to be worked on in the Conroller server:

	CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)
	DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)
	GetCapacity(context.Context, *GetCapacityRequest) (*GetCapacityResponse, error)
	ControllerGetCapabilities(context.Context, *ControllerGetCapabilitiesRequest) (*ControllerGetCapabilitiesResponse, error)
*/
