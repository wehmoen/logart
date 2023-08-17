package main

import (
	"context"
	"flag"
	"github.com/wehmoen/logart/database"
	"log"
	"os"
)

const (
	DEFAULT_DB_URI = "mongodb://localhost:27017"
)

func main() {
	username := flag.String("username", "", "Username")
	if *username == "" {
		log.Fatalf("Username is required via --username")
	}

	dbUriEnv := os.Getenv("LOGART_DB_URI")

	var dbUri string

	if dbUriEnv == "" {
		dbUri = DEFAULT_DB_URI
	} else {
		dbUri = dbUriEnv
	}

	db, err := database.Open(dbUri, "logart")

	if err != nil {
		log.Fatal(err)
	}

	apiKey, err := db.CreateUser(*username)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User %s created with api key %s", *username, apiKey)

	_ = db.Client.Disconnect(context.Background())

}
