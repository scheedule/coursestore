package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coursestoreCmd = &cobra.Command{
	Use:   "coursestore",
	Short: "Course data endpoint",
	Long: "Coursestore is an endpoint that fetches and serves course" +
		"information for UIUC.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var verbose bool
var termURL, servePort, dbHost, dbPort, database, collection string

//Initializes flags
func init() {
	coursestoreCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func addCommands() {
	coursestoreCmd.AddCommand(versionCmd)
	coursestoreCmd.AddCommand(scrapeCmd)
	coursestoreCmd.AddCommand(serveCmd)
}

func initializeConfig() {
	if verbose {
		log.SetLevel(log.InfoLevel)
	}
}

func Execute() {
	addCommands()
	if err := coursestoreCmd.Execute(); err != nil {
		panic(err)
	}
}
