package main

import (
	"log"
	"net/http"
	"os"

	"github.com/imaskm/getir/controllers"
	"github.com/imaskm/getir/database"
)

func main() {

	db, err := database.NewMongoSession()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	ctr := controllers.NewController(db)

	http.HandleFunc("/records", ctr.Records)
	http.HandleFunc("/in-memory", ctr.InMemory)
	log.Println("started server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
