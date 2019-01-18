package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/api/spec"
	"github.com/libopenstorage/openstorage/config"
	"github.com/libopenstorage/openstorage/pkg/grpcserver"
	"github.com/libopenstorage/openstorage/pkg/options"
	"github.com/libopenstorage/openstorage/pkg/util"
	"github.com/libopenstorage/openstorage/volume"
	"github.com/libopenstorage/openstorage/volume/drivers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// VolumeDriver is the string returned in the handshake protocol.
	VolumeDriver = "VolumeDriver"
)

// Implementation of the Docker volumes plugin specification.
type driver struct {
	restBase
	spec.SpecHandler

	sdkUds string
	conn   *grpc.ClientConn
	mu     sync.Mutex
	cookie *CookieOnce
}

type handshakeResp struct {
	Implements []string
}

type volumeRequest struct {
	Name string
	Opts map[string]string
}

type mountRequest struct {
	Name string
	ID   string
}

type volumeResponse struct {
	Err string
}

type volumePathResponse struct {
	Mountpoint string
	volumeResponse
}

type volumeInfo struct {
	Name       string
	Mountpoint string
}

type capabilities struct {
	Scope string
}

type capabilitiesResponse struct {
	Capabilities capabilities
}

func newVolumePlugin(name, sdkUds string) restServer {
	d := &driver{
		restBase:    restBase{name: name, version: "0.3"},
		SpecHandler: spec.NewSpecHandler(),
		sdkUds:      sdkUds,
		cookie:      NewCookieOnce(1 * time.Minute),
	}
	return d
}

func (d *driver) String() string {
	return d.name
}

func volDriverPath(method string) string {
	return fmt.Sprintf("/%s.%s", VolumeDriver, method)
}

func (d *driver) volNotFound(request string, id string, e error, w http.ResponseWriter) error {
	err := fmt.Errorf("Failed to locate volume: " + e.Error())
	if e == volume.ErrDriverInitializing {
		d.logRequest(request, id).Warnln(http.StatusInternalServerError, " ", err.Error())
	} else {
		d.logRequest(request, id).Warnln(http.StatusNotFound, " ", err.Error())
	}
	return err
}

func (d *driver) volNotMounted(request string, id string) error {
	err := fmt.Errorf("volume not mounted")
	d.logRequest(request, id).Debugln(http.StatusNotFound, " ", err.Error())
	return err
}

func (d *driver) Routes() []*Route {
	return []*Route{
		{verb: "POST", path: volDriverPath("Create"), fn: d.create},
		{verb: "POST", path: volDriverPath("Remove"), fn: d.remove},
		{verb: "POST", path: volDriverPath("Mount"), fn: d.mount},
		{verb: "POST", path: volDriverPath("Path"), fn: d.path},
		{verb: "POST", path: volDriverPath("List"), fn: d.list},
		{verb: "POST", path: volDriverPath("Get"), fn: d.get},
		// Run in insecure mode until we figure out how to secure using another model
		{verb: "POST", path: volDriverPath("Unmount"), fn: d.unmountInsecure},
		{verb: "POST", path: volDriverPath("Capabilities"), fn: d.capabilities},
		{verb: "POST", path: "/Plugin.Activate", fn: d.handshake},
		{verb: "GET", path: "/status", fn: d.status},
	}
}

func (d *driver) emptyResponse(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(&volumeResponse{})
}

