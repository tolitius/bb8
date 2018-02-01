package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// StellarCmd is Stellar Mission Control Center's root command.
// Other commands are added to StellarCmd as subcommands.
var StellarCmd = &cobra.Command{
	Use:   "stellar-mc",
	Short: "cli to interact with Stellar network",
	Long:  `stellar is a command line interface to Stellar (https://www.stellar.org/) networks.`,
}

// AddCommands adds sub commands to StellarCmd
func AddCommands() {
	StellarCmd.AddCommand(versionCmd)
	StellarCmd.AddCommand(genKeysCmd)
}

// Execute adds sub commands to StellarCmd and sets all the command line flags
func Execute() {

	StellarCmd.SilenceUsage = true

	AddCommands()

	if c, err := StellarCmd.ExecuteC(); err != nil {
		c.Println("")
		c.Println(c.UsageString())
		os.Exit(-1)
	}
}
