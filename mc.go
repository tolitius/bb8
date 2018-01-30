package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

func readConfig(cpath string) *config {
	//TODO: read config from cpath

	return &config{
		client:  horizon.DefaultTestNetClient,
		network: b.Network{network.TestNetworkPassphrase}}
}

func seedToPair(seed string) keypair.KP {

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatal(err)
	}

	return kp
}

func createNewKeys(fpath string) string {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	fpub, err := os.Create(fpath + ".pub")
	if err != nil {
		log.Fatal(err)
	}

	fseed, err := os.Create(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer fpub.Close()
	defer fseed.Close()

	fmt.Fprint(fpub, pair.Address())
	fmt.Fprint(fseed, pair.Seed())

	fpub.Sync()
	fseed.Sync()

	log.Printf("keys are created and stored in: %s and %s\n", fpub.Name(), fseed.Name())

	return fpath
}

func loadAccount(stellar *horizon.Client, address, message string) horizon.Account {

	account, err := stellar.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%saccount: %s.\nbalances: %+v\nmore details: %s",
		message, address, account.Balances, account.Links.Self.Href)

	return account
}

func fundTestAccount(stellar *horizon.Client, address string) horizon.Account {

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

	return loadAccount(stellar, address, "successfully funded ")
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

func submitTransaction(stellar *horizon.Client, txn *b.TransactionBuilder, seed string) int32 {

	txe, err := txn.Sign(seed)

	if err != nil {
		log.Fatal(err)
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

func (t *tokenPayment) send(conf *config) horizon.Account {

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
	)

	if err != nil {
		log.Fatal(err)
	}

	submitTransaction(conf.client, tx, t.From)

	receiver := loadAccount(conf.client, t.To, "... [payment sent]\n")

	return receiver
}

type newToken struct {
	symbol, issuerSeed, distributorSeed, limit string
}

func (t *newToken) issueNew(conf *config) b.Asset {

	issuer := seedToPair(t.issuerSeed)
	distributor := seedToPair(t.distributorSeed)

	asset := b.CreditAsset(t.symbol, issuer.Address())

	var limit = b.MaxLimit

	if t.limit != "" {
		limit = b.Limit(t.limit)
	}

	tx, err := b.Transaction(
		b.SourceAccount{distributor.Address()},
		b.AutoSequence{conf.client},
		conf.network,
		b.Trust(t.symbol, issuer.Address(), limit),
	)

	if err != nil {
		log.Fatal(err)
	}

	submitTransaction(conf.client, tx, t.distributorSeed)

	loadAccount(conf.client, distributor.Address(), fmt.Sprintf("issued trust for %s to ", t.symbol))

	return asset
}

type txOptions struct {
	HomeDomain   *b.HomeDomain
	MasterWeight *b.MasterWeight
	//TODO: add all transaction options
}

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

func parseOptions(options string) b.SetOptionsBuilder {
	topts := &txOptions{}
	if err := json.Unmarshal([]byte(options), topts); err != nil {
		log.Fatal(err)
	}

	values := structValues(*topts)
	builder := b.SetOptions(values...)

	fmt.Printf("builder: %+v", builder)

	return builder
}

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {
	var fund string
	var keyFpath string
	var txToSubmit string
	var issueToken string
	var sendPayment string
	var txOptions string

	flag.StringVar(&fund, "fund", "", "funds a test account. example: --fund address")
	flag.StringVar(&keyFpath, "gen-keys", "", "creates a pair of keys (in two files \"file-path\" and \"file-path.pub\"). example: --gen-keys file-path")
	flag.StringVar(&txToSubmit, "submit-tx", "", "submits a base64 encoded transaction. example: --submit-tx txn")
	flag.StringVar(&issueToken, "issue-new-token", "", "issue new token (asset). example: --issue-new-token token issuer-seed distributor-seed [limit]")
	flag.StringVar(&sendPayment, "send-payment", "", "send payment from one account to another. example: --send-payment '{\"from\": \"seed\", \"to\": \"address\", \"token\": \"BTC\", \"amount\": \"42.0\", \"issuer-address\": \"address\"}'")
	flag.StringVar(&txOptions, "tx-options", "", "add one or more transaction options. example: --tx-options '{\"homeDomain\": \"stellar.org\", \"maxWeight\": 1}'")

	flag.Parse()

	conf := readConfig("/tmp/todo")

	switch {
	case fund != "":
		fundTestAccount(conf.client, fund)
	case keyFpath != "":
		createNewKeys(keyFpath)
	case txToSubmit != "":
		submitTransactionB64(conf.client, txToSubmit)
	case sendPayment != "":
		payment := &tokenPayment{}
		if err := json.Unmarshal([]byte(sendPayment), payment); err != nil {
			log.Fatal(err)
		}
		payment.send(conf)
	case issueToken != "":

		switch nargs := len(os.Args); nargs {
		case 5:
			(&newToken{
				symbol:          os.Args[2],
				issuerSeed:      os.Args[3],
				distributorSeed: os.Args[4]}).issueNew(conf)
		case 6:
			(&newToken{
				symbol:          os.Args[2],
				issuerSeed:      os.Args[3],
				distributorSeed: os.Args[4],
				limit:           os.Args[5]}).issueNew(conf)
		default:
			log.Fatalf("usage: --issue-new-token token issuer-seed distributor-seed [limit]\nthe arguments given are: %+v", os.Args[1:])
		}

	case len(txOptions) > 0:
		parseOptions(txOptions)
	default:
		flag.PrintDefaults()
	}
}
