package main

import (
	"crypto/tls"

	"encoding/json"

	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/dimfeld/httptreemux"
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

func notAllowed(w http.ResponseWriter, r *http.Request, methods map[string]httptreemux.HandlerFunc) {
	fmt.Println("Incorrect resource request (method): ", r.Method, " ", r.RequestURI, " from ", r.RemoteAddr)
	Error(w, 405, "Method Not Allowed")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incorrect resource request: ", r.Method, " ", r.RequestURI, " from ", r.RemoteAddr)
	Error(w, 404, "API Not Found")
}
func getdomain(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	id := ps["domain"]
	urlStr := "http://10.137.0.6/local/gti/" + id + "/rate"
	fmt.Println("URL IS : ", urlStr)

	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error : ", err)
		return
	}
	fmt.Println("content : ", string(content))
	res := new(resStr)
	if err = json.Unmarshal(content, &res); err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Printf("Response is %+v: ", res.Categories, res.Reputation)
	json.NewEncoder(w).Encode(res)
	return

}

func Serve() bool {

	router := httptreemux.New()
	router.MethodNotAllowedHandler = notAllowed
	router.NotFoundHandler = notFound

	//Mock-gti
	//router.Handle("POST", "/gti", post)
	router.Handle("GET", "/local/gti/:domain/rate", getdomain)

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
