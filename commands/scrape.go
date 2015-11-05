package commands

import (
	"github.com/scheedule/coursestore/db"
	scrape_lib "github.com/scheedule/coursestore/scrape"
	"github.com/scheedule/coursestore/types"
	"github.com/spf13/cobra"
)

var scrape = &cobra.Command{
	Use:   "scrape",
	Short: "Fetch courses",
	Long:  "Fetch courses from from course API by scraping and parsing XML",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		if err := PopulateDB(termUrl, db_host, db_port, database, collection); err != nil {
			panic(err)
		}
	},
}

func init() {
	scrape.Flags().StringVarP(
		&termUrl, "term_url", "t",
		"http://courses.illinois.edu/cisapp/explorer/schedule/2016/spring.xml",
		"URL to term XML.")

	scrape.Flags().StringVarP(
		&db_host, "host", "", "localhost", "Hostname of DB to insert into.")

	scrape.Flags().StringVarP(
		&db_port, "port", "p", "27017", "Port to access DB on.")

	scrape.Flags().StringVarP(
		&database, "db", "", "test", "Database name.")

	scrape.Flags().StringVarP(
		&collection, "collection", "", "classes", "Collection in database to insert classes.")
}

// Populate the selected database with the data scraped from the TERM_URL.
func PopulateDB(term_url, ip, port, db_name, collection_name string) error {
	mydb := db.NewDB(ip, port, db_name, collection_name)

	err := mydb.Init()
	if err != nil {
		return err
	}

	mydb.Purge()

	term, err := scrape_lib.GetXML(term_url)

	course_chan := make(chan types.Class)

	go scrape_lib.DigestAll(term, course_chan)

	for class := range course_chan {
		err = mydb.Put(class)
		if err != nil {
			return err
		}
	}

	return nil
}
