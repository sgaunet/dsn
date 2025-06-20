package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command.
var getScheme = &cobra.Command{
	Use:   "scheme",
	Short: "get scheme of data source name",
	Long:  `get scheme of data source name`,
	Run: func(_ *cobra.Command, _ []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetScheme())
	},
}