func (d *driver) errorResponse(method string, w http.ResponseWriter, err error) {
	if err == volume.ErrDriverInitializing {
		d.sendError(method, "", w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(&volumeResponse{Err: err.Error()})
	}
}

func (d *driver) volFromName(name string) (*api.Volume, error) {
	v, err := volumedrivers.Get(d.name)
	if err != nil {
		return nil, fmt.Errorf("Cannot locate volume driver for %s: %s", d.name, err.Error())
	}
	return util.VolumeFromName(v, name)
}

func (d *driver) volFromNameSdk(ctx context.Context, volumes api.OpenStorageVolumeClient, name string) (*api.Volume, error) {
	// get volume id
	volId, err := d.volIdFromName(ctx, volumes, name)
	if err != nil {
		return nil, err
	}

	// inspect for actual volume
	inspectResp, err := volumes.Inspect(ctx, &api.SdkVolumeInspectRequest{
		VolumeId: volId,
	})
	if err != nil {
		return nil, err
	}
	return inspectResp.Volume, nil
}

func (d *driver) volFromIdSdk(ctx context.Context, volumes api.OpenStorageVolumeClient, volId string) (*api.Volume, error) {
	inspectResp, err := volumes.Inspect(ctx, &api.SdkVolumeInspectRequest{
		VolumeId: volId,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot locate volume with id %s", volId)
	}
	return inspectResp.Volume, nil
}

func (d *driver) volIdFromName(ctx context.Context, volumes api.OpenStorageVolumeClient, name string) (string, error) {
	enumerateResp, err := volumes.EnumerateWithFilters(ctx, &api.SdkVolumeEnumerateWithFiltersRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	} else if len(enumerateResp.VolumeIds) < 1 {
		return "", fmt.Errorf("Cannot locate volume with name %s", name)
	}

	return enumerateResp.VolumeIds[0], nil
}

func (d *driver) decode(method string, w http.ResponseWriter, r *http.Request) (*volumeRequest, error) {
	var request volumeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		e := fmt.Errorf("Unable to decode JSON payload")
		d.sendError(method, "", w, e.Error()+":"+err.Error(), http.StatusBadRequest)
		return nil, e
	}
	d.logRequest(method, request.Name).Debugln("")
	return &request, nil
}

func (d *driver) decodeMount(method string, w http.ResponseWriter, r *http.Request) (*mountRequest, error) {
	var request mountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		e := fmt.Errorf("Unable to decode JSON payload")
		d.sendError(method, "", w, e.Error()+":"+err.Error(), http.StatusBadRequest)
		return nil, e
	}
	d.logRequest(method, request.Name).Debugf("ID: %v", request.ID)
	return &request, nil
}

func (d *driver) handshake(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(&handshakeResp{
		[]string{VolumeDriver},
	})
	if err != nil {
		d.sendError("handshake", "", w, "encode error", http.StatusInternalServerError)
		return
	}
	d.logRequest("handshake", "").Debugln("Handshake completed")
}

func (d *driver) attachToken(ctx context.Context, request *volumeRequest) (context.Context, string) {
	token, tokenInName := d.GetTokenFromString(request.Name)
	if !tokenInName {
		token = request.Opts[api.Token]
	}
	/*
		if len(token) == 0 {
			token, _ = d.cookie.Pop(request.Name)
		}
	*/
	if len(token) == 0 {
		return ctx, ""
		//token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1AcG9ydHdvcnguY29tIiwiZXhwIjoxNTc5MjQxODk3LCJncm91cHMiOlsiKiJdLCJpYXQiOjE1NDc3MDU4OTcsImlzcyI6Im9wZW5zdG9yYWdlLmlvIiwibmFtZSI6Ik1vbml0b3IiLCJyb2xlcyI6WyJzeXN0ZW0udmlldyJdLCJzdWIiOiJiM0JsYm5OMGIzSmhaMlV1YVc4dlRXOXVhWFJ2Y2k5dFFIQnZjblIzYjNKNExtTnZiUT09In0.st6wj68ifx2YxYr-yAE1eCrC04kKFaofG0YV6V8ZgC4"
	}

	md := metadata.New(map[string]string{
		"authorization": "bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md), token
}

func (d *driver) attachTokenMount(ctx context.Context, request *mountRequest) (context.Context, string) {
	token, _ := d.GetTokenFromString(request.Name)
	if len(token) == 0 {
		token, _ = d.cookie.Pop(request.Name)
	}
	if len(token) == 0 {
		return ctx, ""
	}

	md := metadata.New(map[string]string{
		"authorization": "bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md), token
}

func (d *driver) status(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintln("osd plugin", d.version))
}

func (d *driver) mountpath(name string) string {
	return path.Join(volume.MountBase, name)
}

func (d *driver) getConn() (*grpc.ClientConn, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn == nil {
		var err error
		d.conn, err = grpcserver.Connect(
			d.sdkUds,
			[]grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to gRPC handler: %v", err)
		}
	}
	return d.conn, nil
}

func (d *driver) create(w http.ResponseWriter, r *http.Request) {
	method := "create"
	ctx := r.Context()
	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)

	// attach token in context metadata
	ctx, _ = d.attachToken(ctx, request)

	// get spec for volume creation
	specParsed, spec, locator, source, name := d.SpecFromString(request.Name)
	d.logRequest(method, name).Infoln("")
	if !specParsed {
		spec, locator, source, err = d.SpecFromOpts(request.Opts)
		if err != nil {
			d.errorResponse(method, w, err)
			return
		}
	}

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	// clone if exists, create otherwise
	volumes := api.NewOpenStorageVolumeClient(conn)
	if source != nil && len(source.Parent) != 0 {
		// clone
		_, err = volumes.Clone(ctx, &api.SdkVolumeCloneRequest{
			Name:     name,
			ParentId: source.Parent,
		})
	} else {
		// create
		_, err = volumes.Create(ctx, &api.SdkVolumeCreateRequest{
			Name:   name,
			Spec:   spec,
			Labels: locator.GetVolumeLabels(),
		})
	}
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	json.NewEncoder(w).Encode(&volumeResponse{})
}

func (d *driver) remove(w http.ResponseWriter, r *http.Request) {
	method := "remove"
	ctx := r.Context()
	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)

	// attach token in context metadata
	ctx, _ = d.attachToken(ctx, request)

	// get name for deletion
	_, _, _, _, name := d.SpecFromString(request.Name)
	d.logRequest(method, name).Infoln("")

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}
	volumes := api.NewOpenStorageVolumeClient(conn)

	// get volume id to delete
	volId, err := d.volIdFromName(ctx, volumes, name)
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	// delete volume
	_, err = volumes.Delete(ctx, &api.SdkVolumeDeleteRequest{
		VolumeId: volId,
	})
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	json.NewEncoder(w).Encode(&volumeResponse{})
}

