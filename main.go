package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Domain struct {
	Name              string
	StatusCode        string
	CertificateExpiry string
}

func checkDomainStatus(domainName string) string {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://" + domainName)
	if err != nil {
		log.Println("Error:", err)
		return "Error: " + err.Error()
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)

	return strconv.Itoa(statusCode) + " - " + statusText
}

func checkCertificateExpiry(domainName string) string {
	conn, err := tls.Dial("tcp", domainName+":443", nil)
	if err != nil {
		return "Error: " + err.Error()
	}
	defer conn.Close()

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	expiryStr := expiry.Format("January 2, 2006")

	if expiry.After(time.Now()) {
		return expiryStr + " \u2705"
	} else {
		return expiryStr + " \u274C"
	}
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
		// domains := map[string][]Domain{
		// 	"Domains": {
		// 		{
		// 			Name:              "example.com",
		// 			StatusCode:        "200",
		// 			CertificateExpiry: "June 20, 2024",
		// 		},
		// 	},
		// }
		// tmpl.Execute(w, domains)
		tmpl.Execute(w, nil)
	}

	handler2 := func(w http.ResponseWriter, r *http.Request) {
		domainName := r.PostFormValue("domain-name")
		statusCode := checkDomainStatus(domainName)
		certificateExpiry := checkCertificateExpiry(domainName)

		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.ExecuteTemplate(w, "domain-list-element", Domain{Name: domainName, StatusCode: statusCode, CertificateExpiry: certificateExpiry})
	}

	http.HandleFunc("/", handler1)
	http.HandleFunc("/add-domain/", handler2)

	log.Fatal(http.ListenAndServe(port, nil))
}
