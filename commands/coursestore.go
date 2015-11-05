package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CoursestoreCmd = &cobra.Command{
	Use:   "coursestore",
	Short: "Course data endpoint",
	Long: "Coursestore is an endpoint that fetches and serves course" +
		"information for UIUC.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var Verbose bool
var termUrl, serve_port, db_host, db_port, database, collection string

//Initializes flags
func init() {
	CoursestoreCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func AddCommands() {
	CoursestoreCmd.AddCommand(version)
	CoursestoreCmd.AddCommand(scrape)
	CoursestoreCmd.AddCommand(serve)
}

func InitializeConfig() {
	if Verbose {
		log.SetLevel(log.InfoLevel)
	}
}

func Execute() {
	AddCommands()
	if err := CoursestoreCmd.Execute(); err != nil {
		panic(err)
	}
}
