// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var api_pb = require('./api_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');

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

function serialize_openstorage_api_VolumeCreateRequest(arg) {
  if (!(arg instanceof api_pb.VolumeCreateRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeCreateRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeCreateRequest(buffer_arg) {
  return api_pb.VolumeCreateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeCreateResponse(arg) {
  if (!(arg instanceof api_pb.VolumeCreateResponse)) {
    throw new Error('Expected argument of type openstorage.api.VolumeCreateResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeCreateResponse(buffer_arg) {
  return api_pb.VolumeCreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_openstorage_api_VolumeCreateSimpleVolumeRequest(arg) {
  if (!(arg instanceof api_pb.VolumeCreateSimpleVolumeRequest)) {
    throw new Error('Expected argument of type openstorage.api.VolumeCreateSimpleVolumeRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_openstorage_api_VolumeCreateSimpleVolumeRequest(buffer_arg) {
  return api_pb.VolumeCreateSimpleVolumeRequest.deserializeBinary(new Uint8Array(buffer_arg));
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
  // Creates a volume
  create: {
    path: '/openstorage.api.OpenStorageVolume/Create',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeCreateRequest,
    responseType: api_pb.VolumeCreateResponse,
    requestSerialize: serialize_openstorage_api_VolumeCreateRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeCreateRequest,
    responseSerialize: serialize_openstorage_api_VolumeCreateResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeCreateResponse,
  },
  // CreateSimpleVolume provides a simple API to create a volume
  createSimpleVolume: {
    path: '/openstorage.api.OpenStorageVolume/CreateSimpleVolume',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.VolumeCreateSimpleVolumeRequest,
    responseType: api_pb.VolumeCreateResponse,
    requestSerialize: serialize_openstorage_api_VolumeCreateSimpleVolumeRequest,
    requestDeserialize: deserialize_openstorage_api_VolumeCreateSimpleVolumeRequest,
    responseSerialize: serialize_openstorage_api_VolumeCreateResponse,
    responseDeserialize: deserialize_openstorage_api_VolumeCreateResponse,
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
};

exports.OpenStorageVolumeClient = grpc.makeGenericClientConstructor(OpenStorageVolumeService);
