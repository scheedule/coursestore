package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/scheedule/coursestore/db"
	"github.com/scheedule/coursestore/scrape"
	"github.com/scheedule/coursestore/types"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Fetch courses",
	Long:  "Fetch courses from from course API by scraping and parsing XML",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()

		if err := PopulateDB(termURL, dbHost, dbPort, database, collection); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	scrapeCmd.Flags().StringVarP(
		&termURL, "term_url", "t",
		"http://courses.illinois.edu/cisapp/explorer/schedule/2016/spring.xml",
		"URL to term XML.")

	scrapeCmd.Flags().StringVarP(
		&dbHost, "host", "", "localhost", "Hostname of DB to insert into.")

	scrapeCmd.Flags().StringVarP(
		&dbPort, "port", "p", "27017", "Port to access DB on.")

	scrapeCmd.Flags().StringVarP(
		&database, "db", "", "test", "Database name.")

	scrapeCmd.Flags().StringVarP(
		&collection, "collection", "", "classes", "Collection in database to insert classes.")
}

// Populate the selected database with the data scraped from the term URL.
func PopulateDB(termURL, ip, port, dbName, collectionName string) error {
	scrapeDB := db.New(ip, port, dbName, collectionName)

	err := scrapeDB.Init()
	if err != nil {
		return err
	}

	log.Debug("purging database")
	scrapeDB.Purge()

	term, err := scrape.GetXML(termURL)

	courseChan := make(chan types.Class)

	go scrape.DigestAll(term, courseChan)

	for class := range courseChan {
		err = scrapeDB.Put(class)
		if err != nil {
			return err
		}
	}

	log.Debug("finished populating database")

	return nil
}
