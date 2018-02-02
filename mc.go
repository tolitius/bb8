package main

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func structValues(s interface{}) []interface{} {

	v := reflect.ValueOf(s)

	values := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsNil() {
			values = append(values, f.Interface())
		}
	}

	return values
}

type ledgerStreamer struct {
	Seconds int64
	Cursor  horizon.Cursor
	handler horizon.LedgerHandler
}

func (streamer *ledgerStreamer) streamLedger(stellar *horizon.Client) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if streamer.Seconds != 0 {
		go func() {
			time.Sleep(time.Duration(streamer.Seconds) * time.Second)
			cancel()
		}()
	}

	err := stellar.StreamLedgers(ctx, &streamer.Cursor, streamer.handler)

	if err != nil {
		log.Fatal(err)
	}
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

type newTransaction struct {
	Operations    txOperations
	SourceAccount string `json:"source-account"`
	Signers       []string
}

type txOperations struct {
	SourceAccount *b.SourceAccount
	//TODO: add all transaction operations
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

type txOptions struct {
	HomeDomain    *b.HomeDomain    `json:"home-domain"`
	MasterWeight  *b.MasterWeight  `json:"master-weight"`
	InflationDest *b.InflationDest `json:"inflation-destination"`
	//TODO: add all transaction options
}

func parseOptions(options string) b.SetOptionsBuilder {
	topts := &txOptions{}
	if err := json.Unmarshal([]byte(options), topts); err != nil {
		log.Fatal(err)
	}

	values := structValues(*topts)
	return b.SetOptions(values...)
}
