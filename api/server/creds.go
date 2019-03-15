package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/libopenstorage/openstorage/api"
)

func (vd *volAPI) credsEnumerate(w http.ResponseWriter, r *http.Request) {
	method := "credsEnumerate"

	// Get context with auth token
	ctx, err := vd.annotateContext(r)
	if err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get gRPC connection
	conn, err := vd.getConn()
	if err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	credentials := api.NewOpenStorageCredentialsClient(conn)

	// This returns which of the credentials the caller has access to.
	resp, err := credentials.Enumerate(ctx, &api.SdkCredentialEnumerateRequest{})
	if err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusInternalServerError)
	}

	// Now get all the the creds and their data
	allcreds, err := d.CredsEnumerate()
	if err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We cannot return an SDK object in this call (long story). Instead,
	// use these keys as the method to return the data from the golang API.
	creds := make(map[string]interface{})
	for _, credid := range resp.GetCredentialIds() {
		creds[credid] = allcreds[credid]
	}

	json.NewEncoder(w).Encode(creds)
}

func (vd *volAPI) credsCreate(w http.ResponseWriter, r *http.Request) {
	method := "credsCreate"
	var input api.CredCreateRequest
	response := &api.CredCreateResponse{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusBadRequest)
		return
	}
	d, err := vd.getVolDriver(r)
	if err != nil {
		notFound(w, r)
		return
	}

	response.UUID, err = d.CredsCreate(input.InputParams)
	if err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func (vd *volAPI) credsDelete(w http.ResponseWriter, r *http.Request) {
	method := "credsDelete"
	vars := mux.Vars(r)
	uuid, ok := vars["uuid"]
	if !ok {
		vd.sendError(vd.name, method, w, "Could not parse form for uuid", http.StatusBadRequest)
		return
	}

	d, err := vd.getVolDriver(r)
	if err != nil {
		notFound(w, r)
		return
	}
	if err = d.CredsDelete(uuid); err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (vd *volAPI) credsValidate(w http.ResponseWriter, r *http.Request) {
	method := "credsValidate"
	vars := mux.Vars(r)
	uuid, ok := vars["uuid"]
	if !ok {
		vd.sendError(vd.name, method, w, "Could not parse form for uuid", http.StatusBadRequest)
		return
	}
	d, err := vd.getVolDriver(r)
	if err != nil {
		notFound(w, r)
		return
	}

	if err := d.CredsValidate(uuid); err != nil {
		vd.sendError(vd.name, method, w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}
