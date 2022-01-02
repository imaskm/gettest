// +build integration
package integrationtests

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

func TestFetchRecordNilBody(t *testing.T) {
	endpoint := "/records"
	resp, err := http.Post(Server+endpoint, "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("failed record fetch on nil body")
	}
}

func TestFetchRecordValidBody(t *testing.T) {
	endpoint := "/records"

	var jsonStrValidBody = []byte(`{
		"startDate": "2012-01-26",
		"endDate": "2018-02-02",
		"minCount": 200,
		"maxCount": 3000
	}`)

	body := bytes.NewBuffer(jsonStrValidBody)

	resp, err := http.Post(Server+endpoint, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("failed record fetch")
	}

}

func TestFetchRecordInvalidBody(t *testing.T) {
	endpoint := "/records"

	var jsonStrInValidBody = []byte(`{
		"startDate": "2012-0126",
		"endDate": "2018-02-02",
		"minCount": 200,
		"maxCount": 3000
	}`)

	body := bytes.NewBuffer(jsonStrInValidBody)

	resp, err := http.Post(Server+endpoint, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("failed record fetch invalid body")
	}

}
