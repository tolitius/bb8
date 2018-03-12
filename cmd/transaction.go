package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

const (
	accountAddressEnv = "STELLAR_ACCOUNT_ADDRESS"
	accountSeedEnv    = "STELLAR_ACCOUNT_SEED_FILE"
)

var submitTransactionCmd = &cobra.Command{
	Use:   "submit [base64-encoded-transaction]",
	Short: "submit a base64 encoded transaction",
	Long:  `given a base64 encoded Stellar transaction submit it to the network.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//TODO: fetch the right sequence for composed tx operations
		//      to get rid of these 4 lines
		xdr := decodeXDR(args[0])
		envelope := wrapEnvelope(xdr, nil)
		envelope.MutateTX(b.AutoSequence{conf.client})
		encoded, err := envelope.Base64()

		if err != nil {
			log.Fatal(err)
		}

		submitTransactionB64(conf.client, encoded)
	},
	DisableFlagsInUseLine: true,
}

var signTransactionCmd = &cobra.Command{
	Use:   "sign [signers, base64-encoded-transaction]",
	Short: "sign a base64 encoded transaction",
	Long:  `given a set of signers (seeds) and a base64 encoded Stellar transaction sign it with all the signers and encode it back to base64.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var signers []string
		var xdrRaw string

		if len(args) == 1 { // signers are not provided
			xdrRaw = args[0]
			seed, err := resolveSeed("") // resolve the default seed
			if err != nil {
				log.Fatalf("could not sign the transaction because seeds are not explicitely provided, and the default seed is not set")
			}
			signers = []string{seed}
		} else {
			xdrRaw = args[1]
			if err := json.Unmarshal([]byte(args[0]), &signers); err != nil {
				log.Fatal(err)
			}
		}

		xdr := decodeXDR(xdrRaw)
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

	log.Printf("submitting transaction to horizon at %s\n", conf.client.URL)

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

	log.Printf("submitting transaction to horizon at %s\n", conf.client.URL)

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
		err := validateSeed(seed)
		if err != nil {
			log.Fatalf("could not sign the transaction. seed is invalid: %v", err)
		}
		err = signer.MutateTransactionEnvelope(envelope)
		if err != nil {
			log.Fatalf("could not sign the transaction, make sure the seed(s) is provided: %v", err)
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

func submitStandalone(conf *config, sourceAccount string, muts []b.TransactionMutator) int32 {

	seed := sourceAccount
	err := validateSeed(seed)

	if err != nil {
		seed, err = resolveSeed("")

		if err != nil {
			log.Fatalf("could not find a valid seed to submit transaction: %v", err)
		}
	}

	header := txHeader{sourceAccount: seed}.
		newTx(conf)

	ops := append(header, muts...)

	envelope := makeTransactionEnvelope(ops)
	signEnvelope(envelope, seed)

	return submitEnvelope(envelope, conf.client)
}

func makeEnvelope(conf *config, sourceAccount string, muts []b.TransactionMutator) string {

	header := txHeader{sourceAccount: resolveAddress(sourceAccount)}.
		newTx(conf)

	ops := append(header, muts...)
	envelope := makeTransactionEnvelope(ops)
	encoded, err := envelope.Base64()

	if err != nil {
		log.Fatal(err)
	}

	return encoded
}

func composeWithOps(xdr string, muts []b.TransactionMutator) string {
	parent := decodeXDR(xdr)
	envelope := wrapEnvelope(parent, muts)
	encoded, err := envelope.Base64()

	if err != nil {
		log.Fatal(err)
	}

	return encoded
}

func resolveAddress(address string) string {

	if address != "" {
		return address
	}

	account := getEnv(accountAddressEnv, "")

	if account == "" {
		log.Fatalf("can't resolve Stellar account address (a.k.a. source account). you can set it via %s environment variable or provide it as a \"%s\" field of the transaction", accountAddressEnv, "source_account")
	}

	return account
}

func resolveSeed(seed string) (string, error) {

	if seed != "" {
		return seed, nil
	}

	seedFile := getEnv(accountSeedEnv, "")

	if seedFile == "" {
		return "", fmt.Errorf("can't find the account seed. you can either provide it explicitely in the transaction or set \"%v\" environment variable that points to a file with a seed", accountSeedEnv)
	}

	seedBytes, err := ioutil.ReadFile(seedFile)

	if err != nil {
		return "", fmt.Errorf("could not read %v seed file due to %v", seedFile, err)
	}

	resolvedSeed := string(seedBytes)

	err = validateSeed(resolvedSeed)

	if err != nil {
		return "", fmt.Errorf("read a seed from %v file, but it does not appear to be a valid seed: %v", seedFile, err)
	}

	return resolvedSeed, nil
}

func validateSeed(seed string) error {

	if seed == "" {
		return fmt.Errorf("the account seed (private key) is empty")
	}

	kp, err := keypair.Parse(seed)

	if err != nil {
		return fmt.Errorf("could not parse an account seed: %v", err)
	}

	switch v := kp.(type) {
	default:
		return fmt.Errorf("unexpected account seed type: %T", v)
	case *keypair.FromAddress:
		return fmt.Errorf("expected a seed, but a public key (a.k.a. account address) is provided instead: %v", seed)
	case *keypair.Full:
		return nil
	}
}
