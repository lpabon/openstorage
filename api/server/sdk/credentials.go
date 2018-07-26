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
	"reflect"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CredentialServer is an implementation of the gRPC OpenStorageCredential interface
type CredentialServer struct {
	driver volume.VolumeDriver
}

// CreateForAWS method creates credential for AWS S3.
func (s *CredentialServer) CreateForAWS(
	ctx context.Context,
	req *api.SdkCredentialCreateAWSRequest,
) (*api.SdkCredentialCreateAWSResponse, error) {

	if len(req.GetCredential().GetAccessKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Access Key")
	}

	if len(req.GetCredential().GetSecretKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Secret Key")
	}

	if len(req.GetCredential().GetRegion()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Region Key")
	}

	if len(req.GetCredential().GetEndpoint()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Endpoint Key")
	}

	params := make(map[string]string)

	params[api.OptCredType] = "s3"
	params[api.OptCredRegion] = req.GetCredential().GetRegion()
	params[api.OptCredEndpoint] = req.GetCredential().GetEndpoint()
	params[api.OptCredAccessKey] = req.GetCredential().GetAccessKey()
	params[api.OptCredSecretKey] = req.GetCredential().GetSecretKey()

	uuid, err := s.driver.CredsCreate(params)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to create S3 credentials: %v",
			err.Error())
	}

	err = validateAndDeleteIfInvalid(s, uuid)

	if err != nil {
		return nil, err
	}
	return &api.SdkCredentialCreateAWSResponse{CredentialId: uuid}, nil

}

// CreateForAzure method creates credential for Azure.
func (s *CredentialServer) CreateForAzure(
	ctx context.Context,
	req *api.SdkCredentialCreateAzureRequest,
) (*api.SdkCredentialCreateAzureResponse, error) {

	if len(req.GetCredential().GetAccountKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Account Key")
	}

	if len(req.GetCredential().GetAccountName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Account name")
	}

	params := make(map[string]string)

	params[api.OptCredType] = "azure"
	params[api.OptCredAzureAccountKey] = req.GetCredential().GetAccountKey()
	params[api.OptCredAzureAccountName] = req.GetCredential().GetAccountName()

	uuid, err := s.driver.CredsCreate(params)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to create Azure credentials: %v",
			err.Error())
	}

	err = validateAndDeleteIfInvalid(s, uuid)

	if err != nil {
		return nil, err
	}
	return &api.SdkCredentialCreateAzureResponse{CredentialId: uuid}, nil
}

// CreateForGoogle method creates credential for Google.
func (s *CredentialServer) CreateForGoogle(
	ctx context.Context,
	req *api.SdkCredentialCreateGoogleRequest,
) (*api.SdkCredentialCreateGoogleResponse, error) {

	if len(req.GetCredential().GetJsonKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply JSON Key")
	}

	if len(req.GetCredential().GetProjectId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must supply Project ID")
	}

	params := make(map[string]string)

	params[api.OptCredType] = "google"
	params[api.OptCredGoogleProjectID] = req.GetCredential().GetProjectId()
	params[api.OptCredGoogleJsonKey] = req.GetCredential().GetJsonKey()

	uuid, err := s.driver.CredsCreate(params)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to create Google credentials: %v",
			err.Error())
	}

	err = validateAndDeleteIfInvalid(s, uuid)

	if err != nil {
		return nil, err
	}

	return &api.SdkCredentialCreateGoogleResponse{CredentialId: uuid}, nil
}

// Validate validates a specified Credential.
func (s *CredentialServer) Validate(
	ctx context.Context,
	req *api.SdkCredentialValidateRequest,
) (*api.SdkCredentialValidateResponse, error) {

	if len(req.GetCredentialId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must provide credentials uuid")
	}

	validateReq := &api.SdkCredentialValidateRequest{CredentialId: req.GetCredentialId()}

	err := s.driver.CredsValidate(validateReq.GetCredentialId())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to validate credentials: %v",
			err.Error())
	}
	return &api.SdkCredentialValidateResponse{}, nil

}

// Delete deletes a specified credential
func (s *CredentialServer) Delete(
	ctx context.Context,
	req *api.SdkCredentialDeleteRequest,
) (*api.SdkCredentialDeleteResponse, error) {

	if len(req.GetCredentialId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must provide credentials uuid")
	}

	err := s.driver.CredsDelete(req.GetCredentialId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to  delete credentials: %v",
			err.Error())
	}

	return &api.SdkCredentialDeleteResponse{}, nil
}

