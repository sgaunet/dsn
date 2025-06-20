package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command.
var getUser = &cobra.Command{
	Use:   "user",
	Short: "get user of data source name",
	Long:  `get user of data source name`,
	Run: func(_ *cobra.Command, _ []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetUser())
	},
}
