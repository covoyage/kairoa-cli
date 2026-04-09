package main

import (
	"os"

	"github.com/covoyage/kairoa-cli/cmd"
)

// These variables are populated at build time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
