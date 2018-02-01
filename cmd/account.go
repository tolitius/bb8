package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stellar/go/clients/horizon"
)

var loadAccountCmd = &cobra.Command{
	Use:   "load-account [address]",
	Short: "load and return account details",
	Long:  `given an account address reach out to a Stellar network and return details about an account.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(toJSON(loadAccount(conf.client, args[0])))
	},
	DisableFlagsInUseLine: true,
}

func toJSON(foo interface{}) string {
	b, err := json.MarshalIndent(foo, "", "  ")
	if err != nil {
		log.Fatal("error:", err)
	}
	return string(b)
}

func loadAccount(stellar *horizon.Client, address string) horizon.Account {

	account, err := stellar.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	return account
}
