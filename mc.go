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

	//TODO: refactor signing out to a pluggable func to be able to delegate it to external signers such as hardware wallets
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

func (t *tokenPayment) send(conf *config, txOptions b.SetOptionsBuilder) horizon.Account {

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

	submitTransaction(conf.client, tx, t.From)

	receiver := loadAccount(conf.client, t.To, "... [payment sent]\n")

	return receiver
}

type newToken struct {
	IssuerAddress   string `json:"issuer-address"`
	DistributorSeed string `json:"distributor-seed"`
	Code, Limit     string
}

func (t *newToken) issueNew(conf *config, txOptions b.SetOptionsBuilder) b.Asset {

	distributor := seedToPair(t.DistributorSeed)

	asset := b.CreditAsset(t.Code, t.IssuerAddress)

	var limit = b.MaxLimit

	if t.Limit != "" {
		limit = b.Limit(t.Limit)
	}

	tx, err := b.Transaction(
		b.SourceAccount{distributor.Address()},
		b.AutoSequence{conf.client},
		conf.network,
		b.Trust(t.Code, t.IssuerAddress, limit),
		txOptions,
	)

	if err != nil {
		log.Fatal(err)
	}

	submitTransaction(conf.client, tx, t.DistributorSeed)

	loadAccount(conf.client, distributor.Address(), fmt.Sprintf("issued trust for %s to ", t.Code))

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
	return b.SetOptions(values...)
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
	flag.StringVar(&issueToken, "issue-new-token", "", "issue new token/asset. example (\"limit\" param is optional): --issue-new-token '{\"code\": \"XYZ\", \"issuer-address\": \"address\", \"distributor-seed\":\"seed\", \"limit\": \"42.0\"}'")
	flag.StringVar(&sendPayment, "send-payment", "", "send payment from one account to another. example: --send-payment '{\"from\": \"seed\", \"to\": \"address\", \"token\": \"BTC\", \"amount\": \"42.0\", \"issuer\": \"address\"}'")
	flag.StringVar(&txOptions, "tx-options", "", "add one or more transaction options. example: --tx-options '{\"homeDomain\": \"stellar.org\", \"maxWeight\": 1}'")

	flag.Parse()

	conf := readConfig("/tmp/todo")

	var txOptionsBuilder b.SetOptionsBuilder
	if txOptions != "" {
		txOptionsBuilder = parseOptions(txOptions)
	}

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
		payment.send(conf, txOptionsBuilder)
	case issueToken != "":
		nt := &newToken{}
		if err := json.Unmarshal([]byte(issueToken), nt); err != nil {
			log.Fatal(err)
		}
		nt.issueNew(conf, txOptionsBuilder)
	default:
		if txOptions != "" {
			fmt.Printf("options: %+v", txOptionsBuilder)
		} else {
			flag.PrintDefaults()
		}
	}
}
