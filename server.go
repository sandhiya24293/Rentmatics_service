package main

import (
	"crypto/tls"

	"encoding/json"

	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

type resStr struct {
	Reputation int   `json:"reputation"`
	Categories []int `json:"categories"`
}

const (
	contenttypeJSON = "application/json; charset=utf-8"
)

type APIError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

//Error provides error message in JSON format to client (utility)
func Error(w http.ResponseWriter, err int, msg string) {
	e := APIError{}
	e.Error.Code = err
	e.Error.Message = msg
	w.Header().Set("Content-Type", contenttypeJSON)
	w.WriteHeader(err)
	json.NewEncoder(w).Encode(e)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside Rent")

}

func Serve() bool {

	router := mux.NewRouter()

	//Mock-gti
	//router.Handle("POST", "/gti", post)
	router.HandleFunc("/home", home)

	// Default server - non-trusted for debugging
	serverhttp := func() {
		fmt.Println("Server should be available at http", config.Port)
		fmt.Println(http.ListenAndServe(config.Port, router))
	}

	// Setup TLS parameters for trusted server implementation
	if config.SSL && config.Key != "" && config.Cert != "" {
		// Setup TLS parameters
		tlsConfig := &tls.Config{
			ClientAuth:   tls.NoClientCert,
			MinVersion:   tls.VersionTLS12,
			Certificates: make([]tls.Certificate, 1),
		}

		var err error
		// Setup API server private key and certificate
		tlsConfig.Certificates[0], err = tls.X509KeyPair([]byte(config.Cert), []byte(config.Key))
		if err != nil {
			fmt.Println("Error during decoding service key and certificate:", err)
			return false
		}

		tlsConfig.BuildNameToCertificate()

		https := &http.Server{
			Addr:      config.Https_port,
			TLSConfig: tlsConfig,
			Handler:   router,
		}

		// Trusted server implementation
		server := func() {
			fmt.Println("Server should be available at https", config.Https_port)
			fmt.Println(https.ListenAndServeTLS("", ""))
		}
		go server()
	}

	// Schedule API server
	go serverhttp()

	return true
}
