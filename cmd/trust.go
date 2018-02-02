package cmd

import (
	"encoding/json"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ct := &changeTrust{}
		if err := json.Unmarshal([]byte(args[0]), ct); err != nil {
			log.Fatal(err)
		}
		tx := ct.set(conf, parseOptions(txOptionsFlag))
		submitTransaction(conf.client, tx, ct.SourceAccount)
	},
	DisableFlagsInUseLine: true,
}

type changeTrust struct {
	SourceAccount string `json:"source_account"`
	IssuerAddress string `json:"issuer_address"`
	Code, Limit   string
}

func (ct *changeTrust) set(conf *config, txOptions b.SetOptionsBuilder) *b.TransactionBuilder {

	source := seedToPair(ct.SourceAccount)

	var limit = b.MaxLimit
	if ct.Limit != "" {
		limit = b.Limit(ct.Limit)
	}

	tx, err := b.Transaction(
		b.SourceAccount{source.Address()},
		b.AutoSequence{conf.client},
		conf.network,
		b.Trust(ct.Code, ct.IssuerAddress, limit),
		txOptions,
	)

	if err != nil {
		log.Fatal(err)
	}

	return tx
}
