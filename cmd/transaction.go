package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
)

var submitTransactionCmd = &cobra.Command{
	Use:   "submit [base64-encoded-transaction]",
	Short: "submit a base64 encoded transaction",
	Long:  `given a base64 encoded Stellar transaction submit it to the network.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		submitTransactionB64(conf.client, args[0])
	},
	DisableFlagsInUseLine: true,
}

var signTransactionCmd = &cobra.Command{
	Use:   "sign [signers, base64-encoded-transaction]",
	Short: "sign a base64 encoded transaction",
	Long:  `given a set of signers (seeds) and a base64 encoded Stellar transaction sign it with all the signers and encode it back to base64.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		var signers []string
		if err := json.Unmarshal([]byte(args[0]), &signers); err != nil {
			log.Fatal(err)
		}

		xdr := decodeXDR(args[1])
		envelope := wrapEnvelope(xdr, nil)
		envelope.MutateTX(conf.network)

		signEnvelope(envelope, signers...)
		encoded, err := envelope.Base64()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(encoded)
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

func decodeXDR(base64encoded string) *xdr.TransactionEnvelope {

	var tx xdr.TransactionEnvelope

	rawr := strings.NewReader(base64encoded)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	_, err := xdr.Unmarshal(b64r, &tx)

	if err != nil {
		log.Fatalf("could not decode a base64 XDR due to: %v", err)
	}

	return &tx
}

type txHeader struct {
	sourceAccount string
	sequence      *b.Sequence
}

func (header txHeader) newTx(conf *config) (muts []b.TransactionMutator) {

	var seq b.TransactionMutator = header.sequence

	if header.sequence == nil {
		seq = b.AutoSequence{conf.client}
	}

	muts = []b.TransactionMutator{
		b.Defaults{},
		b.SourceAccount{header.sourceAccount},
		seq,
		conf.network}

	return muts
}

func makeTransactionEnvelope(muts []b.TransactionMutator) *b.TransactionEnvelopeBuilder {

	txe := b.TransactionEnvelopeBuilder{}
	txe.Init()

	txe.MutateTX(muts...)
	txe.MutateTX(b.Defaults{})

	return &txe
}

func signEnvelope(envelope *b.TransactionEnvelopeBuilder, seeds ...string) {

	for _, seed := range seeds {
		signer := b.Sign{Seed: seed}
		err := signer.MutateTransactionEnvelope(envelope)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func submitEnvelope(envelope *b.TransactionEnvelopeBuilder, client *horizon.Client) int32 {

	txeB64, err := envelope.Base64()

	if err != nil {
		log.Fatal(err)
	}

	return submitTransactionB64(client, txeB64)
}

func wrapEnvelope(envelope *xdr.TransactionEnvelope, muts []b.TransactionMutator) *b.TransactionEnvelopeBuilder {

	txe := b.TransactionEnvelopeBuilder{E: envelope}
	txe.Init()

	if muts != nil {
		txe.MutateTX(muts...)
		txe.E.Tx.Fee = 0
		txe.MutateTX(b.Defaults{})
	}

	return &txe
}
