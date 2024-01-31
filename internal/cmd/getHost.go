package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getHost = &cobra.Command{
	Use:   "host",
	Short: "get host of data source name",
	Long:  `get host of data source name`,
	Run: func(cmd *cobra.Command, args []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetHost())
	},
}
