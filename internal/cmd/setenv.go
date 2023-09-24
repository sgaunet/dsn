package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"github.com/spf13/cobra"
)

var prefixVar string

// getCmd represents the get command
var setEnvCmd = &cobra.Command{
	Use:   "setenv",
	Short: "generate a bash script to set environment variables",
	Long:  `generate a bash script to set environment variables`,
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
		PrintVar("SCHEME", d.GetScheme())
		PrintVar("DBNAME", d.GetDBName())
		PrintVar("HOST", d.GetHost())
		PrintVar("PORT", d.GetPort(defaultPort))
		PrintVar("USER", d.GetUser())
		PrintVar("PASSWORD", d.GetPassword())
	},
}

func PrintVar(name, value string) {
	fmt.Printf("export %s%s=%s\n", prefixVar, name, value)
}
