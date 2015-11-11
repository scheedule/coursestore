package commands

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/scheedule/coursestore/api"
	"github.com/scheedule/coursestore/db"
)

var myAPI *api.API

// Main command to be executed. Serves coursestore endpoint.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Course Endpoint",
	Long:  "Start serving course data via routes /lookup and /all",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()

		// Create DB Object
		myDB := db.New(dbHost, dbPort, database, collection)
		err := myDB.Init()
		if err != nil {
			log.Fatal("Failed to initialize database connection:", err)
		}

		// API Object
		myAPI = &api.API{myDB}

		http.HandleFunc("/lookup", printURI(myAPI.HandleLookup))
		http.HandleFunc("/all", printURI(myAPI.HandleAll))
		log.Info("Serving on port:", servePort)
		http.ListenAndServe(":"+servePort, nil)
	},
}

func init() {
	serveCmd.Flags().StringVarP(
		&dbHost, "db_host", "", "localhost", "Hostname of DB to insert into.")

	serveCmd.Flags().StringVarP(
		&dbPort, "db_port", "", "27017", "Port to access DB on.")

	serveCmd.Flags().StringVarP(
		&servePort, "serve_port", "", "7819", "Port to serve endpoint on.")

	serveCmd.Flags().StringVarP(
		&database, "db_name", "", "test", "Database name.")

	serveCmd.Flags().StringVarP(
		&collection, "db_collection", "", "classes", "Collection in database to insert classes.")
}

// Middleware to print requests out to the console.
func printURI(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r)
		fn(w, r)
	}
}
