package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Domain struct {
	Name              string
	StatusCode        int
	CertificateExpiry string
}

func main() {
	port := ":8080"
	log.Println("Starting server on port", strings.Split(port, ":")[1])

	healthCheck := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "It works!")
	}
	http.HandleFunc("/healthcheck", healthCheck)

	handler1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		domains := map[string][]Domain{
			"Domains": {
				{
					Name:              "huqas.org",
					StatusCode:        200,
					CertificateExpiry: "June 20, 2024",
				},
			},
		}
		tmpl.Execute(w, domains)
	}

	handler2 := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.ExecuteTemplate(w, "domain-list-element", Domain{Name: name, StatusCode: 200, CertificateExpiry: "June 20, 2024"})
	}

	http.HandleFunc("/", handler1)
	http.HandleFunc("/add-domain", handler2)

	log.Fatal(http.ListenAndServe(port, nil))
}
