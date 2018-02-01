package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	smccVersion = "0.1.2"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of stellar mc",
	Run: func(cmd *cobra.Command, args []string) {
		printStellarMcVersion()
	},
	DisableFlagsInUseLine: true,
}

func printStellarMcVersion() {
	//TODO: get the latest git tag and meta
	fmt.Printf("Stellar Mission Control Center, version %s\n", smccVersion)
}
