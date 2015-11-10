package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/coursestore/api"
	"github.com/scheedule/coursestore/db"
	"github.com/spf13/cobra"
	"net/http"
)

var myapi *api.Api

// Main command to be executed. Serves coursestore endpoint.
var serve = &cobra.Command{
	Use:   "serve",
	Short: "Serve Course Endpoint",
	Long:  "Start serving course data via routes /lookup and /all",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		// Create DB Object
		mydb := db.NewDB(db_host, db_port, database, collection)
		err := mydb.Init()
		if err != nil {
			log.Fatal("Failed to initialize database connection:", err)
		}

		// API Object
		myapi = &api.Api{mydb}

		http.HandleFunc("/lookup", printURI(myapi.HandleLookup))
		http.HandleFunc("/all", printURI(myapi.HandleAll))
		log.Info("Serving on port:", serve_port)
		http.ListenAndServe(":"+serve_port, nil)
	},
}

func init() {
	serve.Flags().StringVarP(
		&db_host, "db_host", "", "localhost", "Hostname of DB to insert into.")

	serve.Flags().StringVarP(
		&db_port, "db_port", "", "27017", "Port to access DB on.")

	serve.Flags().StringVarP(
		&serve_port, "serve_port", "", "7819", "Port to serve endpoint on.")

	serve.Flags().StringVarP(
		&database, "db_name", "", "test", "Database name.")

	serve.Flags().StringVarP(
		&collection, "db_collection", "", "classes", "Collection in database to insert classes.")
}

// Middleware to print requests out to the console.
func printURI(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r)
		fn(w, r)
	}
}
