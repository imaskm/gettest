package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/imaskm/getir/constants"
	"github.com/imaskm/getir/database"
	"github.com/imaskm/getir/inmem"

	"github.com/imaskm/getir/types"
)

type Controller struct {
	db *database.MongoDB
}

func NewController(db *database.MongoDB) *Controller {
	ctr := new(Controller)
	ctr.db = db

	return ctr
}

// SaveKeyValue returns success or failure
// @Param Body {types.Data}
// @Security None
// @Accept  json
// @Produce  json
// @router /in-memory [POST]
// @Success 201 "{types.Data}"
// @failure 400 "{error: error string}"
// @failure 409 "{error: error string}"
// @failure 500 "{error: error string}"
func (ctrl *Controller) SaveKeyValue(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	data := new(types.Data)

	err := decoder.Decode(data)
	if err != nil || data.Key == "" || data.Value == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "invalid input")))
		return
	}

	value, _ := inmem.Get(data.Key)
	if value != nil {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "key already exists")))
		return
	}

	inmem.Set(data)

	result, err := json.Marshal(data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "server error")))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("%v", string(result))))

}

// GetKeyValue returns Key-Value based on key parameter
// @QueryParam key
// @Security None
// @Accept  json
// @Produce  json
// @router /in-memory?key=value [GET]
// @Success 200 "{types.Data}"
// @failure 400 "{error: error string}"
// @failure 404 "{error: error string}"
// @failure 500 "{error: error string}"
func (ctrl *Controller) GetKeyValue(rw http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()["key"]

	if len(keys) != 1 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "bad request")))
		return
	}
	key := keys[0]
	if key == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "invalid input")))
		return
	}

	data, err := inmem.Get(key)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			rw.WriteHeader(http.StatusNotFound)
		}
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", err.Error())))
		return
	}

	result, err := json.Marshal(data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("{\"error\": %s }", "server error")))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("%v", string(result))))

}

func (ctrl *Controller) InMemory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		ctrl.GetKeyValue(rw, r)
	} else if r.Method == http.MethodPost {
		ctrl.SaveKeyValue(rw, r)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// Records fetches records based on filters
// @Param body	body	types.RecordRequest
// @Security None
// @Accept  json
// @Produce  json
// @router /records/ [POST]
// @Success 200 "{code:0, msg:success, records: []types.Record}"
// @failure 400 "{code:1, msg:failure, records: []}"
// @failure 404 "{code:1, msg:failure, records: []}"
// @failure 500 "{code:2, msg:failure, records: []}"

func (ctrl *Controller) Records(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(fmt.Sprintf("{ \"code\": %d, \"msg\": %s, \"records\": [] }", constants.ClientErrorCode, constants.Failure)))
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	requestBody := new(types.RecordRequest)

	err := decoder.Decode(requestBody)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{ \"code\": %d, \"msg\": %s, \"records\": [] }", constants.ClientErrorCode, constants.Failure)))
		return
	}

	records, err := ctrl.db.GetRecords(requestBody)

	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "invalid date") {
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		rw.Write([]byte(fmt.Sprintf("{ \"code\": %d, \"msg\": %s, \"records\": [] }", constants.ServerErrorCode, constants.Failure)))
		return
	}

	recordsJson, err := json.Marshal(records)
	if err != nil {
		log.Fatal(err)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("{ \"code\": %d, \"msg\": %s, \"records\": %v }", constants.SuccessCode, constants.Success, string(recordsJson))))
}
