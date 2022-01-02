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

func (ctrl *Controller) SaveKeyValue(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	data := new(types.Data)

	err := decoder.Decode(data)
	if err != nil || data.Key == "" || data.Value == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"msg\": %s }", constants.Failure)))
		return
	}

	inmem.Set(data)
	rw.WriteHeader(http.StatusNoContent)
	rw.Write([]byte(fmt.Sprintf("{\"msg\": %s }", constants.Success)))
}

func (ctrl *Controller) GetKeyValue(rw http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()["key"]

	if len(keys) != 1 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"msg\": %s }", constants.Failure)))
		return
	}
	key := keys[0]
	if key == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("{\"msg\": %s }", constants.Failure)))
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
