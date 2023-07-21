package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dsn",
	Short: "Tool to extract easily informations of a data source name",
	Long:  `Tool to extract easily informations of a data source name`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVar(&dataSourceName, "d", "", "data source name")

	getCmd.AddCommand(getScheme)
	getCmd.AddCommand(getDBName)
	getCmd.AddCommand(getHost)

	getPort.PersistentFlags().StringVar(&defaultPort, "p", "", "default port to print in case of no port in data source name")
	getCmd.AddCommand(getPort)
	getCmd.AddCommand(getUser)
	getCmd.AddCommand(getPassword)

	setEnvCmd.PersistentFlags().StringVar(&dataSourceName, "d", "", "data source name")
	setEnvCmd.PersistentFlags().StringVar(&defaultPort, "p", "", "default port to print in case of no port in data source name")
	setEnvCmd.PersistentFlags().StringVar(&prefixVar, "pr", "", "prefix for environment variable")
	rootCmd.AddCommand(setEnvCmd)
}
