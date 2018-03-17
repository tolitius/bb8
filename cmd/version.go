package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bb8Version = "0.1.13"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of bb",
	Run: func(cmd *cobra.Command, args []string) {
		printBB8Version()
	},
	DisableFlagsInUseLine: true,
}

func printBB8Version() {
	//TODO: get the latest git tag and meta
	fmt.Printf("BB-8, version %s\n", bb8Version)
}
