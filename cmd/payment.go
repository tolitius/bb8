package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
)

var sendPaymentCmd = &cobra.Command{
	Use:   "send-payment [args]",
	Short: "send payment from one account to another",
	Long: `send payment of any asset from one account to another. this command takes parameters in JSON.

example: send-payment '{"from": "seed", "to": "address", "amount": "42.0"}'
         send-payment '{"from": "seed", "to": "address", "amount": "42.0", "memo": "forty two"}'
         send-payment '{"from": "seed", "to": "address", "token": "BTC", "amount": "42.0", "issuer": "address"}'

         notice there is no issuer when sending XLM since it's a native asset.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payment := &tokenPayment{}
		if err := json.Unmarshal([]byte(args[0]), payment); err != nil {
			log.Fatalf("could not parse transaction JSON data: %v", err)
		}

		if standAloneFlag {
			submitStandalone(conf, payment.From, payment.makeOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, payment.From, payment.makeOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], payment.makeOp())
				fmt.Print(encoded)
			}
		}
	},
}

type tokenPayment struct {
	From, To, Amount, Token, Issuer, Memo string
}

func (t *tokenPayment) makeOp() (muts []b.TransactionMutator) {

	if t.Token == "" {
		t.Token = "XLM"
	}

	from := resolveAddress(t.From)
	to := uniformAddress(t.To)

	if from == "" {
		log.Fatalf("can't send payment, missing a source (a.k.a. account address/seed)")
	}

	log.Printf("sending %s %s from %s to %s", t.Amount, t.Token, seedToPair(from).Address(), t.To)

	var payment b.PaymentBuilder
	var memo = b.MemoText{Value: t.Memo}

	if t.Token == "XLM" && t.Issuer == "" {
		payment = b.Payment(
			b.Destination{to},
			b.NativeAmount{Amount: t.Amount},
		)
	} else {
		payment = b.Payment(
			b.Destination{to},
			b.CreditAmount{t.Token, t.Issuer, t.Amount},
		)
	}

	muts = []b.TransactionMutator{
		b.SourceAccount{from},
		payment,
		memo}

	return muts
}
