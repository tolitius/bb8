package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/stellar/go/clients/federation"
	"github.com/stellar/go/clients/stellartoml"
)

var address string
var account string

var federationCmd = &cobra.Command{
	Use:   "federation",
	Short: "lookup federation addresses",
	Long: `convert federation addresses such as luke_skywalker*scoutship.com into stellar public address and back

example: federation --address luke_skywalker*scoutship.com
         federation --account GBKOPETTWWVE7DM72YOVQ4M2UIY3JCKDYQBTSNLGLHI6L43K7XPDROID`,
	Run: func(cmd *cobra.Command, args []string) {

		if address == "" && account == "" {
			log.Fatal("federation command needs one of the following flags: --address/-a, --account/-s")
		}

		if address != "" {
			resp, err := NewFederationClient().LookupByAddress(address)
			if err != nil {
				log.Fatalf("could not resolve federation address %s due to %v", address, err)
			}

			fmt.Println(resp.AccountID)
		}

		if account != "" {
			resp, err := NewFederationClient().LookupByAccountID(account)
			if err != nil {
				log.Fatalf("could not resolve federation stellar account id %s due to %v", account, err)
			}

			fmt.Println(resp.Address)
		}
	},
}

func init() {
	federationCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "convert federation address to a stellar public account. example: --address luke_skywalker*scoutship.com")
	federationCmd.PersistentFlags().StringVarP(&account, "account", "s", "", "convert stellar account to a federation address. example: --account GBKOPETTWWVE7DM72YOVQ4M2UIY3JCKDYQBTSNLGLHI6L43K7XPDROID")
}

func NewFederationClient() *federation.Client {
	return &federation.Client{
		HTTP:        http.DefaultClient,
		Horizon:     conf.client,
		StellarTOML: stellartoml.DefaultClient,
	}
}
