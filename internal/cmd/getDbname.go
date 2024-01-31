package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getDBName = &cobra.Command{
	Use:   "dbname",
	Short: "get dbname of data source name",
	Long:  `get dbname of data source name`,
	Run: func(cmd *cobra.Command, args []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetDBName())
	},
}
