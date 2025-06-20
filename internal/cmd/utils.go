package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/dsn/v3"
)

// initDsnOrExit creates a DSN object or exits with error code 1 if there's a problem.
// Note: We must return the interface type here since the concrete implementation is unexported.
func initDsnOrExit(dataSourceName string) *dsn.DSN {
	if dataSourceName == "" {
		fmt.Println("data source name is empty")
		os.Exit(1)
	}
	d, err := dsn.New(dataSourceName)
	if err != nil {
		os.Exit(1)
	}
	return d
}
