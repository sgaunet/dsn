package cmd

import (
	"fmt"
	"os"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
)

func initDsnOrExit(dataSourceName string) dsn.DSN {
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
