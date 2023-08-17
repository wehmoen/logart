package main

import (
	"github.com/wehmoen/logart/server"
	"log"
	"os"
)

const (
	DEFAULT_DB_URI = "mongodb://localhost:27017"
)

func main() {

	dbUriEnv := os.Getenv("LOGART_DB_URI")

	var dbUri string

	if dbUriEnv == "" {
		dbUri = DEFAULT_DB_URI
	} else {
		dbUri = dbUriEnv
	}

	logart, err := server.NewLogart(dbUri)
	if err != nil {
		log.Fatal(err)
	}

	logart.Start()
}
