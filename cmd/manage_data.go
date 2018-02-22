package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
)

var manageDataCmd = &cobra.Command{
	Use:   "manage-data [args]",
	Short: "set, modify or delete a Data Entry (name/value pair)",
	Long: `allows you to set, modify or delete a Data Entry (name/value pair) that is attached to a particular account.
an account can have an arbitrary amount of DataEntries attached to it.
each DataEntry increases the minimum balance needed to be held by the account.

this command takes parameters in JSON.

example: manage-data '{"source_account": "seed",
                       "name": "the answer to the ultimate question of life, the universe and everything",
                       "value": 42}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		md := &manageData{}
		if err := json.Unmarshal([]byte(args[0]), md); err != nil {
			log.Fatal(err)
		}

		if standAloneFlag {
			submitStandalone(conf, md.SourceAccount, md.makeOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, md.SourceAccount, md.makeOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], md.makeOp())
				fmt.Print(encoded)
			}
		}
	},
}

type manageData struct {
	SourceAccount string `json:"source_account"`
	Name          string
	Value         string
}

func (md *manageData) makeOp() (muts []b.TransactionMutator) {

	// source := seedToPair(md.SourceAccount)

	muts = []b.TransactionMutator{
		b.SourceAccount{AddressOrSeed: resolveAddress(md.SourceAccount)},
		b.SetData(md.Name, []byte(md.Value))}

	return muts
}