func (d *driver) attachOptionsFromSpec(
	spec *api.VolumeSpec,
) map[string]string {
	if spec.Passphrase != "" {
		opts := make(map[string]string)
		opts[options.OptionsSecret] = spec.Passphrase
		return opts
	}
	return nil
}

func (d *driver) mount(w http.ResponseWriter, r *http.Request) {
	var response volumePathResponse
	ctx := r.Context()
	method := "mount"
	request, err := d.decodeMount(method, w, r)
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)

	// get spec and name from request
	_, spec, _, _, name := d.SpecFromString(request.Name)
	attachOptions := d.attachOptionsFromSpec(spec)

	// attach token in context metadata
	ctx, _ = d.attachTokenMount(ctx, request)

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	//get volume to mount
	volumeClient := api.NewOpenStorageVolumeClient(conn)
	vol, err := d.volFromNameSdk(ctx, volumeClient, name)
	if err != nil {
		d.sendError(method, "", w, err.Error(), http.StatusBadRequest)
		return
	}

	// get and prepare mountpath
	mountpoint := d.mountpath(name)
	response.Mountpoint = mountpoint
	os.MkdirAll(mountpoint, 0755)

	// mount volume
	mountClient := api.NewOpenStorageMountAttachClient(conn)
	_, err = mountClient.Mount(ctx, &api.SdkVolumeMountRequest{
		VolumeId:  vol.Id,
		MountPath: mountpoint,
		Options:   options.NewVolumeAttachOptions(attachOptions),
	})
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	d.logRequest(method, request.Name).Infof("response %v", response.Mountpoint)
	json.NewEncoder(w).Encode(&response)
}

/*
func (d *driver) path(w http.ResponseWriter, r *http.Request) {
	method := "path"
	var response volumePathResponse
	d.logRequest(method, "").Infof("")

	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}

	_, _, _, _, name := d.SpecFromString(request.Name)
	vol, err := d.volFromName(name)
	if err != nil {
		e := d.volNotFound(method, request.Name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	d.logRequest(method, name).Debugf("")
	if len(vol.AttachPath) == 0 || len(vol.AttachPath) == 0 {
		e := d.volNotMounted(method, name)
		d.errorResponse(method, w, e)
		return
	}
	response.Mountpoint = vol.AttachPath[0]
	response.Mountpoint = path.Join(response.Mountpoint, config.DataDir)
	d.logRequest(method, request.Name).Debugf("response %v", response.Mountpoint)
	json.NewEncoder(w).Encode(&response)
}
*/

func (d *driver) list(w http.ResponseWriter, r *http.Request) {
	method := "list"

	d.logRequest(method, "").Infof("")

	v, err := volumedrivers.Get(d.name)
	if err != nil {
		d.logRequest(method, "").Warnf("Cannot locate volume driver: %v", err.Error())
		d.errorResponse(method, w, err)
		return
	}

	vols, err := v.Enumerate(nil, nil)
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	volInfo := make([]volumeInfo, len(vols))
	for i, v := range vols {
		volInfo[i].Name = v.Locator.Name
		if len(v.AttachPath) > 0 || len(v.AttachPath) > 0 {
			volInfo[i].Mountpoint = path.Join(v.AttachPath[0], config.DataDir)
		}
	}
	json.NewEncoder(w).Encode(map[string][]volumeInfo{"Volumes": volInfo})
}

