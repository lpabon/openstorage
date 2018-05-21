package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/libopenstorage/openstorage/objectstore"
)

// swagger:operation GET /cluster/objectstore objectstore objectStoreInspect
//
// Lists Objectstore
//
// This will list current object stores
//
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: success
//     schema:
//      $ref: '#/definitions/ObjectstoreInfo'
func (c *clusterApi) objectStoreInspect(w http.ResponseWriter, r *http.Request) {
	method := "objectStoreInspect"

	objInfo, err := c.ObjectStoreManager.ObjectStoreInspect()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(objInfo)
}

// swagger:operation POST /cluster/objectstore objectstore objectStoreCreate
//
// Create an Object store
//
// This creates the volumes required to run the object store
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
//   description: volume on which object store to run
//   required: true
//   type: string
// responses:
//   '200':
//     description: success
func (c *clusterApi) objectStoreCreate(w http.ResponseWriter, r *http.Request) {
	method := "objectStoreCreate"
	params := r.URL.Query()
	volumeName := params[objectstore.VolumeName]

	if len(volumeName) == 0 || volumeName[0] == "" {
		c.sendError(c.name, method, w, "Missing volume name", http.StatusBadRequest)
		return
	}

	err := c.ObjectStoreManager.ObjectStoreCreate(volumeName[0])
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// swagger:operation PUT /cluster/objectstore objectstore objectStoreUpdate
//
// Updates object store
//
// This will enable/disable object store functionality.
//
// ---
// produces:
// - application/json
// parameters:
// - name: enable
//   in: query
//   description: enable/disable flag for object store
//   required: true
//   type: boolean
// responses:
//   '200':
//     description: success
func (c *clusterApi) objectStoreUpdate(w http.ResponseWriter, r *http.Request) {
	method := "objectStoreUpdate"

	params := r.URL.Query()
	strEnable := params[objectstore.Enable]
	if len(strEnable) == 0 && strEnable[0] == "" {
		c.sendError(c.name, method, w, "enable parameter not set", http.StatusInternalServerError)
		return
	}

	enable, err := strconv.ParseBool(strEnable[0])
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.ObjectStoreManager.ObjectStoreUpdate(enable)
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// swagger:operation DELETE /cluster/objectstore objectstore objectStoreDelete
//
// Delete object store
//
// This will delete object store on node
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: success
func (c *clusterApi) objectStoreDelete(w http.ResponseWriter, r *http.Request) {
	method := "objectStoreDelete"

	err := c.ObjectStoreManager.ObjectStoreDelete()
	if err != nil {
		c.sendError(c.name, method, w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
