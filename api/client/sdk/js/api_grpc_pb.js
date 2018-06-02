// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var api_pb = require('./api_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');
var google_api_annotations_pb = require('./google/api/annotations_pb.js');

function serialize_openstorage_api_ClusterAlertClearRequest(arg) {
  if (!(arg instanceof api_pb.ClusterAlertClearRequest)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertClearRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertClearRequest(buffer_arg) {
  return api_pb.ClusterAlertClearRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterAlertClearResponse(arg) {
  if (!(arg instanceof api_pb.ClusterAlertClearResponse)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertClearResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertClearResponse(buffer_arg) {
  return api_pb.ClusterAlertClearResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterAlertEnumerateRequest(arg) {
  if (!(arg instanceof api_pb.ClusterAlertEnumerateRequest)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertEnumerateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertEnumerateRequest(buffer_arg) {
  return api_pb.ClusterAlertEnumerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterAlertEnumerateResponse(arg) {
  if (!(arg instanceof api_pb.ClusterAlertEnumerateResponse)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertEnumerateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertEnumerateResponse(buffer_arg) {
  return api_pb.ClusterAlertEnumerateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterAlertEraseRequest(arg) {
  if (!(arg instanceof api_pb.ClusterAlertEraseRequest)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertEraseRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertEraseRequest(buffer_arg) {
  return api_pb.ClusterAlertEraseRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterAlertEraseResponse(arg) {
  if (!(arg instanceof api_pb.ClusterAlertEraseResponse)) {
    throw new Error('Expected argument of type openstorage.api.ClusterAlertEraseResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterAlertEraseResponse(buffer_arg) {
  return api_pb.ClusterAlertEraseResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterEnumerateRequest(arg) {
  if (!(arg instanceof api_pb.ClusterEnumerateRequest)) {
    throw new Error('Expected argument of type openstorage.api.ClusterEnumerateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterEnumerateRequest(buffer_arg) {
  return api_pb.ClusterEnumerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterEnumerateResponse(arg) {
  if (!(arg instanceof api_pb.ClusterEnumerateResponse)) {
    throw new Error('Expected argument of type openstorage.api.ClusterEnumerateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterEnumerateResponse(buffer_arg) {
  return api_pb.ClusterEnumerateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterInspectRequest(arg) {
  if (!(arg instanceof api_pb.ClusterInspectRequest)) {
    throw new Error('Expected argument of type openstorage.api.ClusterInspectRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterInspectRequest(buffer_arg) {
  return api_pb.ClusterInspectRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_ClusterInspectResponse(arg) {
  if (!(arg instanceof api_pb.ClusterInspectResponse)) {
    throw new Error('Expected argument of type openstorage.api.ClusterInspectResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_ClusterInspectResponse(buffer_arg) {
  return api_pb.ClusterInspectResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateAWSRequest(arg) {
  if (!(arg instanceof api_pb.CredentialCreateAWSRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateAWSRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateAWSRequest(buffer_arg) {
  return api_pb.CredentialCreateAWSRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateAWSResponse(arg) {
  if (!(arg instanceof api_pb.CredentialCreateAWSResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateAWSResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateAWSResponse(buffer_arg) {
  return api_pb.CredentialCreateAWSResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateAzureRequest(arg) {
  if (!(arg instanceof api_pb.CredentialCreateAzureRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateAzureRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateAzureRequest(buffer_arg) {
  return api_pb.CredentialCreateAzureRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateAzureResponse(arg) {
  if (!(arg instanceof api_pb.CredentialCreateAzureResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateAzureResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateAzureResponse(buffer_arg) {
  return api_pb.CredentialCreateAzureResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateGoogleRequest(arg) {
  if (!(arg instanceof api_pb.CredentialCreateGoogleRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateGoogleRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateGoogleRequest(buffer_arg) {
  return api_pb.CredentialCreateGoogleRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialCreateGoogleResponse(arg) {
  if (!(arg instanceof api_pb.CredentialCreateGoogleResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialCreateGoogleResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialCreateGoogleResponse(buffer_arg) {
  return api_pb.CredentialCreateGoogleResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialDeleteRequest(arg) {
  if (!(arg instanceof api_pb.CredentialDeleteRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialDeleteRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialDeleteRequest(buffer_arg) {
  return api_pb.CredentialDeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialDeleteResponse(arg) {
  if (!(arg instanceof api_pb.CredentialDeleteResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialDeleteResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialDeleteResponse(buffer_arg) {
  return api_pb.CredentialDeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateAWSRequest(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateAWSRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateAWSRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateAWSRequest(buffer_arg) {
  return api_pb.CredentialEnumerateAWSRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateAWSResponse(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateAWSResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateAWSResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateAWSResponse(buffer_arg) {
  return api_pb.CredentialEnumerateAWSResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateAzureRequest(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateAzureRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateAzureRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateAzureRequest(buffer_arg) {
  return api_pb.CredentialEnumerateAzureRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateAzureResponse(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateAzureResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateAzureResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateAzureResponse(buffer_arg) {
  return api_pb.CredentialEnumerateAzureResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateGoogleRequest(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateGoogleRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateGoogleRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateGoogleRequest(buffer_arg) {
  return api_pb.CredentialEnumerateGoogleRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialEnumerateGoogleResponse(arg) {
  if (!(arg instanceof api_pb.CredentialEnumerateGoogleResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialEnumerateGoogleResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialEnumerateGoogleResponse(buffer_arg) {
  return api_pb.CredentialEnumerateGoogleResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialValidateRequest(arg) {
  if (!(arg instanceof api_pb.CredentialValidateRequest)) {
    throw new Error('Expected argument of type openstorage.api.CredentialValidateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialValidateRequest(buffer_arg) {
  return api_pb.CredentialValidateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_CredentialValidateResponse(arg) {
  if (!(arg instanceof api_pb.CredentialValidateResponse)) {
    throw new Error('Expected argument of type openstorage.api.CredentialValidateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_CredentialValidateResponse(buffer_arg) {
  return api_pb.CredentialValidateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_OpenStorageVolumeCreateRequest(arg) {
  if (!(arg instanceof api_pb.OpenStorageVolumeCreateRequest)) {
    throw new Error('Expected argument of type openstorage.api.OpenStorageVolumeCreateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_OpenStorageVolumeCreateRequest(buffer_arg) {
  return api_pb.OpenStorageVolumeCreateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_OpenStorageVolumeCreateResponse(arg) {
  if (!(arg instanceof api_pb.OpenStorageVolumeCreateResponse)) {
    throw new Error('Expected argument of type openstorage.api.OpenStorageVolumeCreateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_OpenStorageVolumeCreateResponse(buffer_arg) {
  return api_pb.OpenStorageVolumeCreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeAttachRequest(arg) {
  if (!(arg instanceof api_pb.VolumeAttachRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeAttachRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeAttachRequest(buffer_arg) {
  return api_pb.VolumeAttachRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeAttachResponse(arg) {
  if (!(arg instanceof api_pb.VolumeAttachResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeAttachResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeAttachResponse(buffer_arg) {
  return api_pb.VolumeAttachResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeCreateFromVolumeIDRequest(arg) {
  if (!(arg instanceof api_pb.VolumeCreateFromVolumeIDRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeCreateFromVolumeIDRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeCreateFromVolumeIDRequest(buffer_arg) {
  return api_pb.VolumeCreateFromVolumeIDRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeCreateFromVolumeIDResponse(arg) {
  if (!(arg instanceof api_pb.VolumeCreateFromVolumeIDResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeCreateFromVolumeIDResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeCreateFromVolumeIDResponse(buffer_arg) {
  return api_pb.VolumeCreateFromVolumeIDResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeDeleteRequest(arg) {
  if (!(arg instanceof api_pb.VolumeDeleteRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeDeleteRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeDeleteRequest(buffer_arg) {
  return api_pb.VolumeDeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeDeleteResponse(arg) {
  if (!(arg instanceof api_pb.VolumeDeleteResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeDeleteResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeDeleteResponse(buffer_arg) {
  return api_pb.VolumeDeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeDetachRequest(arg) {
  if (!(arg instanceof api_pb.VolumeDetachRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeDetachRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeDetachRequest(buffer_arg) {
  return api_pb.VolumeDetachRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeDetachResponse(arg) {
  if (!(arg instanceof api_pb.VolumeDetachResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeDetachResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeDetachResponse(buffer_arg) {
  return api_pb.VolumeDetachResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeEnumerateRequest(arg) {
  if (!(arg instanceof api_pb.VolumeEnumerateRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeEnumerateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeEnumerateRequest(buffer_arg) {
  return api_pb.VolumeEnumerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeEnumerateResponse(arg) {
  if (!(arg instanceof api_pb.VolumeEnumerateResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeEnumerateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeEnumerateResponse(buffer_arg) {
  return api_pb.VolumeEnumerateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeInspectRequest(arg) {
  if (!(arg instanceof api_pb.VolumeInspectRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeInspectRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeInspectRequest(buffer_arg) {
  return api_pb.VolumeInspectRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeInspectResponse(arg) {
  if (!(arg instanceof api_pb.VolumeInspectResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeInspectResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeInspectResponse(buffer_arg) {
  return api_pb.VolumeInspectResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeMountRequest(arg) {
  if (!(arg instanceof api_pb.VolumeMountRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeMountRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeMountRequest(buffer_arg) {
  return api_pb.VolumeMountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeMountResponse(arg) {
  if (!(arg instanceof api_pb.VolumeMountResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeMountResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeMountResponse(buffer_arg) {
  return api_pb.VolumeMountResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotCreateRequest(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotCreateRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotCreateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotCreateRequest(buffer_arg) {
  return api_pb.VolumeSnapshotCreateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotCreateResponse(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotCreateResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotCreateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotCreateResponse(buffer_arg) {
  return api_pb.VolumeSnapshotCreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotEnumerateRequest(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotEnumerateRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotEnumerateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotEnumerateRequest(buffer_arg) {
  return api_pb.VolumeSnapshotEnumerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotEnumerateResponse(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotEnumerateResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotEnumerateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotEnumerateResponse(buffer_arg) {
  return api_pb.VolumeSnapshotEnumerateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotRestoreRequest(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotRestoreRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotRestoreRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotRestoreRequest(buffer_arg) {
  return api_pb.VolumeSnapshotRestoreRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeSnapshotRestoreResponse(arg) {
  if (!(arg instanceof api_pb.VolumeSnapshotRestoreResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeSnapshotRestoreResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeSnapshotRestoreResponse(buffer_arg) {
  return api_pb.VolumeSnapshotRestoreResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeUnmountRequest(arg) {
  if (!(arg instanceof api_pb.VolumeUnmountRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeUnmountRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeUnmountRequest(buffer_arg) {
  return api_pb.VolumeUnmountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeUnmountResponse(arg) {
  if (!(arg instanceof api_pb.VolumeUnmountResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeUnmountResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeUnmountResponse(buffer_arg) {
  return api_pb.VolumeUnmountResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var OpenStorageClusterService = exports.OpenStorageClusterService = {
  // Enumerate lists all the nodes in the cluster.
  enumerate: {
    path: '/openstorage.api.OpenStorageCluster/Enumerate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.ClusterEnumerateRequest,
    responseType: api_pb.ClusterEnumerateResponse,
    requestSerialize: serialize_openstorage_api_ClusterEnumerateRequest,
    requestDeserialize: deserialize_openstorage_api_ClusterEnumerateRequest,
    responseSerialize: serialize_openstorage_api_ClusterEnumerateResponse,
    responseDeserialize: deserialize_openstorage_api_ClusterEnumerateResponse,
  },
  // Inspect the node given a UUID.
  inspect: {
    path: '/openstorage.api.OpenStorageCluster/Inspect',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.ClusterInspectRequest,
    responseType: api_pb.ClusterInspectResponse,
    requestSerialize: serialize_openstorage_api_ClusterInspectRequest,
    requestDeserialize: deserialize_openstorage_api_ClusterInspectRequest,
    responseSerialize: serialize_openstorage_api_ClusterInspectResponse,
    responseDeserialize: deserialize_openstorage_api_ClusterInspectResponse,
  },
  // Get a list of alerts from the storage cluster
  alertEnumerate: {
    path: '/openstorage.api.OpenStorageCluster/AlertEnumerate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.ClusterAlertEnumerateRequest,
    responseType: api_pb.ClusterAlertEnumerateResponse,
    requestSerialize: serialize_openstorage_api_ClusterAlertEnumerateRequest,
    requestDeserialize: deserialize_openstorage_api_ClusterAlertEnumerateRequest,
    responseSerialize: serialize_openstorage_api_ClusterAlertEnumerateResponse,
    responseDeserialize: deserialize_openstorage_api_ClusterAlertEnumerateResponse,
  },
  // Clear the alert for a given resource
  alertClear: {
    path: '/openstorage.api.OpenStorageCluster/AlertClear',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.ClusterAlertClearRequest,
    responseType: api_pb.ClusterAlertClearResponse,
    requestSerialize: serialize_openstorage_api_ClusterAlertClearRequest,
    requestDeserialize: deserialize_openstorage_api_ClusterAlertClearRequest,
    responseSerialize: serialize_openstorage_api_ClusterAlertClearResponse,
    responseDeserialize: deserialize_openstorage_api_ClusterAlertClearResponse,
  },
  // Erases an alert for a given resource
  alertErase: {
    path: '/openstorage.api.OpenStorageCluster/AlertErase',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.ClusterAlertEraseRequest,
    responseType: api_pb.ClusterAlertEraseResponse,
    requestSerialize: serialize_openstorage_api_ClusterAlertEraseRequest,
    requestDeserialize: deserialize_openstorage_api_ClusterAlertEraseRequest,
    responseSerialize: serialize_openstorage_api_ClusterAlertEraseResponse,
    responseDeserialize: deserialize_openstorage_api_ClusterAlertEraseResponse,
  },
};

exports.OpenStorageClusterClient = grpc.makeGenericClientConstructor(OpenStorageClusterService);
var OpenStorageVolumeService = exports.OpenStorageVolumeService = {
  // Creates a new volume
  create: {
    path: '/openstorage.api.OpenStorageVolume/Create',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.OpenStorageVolumeCreateRequest,
    responseType: api_pb.OpenStorageVolumeCreateResponse,
    requestSerialize: serialize_openstorage_api_OpenStorageVolumeCreateRequest,
    requestDeserialize: deserialize_openstorage_api_OpenStorageVolumeCreateRequest,
    responseSerialize: serialize_openstorage_api_OpenStorageVolumeCreateResponse,
    responseDeserialize: deserialize_openstorage_api_OpenStorageVolumeCreateResponse,
  },
  // CreateFromVolumeID creates a new volume cloned from an existing volume
  createFromVolumeID: {
    path: '/openstorage.api.OpenStorageVolume/CreateFromVolumeID',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeCreateFromVolumeIDRequest,
    responseType: api_pb.VolumeCreateFromVolumeIDResponse,
    requestSerialize: serialize_openstorage_api_VolumeCreateFromVolumeIDRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeCreateFromVolumeIDRequest,
    responseSerialize: serialize_openstorage_api_VolumeCreateFromVolumeIDResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeCreateFromVolumeIDResponse,
  },
  // Delete a volume
  delete: {
    path: '/openstorage.api.OpenStorageVolume/Delete',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeDeleteRequest,
    responseType: api_pb.VolumeDeleteResponse,
    requestSerialize: serialize_openstorage_api_VolumeDeleteRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeDeleteRequest,
    responseSerialize: serialize_openstorage_api_VolumeDeleteResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeDeleteResponse,
  },
  // Get information on a volume
  inspect: {
    path: '/openstorage.api.OpenStorageVolume/Inspect',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeInspectRequest,
    responseType: api_pb.VolumeInspectResponse,
    requestSerialize: serialize_openstorage_api_VolumeInspectRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeInspectRequest,
    responseSerialize: serialize_openstorage_api_VolumeInspectResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeInspectResponse,
  },
  // Get a list of volumes
  enumerate: {
    path: '/openstorage.api.OpenStorageVolume/Enumerate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeEnumerateRequest,
    responseType: api_pb.VolumeEnumerateResponse,
    requestSerialize: serialize_openstorage_api_VolumeEnumerateRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeEnumerateRequest,
    responseSerialize: serialize_openstorage_api_VolumeEnumerateResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeEnumerateResponse,
  },
  // Create a snapshot of a volume. This creates an immutable (read-only),
  // point-in-time snapshot of a volume.
  snapshotCreate: {
    path: '/openstorage.api.OpenStorageVolume/SnapshotCreate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeSnapshotCreateRequest,
    responseType: api_pb.VolumeSnapshotCreateResponse,
    requestSerialize: serialize_openstorage_api_VolumeSnapshotCreateRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeSnapshotCreateRequest,
    responseSerialize: serialize_openstorage_api_VolumeSnapshotCreateResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeSnapshotCreateResponse,
  },
  // Restores a volume to a specified snapshot
  snapshotRestore: {
    path: '/openstorage.api.OpenStorageVolume/SnapshotRestore',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeSnapshotRestoreRequest,
    responseType: api_pb.VolumeSnapshotRestoreResponse,
    requestSerialize: serialize_openstorage_api_VolumeSnapshotRestoreRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeSnapshotRestoreRequest,
    responseSerialize: serialize_openstorage_api_VolumeSnapshotRestoreResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeSnapshotRestoreResponse,
  },
  // List the number of snapshots for a specific volume
  snapshotEnumerate: {
    path: '/openstorage.api.OpenStorageVolume/SnapshotEnumerate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeSnapshotEnumerateRequest,
    responseType: api_pb.VolumeSnapshotEnumerateResponse,
    requestSerialize: serialize_openstorage_api_VolumeSnapshotEnumerateRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeSnapshotEnumerateRequest,
    responseSerialize: serialize_openstorage_api_VolumeSnapshotEnumerateResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeSnapshotEnumerateResponse,
  },
  // Attach device to host                                                      
  attach: {
    path: '/openstorage.api.OpenStorageVolume/Attach',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeAttachRequest,
    responseType: api_pb.VolumeAttachResponse,
    requestSerialize: serialize_openstorage_api_VolumeAttachRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeAttachRequest,
    responseSerialize: serialize_openstorage_api_VolumeAttachResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeAttachResponse,
  },
  // Detaches the volume from the node.
  detach: {
    path: '/openstorage.api.OpenStorageVolume/Detach',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeDetachRequest,
    responseType: api_pb.VolumeDetachResponse,
    requestSerialize: serialize_openstorage_api_VolumeDetachRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeDetachRequest,
    responseSerialize: serialize_openstorage_api_VolumeDetachResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeDetachResponse,
  },
  // Attaches the volume to a node.
  mount: {
    path: '/openstorage.api.OpenStorageVolume/Mount',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeMountRequest,
    responseType: api_pb.VolumeMountResponse,
    requestSerialize: serialize_openstorage_api_VolumeMountRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeMountRequest,
    responseSerialize: serialize_openstorage_api_VolumeMountResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeMountResponse,
  },
  // Unmount volume at specified path                                           
  unmount: {
    path: '/openstorage.api.OpenStorageVolume/Unmount',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeUnmountRequest,
    responseType: api_pb.VolumeUnmountResponse,
    requestSerialize: serialize_openstorage_api_VolumeUnmountRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeUnmountRequest,
    responseSerialize: serialize_openstorage_api_VolumeUnmountResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeUnmountResponse,
  },
};

exports.OpenStorageVolumeClient = grpc.makeGenericClientConstructor(OpenStorageVolumeService);
var OpenStorageCredentialsService = exports.OpenStorageCredentialsService = {
  // Provide credentials to OpenStorage and if valid,
  // it will return an identifier to the credentials
  //
  // Create credential for AWS S3 and if valid ,
  // returns a unique identifier
  createForAWS: {
    path: '/openstorage.api.OpenStorageCredentials/CreateForAWS',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialCreateAWSRequest,
    responseType: api_pb.CredentialCreateAWSResponse,
    requestSerialize: serialize_openstorage_api_CredentialCreateAWSRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialCreateAWSRequest,
    responseSerialize: serialize_openstorage_api_CredentialCreateAWSResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialCreateAWSResponse,
  },
  // Create credential for Azure and if valid ,
  // returns a unique identifier
  createForAzure: {
    path: '/openstorage.api.OpenStorageCredentials/CreateForAzure',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialCreateAzureRequest,
    responseType: api_pb.CredentialCreateAzureResponse,
    requestSerialize: serialize_openstorage_api_CredentialCreateAzureRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialCreateAzureRequest,
    responseSerialize: serialize_openstorage_api_CredentialCreateAzureResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialCreateAzureResponse,
  },
  // Create credential for Google and if valid ,
  // returns a unique identifier
  createForGoogle: {
    path: '/openstorage.api.OpenStorageCredentials/CreateForGoogle',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialCreateGoogleRequest,
    responseType: api_pb.CredentialCreateGoogleResponse,
    requestSerialize: serialize_openstorage_api_CredentialCreateGoogleRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialCreateGoogleRequest,
    responseSerialize: serialize_openstorage_api_CredentialCreateGoogleResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialCreateGoogleResponse,
  },
  // EnumerateForAWS lists the configured AWS credentials                      
  enumerateForAWS: {
    path: '/openstorage.api.OpenStorageCredentials/EnumerateForAWS',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialEnumerateAWSRequest,
    responseType: api_pb.CredentialEnumerateAWSResponse,
    requestSerialize: serialize_openstorage_api_CredentialEnumerateAWSRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialEnumerateAWSRequest,
    responseSerialize: serialize_openstorage_api_CredentialEnumerateAWSResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialEnumerateAWSResponse,
  },
  // EnumerateForAzure lists the configured Azure credentials                  
  enumerateForAzure: {
    path: '/openstorage.api.OpenStorageCredentials/EnumerateForAzure',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialEnumerateAzureRequest,
    responseType: api_pb.CredentialEnumerateAzureResponse,
    requestSerialize: serialize_openstorage_api_CredentialEnumerateAzureRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialEnumerateAzureRequest,
    responseSerialize: serialize_openstorage_api_CredentialEnumerateAzureResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialEnumerateAzureResponse,
  },
  // EnumerateForGoogle lists the configured Google credentials                
  enumerateForGoogle: {
    path: '/openstorage.api.OpenStorageCredentials/EnumerateForGoogle',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialEnumerateGoogleRequest,
    responseType: api_pb.CredentialEnumerateGoogleResponse,
    requestSerialize: serialize_openstorage_api_CredentialEnumerateGoogleRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialEnumerateGoogleRequest,
    responseSerialize: serialize_openstorage_api_CredentialEnumerateGoogleResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialEnumerateGoogleResponse,
  },
  // Delete a specified credential                                                 
  credentialDelete: {
    path: '/openstorage.api.OpenStorageCredentials/CredentialDelete',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialDeleteRequest,
    responseType: api_pb.CredentialDeleteResponse,
    requestSerialize: serialize_openstorage_api_CredentialDeleteRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialDeleteRequest,
    responseSerialize: serialize_openstorage_api_CredentialDeleteResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialDeleteResponse,
  },
  // Validate a specified credential
  credentialValidate: {
    path: '/openstorage.api.OpenStorageCredentials/CredentialValidate',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.CredentialValidateRequest,
    responseType: api_pb.CredentialValidateResponse,
    requestSerialize: serialize_openstorage_api_CredentialValidateRequest,
    requestDeserialize: deserialize_openstorage_api_CredentialValidateRequest,
    responseSerialize: serialize_openstorage_api_CredentialValidateResponse,
    responseDeserialize: deserialize_openstorage_api_CredentialValidateResponse,
  },
};

exports.OpenStorageCredentialsClient = grpc.makeGenericClientConstructor(OpenStorageCredentialsService);
