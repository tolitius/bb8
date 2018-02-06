package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
)

var changeTrustCmd = &cobra.Command{
	Use:   "change-trust [args]",
	Short: "create, update, or delete a trustline",
	Long: `create, update, or delete a trustline.
this command takes parameters in JSON and has an optional "limit" param, setting it to "0" removes the trustline.

example: change-trust '{"source_account": "seed", "code": "XYZ", "issuer_address": "address", "limit": "42.0"}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ct := &changeTrust{}
		if err := json.Unmarshal([]byte(args[0]), ct); err != nil {
			log.Fatal(err)
		}

		if standAloneFlag {
			header := txHeader{sourceAccount: ct.SourceAccount}.
				newTx(conf)

			ops := append(header, ct.makeOp()...)

			envelope := makeTransactionEnvelope(ops)
			signEnvelope(envelope, ct.SourceAccount)
			submitEnvelope(envelope, conf.client)
		} else {
			if len(args) == 1 {

				header := txHeader{sourceAccount: ct.SourceAccount}.
					newTx(conf)

				ops := append(header, ct.makeOp()...)
				envelope := makeTransactionEnvelope(ops)
				encoded, err := envelope.Base64()

				if err != nil {
					log.Fatal(err)
				}

				fmt.Print(encoded)

			} else {
				parent := decodeXDR(args[1])
				envelope := wrapEnvelope(parent, ct.makeOp())
				encoded, err := envelope.Base64()

				if err != nil {
					log.Fatal(err)
				}

				fmt.Print(encoded)
			}
		}
	},
}

type changeTrust struct {
	SourceAccount string `json:"source_account"`
	IssuerAddress string `json:"issuer_address"`
	Code, Limit   string
}

func (ct *changeTrust) makeOp() (muts []b.TransactionMutator) {

	// source := seedToPair(ct.SourceAccount)

	var limit = b.MaxLimit
	if ct.Limit != "" {
		limit = b.Limit(ct.Limit)
	}

	muts = []b.TransactionMutator{
		b.SourceAccount{AddressOrSeed: ct.SourceAccount},
		b.Trust(ct.Code, ct.IssuerAddress, limit)}

	return muts
}
