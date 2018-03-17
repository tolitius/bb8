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
	Use:   "create-account [args]",
	Short: "creates a new account",
	Long: `using a source account and a new accounts address create a new account
by sending an (initial) amount of XLM from the source account to the new account address.
hence needs a source account seed to sign this transaction.

example: create-account '{"source_account":"seed", "new_account":"address", "amount":"42.0"}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		naccount := &newAccount{}
		if err := json.Unmarshal([]byte(args[0]), naccount); err != nil {
			log.Fatal(err)
		}

		if standAloneFlag {
			submitStandalone(conf, naccount.SourceAccountSeed, naccount.makeCreateAccountOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, naccount.SourceAccountSeed, naccount.makeCreateAccountOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], naccount.makeCreateAccountOp())
				fmt.Print(encoded)
			}
		}
	},
}

var accountMergeCmd = &cobra.Command{
	Use:   "account-merge [ars]",
	Short: "merges two native (XLM) accounts",
	Long: `transfers the native balance (the amount of XLM an account holds)
to another account and removes the source account from the ledger.

example: account-merge '{"source_account":"seed", "destination":"address"}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		maccount := &accountMerge{}
		if err := json.Unmarshal([]byte(args[0]), maccount); err != nil {
			log.Fatal(err)
		}

		if standAloneFlag {
			submitStandalone(conf, maccount.SourceAccountSeed, maccount.makeAccountMergeOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, maccount.SourceAccountSeed, maccount.makeAccountMergeOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], maccount.makeAccountMergeOp())
				fmt.Print(encoded)
			}
		}
	},
}

func loadAccount(stellar *horizon.Client, address string) horizon.Account {

	uaddr := uniformAddress(address)

	account, err := stellar.LoadAccount(uaddr)
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

func (c *newAccount) makeCreateAccountOp() (muts []b.TransactionMutator) {

	address := uniformAddress(c.NewAccountAddress)

	muts = []b.TransactionMutator{
		b.CreateAccount(
			b.Destination{AddressOrSeed: address},
			b.NativeAmount{Amount: c.Amount},
		)}

	return muts
}

type accountMerge struct {
	SourceAccountSeed string `json:"source_account"`
	Destination       string
}

func (m *accountMerge) makeAccountMergeOp() (muts []b.TransactionMutator) {

	address := uniformAddress(m.Destination)

	muts = []b.TransactionMutator{
		b.AccountMerge(
			b.Destination{address},
		)}

	return muts
}
