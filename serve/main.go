package main

import (
	"fmt"
	"github.com/scheedule/coursestore/db"
	"net/http"
)

var mydb *db.DB

func init() {
	mydb = db.NewDB("mongo", "27017", "test", "classes")
	mydb.Init()
}

func printURI(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/lookup", printURI(HandleLookup(mydb)))
	http.HandleFunc("/all", printURI(HandleAll(mydb)))
	fmt.Println("Listening")
	http.ListenAndServe(":7819", nil)
}
