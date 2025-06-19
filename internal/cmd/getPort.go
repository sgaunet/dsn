package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var defaultPort string

// getCmd represents the get command.
var getPort = &cobra.Command{
	Use:   "port",
	Short: "get port of data source name",
	Long:  `get port of data source name`,
	Run: func(_ *cobra.Command, _ []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetPort(defaultPort))
	},
}