/*
func (d *driver) get(w http.ResponseWriter, r *http.Request) {
	method := "get"

	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)
	parsed, _, _, _, name := d.SpecFromString(request.Name)
	returnName := ""

	// This looks like a bug
	if parsed {
		returnName = request.Name
	} else {
		returnName = name
	}
	d.logRequest(method, "").
		Infof("returnName = %s, request.Name=%s, name=%s",
			returnName,
			request.Name,
			name)
	vol, err := d.volFromName(name)
	if err != nil {
		e := d.volNotFound(method, request.Name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	volInfo := volumeInfo{Name: returnName}
	if len(vol.AttachPath) > 0 || len(vol.AttachPath) > 0 {
		volInfo.Mountpoint = path.Join(vol.AttachPath[0], config.DataDir)
	}

		_, token := d.attachToken(context.Background(), request)
		if len(token) != 0 {
			d.cookie.Push(name, token)
		}

	json.NewEncoder(w).Encode(map[string]volumeInfo{"Volume": volInfo})
}
*/

func (d *driver) path(w http.ResponseWriter, r *http.Request) {
	method := "path"
	ctx := r.Context()
	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)
	var response volumePathResponse

	// attach token in context metadata
	//var token string
	ctx, _ = d.attachToken(ctx, request)

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	volumes := api.NewOpenStorageVolumeClient(conn)
	_, _, _, _, name := d.SpecFromString(request.Name)
	vol, err := d.volFromNameSdk(ctx, volumes, name)
	if err != nil {
		e := d.volNotFound(method, request.Name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	// Save the cookie
	/*
		if len(token) != 0 {
			d.cookie.Push(name, token)
		}
	*/

	d.logRequest(method, name).Debugf("")
	if len(vol.AttachPath) == 0 || len(vol.AttachPath) == 0 {
		e := d.volNotMounted(method, name)
		d.errorResponse(method, w, e)
		return
	}

	response.Mountpoint = vol.AttachPath[0]
	response.Mountpoint = path.Join(response.Mountpoint, config.DataDir)
	d.logRequest(method, request.Name).Debugf("response %v", response.Mountpoint)
	json.NewEncoder(w).Encode(&response)
}

/*

func (d *driver) list(w http.ResponseWriter, r *http.Request) {
	method := "list"
	ctx := r.Context()
	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}

	// attach token in context metadata
	var token string
	ctx, token = d.attachToken(ctx, request)
	d.logRequest(method, "").Infof("Req: %v\nToken: %s", request, token)

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	// get all volumes
	volumes := api.NewOpenStorageVolumeClient(conn)
	enumerateResp, err := volumes.Enumerate(ctx, &api.SdkVolumeEnumerateRequest{})
	if err != nil {
		d.errorResponse(method, w, err)
		d.logRequest(method, "").Infof("List failed")
		return
	}
	d.logRequest(method, "").Infof("Volumes: %+v", enumerateResp.GetVolumeIds())

	volInfo := make([]volumeInfo, len(enumerateResp.GetVolumeIds()))
	for i, id := range enumerateResp.GetVolumeIds() {
		inspectResp, err := volumes.Inspect(ctx, &api.SdkVolumeInspectRequest{
			VolumeId: id,
		})
		if err != nil {
			continue
		}
		volInfo[i].Name = inspectResp.Volume.Locator.Name
		if len(inspectResp.Volume.AttachPath) > 0 || len(inspectResp.Volume.AttachPath) > 0 {
			volInfo[i].Mountpoint = path.Join(inspectResp.Volume.AttachPath[0], config.DataDir)
		}
	}

	type listresp struct {
		Volumes []volumeInfo
		Err     string
	}
	lr := listresp{
		Volumes: volInfo,
	}
	d.logRequest(method, "").Infof("List resp: %v", lr)

	json.NewEncoder(w).Encode(lr)
}
*/

func (d *driver) get(w http.ResponseWriter, r *http.Request) {
	method := "get"
	ctx := r.Context()
	request, err := d.decode(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)

	// attach token in context metadata
	//var token string
	ctx, _ = d.attachToken(ctx, request)

	// get name from the request
	parsed, _, _, _, name := d.SpecFromString(request.Name)
	var returnName string
	if parsed {
		returnName = request.Name
	} else {
		returnName = name
	}

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	// get volume
	volumes := api.NewOpenStorageVolumeClient(conn)
	vol, err := d.volFromNameSdk(ctx, volumes, name)
	if err != nil {
		e := d.volNotFound(method, name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	// Save the cookie
	/*
		if len(token) != 0 {
			d.cookie.Push(name, token)
		}
	*/

	// create response info
	volInfo := volumeInfo{Name: returnName}
	if len(vol.AttachPath) > 0 || len(vol.AttachPath) > 0 {
		volInfo.Mountpoint = path.Join(vol.AttachPath[0], config.DataDir)
	}

	json.NewEncoder(w).Encode(map[string]volumeInfo{"Volume": volInfo})
}

func (d *driver) unmountInsecure(w http.ResponseWriter, r *http.Request) {
	method := "unmount"

	v, err := volumedrivers.Get(d.name)
	if err != nil {
		d.logRequest(method, "").Warnf(
			"Cannot locate volume driver: %v",
			err.Error())
		d.errorResponse(method, w, err)
		return
	}

	request, err := d.decodeMount(method, w, r)
	if err != nil {
		return
	}
	d.logRequest(method, "").Infof("Req: %v", request)

	_, _, _, _, name := d.SpecFromString(request.Name)
	vol, err := d.volFromName(name)
	if err != nil {
		e := d.volNotFound(method, name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	mountpoint := d.mountpath(name)
	id := vol.Id
	if vol.Spec.Scale > 1 {
		id = v.MountedAt(mountpoint)
		if len(id) == 0 {
			err := fmt.Errorf("Failed to find volume mapping for %v",
				mountpoint)
			d.logRequest(method, request.Name).Warnf(
				"Cannot unmount volume %v, %v",
				mountpoint, err)
			d.errorResponse(method, w, err)
			return
		}
	}

	opts := make(map[string]string)
	opts[options.OptionsDeleteAfterUnmount] = "true"

	err = v.Unmount(id, mountpoint, opts)
	if err != nil {
		d.logRequest(method, request.Name).Warnf(
			"Cannot unmount volume %v, %v",
			mountpoint, err)
		d.errorResponse(method, w, err)
		return
	}

	if v.Type() == api.DriverType_DRIVER_TYPE_BLOCK {
		_ = v.Detach(id, nil)
	}
	d.emptyResponse(w)
}

func (d *driver) unmount(w http.ResponseWriter, r *http.Request) {

	method := "unmount"
	ctx := r.Context()
	request, err := d.decodeMount(method, w, r)
	if err != nil {
		return
	}

	// attach token in context metadata
	ctx, _ = d.attachTokenMount(ctx, request)

	// get name from the request
	_, _, _, _, name := d.SpecFromString(request.Name)

	// get grpc connection
	conn, err := d.getConn()
	if err != nil {
		d.errorResponse(method, w, err)
		return
	}

	// get volume to mount
	volumeClient := api.NewOpenStorageVolumeClient(conn)
	vol, err := d.volFromNameSdk(ctx, volumeClient, name)
	if err != nil {
		e := d.volNotFound(method, name, err, w)
		d.errorResponse(method, w, e)
		return
	}

	// create default unmount options
	opts := make(map[string]string)
	opts[options.OptionsDeleteAfterUnmount] = "true"
	opts[options.OptionsWaitBeforeDelete] = "false"

	// Unmount volume
	mountpoint := d.mountpath(name)
	mountClient := api.NewOpenStorageMountAttachClient(conn)
	if _, err = mountClient.Unmount(ctx, &api.SdkVolumeUnmountRequest{
		VolumeId:  vol.Id,
		MountPath: mountpoint,
		Options:   options.NewVolumeUnmountOptions(opts),
	}); err != nil {
		d.logRequest(method, request.Name).Warnf(
			"Cannot unmount volume %v, %v",
			mountpoint, err)
		d.errorResponse(method, w, err)
		return
	}

	d.emptyResponse(w)
}

func (d *driver) capabilities(w http.ResponseWriter, r *http.Request) {
	method := "capabilities"
	var response capabilitiesResponse

	response.Capabilities.Scope = "global"
	d.logRequest(method, "").Infof("response %v", response.Capabilities.Scope)
	json.NewEncoder(w).Encode(&response)
}

type CookieOnce struct {
	lock   sync.Mutex
	cookie map[string]string
	ttl    time.Duration
}

func NewCookieOnce(ttl time.Duration) *CookieOnce {
	return &CookieOnce{
		cookie: make(map[string]string),
		ttl:    ttl,
	}
}

func (c *CookieOnce) Pop(key string) (string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, ok := c.cookie[key]; ok {
		fmt.Printf("** pop %s : %s\n", key, val)
		delete(c.cookie, key)
		return val, true
	}
	return "", false
}

func (c *CookieOnce) Push(key, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cookie[key] = val
	fmt.Printf("** push %s : %s\n", key, val)

	// Start timer to remove key
	go func() {
		time.Sleep(c.ttl)
		c.Pop(key)
	}()
}
