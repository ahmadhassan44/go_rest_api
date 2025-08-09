package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting JSON server!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err.Error())
	}

	server := NewAPIServer(":3000", db)
	server.Listen()
}
