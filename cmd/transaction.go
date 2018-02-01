package cmd

import (
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

var submitTransactionCmd = &cobra.Command{
	Use:   "submit-tx [base64-encoded-transaction]",
	Short: "submit a base64 encoded transaction",
	Long:  `given a base64 encoded Stellar transaction submit it to the network.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		submitTransactionB64(conf.client, args[0])
	},
	DisableFlagsInUseLine: true,
}

func submitTransactionB64(stellar *horizon.Client, base64tx string) int32 {

	resp, err := stellar.SubmitTransaction(base64tx)

	if err != nil {
		log.Println(err)
		herr, isHorizonError := err.(*horizon.Error)
		if isHorizonError {
			resultCodes, err := herr.ResultCodes()
			if err != nil {
				log.Fatalln("failed to extract result codes from horizon response")
			}
			log.Fatalln(resultCodes)
		}
		log.Fatalln("could not submit the transaction")
	}

	return resp.Ledger
}

func submitTransaction(stellar *horizon.Client, txn *b.TransactionBuilder, seed ...string) int32 {

	var txe b.TransactionEnvelopeBuilder
	var err error

	if len(seed) > 0 {
		txe, err = txn.Sign(seed...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//TODO: refactor signing out to a pluggable func to be able to delegate it to external signers such as hardware wallets
		log.Fatal("can't find a seed to sign this transaction, and external / hardware signers are not yet supported")
	}

	txeB64, err := txe.Base64()

	if err != nil {
		log.Fatal(err)
	}

	return submitTransactionB64(stellar, txeB64)
}
