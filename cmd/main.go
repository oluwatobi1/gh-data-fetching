package main

import (
	"log"

	"github.com/oluwatobi1/gh-api-data-fetch/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
