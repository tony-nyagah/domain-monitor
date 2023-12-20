package main

import (
	"io"
	"log"
	"net/http"
)

type Domain struct {
	Name   string
	Status string
}

func main() {
	healthCheck := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "It works!")
	}
	http.HandleFunc("/healthcheck", healthCheck)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
