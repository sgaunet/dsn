package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getPassword = &cobra.Command{
	Use:   "password",
	Short: "get password of data source name",
	Long:  `get password of data source name`,
	Run: func(cmd *cobra.Command, args []string) {
		d := initDsnOrExit(dataSourceName)
		fmt.Println(d.GetPassword())
	},
}
