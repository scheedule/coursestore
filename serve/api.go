package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/scheedule/coursestore/db"
	"net/http"
)

func HandleLookup(db *db.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		department := r.FormValue("department")
		number := r.FormValue("number")

		if department == "" || number == "" {
			http.Error(w, "Specify both department and number", http.StatusBadRequest)
			return
		}

		fmt.Println("Dept:", department)
		fmt.Println("Num:", number)

		// TODO validate

		class, err := db.Lookup(department, number)
		if err != nil {
			http.NotFound(w, r)
		}

		js, err := json.Marshal(class)
		if err != nil {
			http.Error(w, "Oops..", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func HandleAll(db *db.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		classes, err := db.GetAll()
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(classes)
		if err != nil {
			http.Error(w, "Oops..", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gz.Write(js)
	}
}
