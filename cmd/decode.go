package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode [base64-tx]",
	Short: "decodes a base64 encoded transaction",
	Long: `decodes a base64 encoded transaction.
transactions are usually encoded before they are sent to Stellar core.
this command decodes it back and shows it as JSON".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx := decodeXDR(args[0])
		fmt.Println(toJSON(tx))
	},
	DisableFlagsInUseLine: true,
}
