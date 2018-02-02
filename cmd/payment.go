package cmd

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
)

var sendPaymentCmd = &cobra.Command{
	Use:   "send-payment [args]",
	Short: "send payment from one account to another",
	Long: `send payment of any asset from one account to another. this command takes parameters in JSON.
example: send-payment '{"from": "seed", "to": "address", "token": "BTC", "amount": "42.0", "issuer": "address"}'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//TODO: add transactionOptionsFlag
		var txOptionsBuilder b.SetOptionsBuilder
		// if txOptionsFlag != "" {
		// 	txOptionsBuilder = parseOptions(txOptions)
		// }

		payment := &tokenPayment{}
		if err := json.Unmarshal([]byte(args[0]), payment); err != nil {
			log.Fatal(err)
		}
		tx := payment.send(conf, txOptionsBuilder)
		submitTransaction(conf.client, tx, payment.From)
	},
	DisableFlagsInUseLine: true,
}

type tokenPayment struct {
	From, To, Amount, Token, Issuer string
}

func (t *tokenPayment) send(conf *config, txOptions b.SetOptionsBuilder) *b.TransactionBuilder {

	log.Printf("sending %s %s from %s to %s", t.Amount, t.Token, seedToPair(t.From).Address(), t.To)

	asset := b.CreditAsset(t.Token, t.Issuer)

	tx, err := b.Transaction(
		b.SourceAccount{t.From},
		conf.network,
		b.AutoSequence{conf.client},
		b.Payment(
			b.Destination{t.To},
			b.CreditAmount{asset.Code, asset.Issuer, t.Amount},
		),
		txOptions,
	)

	if err != nil {
		log.Fatal(err)
	}

	return tx
}
