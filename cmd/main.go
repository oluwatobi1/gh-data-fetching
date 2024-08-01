package main

import (
	"log"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/cmd/app"
)

func main() {
	app := app.NewAPPServer()
	app.Run()

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
