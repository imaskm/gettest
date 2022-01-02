package integrationtests

import (
	"log"
	"os"
)

var Server string

func init() {
	Server = os.Getenv("SERVER")
	if Server == "" {
		log.Println("server is not exported, using localhost:8081")
		Server = "http://localhost:8081"
	}
}
