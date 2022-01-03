// +build integration
package integrationtests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestInMemoryPostNilBody(t *testing.T) {
	endpoint := "/in-memory"
	resp, err := http.Post(Server+endpoint, "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("in-memory Post failed on nil body")
	}
}

func TestInMemoryPostValidBody(t *testing.T) {
	endpoint := "/in-memory"

	var jsonStrValidBody = []byte(`{
		"key": "wierdKey",
		"value": "Ashwani"
	}`)
	resp, err := http.Post(Server+endpoint, "application/json", bytes.NewBuffer(jsonStrValidBody))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("in-memory Post failed  on valid body")
	}

	jsonStrValidBody = []byte(`{
		"key": "wierdKey",
		"value": "Ashwani12"
	}`)
	resp, err = http.Post(Server+endpoint, "application/json", bytes.NewBuffer(jsonStrValidBody))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusConflict {
		t.Error("in-memory Post failed  on valid existing key")
	}
}

func TestInMemoryPostInvalidBody(t *testing.T) {
	endpoint := "/in-memory"

	var jsonStrValidBody = []byte(`{
		"key": "name",
		"invalidvalue": "Ashwani"
	}`)
	resp, err := http.Post(Server+endpoint, "application/json", bytes.NewBuffer(jsonStrValidBody))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("in-memory Post failed on invalid body")
	}
}

func TestInMemoryGetNoKey(t *testing.T) {
	endpoint := "/in-memory"
	resp, err := http.Get(Server + endpoint)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("in-memory Get failed on no key")
	}
}

func TestInMemoryGetNilKey(t *testing.T) {
	endpoint := "/in-memory?key="
	resp, err := http.Get(Server + endpoint)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("in-memory Post failed on nil key")
	}
}

func TestInMemoryGetAbsentKey(t *testing.T) {
	endpoint := "/in-memory?key=NotPresent"
	resp, err := http.Get(Server + endpoint)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Error("in-memory Post failed on nil body")
	}
}

func TestInMemoryGetPresentKey(t *testing.T) {

	endpoint := "/in-memory"
	var jsonStrValidBody = []byte(`{
		"key": "name",
		"value": "Ashwani"
	}`)
	_, err := http.Post(Server+endpoint, "application/json", bytes.NewBuffer(jsonStrValidBody))
	if err != nil {
		log.Fatal(err)
	}
	endpoint = "/in-memory?key=name"
	resp, err := http.Get(Server + endpoint)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("in-memory Post failed on nil body")
	}
}
