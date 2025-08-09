package main

import "fmt"

func main() {
	fmt.Println("Starting JSON server!")
	server := NewAPIServer(":3000")
	server.Listen()
}
