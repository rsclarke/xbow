// Package main is the entry point for the xbow CLI.
package main

import (
	"os"

	"github.com/rsclarke/xbow/cmd/xbow/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
