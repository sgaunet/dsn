package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/dsn/pkg/dsn"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getUser = &cobra.Command{
	Use:   "user",
	Short: "get user of data source name",
	Long:  `get user of data source name`,
	Run: func(cmd *cobra.Command, args []string) {
		if dataSourceName == "" {
			fmt.Println("data source name is empty")
			os.Exit(1)
		}
		d, err := dsn.New(dataSourceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(d.GetUser())
	},
}