// EnumerateForAWS list credentials for AWS
func (s *CredentialServer) EnumerateForAWS(
	ctx context.Context,
	req *api.SdkCredentialEnumerateAWSRequest,
) (*api.SdkCredentialEnumerateAWSResponse, error) {

	credList, err := s.driver.CredsEnumerate()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials AWS: %v",
			err.Error())
	}

	// By defaultcredList will have all credential details, we will extract for
	// respective cloud provider and return result
	// this may not be expected behaviour, we have to do this since
	// `interface` can't be mapped directly with other lang
	s3Creds, err := getCredentialMap(credList, "s3")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials AWS: %v",
			err.Error())
	}

	// Fill up s3 credential resonse
	creds := []*api.S3Credential{}
	for id, cred := range s3Creds {
		credResp := &api.S3Credential{
			CredentialId: id,
			AccessKey:    cred[api.OptCredAccessKey],
			Endpoint:     cred[api.OptCredEndpoint],
			Region:       cred[api.OptCredRegion],
		}
		creds = append(creds, credResp)
	}

	return &api.SdkCredentialEnumerateAWSResponse{Credential: creds}, nil
}

// EnumerateForAzure list credentials for AWS
func (s *CredentialServer) EnumerateForAzure(
	ctx context.Context,
	req *api.SdkCredentialEnumerateAzureRequest,
) (*api.SdkCredentialEnumerateAzureResponse, error) {
	credList, err := s.driver.CredsEnumerate()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials: %v",
			err.Error())
	}

	// By defaultcredList will have all credential details, we will extract for
	// respective cloud provider and return result
	// this may not be expected behaviour, we have to do this since
	// `interface` can't be mapped directly with other lang
	azureCreds, err := getCredentialMap(credList, "azure")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials Azure: %v",
			err.Error())
	}

	// Fill up azure credential resonse
	creds := []*api.AzureCredential{}
	for id, cred := range azureCreds {
		credResp := &api.AzureCredential{
			CredentialId: id,
			AccountName:  cred[api.OptCredAzureAccountName],
			AccountKey:   cred[api.OptCredAzureAccountKey],
		}
		creds = append(creds, credResp)
	}
	return &api.SdkCredentialEnumerateAzureResponse{Credential: creds}, nil
}

// EnumerateForGoogle list credentials for Google
func (s *CredentialServer) EnumerateForGoogle(
	ctx context.Context,
	req *api.SdkCredentialEnumerateGoogleRequest,
) (*api.SdkCredentialEnumerateGoogleResponse, error) {
	credList, err := s.driver.CredsEnumerate()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials: %v",
			err.Error())
	}

	// By defaultcredList will have all credential details, we will extract for
	// respective cloud provider and return result
	// this may not be expected behaviour, we have to do this since
	// `interface` can't be mapped directly with other lang
	googleCreds, err := getCredentialMap(credList, "google")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to enumerate credentials Azure: %v",
			err.Error())
	}

	// Fill up google credential resonse
	creds := []*api.GoogleCredential{}
	for id, cred := range googleCreds {
		credResp := &api.GoogleCredential{
			CredentialId: id,
			ProjectId:    cred[api.OptCredGoogleProjectID],
		}
		creds = append(creds, credResp)
	}

	return &api.SdkCredentialEnumerateGoogleResponse{Credential: creds}, nil
}

func validateAndDeleteIfInvalid(s *CredentialServer, uuid string) error {
	// Validate if the credentials provided were correct or not
	req := &api.SdkCredentialValidateRequest{CredentialId: uuid}

	validateErr := s.driver.CredsValidate(req.GetCredentialId())

	if validateErr != nil {
		deleteCred := &api.SdkCredentialDeleteRequest{CredentialId: uuid}
		err := s.driver.CredsDelete(deleteCred.GetCredentialId())

		if err != nil {
			return status.Errorf(
				codes.Internal,
				"failed to delete invalid Google credentials: %v",
				err.Error())
		}

		return status.Errorf(
			codes.PermissionDenied,
			"credentials could not be validated: %v",
			validateErr.Error())
	}

	return nil
}

func getCredentialMap(credList map[string]interface{}, credType string) (map[string]map[string]string, error) {
	filtered := make(map[string]map[string]string)

	for k, v := range credList {
		c, ok := v.(map[string]string)
		if !ok {
			return nil, fmt.Errorf("Error parsing credentials of type %v",
				reflect.TypeOf(v).String())
		}

		// Look for only one type, fill up creds with same type array
		if c[api.OptCredType] == credType {
			filtered[k] = c
		}
	}

	return filtered, nil
}
