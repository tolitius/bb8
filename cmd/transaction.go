package cmd

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

var newTransactionCmd = &cobra.Command{
	Use:   "new-tx [args]",
	Short: "build and submit a new transaction",
	Long: `build and submit a new transaction. this command takes parameters in JSON.
"operations" and "signers" are optional, if there are no "signers", the "source_account" seed will be used to sign this transaction.

example: new-tx '{"source_account": "address or seed", {"operations": "trust": {"code": "XYZ", "issuer_address": "address"}}, "signers": ["seed1", "seed2"]}'
         new-tx '{"source_account": "address or seed"}' --set-options '{"home_domain": "stellar.org", "max_weight": 1}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		nt := &newTransaction{}
		if err := json.Unmarshal([]byte(args[0]), nt); err != nil {
			log.Fatal(err)
		}
		nt.Operations.SourceAccount = &b.SourceAccount{nt.SourceAccount}
		tx := nt.Operations.buildTransaction(conf, parseOptions(txOptionsFlag))
		signers := nt.Signers
		if signers == nil {
			signers = []string{nt.SourceAccount}
		}
		submitTransaction(conf.client, tx, signers...)
	},
}

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

type txOperations struct {
	SourceAccount *b.SourceAccount
	//TODO: add all transaction operations
}

type newTransaction struct {
	Operations    txOperations
	SourceAccount string `json:"source_account"`
	Signers       []string
}

func (t *txOperations) toMutators() []b.TransactionMutator {

	values := structValues(*t)
	muts := make([]b.TransactionMutator, len(values))

	for i := 0; i < len(values); i++ {
		switch values[i].(type) {
		case b.TransactionMutator:
			muts[i] = values[i].(b.TransactionMutator)
		default:
			log.Fatalf("%+v is not a valid transaction operation", values[i])
		}
	}

	return muts
}

func (t *txOperations) buildTransaction(
	conf *config,
	options b.SetOptionsBuilder) *b.TransactionBuilder {

	tx, err := b.Transaction(
		t.SourceAccount,
		conf.network,
		b.AutoSequence{conf.client}, //TODO: pass sequence if provided
		options)

	if err != nil {
		log.Fatal(err)
	}

	tx.Mutate(t.toMutators()...)

	return tx
}
