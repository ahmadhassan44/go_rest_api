package main

import (
	"fmt"
	"log"

	"github.com/ahmadhassan44/go_rest_api/api"
	"github.com/ahmadhassan44/go_rest_api/storage"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting JSON server!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", db)
	server.Listen()
}
