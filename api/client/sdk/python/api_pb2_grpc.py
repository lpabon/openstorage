# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import api_pb2 as api__pb2


class OpenStorageClusterStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Enumerate = channel.unary_unary(
        '/openstorage.api.OpenStorageCluster/Enumerate',
        request_serializer=api__pb2.SdkClusterEnumerateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkClusterEnumerateResponse.FromString,
        )
    self.Inspect = channel.unary_unary(
        '/openstorage.api.OpenStorageCluster/Inspect',
        request_serializer=api__pb2.SdkClusterInspectRequest.SerializeToString,
        response_deserializer=api__pb2.SdkClusterInspectResponse.FromString,
        )
    self.AlertEnumerate = channel.unary_unary(
        '/openstorage.api.OpenStorageCluster/AlertEnumerate',
        request_serializer=api__pb2.SdkClusterAlertEnumerateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkClusterAlertEnumerateResponse.FromString,
        )
    self.AlertClear = channel.unary_unary(
        '/openstorage.api.OpenStorageCluster/AlertClear',
        request_serializer=api__pb2.SdkClusterAlertClearRequest.SerializeToString,
        response_deserializer=api__pb2.SdkClusterAlertClearResponse.FromString,
        )
    self.AlertErase = channel.unary_unary(
        '/openstorage.api.OpenStorageCluster/AlertErase',
        request_serializer=api__pb2.SdkClusterAlertEraseRequest.SerializeToString,
        response_deserializer=api__pb2.SdkClusterAlertEraseResponse.FromString,
        )


class OpenStorageClusterServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Enumerate(self, request, context):
    """Enumerate lists all the nodes in the cluster.
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Inspect(self, request, context):
    """Inspect the node given a UUID.
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def AlertEnumerate(self, request, context):
    """Get a list of alerts from the storage cluster
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def AlertClear(self, request, context):
    """Clear the alert for a given resource
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def AlertErase(self, request, context):
    """Erases an alert for a given resource
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_OpenStorageClusterServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Enumerate': grpc.unary_unary_rpc_method_handler(
          servicer.Enumerate,
          request_deserializer=api__pb2.SdkClusterEnumerateRequest.FromString,
          response_serializer=api__pb2.SdkClusterEnumerateResponse.SerializeToString,
      ),
      'Inspect': grpc.unary_unary_rpc_method_handler(
          servicer.Inspect,
          request_deserializer=api__pb2.SdkClusterInspectRequest.FromString,
          response_serializer=api__pb2.SdkClusterInspectResponse.SerializeToString,
      ),
      'AlertEnumerate': grpc.unary_unary_rpc_method_handler(
          servicer.AlertEnumerate,
          request_deserializer=api__pb2.SdkClusterAlertEnumerateRequest.FromString,
          response_serializer=api__pb2.SdkClusterAlertEnumerateResponse.SerializeToString,
      ),
      'AlertClear': grpc.unary_unary_rpc_method_handler(
          servicer.AlertClear,
          request_deserializer=api__pb2.SdkClusterAlertClearRequest.FromString,
          response_serializer=api__pb2.SdkClusterAlertClearResponse.SerializeToString,
      ),
      'AlertErase': grpc.unary_unary_rpc_method_handler(
          servicer.AlertErase,
          request_deserializer=api__pb2.SdkClusterAlertEraseRequest.FromString,
          response_serializer=api__pb2.SdkClusterAlertEraseResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'openstorage.api.OpenStorageCluster', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class OpenStorageVolumeStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Create = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Create',
        request_serializer=api__pb2.SdkVolumeCreateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeCreateResponse.FromString,
        )
    self.CreateFromVolumeId = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/CreateFromVolumeId',
        request_serializer=api__pb2.SdkVolumeCreateFromVolumeIdRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeCreateFromVolumeIdResponse.FromString,
        )
    self.Delete = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Delete',
        request_serializer=api__pb2.SdkVolumeDeleteRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeDeleteResponse.FromString,
        )
    self.Inspect = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Inspect',
        request_serializer=api__pb2.SdkVolumeInspectRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeInspectResponse.FromString,
        )
    self.Enumerate = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Enumerate',
        request_serializer=api__pb2.SdkVolumeEnumerateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeEnumerateResponse.FromString,
        )
    self.SnapshotCreate = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/SnapshotCreate',
        request_serializer=api__pb2.SdkVolumeSnapshotCreateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeSnapshotCreateResponse.FromString,
        )
    self.SnapshotRestore = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/SnapshotRestore',
        request_serializer=api__pb2.SdkVolumeSnapshotRestoreRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeSnapshotRestoreResponse.FromString,
        )
    self.SnapshotEnumerate = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/SnapshotEnumerate',
        request_serializer=api__pb2.SdkVolumeSnapshotEnumerateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeSnapshotEnumerateResponse.FromString,
        )
    self.Attach = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Attach',
        request_serializer=api__pb2.SdkVolumeAttachRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeAttachResponse.FromString,
        )
    self.Detach = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Detach',
        request_serializer=api__pb2.SdkVolumeDetachRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeDetachResponse.FromString,
        )
    self.Mount = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Mount',
        request_serializer=api__pb2.SdkVolumeMountRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeMountResponse.FromString,
        )
    self.Unmount = channel.unary_unary(
        '/openstorage.api.OpenStorageVolume/Unmount',
        request_serializer=api__pb2.SdkVolumeUnmountRequest.SerializeToString,
        response_deserializer=api__pb2.SdkVolumeUnmountResponse.FromString,
        )


class OpenStorageVolumeServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Create(self, request, context):
    """Creates a new volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateFromVolumeId(self, request, context):
    """CreateFromVolumeId creates a new volume cloned from an existing volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Delete(self, request, context):
    """Delete a volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Inspect(self, request, context):
    """Get information on a volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Enumerate(self, request, context):
    """Get a list of volumes
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SnapshotCreate(self, request, context):
    """Create a snapshot of a volume. This creates an immutable (read-only),
    point-in-time snapshot of a volume.
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SnapshotRestore(self, request, context):
    """Restores a volume to a specified snapshot
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SnapshotEnumerate(self, request, context):
    """List the number of snapshots for a specific volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Attach(self, request, context):
    """Attach device to host
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Detach(self, request, context):
    """Detaches the volume from the node.
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Mount(self, request, context):
    """Attaches the volume to a node.
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Unmount(self, request, context):
    """Unmount volume at specified path
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_OpenStorageVolumeServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Create': grpc.unary_unary_rpc_method_handler(
          servicer.Create,
          request_deserializer=api__pb2.SdkVolumeCreateRequest.FromString,
          response_serializer=api__pb2.SdkVolumeCreateResponse.SerializeToString,
      ),
      'CreateFromVolumeId': grpc.unary_unary_rpc_method_handler(
          servicer.CreateFromVolumeId,
          request_deserializer=api__pb2.SdkVolumeCreateFromVolumeIdRequest.FromString,
          response_serializer=api__pb2.SdkVolumeCreateFromVolumeIdResponse.SerializeToString,
      ),
      'Delete': grpc.unary_unary_rpc_method_handler(
          servicer.Delete,
          request_deserializer=api__pb2.SdkVolumeDeleteRequest.FromString,
          response_serializer=api__pb2.SdkVolumeDeleteResponse.SerializeToString,
      ),
      'Inspect': grpc.unary_unary_rpc_method_handler(
          servicer.Inspect,
          request_deserializer=api__pb2.SdkVolumeInspectRequest.FromString,
          response_serializer=api__pb2.SdkVolumeInspectResponse.SerializeToString,
      ),
      'Enumerate': grpc.unary_unary_rpc_method_handler(
          servicer.Enumerate,
          request_deserializer=api__pb2.SdkVolumeEnumerateRequest.FromString,
          response_serializer=api__pb2.SdkVolumeEnumerateResponse.SerializeToString,
      ),
      'SnapshotCreate': grpc.unary_unary_rpc_method_handler(
          servicer.SnapshotCreate,
          request_deserializer=api__pb2.SdkVolumeSnapshotCreateRequest.FromString,
          response_serializer=api__pb2.SdkVolumeSnapshotCreateResponse.SerializeToString,
      ),
      'SnapshotRestore': grpc.unary_unary_rpc_method_handler(
          servicer.SnapshotRestore,
          request_deserializer=api__pb2.SdkVolumeSnapshotRestoreRequest.FromString,
          response_serializer=api__pb2.SdkVolumeSnapshotRestoreResponse.SerializeToString,
      ),
      'SnapshotEnumerate': grpc.unary_unary_rpc_method_handler(
          servicer.SnapshotEnumerate,
          request_deserializer=api__pb2.SdkVolumeSnapshotEnumerateRequest.FromString,
          response_serializer=api__pb2.SdkVolumeSnapshotEnumerateResponse.SerializeToString,
      ),
      'Attach': grpc.unary_unary_rpc_method_handler(
          servicer.Attach,
          request_deserializer=api__pb2.SdkVolumeAttachRequest.FromString,
          response_serializer=api__pb2.SdkVolumeAttachResponse.SerializeToString,
      ),
      'Detach': grpc.unary_unary_rpc_method_handler(
          servicer.Detach,
          request_deserializer=api__pb2.SdkVolumeDetachRequest.FromString,
          response_serializer=api__pb2.SdkVolumeDetachResponse.SerializeToString,
      ),
      'Mount': grpc.unary_unary_rpc_method_handler(
          servicer.Mount,
          request_deserializer=api__pb2.SdkVolumeMountRequest.FromString,
          response_serializer=api__pb2.SdkVolumeMountResponse.SerializeToString,
      ),
      'Unmount': grpc.unary_unary_rpc_method_handler(
          servicer.Unmount,
          request_deserializer=api__pb2.SdkVolumeUnmountRequest.FromString,
          response_serializer=api__pb2.SdkVolumeUnmountResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'openstorage.api.OpenStorageVolume', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class OpenStorageObjectstoreStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.InspectObjectstore = channel.unary_unary(
        '/openstorage.api.OpenStorageObjectstore/InspectObjectstore',
        request_serializer=api__pb2.SdkObjectstoreInspectRequest.SerializeToString,
        response_deserializer=api__pb2.SdkObjectstoreInspectResponse.FromString,
        )
    self.CreateObjectstore = channel.unary_unary(
        '/openstorage.api.OpenStorageObjectstore/CreateObjectstore',
        request_serializer=api__pb2.SdkObjectstoreCreateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkObjectstoreCreateResponse.FromString,
        )
    self.DeleteObjectstore = channel.unary_unary(
        '/openstorage.api.OpenStorageObjectstore/DeleteObjectstore',
        request_serializer=api__pb2.SdkObjectstoreDeleteRequest.SerializeToString,
        response_deserializer=api__pb2.SdkObjectstoreDeleteResponse.FromString,
        )
    self.UpdateObjectstore = channel.unary_unary(
        '/openstorage.api.OpenStorageObjectstore/UpdateObjectstore',
        request_serializer=api__pb2.SdkObjectstoreUpdateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkObjectstoreUpdateResponse.FromString,
        )


class OpenStorageObjectstoreServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def InspectObjectstore(self, request, context):
    """InspectObjectstore returns current status of objectstore 
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateObjectstore(self, request, context):
    """CreateObjectstore creates on specified volume
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DeleteObjectstore(self, request, context):
    """DeleteObjectstore deletes objectstore by id
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def UpdateObjectstore(self, request, context):
    """UpdateObjectstore updates provided objectstore status
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_OpenStorageObjectstoreServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'InspectObjectstore': grpc.unary_unary_rpc_method_handler(
          servicer.InspectObjectstore,
          request_deserializer=api__pb2.SdkObjectstoreInspectRequest.FromString,
          response_serializer=api__pb2.SdkObjectstoreInspectResponse.SerializeToString,
      ),
      'CreateObjectstore': grpc.unary_unary_rpc_method_handler(
          servicer.CreateObjectstore,
          request_deserializer=api__pb2.SdkObjectstoreCreateRequest.FromString,
          response_serializer=api__pb2.SdkObjectstoreCreateResponse.SerializeToString,
      ),
      'DeleteObjectstore': grpc.unary_unary_rpc_method_handler(
          servicer.DeleteObjectstore,
          request_deserializer=api__pb2.SdkObjectstoreDeleteRequest.FromString,
          response_serializer=api__pb2.SdkObjectstoreDeleteResponse.SerializeToString,
      ),
      'UpdateObjectstore': grpc.unary_unary_rpc_method_handler(
          servicer.UpdateObjectstore,
          request_deserializer=api__pb2.SdkObjectstoreUpdateRequest.FromString,
          response_serializer=api__pb2.SdkObjectstoreUpdateResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'openstorage.api.OpenStorageObjectstore', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))


class OpenStorageCredentialsStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.CreateForAWS = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/CreateForAWS',
        request_serializer=api__pb2.SdkCredentialCreateAWSRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialCreateAWSResponse.FromString,
        )
    self.CreateForAzure = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/CreateForAzure',
        request_serializer=api__pb2.SdkCredentialCreateAzureRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialCreateAzureResponse.FromString,
        )
    self.CreateForGoogle = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/CreateForGoogle',
        request_serializer=api__pb2.SdkCredentialCreateGoogleRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialCreateGoogleResponse.FromString,
        )
    self.EnumerateForAWS = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/EnumerateForAWS',
        request_serializer=api__pb2.SdkCredentialEnumerateAWSRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialEnumerateAWSResponse.FromString,
        )
    self.EnumerateForAzure = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/EnumerateForAzure',
        request_serializer=api__pb2.SdkCredentialEnumerateAzureRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialEnumerateAzureResponse.FromString,
        )
    self.EnumerateForGoogle = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/EnumerateForGoogle',
        request_serializer=api__pb2.SdkCredentialEnumerateGoogleRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialEnumerateGoogleResponse.FromString,
        )
    self.CredentialDelete = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/CredentialDelete',
        request_serializer=api__pb2.SdkCredentialDeleteRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialDeleteResponse.FromString,
        )
    self.CredentialValidate = channel.unary_unary(
        '/openstorage.api.OpenStorageCredentials/CredentialValidate',
        request_serializer=api__pb2.SdkCredentialValidateRequest.SerializeToString,
        response_deserializer=api__pb2.SdkCredentialValidateResponse.FromString,
        )


class OpenStorageCredentialsServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def CreateForAWS(self, request, context):
    """Provide credentials to OpenStorage and if valid,
    it will return an identifier to the credentials

    Create credential for AWS S3 and if valid ,
    returns a unique identifier
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateForAzure(self, request, context):
    """Create credential for Azure and if valid ,
    returns a unique identifier
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateForGoogle(self, request, context):
    """Create credential for Google and if valid ,
    returns a unique identifier
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EnumerateForAWS(self, request, context):
    """EnumerateForAWS lists the configured AWS credentials
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EnumerateForAzure(self, request, context):
    """EnumerateForAzure lists the configured Azure credentials
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EnumerateForGoogle(self, request, context):
    """EnumerateForGoogle lists the configured Google credentials
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CredentialDelete(self, request, context):
    """Delete a specified credential
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CredentialValidate(self, request, context):
    """Validate a specified credential
    """
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_OpenStorageCredentialsServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'CreateForAWS': grpc.unary_unary_rpc_method_handler(
          servicer.CreateForAWS,
          request_deserializer=api__pb2.SdkCredentialCreateAWSRequest.FromString,
          response_serializer=api__pb2.SdkCredentialCreateAWSResponse.SerializeToString,
      ),
      'CreateForAzure': grpc.unary_unary_rpc_method_handler(
          servicer.CreateForAzure,
          request_deserializer=api__pb2.SdkCredentialCreateAzureRequest.FromString,
          response_serializer=api__pb2.SdkCredentialCreateAzureResponse.SerializeToString,
      ),
      'CreateForGoogle': grpc.unary_unary_rpc_method_handler(
          servicer.CreateForGoogle,
          request_deserializer=api__pb2.SdkCredentialCreateGoogleRequest.FromString,
          response_serializer=api__pb2.SdkCredentialCreateGoogleResponse.SerializeToString,
      ),
      'EnumerateForAWS': grpc.unary_unary_rpc_method_handler(
          servicer.EnumerateForAWS,
          request_deserializer=api__pb2.SdkCredentialEnumerateAWSRequest.FromString,
          response_serializer=api__pb2.SdkCredentialEnumerateAWSResponse.SerializeToString,
      ),
      'EnumerateForAzure': grpc.unary_unary_rpc_method_handler(
          servicer.EnumerateForAzure,
          request_deserializer=api__pb2.SdkCredentialEnumerateAzureRequest.FromString,
          response_serializer=api__pb2.SdkCredentialEnumerateAzureResponse.SerializeToString,
      ),
      'EnumerateForGoogle': grpc.unary_unary_rpc_method_handler(
          servicer.EnumerateForGoogle,
          request_deserializer=api__pb2.SdkCredentialEnumerateGoogleRequest.FromString,
          response_serializer=api__pb2.SdkCredentialEnumerateGoogleResponse.SerializeToString,
      ),
      'CredentialDelete': grpc.unary_unary_rpc_method_handler(
          servicer.CredentialDelete,
          request_deserializer=api__pb2.SdkCredentialDeleteRequest.FromString,
          response_serializer=api__pb2.SdkCredentialDeleteResponse.SerializeToString,
      ),
      'CredentialValidate': grpc.unary_unary_rpc_method_handler(
          servicer.CredentialValidate,
          request_deserializer=api__pb2.SdkCredentialValidateRequest.FromString,
          response_serializer=api__pb2.SdkCredentialValidateResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'openstorage.api.OpenStorageCredentials', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
