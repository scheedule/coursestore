package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Coursestore",
	Long:  `This is the version of Coursestore`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("There is no version yet!")
	},
}
