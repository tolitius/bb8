package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
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

func seedToPair(seed string) keypair.KP {

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatal(err)
	}

	return kp
}

func toJSON(foo interface{}) string {
	b, err := json.MarshalIndent(foo, "", "  ")
	if err != nil {
		log.Fatal("error:", err)
	}
	return string(b)
}

func loadAccount(stellar *horizon.Client, address string) horizon.Account {

	account, err := stellar.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	return account
}

func fundTestAccount(stellar *horizon.Client, address string) {

	resp, err := http.Get(stellar.URL + "/friendbot?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("could not fund %s, horizon said: %s\n", address, string(body))
	}
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

type changeTrust struct {
	SourceAccount string `json:"source-account"`
	IssuerAddress string `json:"issuer-address"`
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
