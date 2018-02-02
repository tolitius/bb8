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

example: send-payment '{"from": "seed", "to": "address", "token": "XLM", "amount": "42.0"}'
         send-payment '{"from": "seed", "to": "address", "token": "BTC", "amount": "42.0", "issuer": "address"}'

         notice there is no issuer when sending XLM since it's a native asset.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payment := &tokenPayment{}
		if err := json.Unmarshal([]byte(args[0]), payment); err != nil {
			log.Fatal(err)
		}
		tx := payment.send(conf, parseOptions(txOptionsFlag))
		submitTransaction(conf.client, tx, payment.From)
	},
	DisableFlagsInUseLine: true,
}

type tokenPayment struct {
	From, To, Amount, Token, Issuer string
}

func (t *tokenPayment) send(conf *config, txOptions b.SetOptionsBuilder) *b.TransactionBuilder {

	log.Printf("sending %s %s from %s to %s", t.Amount, t.Token, seedToPair(t.From).Address(), t.To)

	var payment b.PaymentBuilder

	if t.Token == "XLM" && t.Issuer == "" {
		payment = b.Payment(
			b.Destination{t.To},
			b.NativeAmount{Amount: t.Amount},
		)
	} else {
		payment = b.Payment(
			b.Destination{t.To},
			b.CreditAmount{t.Token, t.Issuer, t.Amount},
		)
	}

	tx, err := b.Transaction(
		b.SourceAccount{t.From},
		conf.network,
		b.AutoSequence{conf.client},
		payment,
		txOptions,
	)

	if err != nil {
		log.Fatal(err)
	}

	return tx
}
