// Package cmd implements the command-line interface for the dsn tool.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dataSourceName string

// getCmd represents the get command.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Println("please select a subcommand")
		_ = cmd.Help()
	},
}

// func init() {
// 	getCmd.PersistentFlags().StringVar(&dataSourceName, "d", "", "data source name")
// 	// setCmd.Flags().StringVar(&value, "v", "", "value to set")
// 	// setCmd.Flags().BoolVar(&createIniFileIfAbsent, "c", false, "create file if no present")
// }
