package commands

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/scheedule/coursestore/api"
	"github.com/scheedule/coursestore/db"
)

var serveAPI *api.API

// Main command to be executed. Serves coursestore endpoint.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Course Endpoint",
	Long:  "Start serving course data via routes /lookup and /all",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()

		// Create DB Object
		serveDB := db.New(dbHost, dbPort, database, collection)
		err := serveDB.Init()
		if err != nil {
			log.Fatal("Failed to initialize database connection:", err)
		}

		// API Object
		serveAPI = api.New(serveDB)

		http.HandleFunc("/lookup", middleware(serveAPI.HandleLookup, logMiddlewareHandler))
		http.HandleFunc("/all", middleware(serveAPI.HandleAll, logMiddlewareHandler))
		log.Info("Serving on port:", servePort)
		http.ListenAndServe(":"+servePort, nil)
	},
}

func middleware(primaryHandler http.HandlerFunc, middlewareHandlers ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	var next http.HandlerFunc
	for _, mw := range middlewareHandlers {
		next = mw(primaryHandler)
	}
	return next
}

func logMiddlewareHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// stuff
		log.Debug("recieved request:", r.Method, r.URL)

		// Move to next handler
		h.ServeHTTP(w, r)
	}
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
