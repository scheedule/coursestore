// Package main utilizes the api package to serve course queries
//
// Must set the following environment variables
//
// SERVE_PORT: <7819> The port on which to serve endpoint.
package main

import (
	"fmt"
	"github.com/scheedule/coursestore/api"
	"net/http"
	"os"
)

// Port we wish to serve on.
var SERVE_PORT string = func() string {
	if s := os.Getenv("SERVE_PORT"); s != "" {
		return s
	}
	return "7819"
}()

// Middleware to print requests out to the console.
func printURI(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		fn(w, r)
	}
}

// Entrypoint. Start listening.
func main() {
	http.HandleFunc("/lookup", printURI(api.HandleLookup))
	http.HandleFunc("/all", printURI(api.HandleAll))
	fmt.Println("Listening")
	http.ListenAndServe(":"+SERVE_PORT, nil)
}
