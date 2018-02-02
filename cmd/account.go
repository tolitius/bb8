package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
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

var createAccountCmd = &cobra.Command{
	Use:   "create-account [ars]",
	Short: "creates a new account",
	Long: `using a source account and a new accounts address create a new account
by sending an (initial) amount of XML from the source account to the new account address.
hence needs a source account seed to sign this transaction.

example: create-account '{"source_account":"seed", "new_account":"address", "amount":"42.0"}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		naccount := &newAccount{}
		if err := json.Unmarshal([]byte(args[0]), naccount); err != nil {
			log.Fatal(err)
		}
		tx := naccount.create(conf)
		submitTransaction(conf.client, tx, naccount.SourceAccountSeed)
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

type newAccount struct {
	SourceAccountSeed string `json:"source_account"`
	NewAccountAddress string `json:"new_account"`
	Amount            string
}

func (c *newAccount) create(conf *config) *b.TransactionBuilder {

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: c.SourceAccountSeed},
		conf.network,
		b.AutoSequence{conf.client}, //TODO: pass sequence if provided
		b.CreateAccount(
			b.Destination{AddressOrSeed: c.NewAccountAddress},
			b.NativeAmount{Amount: c.Amount},
		))

	if err != nil {
		log.Fatal(err)
	}

	return tx
}
